# Implementation Plan: feature-shaper

## Orden de implementación

El orden importa — el binario es dependencia de todo lo demás.

```
Paso 1: feature-shaper — DB + MCP server (sin TUI)        ← prioridad máxima
Paso 2: SKILL.md + 3 commands
Paso 3: Registro en opencode.json
Paso 4: Testing end-to-end
Paso 5: TUI (Bubble Tea + Bubbles + Lip Gloss)
Paso 6: Actualizar README del repo ai-customizations
```

---

## Paso 1: `feature-shaper` — DB + MCP server

### Archivos a crear

```
tools/feature-shaper/
└── cmd/feature-shaper/main.go
├── internal/db/schema.go
├── internal/db/migrations.go
├── internal/db/queries.go
├── internal/mcp/server.go
├── internal/mcp/handlers.go
├── internal/store/features.go
├── internal/store/projects.go
└── go.mod
```

### Dependencias Go

```
modernc.org/sqlite          — SQLite driver SIN CGO (crítico en Windows)
github.com/mark3labs/mcp-go — MCP server SDK
```

### Criterios de completitud

- `feature-shaper migrate` crea `~/.feature-shaper/features.db` correctamente
- `feature-shaper mcp` arranca y responde al protocolo MCP stdio
- Las 8 tools funcionan: `feature_save`, `feature_get`, `feature_search`, `feature_catalog`, `feature_versions`, `feature_get_version`, `project_register`, `project_list`
- `go build ./cmd/feature-shaper/...` pasa sin errores

### Comando de compilación e instalación (PowerShell)

```powershell
Set-Location tools/feature-shaper
go build ./cmd/feature-shaper/...
go install ./cmd/feature-shaper/...
```

### Verificación rápida

```powershell
feature-shaper migrate
# Debe crear ~/.feature-shaper/features.db sin errores
```

---

## Paso 2: SKILL.md + Commands

### Archivos a crear

```
skills/feature-shaper/SKILL.md
commands/shape.md
commands/shape-refine.md
commands/shape-catalog.md
```

### Referencia de formatos

Ver archivos existentes en el repo:
- `skills/interactive-task/SKILL.md` — patrón de SKILL.md
- `commands/task.md` — patrón de command

### Criterios de completitud

- `SKILL.md` tiene frontmatter YAML válido con `name` y `description`
- El protocolo de 5.5 fases está documentado con suficiente detalle para que el agente lo siga
- Los 3 commands tienen frontmatter YAML válido con `description`
- Los commands cargan el skill y definen el comportamiento especial

---

## Paso 3: Registro en opencode.json

Agregar bajo la clave `"mcp"` en `~/.config/opencode/opencode.json`:

```json
"feature-shaper": {
  "type": "local",
  "command": ["feature-shaper", "mcp"],
  "enabled": true
}
```

**Nota**: El binario debe estar en el PATH antes de agregar esta entrada. Si no está en el PATH, usar ruta absoluta en `command`.

---

## Paso 4: Testing end-to-end

Secuencia de prueba completa:

1. `/shape "quiero un sistema de notificaciones push"`
   - Verificar que el agente conduce las fases correctamente
   - Verificar que el `.md` se genera en `docs/features/`
   - Verificar que la feature aparece en `feature_catalog`

2. `/shape-refine "notificaciones"`
   - Verificar que carga la feature existente
   - Hacer algún cambio (agregar un criterio de aceptación)
   - Verificar que `version` pasó de 1 a 2
   - Verificar que hay un snapshot en `featureVersions`

3. `/shape-catalog`
   - Verificar que muestra la feature con status correcto

4. `feature-shaper tui` (si ya está implementado en este paso)
   - Verificar que muestra el proyecto y la feature

---

## Paso 5: TUI

### Archivos a crear

```
tools/feature-shaper/internal/tui/
├── app.go
├── model.go
├── views/catalog.go
├── views/detail.go
├── views/history.go
├── views/search.go
└── styles.go
```

### Dependencias Go adicionales

```
github.com/charmbracelet/bubbletea  v1.x
github.com/charmbracelet/bubbles    v0.x
github.com/charmbracelet/lipgloss   v1.x
```

### Criterios de completitud

- `feature-shaper tui` muestra la vista de catálogo con dos paneles
- Panel izquierdo: lista de proyectos con contador de features
- Panel derecho: lista de features del proyecto activo con tipo, status y versión coloreados
- Tab alterna entre paneles
- Enter en una feature abre la vista de detalle
- `h` abre el historial de versiones
- `/` activa la búsqueda FTS en vivo
- `q` sale limpiamente

---

## Paso 6: Documentación

- Actualizar `README.md` del repo `ai-customizations` con:
  - Nueva sección para `feature-shaper` skill
  - Instrucciones de instalación del binario `feature-shaper`
  - Descripción de los commands `/shape`, `/shape-refine`, `/shape-catalog`
  - Requisitos previos (Go 1.22+)

---

## Estructura final del repo

```
ai-customizations/
├── commands/
│   ├── shape.md              ← nuevo
│   ├── shape-refine.md       ← nuevo
│   ├── shape-catalog.md      ← nuevo
│   └── [existentes...]
├── skills/
│   ├── feature-shaper/       ← nuevo
│   │   └── SKILL.md
│   └── [existentes...]
├── tools/                    ← nuevo directorio
    └── feature-shaper/
│       ├── cmd/
│       ├── internal/
│       └── go.mod
└── README.md                 ← actualizar
```

---

## Variables de entorno del usuario (referencia)

| Variable | Valor |
|---|---|
| OS | Windows 11 |
| Shell | PowerShell (siempre, nunca bash/linux) |
| Repo | `C:\Users\agust\Desktop\Programacion\Proyectos\ai-customizations` |
| Config OpenCode | `C:\Users\agust\.config\opencode\opencode.json` |
| DB Engram (referencia) | `C:\Users\agust\.engram\engram.db` |
|| DB feature-shaper | `C:\Users\agust\.feature-shaper\features.db` |
| Package manager | preferir `bun` sobre `npm` |
| Go | debe estar instalado |
