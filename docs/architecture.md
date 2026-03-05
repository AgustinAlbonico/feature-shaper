# Architecture: feature-shaper

## Visión general

```
/shape "idea vaga"
        ↓
  SKILL.md (protocolo conversacional — 5.5 fases)
        ↓ feature definition completa
  feature-shaper MCP (persistencia global SQLite)
        ↓
  ~/.feature-shaper/features.db
  docs/features/<slug>.md  (en el repo actual)

Exploración manual:
  feature-shaper tui
```

---

## Componentes

### 1. Binario Go `feature-shaper`

Binario standalone con tres subcomandos:

| Subcomando | Propósito |
|---|---|
|| `feature-shaper mcp` | MCP server (stdio) — expone tools a OpenCode |
|| `feature-shaper tui` | TUI interactiva para navegar el catálogo |
|| `feature-shaper migrate` | Crea/migra la DB (auto-invocado al arrancar) |

**Instalación**: `go install ./cmd/feature-shaper` → binario en el PATH

**Stack tecnológico**:
- `modernc.org/sqlite` — SQLite sin CGO (crítico para Windows)
- `github.com/mark3labs/mcp-go` — MCP server SDK stdio
- `github.com/charmbracelet/bubbletea` — framework TUI (Elm Architecture)
- `github.com/charmbracelet/bubbles` — componentes: list, viewport, textinput
- `github.com/charmbracelet/lipgloss` — estilos y colores

---

### 2. Estructura de directorios

```
tools/feature-shaper/
└── cmd/
    └── feature-shaper/
│       └── main.go              ← entry point, registra subcomandos
├── internal/
│   ├── db/
│   │   ├── schema.go            ← definición del schema SQLite
│   │   ├── migrations.go        ← lógica de migrate
│   │   └── queries.go           ← queries tipadas
│   ├── mcp/
│   │   ├── server.go            ← MCP server stdio
│   │   └── handlers.go          ← implementación de cada tool
│   ├── store/
│   │   ├── features.go          ← lógica de negocio de features
│   │   └── projects.go          ← lógica de negocio de proyectos
│   └── tui/
│       ├── app.go               ← entrada a la TUI (Bubble Tea)
│       ├── model.go             ← model principal (Elm Architecture)
│       ├── views/
│       │   ├── catalog.go       ← vista principal (proyectos + features)
│       │   ├── detail.go        ← vista de feature completa
│       │   ├── history.go       ← vista de historial de versiones
│       │   └── search.go        ← búsqueda FTS en vivo
│       └── styles.go            ← colores y estilos Lip Gloss
├── go.mod
└── README.md
```

---

### 3. Base de datos SQLite

**Ubicación**: `~/.feature-shaper/features.db` (global, igual que `~/.engram/engram.db`)

#### Schema completo

