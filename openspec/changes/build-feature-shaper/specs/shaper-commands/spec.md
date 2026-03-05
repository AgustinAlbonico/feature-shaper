## ADDED Requirements

### Requirement: /shape command
El command `/shape` SHALL cargar el skill feature-shaper y arrancar el proceso de shaping desde Fase 1 con la idea del usuario como argumento.

#### Scenario: Shape with argument
- **WHEN** el usuario ejecuta `/shape "quiero un sistema de notificaciones push"`
- **THEN** el command carga el skill feature-shaper y pasa la idea como contexto inicial

#### Scenario: Shape without argument
- **WHEN** el usuario ejecuta `/shape` sin argumento
- **THEN** el command pregunta "¿Sobre qué feature querés trabajar?"

#### Scenario: Detect similar existing feature
- **WHEN** se detecta una feature similar via feature_search
- **THEN** el command ofrece refinar la existente en lugar de crear una nueva

### Requirement: /shape-refine command
El command `/shape-refine` SHALL cargar el skill feature-shaper y arrancar el proceso de refinamiento desde Fase 3 con la feature existente cargada.

#### Scenario: Refine with argument
- **WHEN** el usuario ejecuta `/shape-refine "invitaciones"`
- **THEN** el command busca la feature, la carga, y arranca el refinamiento

#### Scenario: Refine without argument
- **WHEN** el usuario ejecuta `/shape-refine` sin argumento
- **THEN** el command muestra las últimas 5 features del proyecto actual y pide que el usuario elija

#### Scenario: Multiple matches
- **WHEN** feature_search devuelve más de un resultado
- **THEN** el command presenta la lista para que el usuario elija

### Requirement: /shape-catalog command
El command `/shape-catalog` SHALL mostrar el catálogo de features del proyecto actual usando feature_catalog.

#### Scenario: Catalog without arguments
- **WHEN** el usuario ejecuta `/shape-catalog`
- **THEN** el command detecta el slug del proyecto y muestra las features con status, tipo, versión y fecha

#### Scenario: Catalog with --all flag
- **WHEN** el usuario ejecuta `/shape-catalog --all`
- **THEN** el command lista features de todos los proyectos usando project_list + feature_catalog por proyecto

#### Scenario: Catalog with --status filter
- **WHEN** el usuario ejecuta `/shape-catalog --status ready`
- **THEN** el command muestra solo features con status "ready"

#### Scenario: Catalog with --type filter
- **WHEN** el usuario ejecuta `/shape-catalog --type product`
- **THEN** el command muestra solo features de tipo "product"

#### Scenario: Catalog output format
- **WHEN** se muestran los resultados
- **THEN** el formato usa iconos de estado (✓ ready, ◌ draft, ◎ in-progress) con título, tipo, versión y fecha
