# Adopting Zesh-One-Skills

Guía para integrar esta suite de skills en un proyecto nuevo que use OpenCode como agente AI.

## Prerequisitos

Antes de adoptar, verificá que tenés:

- [ ] [OpenCode](https://opencode.ai) instalado y configurado
- [ ] [Engram](https://github.com/gentleman-programming/engram) disponible en el entorno del agente
- [ ] Node.js ≥ 18 (para ejecutar `npm run verify` localmente si querés validar la suite)

## Paso 1: Copiá el AGENTS.md correspondiente

En `agents-files/` encontrás plantillas de `AGENTS.md` para distintos tipos de proyecto:

| Archivo | Cuándo usarlo |
|---------|---------------|
| `Zesh-custom-AGENTS.md` | Proyectos nuevos bajo el stack Zesh (recomendado como punto de partida) |
| `Prowler-template-AGENTS.md` | Proyectos con estructura tipo Prowler |

**Acción:** Copiá el archivo que aplique a la raíz de tu proyecto como `AGENTS.md`.

```bash
cp agents-files/Zesh-custom-AGENTS.md /ruta/a/tu-proyecto/AGENTS.md
```

> El `AGENTS.md` es el punto de entrada que OpenCode lee al iniciar sesión. Sin él, el agente opera sin contexto del proyecto.

## Paso 2: Revisá qué skills están disponibles

Los skills están organizados por dominio:

- `backend/skills/net8-apirest/` — Arquitectura, data access, endpoints REST para .NET 8
- `shared/skills/github-pr/` — Convenciones de Pull Request y branch workflow
- `vendor/gentleman/` — Skills del ecosistema Gentleman (nextjs-15, react-19, tailwind-4, etc.)

Los skills en `.github/skills/` (vendor Gentleman) son herramientas del agente disponibles bajo demanda — no se auto-invocan, se usan explícitamente.

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
- Los skills vendor (`vendor/gentleman`) no participan en auto-invoke; se referencian explícitamente.
- El alcance principal es backend .NET 8. Los skills de frontend (nextjs-15, react-19, etc.) son referencias del ecosistema Gentleman y se adoptan bajo criterio propio.
