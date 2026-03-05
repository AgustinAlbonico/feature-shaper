# Spec: feature-shaper (Binario Go)

## Módulo: `cmd/feature-shaper/main.go`

Entry point. Parsea subcomandos y delega.

```
Subcomandos:
  mcp      Arranca el MCP server en modo stdio
  tui      Arranca la TUI interactiva
  migrate  Crea o migra la base de datos

Comportamiento de `mcp`:
  1. Llamar migrate() automáticamente
  2. Crear instancia del store
  3. Registrar las 8 MCP tools
  4. Llamar server.Start() — bloquea en stdio

Comportamiento de `migrate`:
  1. Obtener home dir con os.UserHomeDir()
 2. Crear directorio ~/.feature-shaper/ si no existe
 3. Abrir/crear ~/.feature-shaper/features.db
  4. Ejecutar todas las sentencias CREATE TABLE/TRIGGER IF NOT EXISTS
  5. Cerrar conexión y salir con código 0

Comportamiento de `tui`:
  1. Llamar migrate() automáticamente
  2. Conectar a la DB
  3. Llamar tui.Start(db)
```

---

## Módulo: `internal/db/schema.go`

Contiene las sentencias SQL del schema como constantes string. No ejecuta nada — solo define.

```go
const SchemaSQL = `
  CREATE TABLE IF NOT EXISTS projects ( ... );
  CREATE TABLE IF NOT EXISTS features ( ... );
  CREATE TABLE IF NOT EXISTS featureVersions ( ... );
  CREATE VIRTUAL TABLE IF NOT EXISTS featuresFts USING fts5( ... );
  CREATE TRIGGER IF NOT EXISTS features_ai ...;
  CREATE TRIGGER IF NOT EXISTS features_au ...;
  CREATE TRIGGER IF NOT EXISTS features_ad ...;
`
```

---

## Módulo: `internal/db/migrations.go`

```go
// DBPath devuelve la ruta absoluta a la DB según el home del usuario
func DBPath() (string, error)

// Migrate abre (o crea) la DB y ejecuta el schema
func Migrate() (*sql.DB, error)
```

`Migrate()`:
1. Llama `DBPath()` para obtener la ruta
2. Crea el directorio padre si no existe (`os.MkdirAll`)
3. Abre la DB con el driver `modernc.org/sqlite`
4. Ejecuta `pragma journal_mode=WAL` y `pragma foreign_keys=ON`
5. Ejecuta `SchemaSQL` (las sentencias son idempotentes por `IF NOT EXISTS`)
6. Devuelve `*sql.DB` lista para usar

---

## Módulo: `internal/db/queries.go`

Queries SQL tipadas. Cada función recibe `*sql.DB` y devuelve structs o error.

```go
// Tipos de datos
type Project struct {
  ID        int64
  Slug      string
  Name      string
  Path      string
  CreatedAt string
}

type Feature struct {
  ID             int64
  ProjectSlug    string
  Slug           string
  Title          string
  Type           string
  Status         string
  Content        string
  Version        int
  TopicKey       string
  NormalizedHash string
  CreatedAt      string
  UpdatedAt      string
}

type FeatureVersion struct {
  ID        int64
  FeatureID int64
  Version   int
  Content   string
  Changelog string
  CreatedAt string
}

type FeatureSearchResult struct {
  ID          int64
  ProjectSlug string
  Slug        string
  Title       string
  Type        string
  Status      string
  Version     int
  Preview     string  // primeros 200 chars del content
  UpdatedAt   string
}
```

Funciones:
```go
func UpsertProject(db *sql.DB, slug, name, path string) error
func ListProjects(db *sql.DB) ([]Project, error)
func GetProjectFeatureCount(db *sql.DB, slug string) (int, error)

func UpsertFeature(db *sql.DB, projectSlug, slug, title, typ, content, status, changelog string) (*Feature, error)
func GetFeature(db *sql.DB, slug, projectSlug string) (*Feature, error)
func ListFeatures(db *sql.DB, projectSlug, status, typ string) ([]Feature, error)
func SearchFeatures(db *sql.DB, query, projectSlug string) ([]FeatureSearchResult, error)

func ListFeatureVersions(db *sql.DB, featureID int64) ([]FeatureVersion, error)
func GetFeatureVersion(db *sql.DB, featureID int64, version int) (*FeatureVersion, error)
```

Lógica de `UpsertFeature`:
```
topicKey = projectSlug + "/" + slug

SI topicKey ya existe en features:
  1. INSERT INTO featureVersions (featureId, version, content, changelog) — snapshot del estado actual
  2. UPDATE features SET version = version+1, content = newContent, status = newStatus,
     updatedAt = datetime('now'), normalizedHash = hash(newContent)
     WHERE topicKey = topicKey
SINO:
  1. Asegurar que el proyecto existe (INSERT OR IGNORE en projects)
  2. INSERT INTO features con version=1, topicKey, normalizedHash
```

---

## Módulo: `internal/store/features.go`

Capa de negocio sobre `db/queries.go`. Aquí va lógica que no es SQL pura.

