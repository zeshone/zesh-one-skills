# USER_DECISIONS.md

## Propósito

Este archivo centraliza las decisiones de arquitectura y configuración que necesita tomar el usuario del proyecto **Zesh-one-Skills**.

Cada pregunta tiene opciones claras y una recomendación del arquitecto. El usuario responde debajo de cada pregunta y luego avisa `ya podés leer` para que el agente continúe.

---

## Cómo usar este archivo

1. Leé cada pregunta y sus opciones.
2. Escribí tu respuesta debajo de la línea `**Tu respuesta:**`.
3. Cuando terminaste de responder las que podés, avisá: **"ya podés leer"**.
4. El agente leerá este archivo y continuará con el trabajo.

> Podés responder de a una pregunta o varias a la vez. No hace falta contestar todo junto.

---

## Decisiones pendientes

---

### D-01 · Modo de persistencia de artifacts SDD

¿Dónde se guardan los artifacts del proceso SDD (propuestas, specs, diseños, tareas)?

- [X] `engram` — Memoria persistente del agente entre sesiones. Recomendado si trabajás siempre con el mismo agente/herramienta.
- [ ] `openspec` — Archivos en el repositorio (carpeta `openspec/`). Recomendado si querés versionado en Git.
- [ ] `hybrid` — Ambos. Más robusto pero consume más tokens por operación.
- [ ] `none` — Sin persistencia. Los artifacts solo existen en la conversación activa.

> **Recomendación del arquitecto:** `openspec` — los artifacts quedan en el repo, son visibles, versionables y no dependen de la memoria del agente. Ideal para proyectos colaborativos o de largo plazo.

**Tu respuesta:**
engram
---

### D-02 · Metodología de implementación (TDD o estándar)

¿Con qué workflow implementamos el código?

- [X] `TDD` — Primero se escribe el test (RED), luego el código mínimo para pasarlo (GREEN), luego se refactoriza. Requiere suite de tests configurada.
- [ ] `standard` — Se escribe el código directamente. Los tests son opcionales y se agregan después.

> **Recomendación del arquitecto:** `TDD` si el proyecto tiene o va a tener tests. `standard` si el objetivo actual es validar ideas rápido. Para skills de agente, `standard` es suficiente salvo que haya lógica compleja que testear.

**Tu respuesta:**
TDD
---

### D-03 · Idioma de los artifacts SDD

¿En qué idioma se redactan specs, propuestas, diseños y tareas?

- [X] `español` — Todo en español rioplatense. Coherente con el idioma del proyecto y del equipo.
- [ ] `inglés` — Todo en inglés. Más portable si el proyecto es público o tiene colaboradores internacionales.
- [ ] `mixto` — Código y nombres técnicos en inglés, documentación en español.

> **Recomendación del arquitecto:** `español` para proyectos personales o de equipo hispanohablante. `inglés` si el repo va a ser público o open source.

**Tu respuesta:**
español
---

### D-04 · Estrategia de versionado de skills

¿Cómo manejamos las versiones de los skills cuando hay cambios breaking?

- [ ] `semver en carpeta` — Cada versión major tiene su carpeta (`/skill-name/v1/`, `/skill-name/v2/`).
- [ ] `rama git` — Las versiones se manejan como ramas del repositorio.
- [ ] `sin versionado formal` — El skill evoluciona en su lugar, sin historial de versiones explícito.
- [ ] `changelog en SKILL.md` — Un bloque `## Changelog` al final del archivo documenta los cambios.

> **Recomendación del arquitecto:** `changelog en SKILL.md` + Git como historial real. Simple, liviano y sin overhead de carpetas duplicadas. Si hay un breaking change muy grande, ahí sí consideramos carpeta de versión.

**Tu respuesta:**
changelog en skill.md + git
---

### D-05 · Política de skills compartidos entre agentes

¿Los skills de este repo están pensados para un solo agente o para múltiples agentes/herramientas?

- [ ] `un solo agente` — Diseñados para una herramienta específica (ej: OpenCode, Gemini CLI). Pueden usar sintaxis propietaria.
- [ ] `multi-agente` — Deben funcionar con cualquier agente compatible. Requiere formato agnóstico y portabilidad total.
- [ ] `por ahora uno, diseño portable` — Se desarrolla para un agente específico pero respetando estándares para migrar fácil después.

