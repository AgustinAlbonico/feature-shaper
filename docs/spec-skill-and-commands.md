# Spec: feature-shaper Skill + Commands

## Skill: `skills/feature-shaper/SKILL.md`

### Frontmatter y descripción

```yaml
---
name: feature-shaper
description: >
  Conduce una conversación adaptativa de negocio para transformar ideas vagas de features
  en definiciones completas y estructuradas. Usa preguntas de negocio (NO técnicas) organizadas
  en 5.5 fases con checkpoints. Al finalizar, genera una feature definition en formato .md
  y la persiste via feature_save MCP tool. Cargar este skill cuando se invoca /shape o /shape-refine.
---
```

---

## Protocolo de 5.5 Fases

```
FASE 1: Exploración del Contexto
    ↓ checkpoint — usuario confirma contexto entendido
FASE 2: Clasificación de la Feature
    ↓ checkpoint — usuario confirma tipo, capas y complejidad
FASE 3: Definición Adaptiva (preguntas de negocio — NUNCA técnicas)
    ↓ checkpoint — usuario aprueba el borrador de negocio
FASE 3.5: Contexto Técnico [OPCIONAL — el usuario decide]
    ↓ usuario decide si responder o saltear
FASE 4: Especificación Formal (generar el .md)
    ↓ checkpoint — usuario aprueba el .md final
FASE 5: Persistencia y Cierre
```

---

### FASE 1 — Exploración del Contexto

El agente:
1. Lee `package.json` / `go.mod` / `pyproject.toml` para entender el stack
2. Lee `README.md` si existe
3. Llama `feature_catalog` para ver features existentes del proyecto (entiende el dominio)
4. Detecta el slug del proyecto (nombre del directorio raíz o campo `name` del manifest)

Luego formula 2-3 preguntas contextuales basadas en lo encontrado. Si el contexto es claro, puede saltear las preguntas y mostrar directamente el resumen.

**Salida de la fase**:
```
Entendí lo siguiente sobre el proyecto:
- Stack: Next.js 15 + Prisma + PostgreSQL
- Dominio: SaaS de gestión de inventario B2B
- Features existentes: user-auth (ready), product-catalog (draft)

¿Es correcto? ¿Algo importante que agregar antes de continuar?
```

---

### FASE 2 — Clasificación de la Feature

Clasifica con dos dimensiones independientes:

**Dimensión 1 — Tipo**:
- `product` — foco en UX, flujos de usuario, estados de pantalla
- `technical` — foco en contratos, SLAs, dependencias de sistema
- `business` — foco en reglas de negocio complejas, actores, flujos alternativos

**Dimensión 2 — Capas** (detectadas del stack, confirmadas por el usuario):
- `Frontend` — menciona pantallas/UI, o el stack tiene React/Vue/Next/etc.
- `Backend` — casi siempre presente; APIs, lógica, procesos
- `Database` — persistencia, entidades nuevas, cambios de schema

**Complejidad** (determina rondas de preguntas en Fase 3):
- `low` — 1 módulo, flujo lineal → 1 ronda
- `medium` — 2-3 módulos, algunos edge cases → 2 rondas
- `high` — 4+ módulos, reglas de negocio complejas, múltiples actores → 3 rondas

**Presentación al usuario**:
```
Clasifiqué tu idea como:

  Tipo:        business
  Complejidad: alta
  Capas:       ◉ Frontend  ◉ Backend  ◉ Database
  Módulos:     inventory, orders, notifications, users

¿Estás de acuerdo? Podés ajustar las capas o el tipo antes de continuar.
```

---

### FASE 3 — Definición Adaptiva

**Reglas absolutas**:
- Máximo **4 preguntas por ronda** — nunca más
- Cada pregunta numerada ("Pregunta 2/4:") con contexto del por qué importa
- Si una respuesta anterior ya cubre implícitamente una pregunta → se salta
- Si una respuesta genera nueva ambigüedad → se agrega a la siguiente ronda
- **NUNCA preguntar nada técnico** — sin mencionar endpoints, DB, componentes, etc.

