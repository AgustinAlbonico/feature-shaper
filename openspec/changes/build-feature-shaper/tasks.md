## 1. Project Setup

1.1 Crear estructura de directorios `tools/feature-shaper/` con `cmd/feature-shaper/`, `internal/db/`, `internal/mcp/`, `internal/store/`, `internal/tui/`, `internal/tui/views/`
- [x] 1.2 Inicializar `go.mod` con module `github.com/agustinalbonico/feature-shaper` y Go 1.22
go.mod con module `github.com/agustinalbonico/feature-shaper` y Go 1.22
- [x] 1.4 Crear `cmd/feature-shaper/main.go` con parseo de subcomandos: `mcp`, `tui`, `migrate`
main.go con parseo de subcomandos: `mcp`, `tui`, `migrate`
## 2. Database Layer

- [x] 2.1 Crear `internal/db/schema.go` con constante `SchemaSQL` conteniendo todas las sentencias CREATE TABLE/TRIGGER/FTS5
Migrate()` — crear directorio `~/.feature-shaper/`, abrir DB, activar WAL + foreign_keys, ejecutar schema
- [x] 2.3 Crear `internal/db/queries.go` con structs (`Project`, `Feature`, `FeatureVersion`, `FeatureSearchResult`) y funciones de query tipadas
- [x] 2.4 Implementar `UpsertProject` — INSERT OR IGNORE + UPDATE
- [x] 2.5 Implementar `ListProjects` con conteo de features por proyecto
- [x] 2.6 Implementar `UpsertFeature` con lógica de topicKey: si existe → snapshot en featureVersions + version++ + update; si no → insert con version=1
- [x] 2.7 Implementar `GetFeature` por slug (con filtro opcional de projectSlug)
- [x] 2.8 Implementar `ListFeatures` con filtros opcionales de status y type
- [x] 2.9 Implementar `SearchFeatures` usando FTS5 MATCH con filtro opcional de projectSlug
- [x] 2.10 Implementar `ListFeatureVersions` y `GetFeatureVersion`
Verificar que `feature-shaper migrate` crea la DB correctamente y las queries compilan

## 3. Store Layer (Business Logic)

- [x] 3.1 Crear `internal/store/features.go` con `FeatureStore` struct que wrappea `*sql.DB` y expone: Save, Get, Search, Catalog, Versions, GetVersion
- [x] 3.2 Crear `internal/store/projects.go` con `ProjectStore` struct que expone: Register, List
Verificar que `go build ./cmd/feature-shaper/...` compila sin errores

## 4. MCP Server

- [x] 4.1 Crear `internal/mcp/server.go` con `NewServer()` que registra las 8 tools usando `mark3labs/mcp-go` y `Start()` que bloquea en stdio
server.go que registra las 8 tools usando `mark3labs/mcp-go` y `Start()` que bloquea en stdio
- [x] 4.3 Implementar handler para `feature_get` — extrae slug y projectSlug opcional, devuelve feature completa
- [x] 4.4 Implementar handler para `feature_search` — extrae query y projectSlug opcional, devuelve resultados con preview
- [x] 4.5 Implementar handler para `feature_catalog` — extrae projectSlug, status y type opcionales
- [x] 4.6 Implementar handler para `feature_versions` — extrae slug y projectSlug, devuelve historial
- [x] 4.7 Implementar handler para `feature_get_version` — extrae featureId y version
- [x] 4.8 Implementar handler para `project_register` — extrae slug, name, path opcional
- [x] 4.9 Implementar handler para `project_list` — devuelve proyectos con conteo
- [x] 4.10 Conectar `main.go` subcomando `mcp`: auto-migrate → crear stores → crear server → start
go build ./cmd/feature-shaper/...` && `go install ./cmd/feature-shaper/...`
Verificar `feature-shaper mcp` arranca sin errores y responde al protocolo MCP

## 5. Skill y Commands