> **Recomendación del arquitecto:** `por ahora uno, diseño portable`. Evitás over-engineering ahora, pero no te encerrás. Documentás las dependencias de herramienta en el encabezado del SKILL.md para saber qué cambiar si migrás.

**Tu respuesta:**
un solo agente - opencode
---

### D-06 · Estructura de carpetas para nuevos skills

¿Cómo organizamos los skills dentro del repositorio?

- [ ] `plana` — Todos los skills en `skills-files/` directamente, sin subcarpetas.
- [ ] `por dominio` — Subcarpetas por área: `skills-files/sdd/`, `skills-files/testing/`, `skills-files/frontend/`, etc.
- [ ] `por agente` — Subcarpetas por herramienta destino: `skills-files/opencode/`, `skills-files/gemini/`, etc.
- [ ] `por dominio + por agente` — Dos niveles: dominio primero, agente después.

> **Recomendación del arquitecto:** `por dominio`. Es la organización más intuitiva para encontrar y mantener skills. Si en el futuro hay variantes por agente, se agrega un nivel sin romper lo existente.

**Tu respuesta:**
por dominio
---

### D-07 · Política de README y documentación pública

¿Este repositorio va a tener documentación pública dirigida a otros usuarios?

- [ ] `sí, completa` — README principal, guía de uso, ejemplos. El repo está pensado para ser usado por otros.
- [ ] `sí, mínima` — Un README breve que explique qué es y cómo instalar. Sin tutoriales extensos.
- [ ] `no por ahora` — Es un repo de trabajo interno. La documentación puede venir después.

> **Recomendación del arquitecto:** `sí, mínima`. Aunque sea un repo personal, un README mínimo te ahorra tiempo cuando volvés al proyecto después de semanas. Basta con: qué es, cómo usarlo, estructura de carpetas.

**Tu respuesta:**
sí, mínima -> recuerda que son para opencode con las herramientas de gentleman (engram y agents-team-lite)
---

### D-08 · Alcance real del repositorio

¿Qué scope tiene este repo en términos de tecnología cubierta?

- [ ] `A` — Solo .NET / backend. El repo se enfoca exclusivamente en skills y artifacts para proyectos backend con .NET.
- [ ] `B` — Full-stack. Cubre tanto backend (.NET) como frontend (Angular, React, etc.) desde el inicio.
- [ ] `C` — Backend core + frontend opcional. El núcleo es backend, pero se pueden agregar skills de frontend si el contexto lo requiere.

> **Recomendación del arquitecto:** `C` — Empezás con un núcleo sólido de backend sin sobre-comprometerte, y dejás la puerta abierta al frontend cuando sea necesario. Evita el over-engineering inicial.

**Tu respuesta:**
C - Separacion de scope back y front para el front usamos los skills de gentelman
---

### D-09 · Qué hacer con las skills de Gentleman

¿Cómo tratamos los skills existentes del repo de Gentleman como punto de partida?

- [ ] `A` — Mantenerlos como vendor/referencia. Se copian tal cual y se usan sin modificar.
- [ ] `B` — Adaptarlos y reescribirlos como Zesh. Cada skill de Gentleman se toma como base y se refactoriza con identidad propia.
- [ ] `C` — Usarlos solo como inspiración, sin migrarlos tal cual. Se estudian, se extraen ideas, pero el código se escribe desde cero con criterio propio.

> **Recomendación del arquitecto:** `C` para la mayoría — el valor está en entender el patrón, no en copiar el código. `A` de forma selectiva si algún skill ahorra mucho trabajo y no necesita adaptación.

**Tu respuesta:**
A - ya analice los skills y son muy generales por lo tanto me sirven, los ire acoplando y actualizando mas adelante
---

### D-10 · Licencia y destino del repositorio

¿Cuál es la política de visibilidad y licencia de este repo?

- [ ] `A` — Privado por ahora. Sin acceso público hasta nuevo aviso.
- [ ] `B` — Privado ahora, público después. Se abre cuando esté listo y limpio.
- [ ] `C` — Open-source desde el inicio. Licencia MIT o similar, visible desde el primer commit.

> **Recomendación del arquitecto:** `A` o `B` — hasta limpiar todo lo que pueda estar ligado a cliente o tener dependencias que no querés exponer. Después de esa limpieza, `B` es la opción natural.

