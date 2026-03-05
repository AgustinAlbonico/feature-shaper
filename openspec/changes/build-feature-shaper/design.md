## Context

El desarrollo de features hoy va directo de idea vaga a diseño técnico (OpenSpec/SDD), sin una etapa estructurada de claridad de negocio. No existe código previo — este es un proyecto greenfield que introduce una herramienta nueva en el workflow del desarrollador.

El sistema se integra en un ecosistema existente: OpenCode como IDE con soporte MCP, y un patrón ya probado de binarios Go con DB global (referencia: Engram con `~/.engram/engram.db`). La arquitectura sigue este patrón validado.

**Restricciones del entorno**:
- Windows 11 como plataforma principal → sin CGO disponible fácilmente
- OpenCode requiere MCP servers en modo stdio (no HTTP)
- El skill y commands son archivos Markdown que el agente de OpenCode interpreta

## Goals / Non-Goals

**Goals:**
- Binario Go standalone sin runtime externo, instalable con `go install`
- Base de datos SQLite global que centralice features de todos los proyectos del usuario
- Servidor MCP stdio con 8 tools completas para CRUD, búsqueda y versionado
- Skill conversacional que conduzca preguntas de negocio en 5.5 fases con checkpoints
- TUI rica con Bubble Tea para exploración del catálogo
- Feature definitions generadas en formato `.md` estandarizado

**Non-Goals:**
- UI web o desktop (la TUI es suficiente)
- Integración directa con OpenSpec (el output es `.md`, OpenSpec lo consume manualmente)
- Autenticación o multiusuario
- Sincronización entre máquinas
- Exportación a formatos distintos de `.md`

## Decisions

### 1. Go como lenguaje del binario
**Decisión**: Usar Go 1.22+ para `feature-shaper`.
**Alternativa descartada**: Node.js, Python.
**Razón**: Binario standalone sin runtime externo. Compilación estática, distribución simple, patrón ya validado con Engram.

### 2. `modernc.org/sqlite` como driver de SQLite
**Decisión**: Usar la implementación pure-Go de SQLite.
**Alternativa descartada**: `mattn/go-sqlite3` (requiere CGO).
**Razón**: Evita CGO, compilación simple en Windows sin gcc. FTS5 está incluido en el build estándar de modernc.

Base de datos global en `~/.feature-shaper/`
**Decisión**: DB centralizada fuera de cualquier proyecto, en el home del usuario.
**Alternativa descartada**: DB por proyecto (dentro de cada repo).
**Razón**: Catálogo centralizado multi-proyecto. Mismo patrón que Engram (`~/.engram/engram.db`). Permite buscar y comparar features entre proyectos.

### 4. Sistema independiente de Engram
**Decisión**: DB y binario separados de Engram.
**Alternativa descartada**: Usar Engram como backend de persistencia.
**Razón**: Dominio diferente, schema específico con FTS5/versionado/triggers. No contaminar el scope de Engram.

### 5. MCP stdio (no HTTP)
**Decisión**: Servidor MCP en modo stdio.
**Alternativa descartada**: Servidor HTTP.
**Razón**: OpenCode requiere modo stdio para MCPs locales. El MCP se registra en `opencode.json` como proceso local.

### 6. Upsert con versionado automático
**Decisión**: `feature_save` detecta si la feature existe (via `topicKey`) y si existe, guarda un snapshot en `featureVersions` antes de actualizar.
**Razón**: Historial completo sin esfuerzo del usuario. El changelog se captura en cada refinamiento.

### 7. FTS5 para búsqueda
**Decisión**: Virtual table FTS5 con triggers de sincronización automática.
**Alternativa descartada**: LIKE queries manuales.
**Razón**: Búsqueda semántica rápida, ranking por relevancia, soporte nativo en SQLite.

### 8. TUI con Charmbracelet stack
**Decisión**: Bubble Tea (Elm Architecture) + Bubbles (componentes) + Lip Gloss (estilos).
**Alternativa descartada**: TUI frameworks más simples (como Cobra con tablas estáticas).
**Razón**: Experiencia rica e interactiva. Navegación por teclado, búsqueda en vivo, paneles navegables.

### 9. Preguntas de negocio (nunca técnicas) en Fase 3
**Decisión**: El skill hace preguntas exclusivamente de negocio durante la definición. Lo técnico es opcional en Fase 3.5.
**Razón**: Lo técnico viene después con OpenSpec. Mezclar genera features con alcance mal definido — el problema original que se resuelve.

### 10. Estructura de directorios del binario
```
tools/feature-shaper/
├── cmd/feature-shaper/main.go          ← entry point, parsea subcomandos
├── internal/
│   ├── db/
│   │   ├── schema.go                  ← SQL DDL como constantes
│   │   ├── migrations.go              ← lógica de migrate
│   │   └── queries.go                 ← queries tipadas
│   ├── mcp/
│   │   ├── server.go                  ← MCP server stdio
│   │   └── handlers.go                ← handler por tool
│   ├── store/
│   │   ├── features.go                ← lógica de negocio features
│   │   └── projects.go                ← lógica de negocio proyectos
│   └── tui/
│       ├── app.go                     ← entrada TUI
│       ├── model.go                   ← model Elm Architecture
│       ├── views/
│       │   ├── catalog.go             ← vista catálogo (dos paneles)
│       │   ├── detail.go              ← vista detalle feature
│       │   ├── history.go             ← vista historial
│       │   └── search.go              ← búsqueda FTS en vivo
│       └── styles.go                  ← colores y estilos Lip Gloss
└── go.mod
```

## Risks / Trade-offs

| Risk | Likelihood | Mitigation |
|---|---|---|
| `modernc.org/sqlite` sin soporte FTS5 | Low | FTS5 es parte del build estándar de modernc. Verificar en tests iniciales |
| `mark3labs/mcp-go` cambia API | Low | Pinear versión exacta en `go.mod` |
| El skill hace preguntas técnicas por error | Medium | Protocolo explícito con regla "NUNCA técnico en Fase 3", solo en Fase 3.5 opcional |
| Conversación demasiado larga | Medium | Máximo 4 preguntas por ronda, checkpoints para confirmar/salir |
| Compilación en Windows con paths largos | Low | Usar paths relativos cortos en la estructura del proyecto |
| TUI no renderiza bien en ciertos terminales | Low | Usar el stack Charmbracelet que tiene amplio soporte multiplataforma |
