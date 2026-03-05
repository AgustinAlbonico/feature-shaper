# Spec: Formato del Output `.md`

Cada feature definition generada por el skill sigue esta estructura exacta.

## Template completo

```markdown
# [Título de la Feature]

> **Proyecto**: nombre-del-proyecto  
> **Tipo**: product | technical | business  
> **Status**: draft | ready  
> **Versión**: 1  
> **Creado**: YYYY-MM-DD  
> **Actualizado**: YYYY-MM-DD  

---

## Contexto

[Por qué existe esta feature. Qué problema resuelve. Sin soluciones técnicas.]

**Usuario principal**: [quién usa esto]  
**Actores secundarios**: [quién más interactúa — omitir si no aplica]  
**Situación actual**: [cómo se resuelve hoy sin esta feature]  

---

## Alcance

### ✅ Incluido en esta versión
- [item incluido]
- [item incluido]

### ❌ Fuera de alcance (explícito)
- [item excluido] — [razón breve]
- [item excluido] — [razón breve]

### 🔗 Features relacionadas o afectadas
- [feature existente que se modifica o de la que depende — omitir sección si no aplica]

---

## Flujos

### Flujo principal
1. El usuario [acción]
2. El sistema [respuesta/resultado]
3. [...]
4. Resultado: [estado final]

### Flujos alternativos
**[Nombre del flujo alternativo]**  
Condición: [cuándo ocurre]  
1. [paso]
2. [paso]

### Flujos de error
**[Nombre del error]**  
Condición: [qué lo dispara]  
Comportamiento esperado: [qué debe pasar]

---

## Criterios de Éxito

### Criterios de aceptación
- [ ] Dado [contexto], cuando [acción], entonces [resultado esperado]
- [ ] Dado [contexto], cuando [acción], entonces [resultado esperado]
- [ ] [mínimo 3, máximo 8]

### Comportamientos críticos
> Estas cosas deben funcionar sí o sí para que la feature tenga valor:
- [comportamiento no negociable]
- [comportamiento no negociable]

### Edge cases identificados
- [caso borde]: [comportamiento esperado]
- [caso borde]: [comportamiento esperado]

---

## Contexto Técnico

> Esta sección captura intenciones técnicas de alto nivel como insumo para OpenSpec.
> No es un diseño técnico — ese viene después.

**Stack / Tecnologías preferidas**: [valor o "stack actual del proyecto"]  
**Integraciones conocidas**: [valor]  
**Restricciones técnicas**: [valor]  
**Módulos involucrados**: [valor]  
**Otras notas**: [valor]

---

## Notas de Negocio

[Decisiones, restricciones o contexto que no encaja en las secciones anteriores
pero es importante recordar. Omitir sección si no hay nada relevante.]

---

*Generado por feature-shaper v1 · /shape-refine para refinar*
```

---

## Reglas de generación

| Sección | Regla |
|---|---|
| **Contexto** | Máx 3 párrafos cortos. Cero soluciones técnicas. |
| **Alcance — fuera de alcance** | Obligatorio con al menos 2 items. Si el usuario no los mencionó, el agente los propone a partir de lo que quedó implícito. |
| **Flujos alternativos y de error** | Solo si surgieron en el shaping. Omitir las secciones si no aplica. |
| **Criterios de aceptación** | Formato Gherkin simplificado (Dado/Cuando/Entonces). **Siempre se genera** aunque el usuario no los haya mencionado — el agente los deriva de los flujos. |
| **Contexto Técnico** | Solo aparece si el usuario respondió la Fase 3.5. Si la saltea, omitir la sección completa. |
| **Notas de Negocio** | Opcional. |

---

## Ejemplo real: "Sistema de Invitaciones"

```markdown
# Sistema de Invitaciones al Workspace

> **Proyecto**: mi-saas-app  
> **Tipo**: product  
> **Status**: ready  
> **Versión**: 1  
> **Creado**: 2026-03-04  
> **Actualizado**: 2026-03-04  

---

## Contexto

Los usuarios administradores necesitan poder agregar colaboradores a su workspace
sin que esos colaboradores tengan que registrarse de forma independiente. Hoy el
único camino es que el colaborador cree su propia cuenta y luego el admin la vincule
manualmente, lo que genera friccción y errores.

**Usuario principal**: Administrador del workspace  
**Actores secundarios**: Colaborador invitado (recibe email), sistema de email  
**Situación actual**: El admin comparte la URL de registro y luego vincula la cuenta manualmente  

---

## Alcance

### ✅ Incluido en esta versión
- Envío de invitación por email con link de aceptación único (expira en 48hs)
- Flujo de aceptación: el invitado hace click en el link y completa su perfil
- El admin puede ver invitaciones pendientes y cancelarlas
- Límite de 10 invitaciones pendientes simultáneas por workspace (plan free)

### ❌ Fuera de alcance (explícito)
- Invitación masiva por CSV — se agrega en v2
- Roles granulares al momento de invitar — por ahora todos entran como "member"
- Invitaciones por fuera del email (link directo, QR) — fuera de alcance por ahora

### 🔗 Features relacionadas o afectadas
- `user-auth` — el flujo de aceptación reutiliza el onboarding de registro

---

## Flujos

### Flujo principal
1. El admin navega a Configuración → Miembros → Invitar
2. Ingresa el email del colaborador y hace click en "Enviar invitación"
3. El sistema envía el email con el link de aceptación (válido 48hs)
4. El colaborador hace click en el link, completa nombre y contraseña
5. El colaborador queda vinculado al workspace con rol "member"
6. El admin ve al nuevo miembro en la lista

### Flujos alternativos
**Colaborador ya tiene cuenta**  
Condición: el email ya está registrado en la plataforma  
1. El sistema detecta la cuenta existente
2. El email de invitación tiene un link de "Aceptar con tu cuenta existente"
3. El colaborador hace login (si no está logueado) y acepta
4. Queda vinculado al workspace sin crear nueva cuenta

### Flujos de error
**Link expirado**  
Condición: el colaborador hace click en el link después de 48hs  
Comportamiento esperado: Página de error clara con opción de pedir una nueva invitación al admin

**Invitación cancelada**  
Condición: el admin canceló la invitación antes de que el colaborador la acepte  
Comportamiento esperado: Página de error indicando que la invitación ya no es válida

---

## Criterios de Éxito

### Criterios de aceptación
- [ ] Dado que soy admin, cuando envío una invitación, entonces el colaborador recibe el email en menos de 2 minutos
- [ ] Dado que soy colaborador invitado, cuando acepto la invitación, entonces quedo como miembro del workspace sin pasos adicionales
- [ ] Dado que la invitación tiene más de 48hs, cuando el colaborador intenta aceptarla, entonces ve un error claro y puede solicitar una nueva
- [ ] Dado que soy admin, cuando cancelo una invitación pendiente, entonces el link deja de funcionar inmediatamente
- [ ] Dado que el workspace free tiene 10 invitaciones pendientes, cuando intento enviar una más, entonces el sistema me bloquea con un mensaje claro

### Comportamientos críticos
> Estas cosas deben funcionar sí o sí para que la feature tenga valor:
- El link de invitación expira exactamente a las 48hs y no puede ser reutilizado
- Un colaborador no puede unirse a un workspace sin una invitación válida

### Edge cases identificados
- Email con mayúsculas/minúsculas mixtas: tratar como case-insensitive
- El colaborador acepta la invitación desde otro dispositivo que el email: debe funcionar igual
- El admin borra el workspace mientras hay invitaciones pendientes: las invitaciones quedan inválidas

---

*Generado por feature-shaper v1 · /shape-refine para refinar*
```