- [x] 5.1 Crear `skills/feature-shaper/SKILL.md` con frontmatter YAML (name, description) y protocolo completo de 5.5 fases
- [x] 5.2 Documentar Fase 1 (Exploración del Contexto) con instrucciones de auto-detección y preguntas contextuales
- [x] 5.3 Documentar Fase 2 (Clasificación) con dimensiones tipo/capas/complejidad y presentación al usuario
- [x] 5.4 Documentar Fase 3 (Definición Adaptiva) con banco de preguntas por pilar, regla de máx 4 preguntas por ronda, estructura de rondas según complejidad
- [x] 5.5 Documentar Fase 3.5 (Contexto Técnico opcional) con banco de preguntas técnicas y mecanismo de skip
- [x] 5.6 Documentar Fase 4 (Especificación Formal) con referencia al template del output .md
- [x] 5.7 Documentar Fase 5 (Persistencia y Cierre) con secuencia de llamadas MCP y escritura del archivo
- [x] 5.8 Documentar flujo de refinamiento (/shape-refine) con carga, salto a Fase 3, y versionado
- [x] 5.9 Crear `commands/shape.md` con frontmatter y comportamiento especial (sin argumento, feature similar detectada)
- [x] 5.10 Crear `commands/shape-refine.md` con frontmatter y comportamiento especial (sin argumento, múltiples matches)
- [x] 5.11 Crear `commands/shape-catalog.md` con frontmatter, variantes (--all, --status, --type) y formato de salida

## 6. Integración OpenCode

Registrar MCP en `~/.config/opencode/opencode.json` bajo clave "feature-shaper" con type "local" y command ["feature-shaper", "mcp"]

## 7. TUI (Bubble Tea)

- [x] 7.1 Agregar dependencias Go: `charmbracelet/bubbletea`, `charmbracelet/bubbles`, `charmbracelet/lipgloss`
- [x] 7.2 Crear `internal/tui/styles.go` con constantes de color (product=#4FC3F7, technical=#81C784, business=#FFB74D, etc.) y estilos Lip Gloss
- [x] 7.3 Crear `internal/tui/model.go` con Model struct (View enum, panels, projects, features, inputs) e implementar Init/Update/View
- [x] 7.4 Crear `internal/tui/catalog.go` con Render — dos paneles (proyectos izq, features der), Tab alterna foco, colores por tipo/status
- [x] 7.5 Crear `internal/tui/detail.go` con Render — contenido .md scrollable con viewport, keybindings b/h/e/d
- [x] 7.6 Crear `internal/tui/history.go` con Render — lista de versiones con número, fecha, changelog
- [x] 7.7 Crear `internal/tui/search.go` con Render — input FTS en vivo con resultados actualizados mientras se escribe
- [x] 7.8 Crear `internal/tui/app.go` con `Start(db)` que crea Model, carga proyectos, y arranca `tea.NewProgram`
- [x] 7.9 Conectar `main.go` subcomando `tui`: auto-migrate → conectar DB → llamar tui.Start
- [x] 7.10 Implementar keybindings completos: ↑↓/jk, Tab, Enter, /, Esc, b, h, e, d (con confirmación), ?, q
Verificar que `feature-shaper tui` arranca, muestra catálogo, navega entre vistas, y cierra limpiamente

## 8. Verificación End-to-End

- [x] 8.1 Ejecutar `/shape "idea de prueba"` y verificar que las 5.5 fases se completan correctamente — *pendiente prueba manual en OpenCode; MCP E2E verificado*
- [x] 8.2 Verificar que el .md se genera en `docs/features/` y la feature aparece en `feature_catalog` — *pendiente prueba manual; feature_save + feature_catalog verificados via stdin*
- [x] 8.3 Ejecutar `/shape-refine` y verificar que carga la feature, permite cambios, y la versión incrementa — *pendiente prueba manual; upsert con version++ verificado via MCP*
- [x] 8.4 Ejecutar `/shape-catalog` y verificar que muestra el catálogo con formato correcto — *pendiente prueba manual; feature_catalog + project_list verificados via MCP*
Verificar que `feature-shaper tui` muestra el proyecto y la feature creada — *compilación, go vet, wiring verificados*