```sql
-- Proyectos conocidos (auto-registrado al guardar la primera feature)
CREATE TABLE IF NOT EXISTS projects (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  slug      TEXT UNIQUE NOT NULL,   -- "mi-app", "saas-backend"
  name      TEXT NOT NULL,
  path      TEXT,                   -- directorio raíz del proyecto
  createdAt TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Features (una fila = última versión activa)
CREATE TABLE IF NOT EXISTS features (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  projectSlug    TEXT NOT NULL REFERENCES projects(slug),
  slug           TEXT NOT NULL,          -- "user-auth", "payment-flow"
  title          TEXT NOT NULL,
  type           TEXT NOT NULL,          -- product | technical | business
  status         TEXT NOT NULL DEFAULT 'draft',  -- draft | ready | in-progress | done
  content        TEXT NOT NULL,          -- el .md completo de la feature definition
  version        INTEGER NOT NULL DEFAULT 1,
  topicKey       TEXT UNIQUE,            -- "project-slug/feature-slug" para upserts
  normalizedHash TEXT,                   -- dedup de contenido idéntico
  createdAt      TEXT NOT NULL DEFAULT (datetime('now')),
  updatedAt      TEXT NOT NULL DEFAULT (datetime('now')),
  UNIQUE(projectSlug, slug)
);

-- Historial de versiones (snapshot por cada refinamiento)
CREATE TABLE IF NOT EXISTS featureVersions (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  featureId INTEGER NOT NULL REFERENCES features(id) ON DELETE CASCADE,
  version   INTEGER NOT NULL,
  content   TEXT NOT NULL,   -- snapshot completo del .md en esa versión
  changelog TEXT,            -- qué cambió respecto a la versión anterior
  createdAt TEXT NOT NULL DEFAULT (datetime('now'))
);

-- FTS5 para búsqueda semántica
CREATE VIRTUAL TABLE IF NOT EXISTS featuresFts USING fts5(
  title,
  content,
  type,
  status,
  content='features',
  content_rowid='id'
);

-- Trigger: mantener FTS sincronizado con inserts
CREATE TRIGGER IF NOT EXISTS features_ai AFTER INSERT ON features BEGIN
  INSERT INTO featuresFts(rowid, title, content, type, status)
  VALUES (new.id, new.title, new.content, new.type, new.status);
END;

-- Trigger: mantener FTS sincronizado con updates
CREATE TRIGGER IF NOT EXISTS features_au AFTER UPDATE ON features BEGIN
  INSERT INTO featuresFts(featuresFts, rowid, title, content, type, status)
  VALUES ('delete', old.id, old.title, old.content, old.type, old.status);
  INSERT INTO featuresFts(rowid, title, content, type, status)
  VALUES (new.id, new.title, new.content, new.type, new.status);
END;

-- Trigger: mantener FTS sincronizado con deletes
CREATE TRIGGER IF NOT EXISTS features_ad AFTER DELETE ON features BEGIN
  INSERT INTO featuresFts(featuresFts, rowid, title, content, type, status)
  VALUES ('delete', old.id, old.title, old.content, old.type, old.status);
END;
```

---

### 4. MCP Tools (8 tools)

Todas expuestas vía `feature-shaper mcp` en modo stdio.

| Tool | Parámetros requeridos | Parámetros opcionales | Descripción |
|---|---|---|---|
| `feature_save` | projectSlug, slug, title, type, content | status (default: "draft"), changelog | Upsert por topicKey. Si existe → incrementa versión, guarda snapshot |
| `feature_get` | slug | projectSlug | Recupera feature completa (última versión) |
| `feature_search` | query | projectSlug | Búsqueda FTS5. Devuelve título, slug, proyecto, status, preview |
| `feature_catalog` | projectSlug | status, type | Lista features con filtros opcionales |
| `feature_versions` | slug, projectSlug | — | Lista historial: id, versión, changelog, fecha |
| `feature_get_version` | featureId, version | — | Recupera el .md de una versión específica |
| `project_register` | slug, name | path | Upsert de proyecto |
| `project_list` | — | — | Lista proyectos con conteo de features |

#### Lógica de upsert en `feature_save`

```
topicKey = "{projectSlug}/{slug}"

SI topicKey existe en features:
  1. Guardar snapshot actual en featureVersions (version = current_version, content = current_content)
  2. UPDATE features SET version = version + 1, content = nuevo_content, updatedAt = now(), changelog = nuevo_changelog
SINO:
  1. Asegurar que projects tiene el projectSlug (auto-register si no existe)
  2. INSERT INTO features con version = 1
```

---

### 5. Registro en opencode.json

```json
"feature-shaper": {
  "type": "local",
  "command": ["feature-shaper", "mcp"],
  "enabled": true
}
```

Referencia del patrón existente (entrada Engram):
```json
"engram": {
  "type": "local",
  "command": ["engram", "mcp", "--tools=agent"],
  "enabled": true
}
```

---

### 6. TUI

**Punto de entrada**: `feature-shaper tui` (standalone).

#### Vista Principal — catálogo (dos paneles)