**Al final de cada ronda excepto la última**:
```
Ronda 1 completada. Hasta ahora entiendo:
- [bullet]
- [bullet]
- [bullet]

¿Hay algo incorrecto o querés agregar algo antes de continuar?
```

#### Banco de preguntas — 4 pilares

**Pilar 1 — Contexto / Por qué** (Ronda 1):
- ¿Qué problema concreto resuelve esta feature para el usuario?
- ¿Qué pasa hoy sin esta feature? ¿Cómo lo resuelven actualmente?
- ¿Quién es el usuario principal? ¿Hay otros actores involucrados?
- ¿Hay alguna restricción de negocio que la feature deba respetar?

**Pilar 2 — Alcance** (Ronda 1):
- ¿Qué es lo mínimo que tiene que hacer esta feature para tener valor real?
- ¿Qué cosas relacionadas quedan explícitamente fuera de esta versión?
- ¿Hay features existentes que se ven afectadas o modificadas?
- ¿Esta feature tiene fases o se entrega completa de una vez?

**Pilar 3 — Flujos** (Ronda 2):
- Describí el flujo principal: el usuario hace X, luego Y, luego Z. ¿Cómo termina?
- ¿Qué puede hacer el usuario si algo falla en el medio del flujo?
- ¿Hay flujos alternativos válidos además del camino feliz?
- ¿Qué pasa con los datos o el estado si el usuario abandona a mitad del proceso?
- ¿Hay actores secundarios que interactúan en algún punto del flujo?

**Pilar 4 — Criterios de éxito / AC** (Ronda 2-3):
- ¿Cómo sabés que esta feature funciona correctamente? Dame 2-3 escenarios concretos.
- ¿Hay algún comportamiento que si no funciona, la feature directamente no sirve?
- ¿Qué métricas o indicadores muestran que la feature tuvo el impacto esperado?
- ¿Hay restricciones de performance o confiabilidad que el negocio requiere?

**Edge cases** (solo si `complexity=high`):
- ¿Qué pasa si dos usuarios hacen la misma acción al mismo tiempo?
- ¿Hay límites? (máximo de items, vencimientos, cuotas, restricciones por plan)
- ¿Qué pasa con los datos históricos cuando se activa la feature?
- ¿Hay casos donde la feature debería bloquearse o restringirse?

**Énfasis adicional por tipo** (1-2 preguntas extra mezcladas en las rondas):

`product`:
- ¿Cuál es el happy path ideal desde la perspectiva del usuario?
- ¿Qué pasa si el usuario sale a la mitad del flujo? (draft, pérdida de datos, warning)

`technical`:
- ¿Hay restricciones de retrocompatibilidad que debés respetar?
- ¿Cuál es la estrategia si algo sale mal en producción?

`business`:
- ¿Cuáles son los flujos alternativos y de excepción más importantes?
- ¿Hay casos límite conocidos que históricamente generaron bugs o confusión?

#### Estructura de rondas

```
Ronda 1 → Pilar 1 (contexto) + Pilar 2 (alcance) + 1 pregunta de tipo
           máx 4 preguntas, las más abiertas

Ronda 2 → Pilar 3 (flujos) + inicio Pilar 4 (AC)
           máx 4 preguntas, adaptadas a respuestas de Ronda 1
           (solo si complexity=medium o high)

Ronda 3 → Pilar 4 (AC) completo + edge cases detectados
           máx 4 preguntas
           (solo si complexity=high)
```

---

### FASE 3.5 — Contexto Técnico [OPCIONAL]

El agente siempre ofrece esta fase. El usuario decide si responder o saltear.

**Presentación**:
```
La parte de negocio está completa. ¿Querés agregar contexto técnico de alto nivel
para que OpenSpec tenga más información cuando arranque el diseño?

Son 2-4 preguntas cortas, completamente opcionales.
[Sí, agregar contexto técnico]  [No, pasar directo a la especificación]
```

