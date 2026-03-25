# Zesh-One-Skills

> Suite de skills para agentes AI — Backend .NET 8 REST API, con soporte a herramientas [OpenCode](https://opencode.ai) + [Engram](https://github.com/gentleman-programming/engram) + [Agents Team Lite](https://github.com/gentleman-programming).

## Qué es esto

Repositorio privado de skills e instrucciones de contexto para agentes AI utilizados en proyectos backend .NET 8. Centraliza las convenciones, reglas de arquitectura y guías de desarrollo que el agente debe conocer para trabajar correctamente en proyectos Zesh.

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
├── agents-files/          # Plantillas de AGENTS.md por tipo de proyecto
├── backend/               # Assets específicos del dominio backend
│   ├── skills/            # Skills backend (net8-apirest, etc.)
│   └── AGENTS.md          # Instrucciones de agente para proyectos backend
├── frontend/              # Assets específicos del dominio frontend
│   └── AGENTS.md          # Instrucciones de agente para proyectos frontend
├── vendor/                # Skills de fuentes externas
│   └── gentleman/         # Skills de gentleman-programming (nextjs-15, react-19, etc.)
├── shared/                # Utilities cross-domain
│   └── skills/            # Skills compartidos (github-pr, etc.)
├── rules-to-skills/       # Documento auxiliar de reglas .NET (referencia)
├── .github/skills/        # Skills de agente (vendor: gentleman-programming)
│   └── skill-creator/     # Skill para crear nuevas skills
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
