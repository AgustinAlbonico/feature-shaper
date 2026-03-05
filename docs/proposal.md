# Proposal: feature-shaper

## Intent

El workflow actual de desarrollo va directo de "idea vaga" a OpenSpec/SDD, saltando la etapa de claridad de negocio. Esto genera especificaciones técnicas con alcance mal definido, criterios de aceptación vagos y flujos alternativos sin contemplar.

`feature-shaper` es una herramienta que actúa como etapa intermedia obligatoria: toma una idea vaga, conduce una conversación adaptativa de **negocio** (no técnica), y produce una `feature definition` completa y estructurada, lista para alimentar OpenSpec/SDD.

El sistema se compone de:
1. Un **binario Go standalone** (`feature-shaper`) con MCP server + DB global SQLite + TUI
2. Un **skill conversacional** (`feature-shaper`) para OpenCode con protocolo de 5.5 fases
3. **Tres commands** (`/shape`, `/shape-refine`, `/shape-catalog`) como puntos de entrada

## Scope

### In Scope
- Binario Go `feature-shaper` con subcomandos `mcp`, `tui`, `migrate`
- Base de datos SQLite global en `~/.feature-shaper/features.db`
- 8 MCP tools expuestas a OpenCode via stdio
- FTS5 para búsqueda semántica de features
- Historial de versiones por feature (snapshots en `featureVersions`)
- Skill `feature-shaper/SKILL.md` con protocolo conversacional de 5.5 fases
- Preguntas de negocio (4 pilares: contexto, alcance, flujos, AC)
- Fase 3.5 de contexto técnico opcional
- Generación de `feature definition` en formato `.md` estandarizado
- TUI con Bubble Tea para explorar el catálogo (vista catálogo, detalle, historial, búsqueda)
- Commands `/shape`, `/shape-refine`, `/shape-catalog`
- Registro del MCP server en `opencode.json`

### Out of Scope
- Integración directa con OpenSpec (el output es un `.md`; OpenSpec lo consume manualmente)
- UI web o desktop (la TUI es suficiente para exploración)
- Autenticación o multiusuario (DB local por usuario)
- Exportación a formatos distintos de `.md`
- Sincronización entre máquinas

## Approach

Binario Go standalone sin runtime externo. La DB es global (`~/.feature-shaper/`) igual que Engram, lo que permite un catálogo centralizado de todos los proyectos del usuario. El MCP server corre en modo stdio para integrarse con OpenCode. La TUI usa el stack Charmbracelet (Bubble Tea + Bubbles + Lip Gloss) para una experiencia rica y colorida.

El skill conduce la conversación en fases con checkpoints explícitos. Las preguntas son de negocio — lo técnico viene después con OpenSpec. Una fase 3.5 opcional permite capturar contexto técnico de alto nivel si el usuario quiere dárselo a OpenSpec como insumo adicional.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
|| `tools/feature-shaper/` | New | Binario Go completo (DB + MCP + TUI) |
| `skills/feature-shaper/SKILL.md` | New | Skill conversacional con protocolo de 5.5 fases |
| `commands/shape.md` | New | Command de entrada para shaping desde cero |
| `commands/shape-refine.md` | New | Command para refinar features existentes |
| `commands/shape-catalog.md` | New | Command para listar el catálogo del proyecto |
|| `~/.feature-shaper/features.db` | New | Base de datos SQLite global |
|| `~/.config/opencode/opencode.json` | Modified | Agregar entrada `feature-shaper` en `mcp` |
| `docs/features/<slug>.md` | New (por feature) | Feature definitions generadas en el repo activo |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| `modernc.org/sqlite` sin soporte FTS5 | Low | Verificar en tests iniciales; FTS5 es parte del build estándar de modernc |
| `mark3labs/mcp-go` cambia API | Low | Pinear versión exacta en go.mod |
| Compilación CGO requerida en Windows | Low | Usar `modernc.org/sqlite` (sin CGO) |
| El skill hace preguntas técnicas en lugar de de negocio | Med | Protocolo explícito con regla "NUNCA técnico en Fase 3" |
| La conversación es demasiado larga (fatiga del usuario) | Med | Máximo 4 preguntas por ronda, checkpoints para salir |

## Rollback Plan

- El binario es standalone — borrarlo del PATH lo desactiva completamente
- Remover la entrada `feature-shaper` de `opencode.json` desactiva el MCP
- La DB en `~/.feature-shaper/` es independiente — se puede borrar sin afectar nada más
- Los files `.md` generados en `docs/features/` son archivos normales — no hay dependencia de runtime

## Dependencies

- Go 1.22+ instalado en el sistema
- `modernc.org/sqlite` — SQLite driver sin CGO
- `github.com/mark3labs/mcp-go` — MCP server SDK
- `github.com/charmbracelet/bubbletea` + `bubbles` + `lipgloss` — TUI stack
- OpenCode con soporte MCP local

## Success Criteria

- [ ] `feature-shaper migrate` crea `~/.feature-shaper/features.db` con schema correcto
- [ ] `feature-shaper mcp` arranca y el MCP server responde por stdio
- [ ] Las 8 MCP tools funcionan (verificado con cliente MCP)
- [ ] `/shape "idea"` conduce la conversación completa y genera el `.md`
- [ ] El `.md` generado se guarda en `docs/features/` y en la DB
- [ ] `/shape-refine "nombre"` carga la feature existente y permite refinarla
- [ ] La versión se incrementa y hay snapshot en `featureVersions`
- [ ] `/shape-catalog` muestra el catálogo del proyecto actual
- [ ] `feature-shaper tui` muestra la vista de catálogo con proyectos y features
- [ ] La TUI responde a todos los keybindings documentados
