## ADDED Requirements

### Requirement: Database schema creation
El sistema SHALL crear el schema SQLite completo al ejecutar `feature-shaper migrate`, incluyendo las tablas `projects`, `features`, `featureVersions`, la virtual table FTS5 `featuresFts`, y los 3 triggers de sincronización.

#### Scenario: First-time migration
- **WHEN** el usuario ejecuta `feature-shaper migrate` por primera vez
- **THEN** el sistema crea `~/.feature-shaper/features.db` con todas las tablas, virtual table FTS5, y triggers

#### Scenario: Idempotent migration
- **WHEN** el usuario ejecuta `feature-shaper migrate` con una DB existente
- **THEN** el sistema no falla y no pierde datos (todas las sentencias usan `IF NOT EXISTS`)

#### Scenario: WAL mode enabled
- **WHEN** la DB se abre
- **THEN** el sistema activa `pragma journal_mode=WAL` y `pragma foreign_keys=ON`

### Requirement: Auto-migration on MCP and TUI start
El sistema SHALL ejecutar migrate automáticamente al arrancar tanto `feature-shaper mcp` como `feature-shaper tui`.

#### Scenario: MCP auto-migrates
el usuario ejecuta `feature-shaper mcp`
- **THEN** el sistema ejecuta migrate antes de iniciar el servidor MCP

#### Scenario: TUI auto-migrates
el usuario ejecuta `feature-shaper tui`
- **THEN** el sistema ejecuta migrate antes de iniciar la TUI

### Requirement: Projects table schema
La tabla `projects` SHALL tener las columnas: `id` (PK autoincrement), `slug` (UNIQUE NOT NULL), `name` (NOT NULL), `path` (nullable), `createdAt` (default datetime('now')).

#### Scenario: Project uniqueness
- **WHEN** se intenta insertar un proyecto con un slug que ya existe
- **THEN** el sistema actualiza los campos en lugar de duplicar (upsert)

### Requirement: Features table schema
La tabla `features` SHALL tener las columnas: `id` (PK autoincrement), `projectSlug` (FK a projects.slug), `slug`, `title`, `type`, `status` (default 'draft'), `content`, `version` (default 1), `topicKey` (UNIQUE), `normalizedHash`, `createdAt`, `updatedAt`, y constraint UNIQUE(projectSlug, slug).

#### Scenario: Feature uniqueness per project
- **WHEN** se intenta insertar una feature con mismo projectSlug y slug
- **THEN** el sistema detecta el duplicado via topicKey y ejecuta el flujo de upsert

### Requirement: Feature versions table schema
La tabla `featureVersions` SHALL tener: `id` (PK autoincrement), `featureId` (FK a features.id con ON DELETE CASCADE), `version`, `content`, `changelog` (nullable), `createdAt`.

#### Scenario: Cascade delete
- **WHEN** se elimina una feature de la tabla features
- **THEN** todas sus versiones en featureVersions se eliminan automáticamente

### Requirement: FTS5 search index
El sistema SHALL mantener una virtual table FTS5 `featuresFts` sincronizada con la tabla `features` mediante triggers de INSERT, UPDATE y DELETE.

#### Scenario: FTS synced on insert
- **WHEN** se inserta una nueva feature
- **THEN** el trigger `features_ai` inserta los campos title, content, type, status en featuresFts

#### Scenario: FTS synced on update
- **WHEN** se actualiza una feature
- **THEN** el trigger `features_au` elimina la entrada vieja e inserta la nueva en featuresFts

#### Scenario: FTS synced on delete
- **WHEN** se elimina una feature
- **THEN** el trigger `features_ad` elimina la entrada correspondiente de featuresFts
