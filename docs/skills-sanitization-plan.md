## Skills Sanitization Plan (P1/P2/P3)

Objetivo: alinear `skills/**/SKILL.md` al rubric (`docs/skills-review-rubric.md`) con foco en skills operativas (sazón del repo), no manuales genéricos.

### Orden sugerido

1. **P1** (alto impacto, rewrite)
2. **P2** (impacto medio, trim + recalibración)
3. **P3** (mejoras menores y consistencia)

---

## P1 — Rewrites de alto impacto

**Qué entra**
- Skills con deriva fuerte a handbook/enciclopedia.
- Skills con exceso de longitud y ruido histórico en el core operativo.
- Skills con defaults ambiguos o múltiples opciones sin jerarquía.

**Modelo de referencia**
- `skills/backend/dataaccess/SKILL.md` (**skill modelo** para estructura, tono y densidad de reglas).

**Entregables P1**
- Reescritura completa orientada a:
  - `When to Use`
  - `Critical Patterns`
  - `Constraints & Tradeoffs`
  - `Anti-Patterns`
  - `Progressive Disclosure`
  - `Resources`
  - `Changelog`

**Definition of Done (P1)**
- Layer 1: **PASS** (spec completa, `allowed-tools` correcto).
- Layer 2: **GOOD** o **TRIM** (sin contenido enciclopédico en el core).
- Layer 3: **ALIGNED** (convenciones/gotchas reales del repo).
- Longitud objetivo alcanzada (ideal 80–140; aceptable hasta ~180–220 según rubric).

---

## P2 — Trim de impacto medio

**Qué entra**
- Skills que pasan spec pero tienen ruido moderado.
- Skills útiles pero con snippets de más o secciones redundantes.

**Acciones típicas**
- Recortar teoría genérica y onboarding de framework.
- Consolidar reglas duplicadas en reglas atómicas (Do/Don’t + Why).
- Mover detalle no crítico a `references/` cuando aplique.

**Definition of Done (P2)**
- Sin bloqueos de Layer 1.
- Decisión final esperada: **KEEP** o **TRIM** (no REWRITE pendiente).
- Se mantiene activación clara y default principal inequívoco.

---

## P3 — Mejoras menores y homogenización

**Qué entra**
- Ajustes de consistencia editorial y cross-references.
- Limpieza de frases ambiguas, ejemplos innecesarios y formato.

**Acciones típicas**
- Unificar estilo de triggers y secciones.
- Validar links internos y referencias cruzadas.
- Reducir fricción de lectura (menos ruido, más accionable).

**Definition of Done (P3)**
- Todas las skills en estado **operable** y consistentes en estructura base.
- Sin historial de auditoría interna en el core (`[DEFICIENCY FIX]`, etc.).
- Señal clara de “skill de cocina del repo” (no manual genérico).

---

## Regla de operación por lote

- Ejecutar por lotes pequeños, validar lint/tests al cierre de cada lote.
- No mezclar refactors estructurales de tools con sanitización de skills.
- Si una skill falla Layer 1, se corrige primero y recién luego se optimiza contenido.
