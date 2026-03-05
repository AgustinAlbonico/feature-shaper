# Guía de Usuario — feature-shaper

feature-shaper es una herramienta que te ayuda a transformar ideas vagas en definiciones de features completas y estructuradas, listas para implementar.

---

## ¿Qué es feature-shaper?

feature-shaper actúa como una etapa intermedia entre "tengo una idea" y "vamos a codear". Te guía a través de preguntas de **negocio** (no técnicas) para clarificar:

- **Qué** quieres construir
- **Por qué** es importante
- **Para quién** es
- **Cómo** se verá el éxito
- **Qué** casos edge hay que considerar

El resultado es un documento `.md` estructurado que puedes usar como base para diseño técnico, planning, o compartir con tu equipo.

---

## Formas de Uso

feature-shaper tiene tres interfaces:

| Interfaz | Comando | Descripción |
|----------|---------|-------------|
| **Conversacional** | `/shape` | Diálogo guiado con el agente de OpenCode |
| **TUI** | `feature-shaper tui` | Interfaz visual en terminal |
| **MCP Tools** | Directas | Herramientas programáticas |

---

## Uso Conversacional (Recomendado)

### Crear una feature desde cero

```
/shape "quiero agregar sistema de invitaciones al workspace"
```

El agente te guiará a través de **5.5 fases**:

#### Fase 1: Exploración del Contexto

El agente detecta automáticamente el contexto de tu proyecto y hace preguntas iniciales:

- ¿En qué proyecto estás trabajando?
- ¿Hay features relacionadas?
- ¿Cuál es el objetivo principal?

#### Fase 2: Clasificación

El agente clasifica tu feature en tres dimensiones:

| Dimensión | Valores | Descripción |
|-----------|---------|-------------|
| **Tipo** | `product`, `technical`, `business` | Naturaleza de la feature |
| **Capas** | Frontend, Backend, DB, API, etc. | Componentes afectados |
| **Complejidad** | `low`, `medium`, `high` | Nivel de detalle necesario |

Ejemplo de clasificación:

```
📊 Clasificación detectada:
   Tipo: product
   Capas: Frontend, Backend, DB
   Complejidad: medium

¿Es correcta esta clasificación?
```

#### Fase 3: Definición Adaptiva

El agente te hace preguntas de negocio agrupadas en pilares:

**Pilares de preguntas:**

| Pilar | Preguntas típicas |
|-------|-------------------|
| **Usuarios** | ¿Quién puede invitar? ¿Quién puede ser invitado? |
| **Acciones** | ¿Qué pasa cuando se invita? ¿Cómo se acepta? |
| **Datos** | ¿Qué información se necesita? ¿Qué se almacena? |
| **Reglas** | ¿Hay límites? ¿Validaciones? |
| **Flujos** | ¿Qué pasa si X falla? ¿Casos alternativos? |
| **Éxito** | ¿Cómo sabemos que funciona? ¿Métricas? |

**Regla importante:** El agente hace **máximo 4 preguntas por ronda** para no abrumarte.

#### Fase 3.5: Contexto Técnico (Opcional)

Si tu feature tiene consideraciones técnicas importantes:

```
¿Hay consideraciones técnicas que debamos discutir?
[1] Sí, hablemos de arquitectura
[2] No, saltemos directo a la especificación
```

Si eliges "Sí", el agente te pregunta sobre:
- Integraciones existentes
- Restricciones de performance
- Consideraciones de seguridad
- Dependencias

#### Fase 4: Especificación Formal

El agente genera un documento estructurado:

```markdown
# Sistema de Invitaciones al Workspace

## Contexto
Los workspaces necesitan permitir que miembros existentes inviten a nuevos usuarios...

## Alcance
### Incluye
- Envío de invitaciones por email
- Aceptación de invitaciones
- Gestión de roles

### Excluye
- Invitaciones masivas
- Integración con LDAP

## Flujos
### Flujo Principal
1. Usuario hace clic en "Invitar"
2. Ingresa email del invitado
3. Sistema envía email con link único
4. Invitado hace clic en el link
5. Invitado crea cuenta o inicia sesión
6. Invitado se une al workspace

### Flujos Alternativos
- Link expirado → Mostrar mensaje, ofrecer reenvío
- Email ya registrado → Ofrecer unirse directamente

## Criterios de Éxito
- Usuario puede invitar en < 30 segundos
- 90% de invitaciones aceptadas en 48hs
- Zero invitaciones duplicadas
```

#### Fase 5: Persistencia

El agente guarda todo automáticamente:

```
✅ Feature guardada:
   - Base de datos: ~/.feature-shaper/features.db
   - Archivo: docs/features/sistema-invitaciones-workspace.md
   - Versión: 1
```

---

### Refinar una feature existente

Si ya creaste una feature y quieres mejorarla:

```
/shape-refine "invitaciones"
```

El agente:
1. Busca features que coincidan con "invitaciones"
2. Carga la definición actual
3. Salta directo a Fase 3 (preguntas de negocio)
4. Mantiene el contexto previo
5. Incrementa la versión al guardar

