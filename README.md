# Zesh-One-Skills

> Suite de skills para agentes de IA en backend (.NET 8 REST), frontend (Next.js 15 + React 19) y aplicaciones móviles (Ionic + Capacitor), hecha para flujos de trabajo con OpenCode + Engram + Agent Teams Lite.

## Qué es este repositorio

`zesh-one-skills` centraliza guías de ejecución específicas de dominio para agentes de IA.

La suite es intencionalmente **operativa**:
- valores predeterminados concisos,
- restricciones explícitas,
- antipatrones y divulgación progresiva,
- ruido de contexto mínimo.

Todos los archivos de skills y los archivos de dominio `AGENTS.*.md` siguen una política **English-only** (solo inglés).

## Herramientas del entorno

| Herramienta | Función |
|------|------|
| [OpenCode](https://opencode.ai) | Runtime del agente |
| [Engram](https://github.com/gentleman-programming/engram) | Memoria persistente entre sesiones |
| Agent Teams Lite | Orquestación y delegación multiagente |
| Node.js | Canal de validación local (`npm run verify`) |

## Estructura del repositorio

```text
zesh-one-skills/
├── AGENTS.backend.md
├── AGENTS.frontend.md
├── AGENTS.apps.md
├── skills/
│   ├── backend/
│   │   ├── dataaccess/
│   │   ├── general/
│   │   ├── logging/
│   │   ├── mapping/
│   │   ├── performance/
│   │   ├── requests/
│   │   ├── responses/
│   │   ├── security/
│   │   ├── testing-unit/
│   │   └── validations/
│   ├── frontend/
│   │   ├── nextjs-15/
│   │   ├── react-19/
│   │   ├── security/
│   │   ├── shadcn-ui/
│   │   ├── tanstack-query/
│   │   ├── tailwind-4/
│   │   ├── testing/
│   │   ├── typescript/
│   │   ├── zod-4/
│   │   └── zustand-5/
│   ├── apps/
│   │   ├── capacitor/
│   │   └── ionic-angular/
│   └── shared/
│       ├── accessibility/
│       ├── cli-design/
│       ├── diagnosing-bugs/
│       └── github-pr/
├── docs/
│   ├── MAINTAINING.md
│   ├── skills-review-rubric.md
│   └── skills-sanitization-plan.md
├── tools/
├── ADOPTING.md
├── CHANGELOG.md
└── CREDITS.md
```

## Instalador de la CLI en el proyecto

Ejecuta `zesh-skills` desde la raíz de un proyecto consumidor para elegir skills por alcance y tema, con descripciones y sin paquetes obligatorios. Es compatible con Claude Code, Codex, OpenCode y Pi, conserva las instrucciones existentes y lee cada ejecución el catálogo y el contenido más recientes de `main` desde un único commit fijado.

El catálogo ya está publicado en `main`. Descarga el binario de tu plataforma y `SHA256SUMS` desde [v2.9.0](https://github.com/zeshone/zesh-one-skills/releases/tag/v2.9.0); consulta los [enlaces de descarga por plataforma, la verificación del checksum y las instrucciones de `chmod` en Unix y PowerShell en Windows](cli/README.md#descarga-de-v290). Por ejemplo, Linux x64 usa [zesh-skills-linux-amd64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.9.0/zesh-skills-linux-amd64); verifica el archivo contra [SHA256SUMS](https://github.com/zeshone/zesh-one-skills/releases/download/v2.9.0/SHA256SUMS) antes de ejecutar `--version` (valor esperado: `2.9.0`). Si esa página de release todavía no es pública, puedes reproducir los mismos archivos con `npm run build:zesh-cli -- --release 2.9.0`.

El ejecutable no requiere Go, Node, Git ni el runtime de .NET; los skills permanecen en el proyecto. Ejecuta `zesh-skills --install-cli` para copiar el ejecutable a una ubicación del usuario y agregarlo a la `PATH` del usuario sin permisos de administrador; conserva el archivo original y no instala skills. Consulta [instalar la CLI para tu usuario](cli/README.md#instala-la-cli-para-tu-usuario). El asistente usa inglés de forma predeterminada y español con `--es`, con 35 descripciones traducidas incluidas y sin cambiar los identificadores en inglés ni los formatos instalados. Consulta [uso de la CLI, elección de arquitectura y límites de la evidencia](cli/README.md). Las compilaciones para macOS y ARM solo tienen evidencia de compilación, no una verificación de ejecución.

Compila localmente con `npm run build:zesh-cli` (versión `dev`), o prepara los seis archivos versionados con `npm run build:zesh-cli -- --release 2.9.0` en `dist/zesh-skills-release/v2.9.0/`. Ninguno de los dos comandos publica ni instala de forma global.

Mantenedores: ejecuta `npm run catalog:generate` después de cambiar skills. Las comprobaciones actuales, solo con Node, validan la vigencia del catálogo; las comprobaciones ejecutables de Go se ejecutan por separado con `npm run test:zesh-cli`.

## Adopción

Usa la guía de incorporación:

-> **[ADOPTING.md](./ADOPTING.md)**

Para cambios en el repositorio, usa el [flujo de trabajo del mantenedor](docs/MAINTAINING.md).

## Validación

```bash
npm install
npm run verify
```

`verify` ejecuta las pruebas de fixtures, las comprobaciones semánticas, el lint de skills y las comprobaciones de preparación de release. Todo sigue siendo solo con Node y no ejecuta una aplicación consumidora.

Para trabajo en la interfaz de una CLI, carga [cli-design](skills/shared/cli-design/SKILL.md) junto con el dominio que corresponda. `npm run test:cli-instructions` ejecuta sus comprobaciones de instrucciones y rutas, incluidas también en `test` y `verify`; aprobarlas no demuestra obediencia del modelo ni comportamiento ejecutable de la CLI.

Con un SDK de .NET 10 y sus paquetes de referencia ya instalados, `npm run test:cli` ejecuta el [espécimen de aceptación de CLI ejecutable](skills/shared/cli-design/assets/line-tally/README.md). Usa estado desechable del proyecto y ninguna fuente de paquetes externa; no forma parte de la verificación predeterminada y no instala un SDK ni paquetes.

Para elegir pruebas de una React SPA y comportamiento observable de UI, formularios y API, carga [testing de frontend](skills/frontend/testing/SKILL.md). `npm run test:frontend-instructions` comprueba su contrato de instrucciones y su ruta condicional; estas comprobaciones, solo con Node, no renderizan React ni demuestran comportamiento de navegador. Reutiliza las herramientas del consumidor; si faltan requisitos previos se requiere autorización explícita, nunca una instalación automática. La [prueba opcional de React DOM](tools/fixtures/frontend-dom/README.md) se ejecuta con `npm run test:frontend-dom` después de autorizar por separado su preparación aislada; no forma parte de la verificación predeterminada ni es evidencia de navegador real.

Los perfiles existentes de [pruebas unitarias de .NET 8](skills/backend/testing-unit/SKILL.md) y [pruebas de .NET 10](skills/dotnet-10/testing/SKILL.md) cubren límites aplicables, efectos al rechazar y tiempo y cancelación deterministas, sin unificar sus librerías ni agregar infraestructura. `npm run test:dotnet-instructions` comprueba estos contratos de instrucciones; no ejecuta xUnit. La [prueba aislada opcional con xUnit](tools/fixtures/dotnet-testing/README.md) se ejecuta con `npm run test:dotnet-unit` después de autorizar explícitamente la preparación local de NuGet; la ejecución ordinaria nunca descarga paquetes.

Para investigar defectos en todos los perfiles, carga [diagnosing-bugs](skills/shared/diagnosing-bugs/SKILL.md): reproduce, separa hipótesis de observaciones y verifica la corrección mínima. `npm run test:debugging-instructions` comprueba sus instrucciones y rutas, no un defecto del consumidor ni el cumplimiento del agente; consulta la [evidencia de la propuesta acotada](odd/evidence/debugging-accessibility-adoption.md).

Para cambios de interacción web o híbrida, carga [accessibility](skills/shared/accessibility/SKILL.md) a través de la guía de frontend o apps. `npm run test:accessibility-instructions` comprueba las instrucciones semánticas, de teclado y foco, y de retroalimentación, además de las rutas; no ejecuta un navegador, una tecnología de asistencia ni una auditoría de cumplimiento. Conserva las librerías de UI existentes y usa recorridos reales del consumidor para esas comprobaciones; ninguna instalación de herramientas es automática.

## Changelog

-> **[CHANGELOG.md](./CHANGELOG.md)**

## Créditos

-> **[CREDITS.md](./CREDITS.md)**
