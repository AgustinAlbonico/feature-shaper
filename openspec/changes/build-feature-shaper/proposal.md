## Why

El workflow actual de desarrollo va directo de "idea vaga" a OpenSpec/SDD, saltando la etapa de claridad de negocio. Esto genera especificaciones técnicas con alcance mal definido, criterios de aceptación vagos y flujos alternativos sin contemplar. Se necesita una herramienta que fuerce la claridad de negocio antes de arrancar el diseño técnico.

## What Changes

Nuevo binario Go `feature-shaper`": servidor MCP stdio + base de datos SQLite global + TUI interactiva para explorar el catálogo de features
- **Nuevo skill `feature-shaper`**: protocolo conversacional de 5.5 fases que conduce preguntas de negocio (nunca técnicas) para transformar ideas vagas en definiciones completas
- **3 nuevos commands**: `/shape` (crear feature desde cero), `/shape-refine` (refinar existente), `/shape-catalog` (listar catálogo)
- **8 MCP tools** expuestas por el binario: `feature_save`, `feature_get`, `feature_search`, `feature_catalog`, `feature_versions`, `feature_get_version`, `project_register`, `project_list`
Nuevo skill `feature-shaper`": protocolo conversacional de 5.5 fases que conduce preguntas de negocio (nunca técnicas) para transformar ideas vagas en definiciones completas
- **Formato estandarizado `.md`** para feature definitions generadas, con secciones de Contexto, Alcance, Flujos, Criterios de Éxito y Contexto Técnico opcional
- **Registro del MCP** en `opencode.json` para integración con OpenCode

## Capabilities

### New Capabilities
`feature-shaper-db`: Schema SQLite con 3 tablas (projects, features, featureVersions), FTS5, triggers de sincronización, lógica de upsert con versionado automático
`feature-shaper-mcp`: Servidor MCP stdio con 8 tools para CRUD de features y proyectos, búsqueda FTS5, y gestión de versiones
- `feature-shaper-tui`: TUI interactiva con Bubble Tea — vista catálogo (dos paneles), detalle de feature, historial de versiones, búsqueda FTS en vivo
- `shaper-skill`: Skill conversacional con protocolo de 5.5 fases — exploración de contexto, clasificación, definición adaptiva con banco de preguntas por pilares, contexto técnico opcional, especificación formal, persistencia
- `shaper-commands`: Tres commands de OpenCode (/shape, /shape-refine, /shape-catalog) como puntos de entrada al skill
- `feature-output-format`: Formato estandarizado del .md generado con template, reglas de generación por sección, y ejemplo de referencia

### Modified Capabilities
<!-- No hay capabilities existentes — proyecto greenfield -->

## Impact

Nuevo directorio `tools/feature-shaper/` con código Go completo (~15 archivos)
- **Nuevo directorio** `skills/feature-shaper/` con SKILL.md
- **Nuevos archivos** `commands/shape.md`, `commands/shape-refine.md`, `commands/shape-catalog.md`
- **Archivos generados** por el workflow: `docs/features/<slug>.md` en cada proyecto del usuario
Dependencias Go
Configuración": Nueva entrada `feature-shaper` en `~/.config/opencode/opencode.json` bajo `mcp`
- **Sistema de archivos global**: Nuevo directorio `~/.feature-shaper/` con `features.db`
