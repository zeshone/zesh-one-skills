# Adopting Zesh-One-Skills

Guía para integrar esta suite de skills en un proyecto nuevo que use OpenCode como agente AI.

## Prerequisitos

Antes de adoptar, verificá que tenés:

- [ ] [OpenCode](https://opencode.ai) instalado y configurado
- [ ] [Engram](https://github.com/gentleman-programming/engram) disponible en el entorno del agente
- [ ] Node.js ≥ 18 (para ejecutar `npm run verify` localmente si querés validar la suite)

## Paso 1: Copiá el AGENTS.md correspondiente

En la raíz del repositorio encontrás los archivos de instrucciones de agente por dominio:

| Archivo | Cuándo usarlo |
|---------|---------------|
| `AGENTS.backend.md` | Proyectos backend .NET 8 REST API (recomendado como punto de partida) |
| `AGENTS.dotnet-10.md` | Proyectos backend .NET 10 / C# 14 Minimal API |
| `AGENTS.frontend.md` | Proyectos frontend (Next.js 15, React 19, etc.) |

**Acción:** Copiá el archivo que aplique a la raíz de tu proyecto y renombralo a `AGENTS.md`.

```bash
# Para proyectos backend:
cp AGENTS.backend.md /ruta/a/tu-proyecto/AGENTS.md

# Para proyectos frontend:
cp AGENTS.frontend.md /ruta/a/tu-proyecto/AGENTS.md
```

> El `AGENTS.md` es el punto de entrada que OpenCode lee al iniciar sesión. Sin él, el agente opera sin contexto del proyecto.

## Paso 2: Copiá los skills relevantes

Los skills están organizados en cinco dominios bajo `skills/`:

- `skills/backend/` — Arquitectura, data access, logging, responses, security, testing y más para .NET 8 REST API
- `skills/dotnet-10/` — Minimal APIs, arquitectura, persistencia, seguridad y pruebas para .NET 10 / C# 14
- `skills/frontend/` — Next.js 15, React 19, Tailwind 4, TypeScript, Zod 4, Zustand 5
- `skills/apps/` — Mobile native layer (Ionic Angular, Capacitor)
- `skills/shared/` — Skills cross-domain (github-pr, convenciones de Pull Request)

**Acción:** Copiá el directorio del dominio que necesitás (más `skills/shared/`) a tu proyecto o referenciá este repositorio directamente desde tu `AGENTS.md`.

### Seguridad diaria en backend

Conserva la estructura `skills/backend/` o `skills/dotnet-10/` junto con `skills/shared/security/`. Ambas guías backend indican cargar `skills/shared/security/baseline.md` al iniciar cada tarea; las referencias `application-checks.md` y `delivery-checks.md` se consultan solo según el cambio. Si usas una instalación externa, ajusta las rutas de la guía para que resuelvan al repositorio de skills.

Esto conserva un núcleo común sin mezclar las convenciones de .NET 8 y .NET 10. No requiere instalar Cyber Neo, ejecutar sus scripts ni realizar una auditoría completa. La carga depende del runtime que use `AGENTS.md`; los enlaces no garantizan activación automática ni aislamiento técnico.

## Paso 3: Validá que la suite esté funcional

Desde la raíz de **este repositorio** (zesh-one-skills), ejecutá:

```bash
npm install
npm run verify
```

Si todo pasa sin errores, la suite está lista. Si hay fallos, revisá los mensajes — suelen indicar skills con estructura inválida o fixtures rotos.

## Troubleshooting

**El agente no encuentra las skills**
→ Verificá que el `AGENTS.md` copiado en tu proyecto referencia las rutas correctas a este repositorio.

**`npm run verify` falla con errores de YAML**
→ Revisá que los archivos `.yaml` de los skills no tengan indentación mixta (tabs vs espacios).

**El agente no recuerda decisiones de sesiones anteriores**
→ Engram debe estar activo. Verificá la configuración del runtime de OpenCode.

## Notas

- Esta suite está diseñada para un solo agente (OpenCode). No es multi-agente por defecto.
- El alcance principal es backend .NET 8, frontend (Next.js 15, React 19) y apps móviles (Ionic Angular, Capacitor). Los skills están separados por dominio en `skills/backend/`, `skills/frontend/` y `skills/apps/`.
- `skills/shared/` contiene skills cross-domain — inclúyelos siempre junto con el dominio principal.

## Uso en monorepos o proyectos multi-dominio

Si tu repositorio tiene múltiples capas (por ejemplo, `/backend` y `/frontend` en el mismo repo), podés:

1. **Tener un único `AGENTS.md` en la raíz** que combine las secciones relevantes de `AGENTS.backend.md` y `AGENTS.frontend.md`.
2. **O usar un `AGENTS.md` por subdirectorio** — OpenCode respeta el `AGENTS.md` más cercano al directorio de trabajo activo.

En ambos casos, las rutas a los skills en `AGENTS.md` deben ser **relativas al archivo** o **absolutas al repositorio de skills**. Ajustá según cómo montés este repositorio (submódulo, dependencia copiada, o referencia directa).
