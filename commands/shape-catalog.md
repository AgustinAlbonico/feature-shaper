---
description: Muestra el catálogo de features del proyecto actual con status, tipo y versión
---

Mostrá el catálogo de features del proyecto actual usando feature_catalog.

Detectá el slug del proyecto desde el directorio actual (nombre del directorio o campo name del package.json/go.mod).

Variantes soportadas:
- Sin argumentos: features del proyecto actual
- --all: features de todos los proyectos (usar project_list + feature_catalog por proyecto)
- --status <valor>: filtrar por status (draft/ready/in-progress/done)
- --type <valor>: filtrar por tipo (product/technical/business)

Formato de salida:
📋 Features — <proyecto> (N features)

  ✓ ready       <título>   <tipo>   v<n>  · <fecha>
  ◌ draft       <título>   <tipo>   v<n>  · <fecha>
  ◎ in-progress <título>   <tipo>   v<n>  · <fecha>

Usá /shape-refine "nombre" para refinar. Usá /shape "idea" para agregar una nueva.
