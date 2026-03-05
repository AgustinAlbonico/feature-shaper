---
description: Refina una feature existente cargando su definición actual y conduciendo una conversación de actualización
---

Cargá el skill "feature-shaper" y comenzá el proceso de refinamiento.

El nombre o slug de la feature está disponible como argumento del comando.

Comportamiento especial:
- Si no se pasó argumento: mostrá las últimas 5 features del proyecto actual y pedí que el usuario elija
- Si hay múltiples matches en feature_search: mostrá lista para elegir
- Saltá directamente a Fase 3 del protocolo con el contexto de la feature ya cargado