**Nota:** Si hay múltiples matches, el agente te pedirá que elijas.

---

### Ver el catálogo del proyecto

Para ver todas las features de tu proyecto:

```
/shape-catalog
```

**Con filtros:**

```
/shape-catalog --status draft
/shape-catalog --status ready
/shape-catalog --type product
/shape-catalog --type technical
/shape-catalog --all
```

**Salida ejemplo:**

```
📁 Proyecto: mi-app

┌─────────────────────────────────────────────────────────────┐
│ ◆ Sistema de Invitaciones                        v2 ready  │
│   Permite invitar usuarios al workspace via email           │
├─────────────────────────────────────────────────────────────┤
│ ◆ Notificaciones Push                           v1 draft   │
│   Notificaciones en tiempo real para eventos críticos       │
├─────────────────────────────────────────────────────────────┤
│ ⚙ API de Webhooks                               v1 ready   │
│   Endpoints para integraciones externas                     │
└─────────────────────────────────────────────────────────────┘

Total: 3 features (2 ready, 1 draft)
```

---

## Uso de la TUI

La Terminal User Interface te permite explorar tu catálogo visualmente:

```powershell
feature-shaper tui
```

### Pantalla Principal

```
┌──────────────────────────────────────────────────────────────────────┐
│ ◈ feature-shaper          3 proyectos • 12 features   [/] [?] [q]   │
├────────────────────┬─────────────────────────────────────────────────┤
│ 📁 PROYECTOS       │ ✨ FEATURES — mi-app                            │
│ ──────────────────│ ─────────────────────────────────────────────────│
│ ▶ mi-app (5)       │ ▶ ◆ Sistema de Invitaciones  product  ● ready  │
│   otro-proyecto (4)│   ◆ Notificaciones Push     product  ○ draft   │
│   tercer-app (3)   │   ⚙ API de Webhooks        technical ● ready  │
│                    │   ◈ Reportes de Ventas     business ◐ in-prog  │
└────────────────────┴─────────────────────────────────────────────────┘
│ [enter] abrir  [←] proyectos  [d] eliminar  [e] exportar  [q] salir │
└──────────────────────────────────────────────────────────────────────┘
```

### Navegación

| Tecla | Acción |
|-------|--------|
| `↑` / `k` | Mover cursor arriba |
| `↓` / `j` | Mover cursor abajo |
| `Tab` / `←` / `→` | Cambiar entre paneles |
| `Enter` | Abrir feature / seleccionar proyecto |
| `/` | Activar búsqueda |
| `?` | Mostrar ayuda |
| `q` | Salir |

### Vista de Detalle

Al presionar `Enter` sobre una feature:

```
┌──────────────────────────────────────────────────────────────────────┐
│ Sistema de Invitaciones al Workspace                                 │
│ [◆ product] [● ready] [v2]                        [b] [h] [e]        │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  # Sistema de Invitaciones                                          │
│                                                                      │
│  ## Contexto                                                        │
│  Los workspaces necesitan permitir que miembros existentes...       │
│                                                                      │
│  ## Alcance                                                         │
│  ### Incluye                                                        │
│  - Envío de invitaciones por email                                  │
│  ...                                                                │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ [↑↓/jk] scroll  [h] historial  [e] exportar  45%  [b] volver        │
└──────────────────────────────────────────────────────────────────────┘
```

### Historial de Versiones

Presiona `h` para ver el historial:

```
┌──────────────────────────────────────────────────────────────────────┐
│ 📜 Historial — Sistema de Invitaciones              [b] volver       │
│ 3 versiones guardadas                                                │
├──────────────────────────────────────────────────────────────────────┤
│  VER       FECHA         CAMBIOS                                     │
│ ────────────────────────────────────────────────────────────────────│
│ ▶ v3       2026-03-05   Agregado flujo de expiración               │
│   v2       2026-03-04   Refinado alcance, agregados casos edge     │
│   v1       2026-03-03   Versión inicial                            │
└──────────────────────────────────────────────────────────────────────┘
```

### Búsqueda FTS

Presiona `/` para buscar en todas las features:

```
┌──────────────────────────────────────────────────────────────────────┐
│ 🔍 Búsqueda de features                         2 resultados  [esc]  │
├──────────────────────────────────────────────────────────────────────┤
│ [invitacion █]                                                      │
├──────────────────────────────────────────────────────────────────────┤
│   ◆ Sistema de Invitaciones  [mi-app]  product  ● ready  v2        │
│     Permite invitar usuarios al workspace via email...              │
│                                                                      │
│   ◆ Invitaciones Masivas    [mi-app]  product  ○ draft   v1        │
│     Enviar múltiples invitaciones de una vez...                     │
└──────────────────────────────────────────────────────────────────────┘
```

### Acciones