**Tu respuesta:**
A
---

### D-11 · Firma y autoría canónica

¿Bajo qué nombre/handle se firma la autoría de este repo y sus skills?

- [ ] `A` — `zesh`
- [ ] `B` — `zesh-one`
- [ ] `C` — Otro nombre (especificar debajo).

> **Recomendación del arquitecto:** `B` — `zesh-one` es más específico, evita colisiones con handles genéricos y mantiene coherencia con el nombre del proyecto.

**Tu respuesta:**
B - respetamos las firmadas por gentleman
---

### D-12 · Qué significa TDD para este repositorio

¿Cómo se entiende y aplica TDD en el contexto concreto de este repo de skills?

- [ ] `A` — Casos de prueba documentados en Markdown. Los "tests" son escenarios descritos en texto, sin ejecución automática.
- [ ] `B` — Smoke tests con prompts / agente. Se prueba que el skill funciona ejecutándolo manualmente con un agente.
- [ ] `C` — Spec-first + smoke tests. Primero se escribe la spec del skill, luego se valida con un smoke test ejecutado con el agente.
- [ ] `D` — Otro enfoque definido por el usuario (especificar debajo).

> **Recomendación del arquitecto:** `C` — El spec-first te obliga a pensar qué debe hacer el skill antes de escribirlo, y el smoke test valida que realmente funciona. Es liviano, no requiere infraestructura de tests, y es coherente con el flujo SDD del repo.

**Tu respuesta:**
C
---

### D-13 · ¿Dónde viven las skills de Gentleman?

¿Cuál es la ubicación física de las skills que vienen del repo de Gentleman dentro de este repositorio?

- [ ] `A` — Dentro del repo en `vendor/gentleman`. Se copian tal cual en una carpeta vendor claramente separada.
- [ ] `B` — Fuera del repo como referencia externa. Solo se referencia el repositorio original, sin copiar archivos.
- [ ] `C` — Dentro del repo reorganizadas por dominio. Se integran directamente en la estructura del proyecto.

> **Recomendación del arquitecto:** `A` — Tener las skills en `vendor/gentleman` las mantiene accesibles sin mezclarlas con las propias. Queda claro cuál es vendor y cuál es Zesh. Si en algún momento cambia la fuente, sabés exactamente qué actualizar.

**Tu respuesta:**
A
---

### D-14 · ¿Cuándo una skill deja de ser vendor y pasa a ser Zesh?

¿Cuál es el criterio para "graduarla" de vendor a skill propia del proyecto?

- [ ] `A` — Cuando solo se traduce o cambia el formato. Cualquier modificación superficial la convierte en Zesh.
- [ ] `B` — Cuando cambia contenido sustancial: reglas, ejemplos, flujo de trabajo. El criterio es el cambio de fondo.
- [ ] `C` — Cuando se toca cualquier cosa. Cualquier edición implica que ya no es la skill original.
- [ ] `D` — Otra regla definida por el usuario (especificar debajo).

> **Recomendación del arquitecto:** `B` — Traducciones y ajustes de formato son cosmética; el skill sigue siendo del autor original. Cuando cambiás reglas, ejemplos o el flujo de trabajo, ahí sí tiene identidad propia. Es el criterio más honesto y evita inflación de "skills propias".

**Tu respuesta:**
B - mantendre el merito de gentleman hasta que haga cambios mayores como bien lo mencionas
---

### D-15 · ¿Cómo manejar créditos y atribución?

¿Cómo documentamos la autoría cuando una skill deriva de otra?

- [ ] `A` — Mantener el frontmatter original en vendor y listo. Sin cambios adicionales.
- [ ] `B` — Agregar sección `Credits` cuando una skill derive de otra. La atribución vive dentro del archivo.
- [ ] `C` — Archivo central `CREDITS.md`. Un único registro de todas las atribuciones del repositorio.
- [ ] `D` — Combinación de `B` + `C`. Sección `Credits` en el archivo Y entrada en `CREDITS.md`.

> **Recomendación del arquitecto:** `D` — La sección `Credits` en el archivo es inmediata y auto-contenida; el `CREDITS.md` central es el índice global. Juntos dan trazabilidad completa sin overhead excesivo.

**Tu respuesta:**
C
---

### D-16 · Primer skill a trabajar

¿Con cuál skill empezamos el desarrollo concreto del repositorio?

