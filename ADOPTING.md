# Adopción de Zesh-One-Skills

Guía para integrar esta suite de skills en un proyecto nuevo que use OpenCode como agente de IA.

## Requisitos previos

Antes de adoptar, verifica que tienes:

- [ ] [OpenCode](https://opencode.ai) instalado y configurado
- [ ] [Engram](https://github.com/gentleman-programming/engram) disponible en el entorno del agente
- [ ] Node.js ≥ 18 (para ejecutar `npm run verify` localmente si quieres validar la suite)

## Paso 1: Copia el AGENTS.md correspondiente

En la raíz del repositorio encontrarás los archivos de instrucciones de agente por dominio:

| Archivo | Cuándo usarlo |
|---------|---------------|
| `AGENTS.backend.md` | Proyectos backend .NET 8 REST API (recomendado como punto de partida) |
| `AGENTS.dotnet-10.md` | .NET 10 / C# 14: API, CLI o worker; recetas HTTP solo para APIs |
| `AGENTS.frontend.md` | React SPA o Next.js; conserva el runtime instalado y su flujo de datos |
| `AGENTS.apps.md` | Angular/Ionic/Capacitor; verifica las versiones y convenciones del proyecto |

**Acción:** copia el archivo que aplique a la raíz de tu proyecto y cámbiale el nombre a `AGENTS.md`.

```bash
# Para proyectos backend:
cp AGENTS.backend.md /ruta/a/tu-proyecto/AGENTS.md

# Para proyectos frontend:
cp AGENTS.frontend.md /ruta/a/tu-proyecto/AGENTS.md
```

> El `AGENTS.md` es el punto de entrada que OpenCode lee al iniciar sesión. Sin él, el agente opera sin contexto del proyecto.

## Paso 2: Copia los skills relevantes

Los skills están organizados en cinco dominios bajo `skills/`:

- `skills/backend/` — Arquitectura, acceso a datos, logging, responses, seguridad, pruebas y más para .NET 8 REST API
- `skills/dotnet-10/` — Minimal APIs, arquitectura, persistencia, seguridad y pruebas para .NET 10 / C# 14
- `skills/frontend/` — Next.js 15, React 19, Tailwind 4, TypeScript, Zod 4, Zustand 5
- `skills/apps/` — Capa nativa móvil (Ionic Angular, Capacitor)
- `skills/shared/` — Skills transversales (`cli-design`, `github-pr`) y referencias compartidas de arquitectura y seguridad

**Acción:** copia el directorio del dominio que necesitas (más `skills/shared/`) a tu proyecto, o referencia este repositorio directamente desde tu `AGENTS.md`.

### Elección de arquitectura según el contexto

Incluye `skills/shared/architecture-selection.md` junto con el dominio. Los skills de arquitectura existentes usan esa referencia para inspeccionar el tipo de aplicación, las versiones y las decisiones previas; recomendar la estructura mínima con alternativas y costos; y esperar tu elección antes de crear o reorganizar proyectos. Una decisión vigente no se vuelve a pedir en cada funcionalidad.

Una CLI pequeña puede usar un proyecto de aplicación más pruebas, sin capas vacías, HTTP, base de datos ni autenticación inventados. Un segundo contexto no obliga a separar ensamblados. La arquitectura existente y consistente se conserva; Clean/Hexagonal, Vertical Slice y los límites físicos son decisiones combinables, no escalones obligatorios por tamaño. Las convenciones de .NET 8 y .NET 10 permanecen separadas.

### Interfaces de CLI

Cuando diseñes o cambies comandos, argumentos, ayuda, códigos de salida o modos de automatización, carga [cli-design](skills/shared/cli-design/SKILL.md). Las guías de .NET incluyen esa ruta condicional; en otros lenguajes, agrégala al punto de entrada equivalente. Conserva el selector de arquitectura, la seguridad compartida y las convenciones del lenguaje. El skill no obliga a instalar un parser, incorporar infraestructura ni rediseñar comandos existentes.

`npm run test:cli-instructions` comprueba el contrato de instrucciones y sus rutas; forma parte de `test` y `verify`. La suite habitual sigue requiriendo solo Node y no ejecuta una aplicación consumidora.

La [prueba ejecutable opcional](skills/shared/cli-design/assets/line-tally/README.md) se ejecuta con `npm run test:cli` y requiere un SDK de .NET 10 con sus paquetes de referencia ya instalados. Compila una copia desechable local, sin fuentes NuGet externas ni instalaciones, y comprueba el proceso público con datos sintéticos. Es aceptación de ese espécimen, no cobertura xUnit ni garantía de obediencia del modelo. Si el índice no muestra `cli-design`, usa la ruta explícita de la guía; no supongas descubrimiento automático.

### Perfiles de frontend y apps

Una React SPA mantiene su cliente autenticado y su router: no necesita Next.js, Server Components, Server Actions ni SSR. En el perfil Next.js se conserva el límite Client → Server Action → .NET API y sus controles. TanStack Query y Zustand solo se aplican si ya están instalados o se eligen explícitamente, y se conserva la arquitectura de formularios existente. `revalidatePath`/`revalidateTag` no invalidan la caché de Query en el navegador. `npm run test:react-stack-instructions` comprueba estas instrucciones, no ejecuta la aplicación. En Angular/Ionic, verifica las versiones instaladas y conserva formularios, HttpClient, navegación y ciclo de vida válidos; no copies reglas de otra integración o versión. No se instalan plugins ni se ejecuta sincronización nativa para una tarea que solo es de UI.

En apps existentes con NgModule o RxJS no se exige migrar a standalone ni a signals. Al reingresar a una página, los datos se actualizan según su vigencia, y el trabajo de una página oculta se detiene al salir de ella. Capacitor distingue plataforma, disponibilidad del plugin y permisos propios de cada plugin. `npm run test:apps-instructions` comprueba estas instrucciones, no el comportamiento en un dispositivo.

Para seleccionar o escribir pruebas de una React SPA, carga [testing de frontend](skills/frontend/testing/SKILL.md) mediante la ruta condicional de la guía. Conserva el runner instalado, el cliente autenticado y el flujo de formularios; no migres a Next ni a otro sistema de formularios por una tarea de pruebas. La guía distingue pruebas unitarias, de componente, de integración y E2E, con estados de UI y efectos observables.

`npm run test:frontend-instructions` comprueba instrucciones y rutas, no ejecuta React ni un navegador; forma parte de la suite habitual, solo con Node. Si faltan herramientas para una prueba ejecutable, informa la evidencia que no se puede obtener y solicita autorización antes de instalar. Los enlaces explícitos no garantizan el descubrimiento en un índice externo.

La [prueba opcional de React DOM](tools/fixtures/frontend-dom/README.md) usa dependencias propias fijadas y datos sintéticos. Instala sus requisitos solo con autorización explícita, dentro de su proyecto y caché local y con los scripts de instalación desactivados; después ejecuta `npm run test:frontend-dom`. La suite habitual no necesita esas dependencias. jsdom valida un DOM simulado: no demuestra comportamiento de navegador, E2E, accesibilidad completa ni integración con una API real.

### Pruebas de .NET según el perfil existente

Conserva [testing-unit de backend .NET 8](skills/backend/testing-unit/SKILL.md) o [testing de .NET 10](skills/dotnet-10/testing/SKILL.md), según corresponda; no transfieras sus librerías ni nombres de pruebas entre perfiles. Las guías refuerzan límites aplicables, resultados y ausencia de efectos prohibidos al rechazar, además de escenarios deterministas de tiempo y cancelación solo cuando existen esas capacidades. No duplican ni sustituyen las obligaciones de CLI, WAF, base de datos real y seguridad.

`npm run test:dotnet-instructions` comprueba instrucciones y preservación de perfiles; forma parte de la suite Node habitual. No es una ejecución de xUnit ni autoriza descargas de NuGet o cambios en las dependencias de un proyecto consumidor.

El [ensayo opcional con xUnit](tools/fixtures/dotnet-testing/README.md) verifica una regla pura de capacidad y sus efectos públicos con `npm run test:dotnet-unit`. Sus requisitos de NuGet se restauran por separado solo con autorización, dentro del repositorio; la ejecución normal usa fuentes vacías y paquetes locales fijados. No agrega WAF, base de datos ni paquetes de reloj a una regla que no los necesita; tampoco valida por ejecución las nuevas instrucciones de tiempo y cancelación.

### Diagnóstico de errores

Para investigar un defecto, carga [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md) mediante las rutas condicionales de las cuatro guías. Reproduce con datos sintéticos, contrasta hipótesis y verifica la corrección mínima con el perfil de pruebas existente. No instala herramientas ni autoriza repetir datos de producción. `npm run test:debugging-instructions` comprueba instrucciones y rutas, no ejecuta un defecto de una aplicación ni demuestra obediencia del modelo.

### Accesibilidad proporcional al cambio

Para cambios de interacción web o híbrida, carga [accessibility](skills/shared/accessibility/SKILL.md) desde las guías de frontend y apps. Conserva la librería de UI, los formularios y la navegación existentes. Revisa semántica y nombres, teclado y foco, errores y anuncios, contraste, zoom y movimiento en el recorrido afectado; no impone shadcn, plugins ni una auditoría global.

`npm run test:accessibility-instructions` comprueba el contrato y las rutas; forma parte de `test` y `verify` sin dependencias nuevas. Una prueba de jsdom o un escáner no demuestran accesibilidad completa: documenta por separado las pruebas reales de navegador, tecnología de asistencia y dispositivo que se ejecutaron o que faltan. No se instala ninguna herramienta sin autorización.

### Seguridad cotidiana en todos los perfiles

Conserva el dominio elegido junto con `skills/shared/security/`. Las cuatro guías indican cargar `skills/shared/security/baseline.md` al iniciar cada tarea; `application-checks.md` y `delivery-checks.md` se consultan solo según el cambio. El skill de seguridad de frontend contiene recetas de Next.js: una SPA o una app móvil usa los controles compartidos y su diseño de autenticación existente, sin inventar infraestructura de servidor. Si usas una instalación externa, ajusta las rutas para que resuelvan al repositorio de skills.

Esto conserva un núcleo común sin mezclar las convenciones de .NET 8 y .NET 10. No requiere instalar Cyber Neo, ejecutar sus scripts ni realizar una auditoría completa. La carga depende del runtime que use `AGENTS.md`; los enlaces no garantizan activación automática ni aislamiento técnico.

## Paso 3: Valida que la suite funcione

Desde la raíz de **este repositorio** (zesh-one-skills), ejecuta:

```bash
npm install
npm run verify
```

Si todo pasa sin errores, las comprobaciones de estructura e instrucciones de la suite pasan. Si hay fallos, revisa los mensajes. `npm run test:pilot` ejecuta las regresiones de selección y aplicabilidad; también forma parte de `test` y `verify`. Estas comprobaciones no garantizan obediencia del modelo ni seguridad de una aplicación: consulta los [escenarios observados y sus límites](odd/evidence/architecture-adoption-pilot.md).

## Solución de problemas

**El agente no encuentra los skills**
→ Verifica que el `AGENTS.md` copiado en tu proyecto referencia las rutas correctas a este repositorio.

**`npm run verify` falla con errores de YAML**
→ Revisa que los archivos `.yaml` de los skills no tengan indentación mixta (tabulaciones frente a espacios).

**El agente no recuerda decisiones de sesiones anteriores**
→ Engram debe estar activo. Verifica la configuración del runtime de OpenCode.

## Notas

- Esta suite está diseñada para un solo agente (OpenCode). No es multiagente por defecto.
- El alcance principal es backend .NET 8, frontend (Next.js 15, React 19) y apps móviles (Ionic Angular, Capacitor). Los skills están separados por dominio en `skills/backend/`, `skills/frontend/` y `skills/apps/`.
- `skills/shared/` contiene skills transversales; inclúyelos siempre junto con el dominio principal.

## Uso en monorepos o proyectos con varios dominios

Si tu repositorio tiene varias capas (por ejemplo, `/backend` y `/frontend` en el mismo repositorio), puedes:

1. **Tener un único `AGENTS.md` en la raíz** que combine las secciones relevantes de `AGENTS.backend.md` y `AGENTS.frontend.md`.
2. **O usar un `AGENTS.md` por subdirectorio** — OpenCode respeta el `AGENTS.md` más cercano al directorio de trabajo activo.

En ambos casos, las rutas a los skills en `AGENTS.md` deben ser **relativas al archivo** o **absolutas al repositorio de skills**. Ajusta según cómo montes este repositorio (submódulo, dependencia copiada o referencia directa).
