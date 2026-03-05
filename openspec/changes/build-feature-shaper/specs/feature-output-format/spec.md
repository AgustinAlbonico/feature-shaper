## ADDED Requirements

### Requirement: Standardized feature definition template
El output `.md` generado SHALL seguir un template estandarizado con secciones: metadata, Contexto, Alcance, Flujos, Criterios de Éxito, Contexto Técnico (opcional), y Notas de Negocio (opcional).

#### Scenario: Complete template structure
- **WHEN** el skill genera un .md
- **THEN** incluye header con metadata (Proyecto, Tipo, Status, Versión, fechas), y secciones Contexto, Alcance (incluido/excluido/relacionadas), Flujos (principal/alternativos/error), Criterios de Éxito (AC/críticos/edge cases)

#### Scenario: Metadata header format
- **WHEN** se genera el header
- **THEN** usa formato blockquote con campos: Proyecto, Tipo (product|technical|business), Status (draft|ready), Versión, Creado (YYYY-MM-DD), Actualizado (YYYY-MM-DD)

### Requirement: Context section rules
La sección Contexto SHALL tener máximo 3 párrafos cortos y cero soluciones técnicas. MUST incluir usuario principal, actores secundarios (si aplica), y situación actual.

#### Scenario: No technical solutions in context
- **WHEN** el skill genera la sección Contexto
- **THEN** describe el problema y usuarios sin mencionar implementación técnica

### Requirement: Scope section with mandatory exclusions
La sección Alcance SHALL incluir items incluidos, items fuera de alcance (mínimo 2), y features relacionadas (si aplica).

#### Scenario: Agent proposes exclusions
- **WHEN** el usuario no mencionó exclusiones explícitamente
- **THEN** el skill propone al menos 2 items fuera de alcance derivados de lo que quedó implícito en la conversación

### Requirement: Flows section
La sección Flujos SHALL documentar el flujo principal como pasos numerados. Flujos alternativos y de error solo se incluyen si surgieron en el shaping.

#### Scenario: Principal flow format
- **WHEN** se genera el flujo principal
- **THEN** usa pasos numerados con actor ("El usuario...") y sistema ("El sistema...")

#### Scenario: Optional alternative flows
- **WHEN** no surgieron flujos alternativos en la conversación
- **THEN** la subsección de flujos alternativos se omite

### Requirement: Acceptance criteria in Gherkin format
Los criterios de aceptación SHALL usar formato Gherkin simplificado (Dado/Cuando/Entonces) y SIEMPRE se generan aunque el usuario no los haya mencionado.

#### Scenario: Agent derives AC from flows
- **WHEN** el usuario no proporcionó criterios de aceptación explícitos
- **THEN** el skill los deriva de los flujos documentados, con mínimo 3 y máximo 8

#### Scenario: Critical behaviors always present
- **WHEN** se genera la sección de Criterios de Éxito
- **THEN** incluye "Comportamientos críticos" (cosas que deben funcionar sí o sí)

### Requirement: Technical context section is conditional
La sección Contexto Técnico SHALL aparecer solo si el usuario respondió la Fase 3.5.

#### Scenario: User skipped Phase 3.5
- **WHEN** el usuario saltó la Fase 3.5
- **THEN** el .md no incluye la sección "Contexto Técnico"

#### Scenario: User answered Phase 3.5
- **WHEN** el usuario respondió las preguntas técnicas
- **THEN** el .md incluye la sección con stack/tecnologías, integraciones, restricciones, módulos, y notas

### Requirement: Footer signature
El .md generado SHALL terminar con la línea: `*Generado por feature-shaper v1 · /shape-refine para refinar*`

#### Scenario: Footer present
- **WHEN** se genera el .md
- **THEN** la última línea es la firma del generador