```
┌─────────────────────────────────────────────────────────────────────┐
│  ◆ feature-shaper                              [/] buscar  [?] ayuda │
├──────────────────┬──────────────────────────────────────────────────┤
│  PROYECTOS       │  FEATURES — mi-saas-app                          │
│  ▶ mi-saas-app 3 │  ● user-auth      product   ✓ ready    v3        │
│    backend-api 7 │  ● payment-flow   business  ◌ draft    v1        │
│    mobile-app  2 │  ● notifications  technical ◎ in-progress v2     │
│                  │                                                   │
├──────────────────┴──────────────────────────────────────────────────┤
│  [enter] abrir  [d] eliminar  [e] exportar  [h] historial  [q] salir│
└─────────────────────────────────────────────────────────────────────┘
```

#### Vista de Feature (detalle)

```
┌─────────────────────────────────────────────────────────────────────┐
│  ◆ user-auth  [product] [ready]  v3                      [b] volver │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  # Sistema de Autenticación de Usuarios                              │
│  ...contenido scrollable con j/k o ↑↓...                            │
│                                                                      │
├─────────────────────────────────────────────────────────────────────┤
│  [h] historial  [e] exportar .md  [d] eliminar  [b] volver          │
└─────────────────────────────────────────────────────────────────────┘
```

#### Vista de Historial

```
┌─────────────────────────────────────────────────────────────────────┐
│  ◆ user-auth — Historial                               [b] volver   │
├─────────────────────────────────────────────────────────────────────┤
│  v3  2026-03-04  Agregado 2FA y flujo de recovery                    │
│  v2  2026-02-28  Refinado criterios de aceptación                    │
│  v1  2026-02-20  Definición inicial                                  │
├─────────────────────────────────────────────────────────────────────┤
│  [enter] ver versión  [b] volver                                     │
└─────────────────────────────────────────────────────────────────────┘
```

#### Búsqueda en vivo

Al presionar `/`, el panel derecho se convierte en input FTS5. Resultados se actualizan mientras se escribe, abarcan todos los proyectos. `Esc` cancela.

#### Esquema de colores

| Elemento | Color |
|---|---|
| `product` | Azul `#4FC3F7` |
| `technical` | Verde `#81C784` |
| `business` | Naranja `#FFB74D` |
| `ready` | Verde brillante `#69F0AE` |
| `draft` | Gris `#90A4AE` |
| `in-progress` | Amarillo `#FFD54F` |
| `done` | Verde apagado `#A5D6A7` |
| Proyecto activo | Highlight violeta `#CE93D8` |
| Border activo | Violeta `#9C27B0` |

#### Keybindings

| Tecla | Acción |
|---|---|
| `↑↓` / `j k` | Navegar lista |
| `Tab` | Alternar panel izquierdo/derecho |
| `Enter` | Seleccionar / abrir |
| `/` | Activar búsqueda FTS en vivo |
| `Esc` | Cancelar búsqueda / volver |
| `b` | Volver a vista anterior |
| `h` | Ver historial de versiones |
| `e` | Exportar feature a .md en proyecto actual |
| `d` | Eliminar feature (con confirmación) |
| `?` | Toggle ayuda |
| `q` | Salir |

---

### 7. Decisiones de diseño

| Decisión | Alternativa descartada | Razón |
|---|---|---|
| Binario Go | Node.js / Python | Binario standalone sin runtime externo |
|| DB global `~/.feature-shaper/` | DB por proyecto | Catálogo centralizado, misma arquitectura que Engram |
| Sistema independiente | Usar Engram como backend | Dominio diferente, schema específico, no contaminar Engram |
| TUI standalone | TUI integrada al workflow | El workflow ya tiene checkpoints conversacionales; la TUI es para exploración |
| Preguntas de negocio únicamente en Fase 3 | Preguntas técnicas mezcladas | Lo técnico viene después con OpenSpec |
| Fase 3.5 opcional | Técnico siempre o nunca | Flexible según el estado mental del usuario |
| `modernc.org/sqlite` | `mattn/go-sqlite3` | Evita CGO — compilación simple en Windows sin gcc |
| MCP stdio | MCP HTTP | OpenCode requiere stdio para MCPs locales |