- [ ] `A` — `net8-apirest-dataaccess`
- [ ] `B` — `net8-apirest-general`
- [ ] `C` — `net8-testing-unit`
- [ ] `D` — `AGENTS.md` root
- [ ] `E` — Otro skill definido por el usuario (especificar debajo).

> **Recomendación del arquitecto:** `A` primero (`net8-apirest-dataaccess`), luego `C` (`net8-testing-unit`). El acceso a datos es el núcleo funcional de cualquier API REST y tener la capa de data access bien definida permite que el resto de los skills la referencien. Los tests unitarios vienen inmediatamente después para validar esa capa.

**Tu respuesta:**
A, C, B
---

### D-17 · Nombre final de la carpeta raíz de skills propias

¿Cómo se llama la carpeta raíz donde viven las skills de autoría Zesh dentro del repositorio?

- [ ] `A` — `skills/` — Nombre simple, directo y estándar.
- [ ] `B` — `zesh-skills/` — Nombre con identidad de proyecto, evita ambigüedad con vendor.
- [ ] `C` — Otro nombre definido por el usuario (especificar debajo).

> **Recomendación del arquitecto:** `A` — `skills/` es limpio, predecible y suficiente si ya existe `vendor/gentleman` como separación explícita. No hace falta redundar el nombre del proyecto en la carpeta.

**Tu respuesta:**
C - voy a agregar un nivel adicional
backend/skills - hasta ahora skills propias de .NET
forntend/skills - en espera de hacer mis propias skills para front
vendor/gentleman - separacion de autoria
---

### D-18 · ¿Querés AGENTS.md root desde el arranque?

¿Incluimos el archivo `AGENTS.md` en la raíz del repositorio desde el inicio del proyecto?

- [ ] `A` — Sí — Se crea `AGENTS.md` root desde el arranque para establecer contexto global del agente desde el primer momento.
- [ ] `B` — No todavía — Se agrega más adelante cuando el repo tenga más contenido y contexto para documentar.

> **Recomendación del arquitecto:** `A` — El `AGENTS.md` root es la carta de presentación del repo para cualquier agente que lo lea. Crearlo desde el inicio evita que el agente opere sin contexto durante las primeras sesiones, que son las más importantes.

**Tu respuesta:**
como estoy separando skills para back y fornt tambien quiero separar el agents.md
ya que no necesito los autoinvoke de front en proyectos back y viceversa
---

### D-19 · ¿vendor/gentleman queda fuera del auto-invoke principal?

¿Las skills en `vendor/gentleman` se excluyen del mecanismo de auto-invoke y se referencian solo al inicio de sesión?

- [ ] `A` — Sí, solo referencia al principio — `vendor/gentleman` se menciona como contexto inicial pero no participa en el auto-invoke automático del agente.
- [ ] `B` — No, también disponible para auto-invoke — Las skills de Gentleman se comportan igual que las propias y el agente las puede invocar automáticamente.

> **Recomendación del arquitecto:** `A` — Mantener vendor fuera del auto-invoke evita que el agente mezcle skills de Gentleman con las de Zesh en tiempo de ejecución. El contexto se da explícitamente, no de forma automática.

**Tu respuesta:**
No se toman encuenta las skills de vendor
---

### D-20 · ¿Querés preparar esto con SDD formal o con plan maestro + ejecución?

¿Cómo organizamos el trabajo de desarrollo de este repositorio a partir de ahora?

- [ ] `A` — SDD formal — Usamos el flujo completo: propuesta → spec → diseño → tareas → apply → verify → archive. Máxima trazabilidad.
- [ ] `B` — Plan maestro + ejecución — Definimos un plan de alto nivel y ejecutamos directamente sin artifacts SDD formales.

> **Recomendación del arquitecto:** `A` — El flujo SDD formal es exactamente para esto: proyectos propios donde querés construir con criterio, tener historial de decisiones y poder retomar en cualquier punto. El overhead inicial se paga con creces en consistencia y velocidad a largo plazo.

**Tu respuesta:**
A
---

## Decisiones futuras

> Esta sección es para agregar nuevas preguntas a medida que surjan durante el desarrollo del proyecto.
> Formato: copiar el bloque de cualquier decisión de arriba y adaptar.

<!-- Agregar nuevas decisiones aquí -->
