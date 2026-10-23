# Credits

Atribución de assets, skills y herramientas de terceros presentes en este repositorio.

## Assets Vendor

| Asset | Autor | Licencia | Ubicación |
|-------|-------|----------|-----------|
| `github-pr` | [gentleman-programming](https://github.com/gentleman-programming) | — ¹ | `skills/shared/github-pr/` |
| `skill-creator` | [gentleman-programming](https://github.com/gentleman-programming) | — ¹ | `.github/skills/skill-creator/` |

> ¹ **Nota sobre licencia vendor**: La licencia de los assets originales de `gentleman-programming` no está confirmada al momento de esta release interna. Se usan con atribución explícita. Esta nota deberá actualizarse cuando se confirme la licencia upstream.

## Autoría Zesh

Los siguientes skills son de autoría propia del proyecto `zesh-one`:

| Dominio | Ubicación |
|---------|-----------|
| Backend (.NET 8 REST API) | `skills/backend/` |
| Frontend (Next.js 15, React 19, Tailwind 4, Zod 4, Zustand 5, TanStack Query, shadcn/ui, security, TypeScript) | `skills/frontend/` |
| Apps móviles (Ionic Angular, Capacitor) | `skills/apps/` |
| Contratos de interacción CLI | `skills/shared/cli-design/` |
| Diagnóstico de errores | `skills/shared/diagnosing-bugs/` |
| Accesibilidad web e híbrida | `skills/shared/accessibility/` |

Los skills de `skills/frontend/` tienen inspiración conceptual en los originales de `gentleman-programming` pero son obras derivadas con modificaciones sustanciales — arquitectura, convenciones, ejemplos y reglas propias de Zesh-One.

La guía `frontend/testing` es de autoría del proyecto: criterios de evidencia y comportamiento observable adaptados a sus rutas SPA, seguridad y convenciones existentes; no importa una colección externa ni obliga a instalar un runner.

La guía `cli-design` es de autoría del proyecto y reutiliza sus límites de arquitectura, seguridad y pruebas. No incorpora contenido, scripts ni paquetes de una colección CLI externa.

La guía `diagnosing-bugs` es propia del proyecto y articula sus reglas existentes de seguridad, diagnóstico y regresión. No se importó una skill externa ni se inspeccionaron nuevas fuentes externas en esta fase.

La guía `accessibility` también es propia del proyecto: reúne criterios de interacción y límites de evidencia ya presentes en sus guías UI/testing, sin importar contenido externo, añadir herramientas ni certificar conformidad.

## Referencias conceptuales del piloto de adopción

El piloto adapta ideas puntuales a las reglas existentes; no importa colecciones, scripts ni dependencias de estas fuentes.

| Referencia | Uso acotado | Licencia declarada |
|------------|-------------|--------------------|
| [dotnet-claude-kit / architecture-advisor](https://github.com/codewithmukesh/dotnet-claude-kit) — Mukesh Murugan | Contexto, alternativas y costos antes de elegir arquitectura; sin umbrales rígidos | MIT |
| [mattpocock/skills](https://github.com/mattpocock/skills) — Matt Pocock | Abstracciones justificadas y referencias de instrucciones según contexto | MIT |
| [Cyber Neo](https://github.com/Hainrixz/cyber-neo) — Hainrixz | Inspiración del enriquecimiento de seguridad ya integrado; sin importar código ni flujo de auditoría | Sin nueva conclusión de licencia en este piloto |
| [midudev/autoskills](https://github.com/midudev/autoskills) | Descubrimiento de candidatos, sin incorporar contenido | CC BY-NC 4.0 en la raíz; las fuentes registradas varían |

Las revisiones consultadas, prácticas aceptadas/rechazadas y límites de evidencia están en el [registro del piloto](odd/evidence/architecture-adoption-pilot.md).

El enriquecimiento de las guías de testing .NET es propio del proyecto: escenarios de límites, efectos y tiempo/cancelación sobre sus reglas existentes. No importa otra colección ni unifica las convenciones de .NET 8 y .NET 10.

## Inspiraciones de las skills de apps

El metadato `inspired-by` de las skills de apps registra estas etiquetas. No se registró una URL de origen verificada, por lo que se conservan como identificadores sin enlace.

| Referencia | Uso registrado | Revisión inspeccionada | Licencia upstream |
|------------|----------------|------------------------|-------------------|
| `capawesome-team/ionic-angular` | Inspiración nombrada en [ionic-angular](skills/apps/ionic-angular/SKILL.md) | not recorded | not recorded |
| `ionic-capacitor/angular.md` | Etiqueta adicional en [ionic-angular](skills/apps/ionic-angular/SKILL.md) | not recorded | not recorded |
| `capawesome-team/ionic-capacitor` | Inspiración nombrada en [capacitor](skills/apps/capacitor/SKILL.md) | not recorded | not recorded |

El metadato no establece si el uso fue copiado o conceptual, ni la fuente exacta de `ionic-capacitor/angular.md`. Estas etiquetas y la declaración Apache-2.0 de las skills locales no permiten inferir una revisión upstream, una conclusión de licencia ni una inspección histórica.

## Dependencias del ensayo React DOM

El espécimen propio de `tools/fixtures/frontend-dom/` utiliza React/react-dom, Testing Library (React/DOM) y jsdom solo en su proyecto de pruebas opcional. Sus versiones, dependencias transitivas, licencias declaradas e integridad de descarga quedan en su `package-lock.json`; no se incorporan sus fuentes a la colección ni se añaden al paquete raíz. Esto no modifica la atribución de las skills ni implica una auditoría de licencias o seguridad de dependencias.

## Dependencias del ensayo xUnit

El espécimen propio de `tools/fixtures/dotnet-testing/` utiliza xunit, xunit.runner.visualstudio y Microsoft.NET.Test.Sdk, con sus dependencias transitivas, solo en el ensayo opcional. Su lockfile registra versiones y hashes; no incorpora sus fuentes ni modifica las dependencias raíz. La restauración desde nuget.org fue autorizada por separado, sin auditoría de seguridad o licencias implícita.

## Herramientas del ecosistema

| Herramienta | Autor | Uso en este repo |
|-------------|-------|-----------------|
| OpenCode | opencode.ai | Runtime del agente |
| Engram | gentleman-programming | Memoria persistente entre sesiones |
| Agents Team Lite | gentleman-programming | Orquestación de sub-agentes |
