## ADDED Requirements

### Requirement: Five-and-a-half phase conversational protocol
El skill SHALL conducir una conversación adaptativa de 5.5 fases con checkpoints explícitos entre cada fase.

#### Scenario: Complete shaping flow
- **WHEN** el usuario invoca /shape con una idea
- **THEN** el skill ejecuta las fases en orden: Exploración del Contexto → Clasificación → Definición Adaptiva → Contexto Técnico (opcional) → Especificación Formal → Persistencia y Cierre

### Requirement: Phase 1 — Context exploration
El skill SHALL explorar el contexto del proyecto leyendo manifest files, README, y features existentes, y luego formular 2-3 preguntas contextuales.

#### Scenario: Auto-detect project context
- **WHEN** el skill arranca la Fase 1
- **THEN** lee package.json/go.mod/pyproject.toml, README.md, y llama feature_catalog para detectar el dominio y slug del proyecto

#### Scenario: Context summary checkpoint
- **WHEN** el skill tiene suficiente contexto
- **THEN** presenta un resumen ("Entendí lo siguiente sobre el proyecto...") y pide confirmación antes de continuar

### Requirement: Phase 2 — Feature classification
El skill SHALL clasificar la feature en dos dimensiones: tipo (product/technical/business) y capas (Frontend/Backend/Database), más un nivel de complejidad (low/medium/high).

#### Scenario: Classification presentation
- **WHEN** el skill completa la clasificación
- **THEN** presenta tipo, complejidad, capas y módulos detectados, y permite al usuario ajustar antes de continuar

#### Scenario: Complexity determines rounds
- **WHEN** la complejidad es low
- **THEN** la Fase 3 tiene 1 ronda de preguntas
- **WHEN** la complejidad es medium
- **THEN** la Fase 3 tiene 2 rondas
- **WHEN** la complejidad es high
- **THEN** la Fase 3 tiene 3 rondas

### Requirement: Phase 3 — Adaptive definition with business questions only
El skill SHALL hacer preguntas exclusivamente de negocio organizadas en 4 pilares: Contexto/Por qué, Alcance, Flujos, Criterios de éxito. NUNCA preguntas técnicas.

#### Scenario: Maximum 4 questions per round
- **WHEN** el skill formula preguntas en cualquier ronda
- **THEN** no excede 4 preguntas, cada una numerada ("Pregunta 2/4:")

#### Scenario: Skip already-answered questions
- **WHEN** una respuesta anterior cubre implícitamente una pregunta planificada
- **THEN** el skill la salta sin presentarla

#### Scenario: Round summary checkpoint
- **WHEN** el skill completa una ronda (excepto la última)
- **THEN** presenta un resumen de lo entendido hasta ahora y pide confirmación

#### Scenario: Never ask technical questions
- **WHEN** el skill está en Fase 3
- **THEN** NUNCA menciona endpoints, DB, componentes, APIs, o detalles de implementación

### Requirement: Phase 3.5 — Optional technical context
El skill SHALL ofrecer una fase opcional de 2-4 preguntas técnicas de alto nivel. El usuario decide si responder o saltear.

#### Scenario: Offer technical context
- **WHEN** la Fase 3 está completa
- **THEN** el skill pregunta "¿Querés agregar contexto técnico de alto nivel?" con opciones de aceptar o saltear

#### Scenario: Skip technical context
- **WHEN** el usuario elige saltear
- **THEN** el skill pasa directo a Fase 4 sin generar la sección "Contexto Técnico" en el .md

### Requirement: Phase 4 — Formal specification generation
El skill SHALL generar el .md completo siguiendo el template estandarizado y presentarlo para aprobación.

#### Scenario: Generate and present
- **WHEN** el skill tiene toda la información de negocio (y técnica opcional)
- **THEN** genera el .md completo con todas las secciones y lo presenta al usuario

#### Scenario: User requests adjustments
- **WHEN** el usuario pide cambios al .md generado
- **THEN** el skill incorpora los ajustes y vuelve a presentar

### Requirement: Phase 5 — Persistence and closure
El skill SHALL persistir la feature llamando a project_register, feature_save, y escribiendo el .md en docs/features/.

#### Scenario: Save and confirm
- **WHEN** el usuario aprueba el .md final
- **THEN** el skill llama project_register, luego feature_save, luego escribe docs/features/<slug>.md, y muestra confirmación con versión y ruta del archivo

### Requirement: Refine flow
El skill SHALL soportar refinamiento de features existentes cargando la versión actual y saltando a Fase 3 con contexto preservado.

#### Scenario: Load existing feature for refinement
- **WHEN** el usuario invoca /shape-refine con un nombre
- **THEN** el skill busca con feature_search, carga con feature_get, presenta resumen de la versión actual, y pregunta qué cambiar

#### Scenario: Version increment on refine
- **WHEN** el usuario aprueba un refinamiento
- **THEN** el skill guarda via feature_save que incrementa la versión y crea snapshot en featureVersions con changelog