**Banco de preguntas técnicas** (el agente elige 2-4 según lo detectado):

Siempre disponibles:
- ¿Tenés alguna preferencia o restricción sobre el stack o tecnologías a usar?
- ¿Hay alguna integración con sistemas externos que ya sabés que va a ser necesaria?
- ¿Hay alguna restricción técnica que el negocio impone? (performance, disponibilidad, seguridad)
- ¿Hay algún módulo o parte del sistema existente que claramente va a estar involucrado?

Si la feature tiene flujos complejos o múltiples actores:
- ¿Hay algún proceso que debería ocurrir en background o de forma asíncrona?
- ¿Hay alguna consideración de escala que ya sepas?

Si la feature modifica features existentes:
- ¿Hay restricciones de retrocompatibilidad que debés respetar?
- ¿Hay deuda técnica en esa área que podría complicar la implementación?

Si la feature involucra datos sensibles:
- ¿Hay consideraciones sobre privacidad o seguridad de los datos?
- ¿Los datos tienen restricciones legales o de compliance?

---

### FASE 4 — Especificación Formal

El agente genera el `.md` completo (ver `spec-output-format.md`). Presenta el documento y pregunta:

```
¿Aprobás esta especificación para guardarla?
Podés pedirme ajustes antes de confirmar.
```

El usuario puede pedir ajustes. El agente los incorpora y vuelve a presentar.

---

### FASE 5 — Persistencia y Cierre

1. Llama `project_register` para asegurar que el proyecto existe en la DB
2. Llama `feature_save` con todos los datos
3. Escribe el archivo `docs/features/<slug>.md` en el repo actual
4. Muestra confirmación:

```
✓ Feature guardada: Sistema de Invitaciones (v1)
✓ Archivo creado: docs/features/workspace-invites.md

¿Querés arrancar OpenSpec ahora para especificar esta feature?
(Podés hacerlo después con /sdd-propose)
```

---

### Flujo de `/shape-refine`

1. Llama `feature_search` con el argumento para encontrar la feature
2. Si hay más de un match → presenta lista para elegir
3. Llama `feature_get` para cargar el contenido actual
4. Presenta resumen de la versión actual:

```
Refinando: Sistema de Invitaciones (v2, ready)
Última actualización: 2026-02-28

Resumen actual:
- Contexto: Permite a admins invitar colaboradores al workspace por email
- Alcance: 4 items incluidos, 2 excluidos explícitamente
- Flujos: flujo principal + 2 alternativos + 2 flujos de error
- AC: 5 criterios definidos

¿Qué querés cambiar o agregar?
```

5. Salta directo a **Fase 3** con todo el contexto cargado
6. Al guardar: `version++`, snapshot en `featureVersions` con `changelog`

---

## Commands

### `commands/shape.md`

```markdown
---
description: Transforma una idea de feature en una definición completa de negocio, lista para OpenSpec
---

Cargá el skill "feature-shaper" y comenzá el proceso de shaping.

La idea inicial del usuario ya está disponible como argumento del comando.

Comportamiento especial:
- Si no se pasó argumento: preguntá "¿Sobre qué feature querés trabajar?"
- Si se detecta una feature similar via feature_search: ofrecé refinar en lugar de crear nueva
- Arrancá siempre desde Fase 1 del protocolo (Exploración del Contexto)
```

---

### `commands/shape-refine.md`

```markdown
---
description: Refina una feature existente cargando su definición actual y conduciendo una conversación de actualización
---

Cargá el skill "feature-shaper" y comenzá el proceso de refinamiento.

El nombre o slug de la feature está disponible como argumento del comando.

Comportamiento especial:
- Si no se pasó argumento: mostrá las últimas 5 features del proyecto actual y pedí que el usuario elija
- Si hay múltiples matches en feature_search: mostrá lista para elegir
- Saltá directamente a Fase 3 del protocolo con el contexto de la feature ya cargado
```

---

### `commands/shape-catalog.md`

```markdown
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
```
