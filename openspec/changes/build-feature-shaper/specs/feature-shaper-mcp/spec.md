## ADDED Requirements

### Requirement: MCP server stdio mode
El sistema SHALL arrancar un servidor MCP en modo stdio al ejecutar `feature-shaper mcp`, exponiendo las 8 tools registradas.

#### Scenario: Server starts and blocks on stdio
el usuario ejecuta `feature-shaper mcp`
- **THEN** el servidor arranca, ejecuta auto-migrate, registra las 8 tools, y bloquea leyendo stdin

### Requirement: feature_save tool
La tool `feature_save` SHALL hacer upsert de una feature. Si el topicKey (`projectSlug/slug`) ya existe, guarda snapshot en featureVersions e incrementa versión. Si no existe, crea la feature con version=1.

#### Scenario: Save new feature
- **WHEN** se llama `feature_save` con projectSlug, slug, title, type, content
- **THEN** el sistema auto-registra el proyecto si no existe, crea la feature con version=1, y devuelve la feature creada

#### Scenario: Update existing feature
- **WHEN** se llama `feature_save` con un slug que ya existe para ese projectSlug
- **THEN** el sistema guarda snapshot del contenido actual en featureVersions, incrementa version, actualiza content y updatedAt

#### Scenario: Save with changelog
- **WHEN** se llama `feature_save` con parámetro changelog
- **THEN** el changelog se asocia al snapshot de la versión anterior en featureVersions

### Requirement: feature_get tool
La tool `feature_get` SHALL recuperar la última versión de una feature por slug.

el usuario ejecuta `feature-shaper tui`
- **WHEN** se llama `feature_get` con slug existente
- **THEN** el sistema devuelve la feature completa incluyendo content, version, status, type

#### Scenario: Get with projectSlug filter
- **WHEN** se llama `feature_get` con slug y projectSlug
- **THEN** el sistema busca solo en el proyecto especificado

### Requirement: feature_search tool
La tool `feature_search` SHALL buscar features usando FTS5 con el query proporcionado.

#### Scenario: Search across all projects
- **WHEN** se llama `feature_search` con query sin projectSlug
- **THEN** el sistema busca en todas las features de todos los proyectos

#### Scenario: Search within project
- **WHEN** se llama `feature_search` con query y projectSlug
- **THEN** el sistema busca solo en las features del proyecto especificado

#### Scenario: Search result format
- **WHEN** la búsqueda encuentra resultados
- **THEN** cada resultado incluye id, projectSlug, slug, title, type, status, version, preview (primeros 200 chars), updatedAt

### Requirement: feature_catalog tool
La tool `feature_catalog` SHALL listar features de un proyecto con filtros opcionales de status y type.

#### Scenario: List all features of a project
- **WHEN** se llama `feature_catalog` con projectSlug
- **THEN** el sistema devuelve todas las features del proyecto

#### Scenario: Filter by status
- **WHEN** se llama `feature_catalog` con projectSlug y status="ready"
- **THEN** el sistema devuelve solo features con status "ready"

#### Scenario: Filter by type
- **WHEN** se llama `feature_catalog` con projectSlug y type="product"
- **THEN** el sistema devuelve solo features de tipo "product"

### Requirement: feature_versions tool
La tool `feature_versions` SHALL listar el historial de versiones de una feature.

#### Scenario: List version history
- **WHEN** se llama `feature_versions` con slug y projectSlug
- **THEN** el sistema devuelve la lista de versiones con id, version, changelog, createdAt

### Requirement: feature_get_version tool
La tool `feature_get_version` SHALL recuperar el contenido de una versión específica de una feature.

#### Scenario: Get specific version
- **WHEN** se llama `feature_get_version` con featureId y version
- **THEN** el sistema devuelve el content completo del snapshot de esa versión

### Requirement: project_register tool
La tool `project_register` SHALL hacer upsert de un proyecto por slug.

#### Scenario: Register new project
- **WHEN** se llama `project_register` con slug y name
- **THEN** el sistema crea el proyecto si no existe

#### Scenario: Update existing project
- **WHEN** se llama `project_register` con un slug que ya existe
- **THEN** el sistema actualiza name y path sin duplicar

### Requirement: project_list tool
La tool `project_list` SHALL listar todos los proyectos registrados con conteo de features.

#### Scenario: List projects
- **WHEN** se llama `project_list`
- **THEN** el sistema devuelve todos los proyectos con slug, name, path, createdAt, y featureCount
