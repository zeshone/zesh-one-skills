# Zesh-One-Skills

> Suite de skills para agentes AI — Backend .NET 8 REST API y Frontend, con soporte a herramientas [OpenCode](https://opencode.ai) + [Engram](https://github.com/gentleman-programming/engram) + [Agents Team Lite](https://github.com/gentleman-programming).

## Qué es esto

Repositorio privado de skills e instrucciones de contexto para agentes AI utilizados en proyectos backend .NET 8 y frontend. Centraliza las convenciones, reglas de arquitectura y guías de desarrollo que el agente debe conocer para trabajar correctamente en proyectos Zesh.

Los skills están diseñados para **OpenCode** con el stack de herramientas de Gentleman Programming: engram (memoria persistente) y agents-team-lite (orquestación de sub-agentes).

## Stack y herramientas

| Herramienta | Rol |
|-------------|-----|
| [OpenCode](https://opencode.ai) | Runtime del agente |
| [Engram](https://github.com/gentleman-programming/engram) | Memoria persistente entre sesiones |
| Agents Team Lite | Orquestación y delegación de sub-agentes |
| Node.js | Pipeline de validación local (`npm run verify`) |

## Estructura del repo

```
zesh-one-skills/
├── AGENTS.backend.md      # Instrucciones de agente para proyectos backend (distribuir como AGENTS.md)
├── AGENTS.frontend.md     # Instrucciones de agente para proyectos frontend (distribuir como AGENTS.md)
├── skills/
│   ├── backend/           # Skills de dominio backend (.NET 8 REST API)
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
│   ├── frontend/          # Skills de dominio frontend (Next.js 15, React 19, etc.)
│   │   ├── nextjs-15/
│   │   ├── react-19/
│   │   ├── tailwind-4/
│   │   ├── typescript/
│   │   ├── zod-4/
│   │   └── zustand-5/
│   └── shared/            # Skills cross-domain
│       └── github-pr/
├── tools/                 # Scripts de validación (lint, test, verify)
├── ADOPTING.md            # Guía para adoptar esta suite en un proyecto nuevo
├── CREDITS.md             # Atribución de assets vendor
└── USER_DECISIONS.md      # Decisiones de arquitectura del proyecto
```

## Quick Start

Para adoptar esta suite en un proyecto nuevo, seguí la guía paso a paso:

→ **[ADOPTING.md](./ADOPTING.md)**

## Verificación

Para validar que la suite está consistente y los skills tienen el formato correcto:

```bash
npm install
npm run verify
```

El comando ejecuta lint de estructura de skills, prueba fixtures y valida logging specs. Debe terminar con 0 errores.

## Créditos

→ **[CREDITS.md](./CREDITS.md)**
