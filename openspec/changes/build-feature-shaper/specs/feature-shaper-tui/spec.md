## ADDED Requirements

### Requirement: TUI catalog view with two panels
La TUI SHALL mostrar una vista de catálogo con panel izquierdo (proyectos) y panel derecho (features del proyecto activo).

#### Scenario: Initial load
- **WHEN** el usuario ejecuta `feature-shaper tui`
el usuario ejecuta `feature-shaper tui`

#### Scenario: Panel navigation with Tab
- **WHEN** el usuario presiona Tab
- **THEN** el foco alterna entre panel izquierdo (proyectos) y panel derecho (features)

#### Scenario: Project selection updates features
- **WHEN** el usuario selecciona un proyecto diferente en el panel izquierdo
- **THEN** el panel derecho se actualiza mostrando las features de ese proyecto

### Requirement: Feature detail view
La TUI SHALL mostrar el contenido completo de una feature en una vista scrollable.

#### Scenario: Open feature detail
- **WHEN** el usuario presiona Enter en una feature del catálogo
- **THEN** la TUI navega a la vista de detalle mostrando el .md completo con scroll

#### Scenario: Scroll navigation
- **WHEN** el usuario está en la vista de detalle
- **THEN** puede hacer scroll con j/k o ↑↓

#### Scenario: Return to catalog
- **WHEN** el usuario presiona b en la vista de detalle
- **THEN** la TUI vuelve a la vista de catálogo

### Requirement: Version history view
La TUI SHALL mostrar el historial de versiones de una feature con changelog.

#### Scenario: Open history
- **WHEN** el usuario presiona h en la vista de catálogo o detalle
- **THEN** la TUI muestra la lista de versiones con número, fecha, y changelog

#### Scenario: View specific version
- **WHEN** el usuario presiona Enter en una versión del historial
- **THEN** la TUI muestra el contenido completo de esa versión

### Requirement: FTS live search
La TUI SHALL ofrecer búsqueda FTS en vivo al presionar `/`.

#### Scenario: Activate search
- **WHEN** el usuario presiona /
- **THEN** el panel derecho se convierte en input de búsqueda con resultados en vivo

#### Scenario: Live results
- **WHEN** el usuario escribe en el input de búsqueda
- **THEN** los resultados FTS se actualizan mientras escribe, abarcando todos los proyectos

#### Scenario: Cancel search
- **WHEN** el usuario presiona Esc durante la búsqueda
- **THEN** la TUI vuelve a la vista de catálogo

### Requirement: Color scheme
La TUI SHALL usar un esquema de colores consistente para tipos y estados.

#### Scenario: Type colors
- **WHEN** se muestra una feature
- **THEN** el tipo usa colores: product=#4FC3F7 (azul), technical=#81C784 (verde), business=#FFB74D (naranja)

#### Scenario: Status colors
- **WHEN** se muestra una feature
- **THEN** el status usa colores: ready=#69F0AE, draft=#90A4AE, in-progress=#FFD54F, done=#A5D6A7

### Requirement: Keybindings
La TUI SHALL responder a los keybindings documentados: ↑↓/jk (navegar), Tab (alternar panel), Enter (seleccionar), / (buscar), Esc (cancelar/volver), b (volver), h (historial), e (exportar), d (eliminar con confirmación), ? (toggle ayuda), q (salir).

#### Scenario: Delete with confirmation
- **WHEN** el usuario presiona d en una feature
- **THEN** la TUI pide confirmación antes de eliminar

#### Scenario: Export to markdown
- **WHEN** el usuario presiona e en una feature
- **THEN** la TUI exporta el .md al directorio docs/features/ del proyecto actual

#### Scenario: Quit
- **WHEN** el usuario presiona q
- **THEN** la TUI se cierra limpiamente