| Tecla | Acción | Descripción |
|-------|--------|-------------|
| `e` | Exportar | Guarda la feature como `.md` en `docs/features/` |
| `d` | Eliminar | Elimina la feature (requiere confirmar con `d` dos veces) |

---

## Tipos de Feature

| Tipo | Icono | Color | Descripción |
|------|-------|-------|-------------|
| `product` | ◆ | Azul | Funcionalidad de usuario, UI/UX |
| `technical` | ⚙ | Verde | Infraestructura, APIs, herramientas |
| `business` | ◈ | Naranja | Lógica de negocio, reglas, procesos |

---

## Estados de Feature

| Estado | Icono | Color | Descripción |
|--------|-------|-------|-------------|
| `draft` | ○ | Gris | En desarrollo, no lista |
| `in-progress` | ◐ | Amarillo | Siendo refinada activamente |
| `ready` | ● | Verde | Lista para implementar |
| `done` | ✓ | Verde claro | Implementada |

---

## MCP Tools (Avanzado)

Si necesitas usar las herramientas directamente:

### feature_save

Guarda o actualiza una feature:

```json
{
  "tool": "feature_save",
  "arguments": {
    "projectSlug": "mi-app",
    "slug": "sistema-invitaciones",
    "title": "Sistema de Invitaciones",
    "type": "product",
    "status": "ready",
    "content": "# ... markdown content ...",
    "changelog": "Agregado flujo de expiración"
  }
}
```

### feature_get

Obtiene una feature por slug:

```json
{
  "tool": "feature_get",
  "arguments": {
    "slug": "sistema-invitaciones",
    "projectSlug": "mi-app"
  }
}
```

### feature_search

Búsqueda full-text:

```json
{
  "tool": "feature_search",
  "arguments": {
    "query": "invitaciones email",
    "projectSlug": "mi-app"
  }
}
```

### feature_catalog

Lista features del proyecto:

```json
{
  "tool": "feature_catalog",
  "arguments": {
    "projectSlug": "mi-app",
    "status": "ready",
    "type": "product"
  }
}
```

### feature_versions

Historial de versiones:

```json
{
  "tool": "feature_versions",
  "arguments": {
    "slug": "sistema-invitaciones",
    "projectSlug": "mi-app"
  }
}
```

### project_register

Registra un nuevo proyecto:

```json
{
  "tool": "project_register",
  "arguments": {
    "slug": "nuevo-proyecto",
    "name": "Nuevo Proyecto",
    "path": "/ruta/al/proyecto"
  }
}
```

### project_list

Lista todos los proyectos:

```json
{
  "tool": "project_list",
  "arguments": {}
}
```

---

## Flujo de Trabajo Recomendado

### 1. Ideación

```
/shape "tengo una idea para..."
```

Deja que el agente te guíe. No intentes tener todo claro desde el inicio.

### 2. Refinamiento

Después de una reunión o nuevo contexto:

```
/shape-refine "la-feature"
```

El agente recordará todo lo discutido anteriormente.

### 3. Revisión

Abre la TUI para revisar el catálogo:

```powershell
feature-shaper tui
```

### 4. Exportación

Desde la TUI, presiona `e` para exportar a `.md`.

### 5. Handoff

Comparte el archivo `.md` con tu equipo o úsalo como base para diseño técnico.

---

## Tips y Buenas Prácticas

### ✅ Hacer

- **Sé específico en el título:** "Sistema de invitaciones con expiración" > "Invitaciones"
- **Describe el problema, no la solución:** Deja que el agente te ayude a encontrar la mejor solución
- **Usa `/shape-refine` frecuentemente:** Es mejor iterar que intentar hacer todo perfecto la primera vez
- **Exporta regularmente:** Guarda versiones en `docs/features/` para compartir

### ❌ Evitar

- **No intentes responder todo de una vez:** El agente hace 4 preguntas a la vez por diseño
- **No saltees fases:** Cada fase tiene un propósito
- **No asumas que el agente sabe todo:** Dale contexto cuando te lo pida

---

## Preguntas Frecuentes

### ¿Dónde se guardan mis features?

- **Base de datos:** `~/.feature-shaper/features.db`
- **Exports:** `docs/features/<slug>.md` (en cada proyecto)

### ¿Puedo tener features en múltiples proyectos?

Sí. Cada feature está asociada a un proyecto. Usa `/shape-catalog` para ver features del proyecto actual, o `feature-shaper tui` para ver todos.

### ¿Qué pasa si quiero cambiar algo que ya guardé?

Usa `/shape-refine`. El sistema guarda un historial de versiones.

### ¿Puedo eliminar una feature?

Sí. En la TUI, selecciona la feature y presiona `d` dos veces.

### ¿El agente hace preguntas técnicas?

Solo si tú quieres. En la Fase 3.5 puedes elegir saltar las preguntas técnicas.

---

## Próximos Pasos

1. **Prueba el flujo:** `/shape "mi primera feature"`
2. **Explora la TUI:** `feature-shaper tui`
3. **Integra con tu workflow:** Usa las features exportadas como base para planning
