# feature-shaper

Herramienta que transforma ideas vagas en definiciones de features completas y estructuradas.

## ¿Qué hace?

`feature-shaper` te guía a través de preguntas de **negocio** (no técnicas) para clarificar qué quieres construir, para quién, y cómo se verá el éxito. El resultado es un documento estructurado listo para compartir con tu equipo o usar como base para diseño técnico.

```
/shape "quiero un sistema de notificaciones push"
        ↓
  Conversación adaptativa (5.5 fases)
        ↓
  Feature definition completa en .md
        ↓
  Guardada en DB global + docs/features/<slug>.md
```

## Instalación

### Desde Releases (Recomendado)

Descarga el ejecutable desde [GitHub Releases](https://github.com/AgustinAlbonico/feature-shaper/releases):

```powershell
# Windows: descargar feature-shaper-windows-amd64.zip
# Descomprimir y agregar al PATH
```

### Desde Código Fuente

```powershell
git clone https://github.com/AgustinAlbonico/feature-shaper.git
cd feature-shaper/tools/feature-shaper
go install ./cmd/feature-shaper/...
```

## Configuración

```powershell
# 1. Crear la base de datos
feature-shaper migrate

# 2. Registrar el MCP en opencode.json
# Agregar en ~/.config/opencode/opencode.json:
# "feature-shaper": { "type": "local", "command": ["feature-shaper", "mcp"], "enabled": true }

# 3. Instalar el skill
cd tools/feature-shaper
npx skills add . --skill feature-shaper --agent opencode -y
```

## Uso

```
/shape "mi primera feature"
```

## Interfaces

| Interfaz | Comando | Descripción |
|----------|---------|-------------|
| Conversacional | `/shape "idea"` | Diálogo guiado con OpenCode |
| TUI | `feature-shaper tui` | Interfaz visual en terminal |
| MCP Tools | Directas | Herramientas programáticas |

## Comandos Principales

```
/shape "idea"              # Crear feature desde cero
/shape-refine "slug"       # Refinar feature existente
/shape-catalog             # Ver catálogo del proyecto
/shape-catalog --status draft
/shape-catalog --type product
```

## TUI

```powershell
feature-shaper tui
```

Interfaz visual con:
- Navegación entre proyectos y features
- Vista de detalle con scroll
- Historial de versiones
- Búsqueda full-text
- Exportación a `.md`

**Keybindings:** `↑↓` navegar • `Tab` cambiar panel • `Enter` abrir • `/` buscar • `?` ayuda • `q` salir

## Documentación

| Documento | Descripción |
|-----------|-------------|
| [Guía de Instalación](docs/INSTALLATION.md) | Instalación paso a paso desde cero |
| [Guía de Usuario](docs/USER_GUIDE.md) | Uso completo de todas las funcionalidades |
| [Arquitectura](docs/architecture.md) | Diseño técnico del sistema |
| [Formato de Output](docs/spec-output-format.md) | Estructura del `.md` generado |

## Componentes

| Componente | Descripción |
|------------|-------------|
| `feature-shaper` binario | Servidor MCP + DB SQLite + TUI |
| `feature-shaper` skill | Protocolo conversacional de 5.5 fases |
| `/shape` | Crear feature desde cero |
| `/shape-refine` | Refinar feature existente |
| `/shape-catalog` | Ver catálogo del proyecto |

## Requisitos

- Go 1.22+
- OpenCode con soporte MCP local

## Licencia

MIT