```go
type FeatureStore struct {
  db *sql.DB
}

func NewFeatureStore(db *sql.DB) *FeatureStore

func (s *FeatureStore) Save(projectSlug, slug, title, typ, content, status, changelog string) (*Feature, error)
func (s *FeatureStore) Get(slug, projectSlug string) (*Feature, error)
func (s *FeatureStore) Search(query, projectSlug string) ([]FeatureSearchResult, error)
func (s *FeatureStore) Catalog(projectSlug, status, typ string) ([]Feature, error)
func (s *FeatureStore) Versions(slug, projectSlug string) ([]FeatureVersion, error)
func (s *FeatureStore) GetVersion(featureID int64, version int) (*FeatureVersion, error)
```

---

## Módulo: `internal/store/projects.go`

```go
type ProjectStore struct {
  db *sql.DB
}

func NewProjectStore(db *sql.DB) *ProjectStore

func (s *ProjectStore) Register(slug, name, path string) error
func (s *ProjectStore) List() ([]ProjectWithCount, error)
// ProjectWithCount = Project + FeatureCount int
```

---

## Módulo: `internal/mcp/server.go`

```go
// NewServer crea el MCP server con todas las tools registradas
func NewServer(features *store.FeatureStore, projects *store.ProjectStore) *mcp.Server

// Start arranca el server en modo stdio (bloquea)
func Start(server *mcp.Server) error
```

Patrón de registro de tool con `mark3labs/mcp-go`:
```go
server.AddTool(mcp.NewTool("feature_save",
  mcp.WithDescription("Guarda o actualiza una feature definition"),
  mcp.WithString("projectSlug", mcp.Required(), mcp.Description("...")),
  mcp.WithString("slug",        mcp.Required(), mcp.Description("...")),
  // ...
), handlers.FeatureSave(features))
```

---

## Módulo: `internal/mcp/handlers.go`

Una función handler por tool. Cada handler:
1. Extrae parámetros del `mcp.CallToolRequest`
2. Llama al store correspondiente
3. Serializa resultado a JSON y devuelve `mcp.CallToolResult`

```go
func FeatureSave(s *store.FeatureStore) mcp.ToolHandlerFunc
func FeatureGet(s *store.FeatureStore) mcp.ToolHandlerFunc
func FeatureSearch(s *store.FeatureStore) mcp.ToolHandlerFunc
func FeatureCatalog(s *store.FeatureStore) mcp.ToolHandlerFunc
func FeatureVersions(s *store.FeatureStore) mcp.ToolHandlerFunc
func FeatureGetVersion(s *store.FeatureStore) mcp.ToolHandlerFunc
func ProjectRegister(s *store.ProjectStore) mcp.ToolHandlerFunc
func ProjectList(s *store.ProjectStore) mcp.ToolHandlerFunc
```

---

## Módulo: `internal/tui/`

### `styles.go`

Define todas las constantes de color y estilos Lip Gloss. No contiene lógica.

```go
var (
  ColorProduct    = lipgloss.Color("#4FC3F7")
  ColorTechnical  = lipgloss.Color("#81C784")
  ColorBusiness   = lipgloss.Color("#FFB74D")
  ColorReady      = lipgloss.Color("#69F0AE")
  ColorDraft      = lipgloss.Color("#90A4AE")
  ColorInProgress = lipgloss.Color("#FFD54F")
  ColorDone       = lipgloss.Color("#A5D6A7")
  ColorActive     = lipgloss.Color("#CE93D8")
  ColorBorder     = lipgloss.Color("#9C27B0")

  StyleTitle      = lipgloss.NewStyle().Bold(true).Foreground(ColorActive)
  StyleBorderActive = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(ColorBorder)
  StyleBorderNormal = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#444"))
  // ...
)
```

### `model.go`

Model principal (Elm Architecture de Bubble Tea).

```go
type View int
const (
  ViewCatalog View = iota
  ViewDetail
  ViewHistory
  ViewSearch
)

type Model struct {
  db           *sql.DB
  activeView   View
  activePanel  int  // 0=proyectos, 1=features
  projects     []store.ProjectWithCount
  features     []db.Feature
  selectedProj int
  selectedFeat int
  searchInput  textinput.Model
  viewport     viewport.Model
  // ...
}

func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m Model) View() string
```

### `app.go`

```go
func Start(database *sql.DB) error
// Crea el Model inicial, carga proyectos, y arranca tea.NewProgram(model).Run()
```

### `views/catalog.go`, `views/detail.go`, `views/history.go`, `views/search.go`

Cada archivo expone una función `Render(m Model) string` que devuelve el string listo para imprimir. El `View()` del model delega aquí según `m.activeView`.

---

## go.mod

```
module github.com/agustinalbonico/feature-shaper

go 1.22

require (
  modernc.org/sqlite          v1.x.x
  github.com/mark3labs/mcp-go v0.x.x
  github.com/charmbracelet/bubbletea  v1.x.x
  github.com/charmbracelet/bubbles    v0.x.x
  github.com/charmbracelet/lipgloss   v1.x.x
)
```
