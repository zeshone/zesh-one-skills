# Instala solo los skills que tu proyecto necesita

`zesh-skills` es un instalador de skills independiente y local al proyecto para Claude Code, Codex, OpenCode y Pi. El ejecutable puede vivir en tu `PATH`; los skills nunca se instalan en la configuración de agentes de tu directorio personal.

## Descarga de v2.8.0

Descarga un ejecutable sin empaquetar y el archivo [SHA256SUMS](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/SHA256SUMS) desde el release [v2.8.0](https://github.com/zeshone/zesh-one-skills/releases/tag/v2.8.0). Si esa página todavía no es pública, puedes reproducir los mismos seis archivos desde el commit de release con `npm run build:zesh-cli -- --release 2.8.0`; una compilación de desarrollo informa `dev`, así que revisa `--version` antes de confiar en un binario.

| Plataforma | Archivo |
| --- | --- |
| Windows x64 | [zesh-skills-windows-amd64.exe](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-windows-amd64.exe) |
| Windows ARM64 | [zesh-skills-windows-arm64.exe](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-windows-arm64.exe) |
| Linux x64 | [zesh-skills-linux-amd64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-linux-amd64) |
| Linux ARM64 | [zesh-skills-linux-arm64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-linux-arm64) |
| macOS Intel | [zesh-skills-darwin-amd64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-darwin-amd64) |
| macOS Apple Silicon | [zesh-skills-darwin-arm64](https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0/zesh-skills-darwin-arm64) |

Ejemplo para Linux (reemplaza el nombre del archivo por el de tu plataforma):

```bash
base=https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0
curl -fLO "$base/zesh-skills-linux-amd64"
curl -fLO "$base/SHA256SUMS"
grep '  zesh-skills-linux-amd64$' SHA256SUMS | sha256sum --check -
# Continúa solo si el checksum coincide.
chmod +x zesh-skills-linux-amd64
./zesh-skills-linux-amd64 --version  # valor esperado: 2.8.0
```

En macOS usa `shasum -a 256 -c` en lugar de `sha256sum --check -`, con el nombre de archivo de macOS seleccionado. Verifica solo la línea del archivo descargado; el manifiesto también lista las otras cinco plataformas.

Ejemplo para Windows PowerShell:

```powershell
$base = 'https://github.com/zeshone/zesh-one-skills/releases/download/v2.8.0'
$name = 'zesh-skills-windows-amd64.exe'
Invoke-WebRequest "$base/$name" -OutFile $name
Invoke-WebRequest "$base/SHA256SUMS" -OutFile SHA256SUMS
$line = Get-Content SHA256SUMS | Where-Object { $_.EndsWith("  $name") }
if (@($line).Count -ne 1) { throw 'Missing or duplicate checksum' }
$expected = ($line -split '  ')[0]
if ((Get-FileHash $name -Algorithm SHA256).Hash.ToLowerInvariant() -ne $expected) { throw 'Checksum mismatch' }
& ".\$name" --version  # valor esperado: 2.8.0
```

Los checksums detectan bytes que no coinciden, no son una firma independiente de la identidad del publicador. Después de verificar, opcionalmente cambia el nombre del ejecutable a `zesh-skills` (`zesh-skills.exe` en Windows) y sigue la guía manual de `PATH` que aparece más abajo. No se necesita Go, Node, Git ni el runtime de .NET para ejecutarlo. Los skills permanecen en el proyecto. Los archivos de macOS y ARM solo tienen evidencia de compilación: la matriz de CI ejecuta las pruebas de Go en las tres plataformas, no los binarios publicados.

## Ruta rápida

```text
cd tu-proyecto
zesh-skills
# Interfaz en español; el inglés sigue siendo el idioma predeterminado sin esta opción.
zesh-skills --es
```

### Idioma de la interfaz

Usa `--es` para traducir al español los mensajes, las opciones de alcance y arquitectura, las descripciones, los controles del selector, las vistas previas, las confirmaciones y los diagnósticos de la aplicación. No hay una pregunta adicional al iniciar ni una preferencia de idioma guardada. Las confirmaciones en español aceptan `s`, `si`, `sí`, `y` o `yes`; el comportamiento en inglés no cambia. `--es=false` mantiene explícitamente el inglés.

Los nombres de arquitectura siguen siendo reconocibles en español: permanecen en inglés o acompañan su traducción, por ejemplo `Arquitectura Limpia (Clean Architecture)`. Los nombres híbridos como `CQRS + event sourcing` se mantienen en inglés para que el menú quede compacto. Las explicaciones siguen en español; términos conocidos como UI, view models, focus, feedback y triggers conservan su nombre en inglés.

En el catálogo de skills (`--list` y el selector interactivo), **solo las descripciones** se traducen: los nombres e identificadores de skills y los encabezados y etiquetas de tema conservan su inglés original. Con `--es`, tiene prioridad un `descriptionEs` remoto validado. El binario también incluye 35 traducciones al español para el catálogo original, usadas solo cuando coinciden el identificador y la descripción original exacta en inglés. Los skills desconocidos o las descripciones modificadas sin traducción vuelven al inglés de forma honesta y siguen siendo detectables. Sin `--es` (o con `--es=false`), las descripciones quedan exactamente como se publicaron en inglés.

La opción de idioma no traduce flags, formatos JSON, valores personalizados del usuario, contenido de skills descargado ni documentos de agente generados. Las descripciones traducidas son solo de presentación: el manifiesto y las tablas de instrucciones instaladas conservan los metadatos originales canónicos, sin `descriptionEs`. Los detalles técnicos del sistema operativo pueden conservar su idioma original y se identifican como tales. Los perfiles de arquitectura predefinidos se guardan por identificador estable, sea cual sea el idioma mostrado.

El asistente confirma el directorio actual como raíz del proyecto, detecta manifiestos de framework y archivos de instrucciones existentes, pregunta la arquitectura y luego abre un selector de casillas controlado por teclado. Cada skill relevante sigue disponible con su tema y una descripción breve; no hace falta escribir identificadores de skills.

| Tecla | Acción |
| --- | --- |
| Arriba/Abajo o `j`/`k` | Recorrer la lista desplazable de skills |
| Espacio | Marcar o desmarcar el skill enfocado |
| Enter | Aceptar la selección actual |
| `a` / `n` / `r` | Seleccionar todo / limpiar la selección / restaurar sugerencias |
| `q` o Ctrl-C | Cancelar sin instalar nada |

Un conjunto pequeño de sugerencias aparece marcado y señalado con `*`. Las sugerencias incluyen diagnóstico independiente del framework, accesibilidad para frontend y apps, y fundamentos de framework y diseño cuando las declaraciones de paquetes aportan evidencia. Se sugieren como máximo cinco. Los skills versionados requieren la misma versión mayor declarada del paquete; los rangos desconocidos y las versiones conflictivas de un monorepo no seleccionan un perfil versionado. La falta de evidencia nunca fuerza .NET 8/10 ni una nueva librería de UI. Todas las sugerencias se pueden desmarcar; los skills de seguridad y acceso a datos no son obligatorios. Para un mockup de landing page, conserva solo los skills de diseño que quieras. Las descripciones y los motivos de sugerencia permanecen visibles, y la pantalla, el cursor y el modo de entrada se restauran al confirmar, cancelar, ante un fallo de entrada o ante una terminación SIGTERM capturable. Una terminación forzada del proceso no puede ejecutar la limpieza de la terminal. El alcance usa el mismo menú de casillas y puede combinar frontend, backend y app. La arquitectura usa un menú de selección exclusiva; la configuración del agente conserva su pregunta breve.

### Perfiles de arquitectura

El menú de arquitectura incluye 20 opciones explicadas: monolito mínimo, por capas/N-Tier, vertical slices, Clean, Hexagonal, Onion, monolito modular, microservicios, dirigida por eventos, funciones serverless, MVC, MVVM, UI basada en componentes y UI basada en funcionalidades. Los híbridos habituales incluyen monolito modular + vertical slices, Clean + vertical slices, Hexagonal + DDD, microservicios + eventos y CQRS + event sourcing. **Arquitectura personalizada** acepta tu propio texto.

Cada opción explica su alcance y sus costos. Los estilos de despliegue y las combinaciones de patrones se identifican como tales, no se presentan como diseños de dominio intercambiables. Una selección integrada guarda el identificador estable de su perfil; el texto personalizado se conserva. Elegir una opción no genera estructura de código, no agrega infraestructura ni impone un paquete de skills. Los perfiles se mantienen de forma centralizada en [`profiles.go`](profiles.go).

La vista previa final lista el commit de origen, los identificadores seleccionados, los directorios de destino y los archivos de instrucciones. Responde `y` para instalar o `n` para dejar el proyecto sin cambios. Confirmar una selección vacía en el selector también cancela. La arquitectura acepta las opciones mostradas o tu propio texto; no genera ni altera la arquitectura de tu aplicación. Limpiar las opciones no es una desinstalación: las selecciones instaladas antes se conservan.

```bash
zesh-skills --scope frontend --list
zesh-skills --help

# Automatización: opciones explícitas; sin preguntas, incluso ante conflictos locales.
zesh-skills --non-interactive --agents pi --architecture mockup \
  --skills frontend/tailwind-4 --yes

# Proyecto mixto; varios destinos de agente y un archivo de instrucciones adicional.
zesh-skills --scope frontend,backend --agents claude,codex \
  --architecture "modular monolith" --skills frontend/react-19,backend/general \
  --instructions GEMINI.md --yes
```

`--root` reemplaza el directorio actual. `--yes` aprueba la selección mostrada, no el reemplazo de ediciones locales. En una terminal todavía puede preguntar por archivos modificados localmente; en automatización, un conflicto falla sin cambiar los datos instalados. `--list`, `--help` y `--version` nunca instalan nada; ayuda y versión no acceden a la red. Códigos de salida: `0` éxito o cancelación deliberada, `1` fallo de red, de sistema de archivos o conflicto, `2` argumentos inválidos o entrada incompleta.

## Instala la CLI para tu usuario

`--install-cli` copia el ejecutable en ejecución a un directorio del usuario y deja el comando disponible desde cualquier terminal, sin permisos de administrador:

```bash
zesh-skills --install-cli          # vista previa y luego confirmación
zesh-skills --install-cli --yes    # automatización: confirma la vista previa mostrada arriba
zesh-skills --es --install-cli     # vista previa y mensajes en español
```

| Plataforma | Ejecutable | Configuración persistente |
| --- | --- | --- |
| Windows | `%LOCALAPPDATA%\Programs\zesh-skills\zesh-skills.exe` | `PATH` del usuario en `HKCU\Environment`, conservando el valor, el orden y el tipo de valor del registro existentes |
| Linux/macOS | `~/.local/bin/zesh-skills` | Archivo de inicio detectado de `sh`, `bash`, `zsh` o `fish`, con un bloque marcado agregado que mantiene intacto el resto del archivo |

Solo cambia el nivel del usuario: no hay `PATH` de toda la máquina, ni elevación, ni eliminación o movimiento del ejecutable original, ni skills instalados globalmente. `--install-cli` nunca descarga nada, así que funciona sin conexión y no lee el catálogo de un proyecto.

La vista previa lista el ejecutable de origen, el destino, los archivos de inicio y si cambia la `PATH` del usuario. No se escribe nada antes de la confirmación. `--yes` confirma esa vista previa, pero nunca reemplaza un ejecutable existente **distinto**: ese caso requiere una confirmación interactiva de reemplazo y, si no la hay, termina con código `1` y el destino intacto. Volver a ejecutar con el mismo ejecutable es idempotente y no escribe una entrada duplicada en la `PATH`.

Se rechazan antes de cualquier escritura: un origen ausente o que no sea un archivo regular, rutas de destino relativas o con delimitadores, destinos o directorios superiores que sean enlaces simbólicos, caracteres de control, una `SHELL` distinta de `sh`, `bash`, `zsh` o `fish`, un `ZDOTDIR`/`XDG_CONFIG_HOME` fuera del directorio personal, y un valor de `PATH` del usuario cuyo tipo no sea una cadena. Los archivos que cambian entre la vista previa y la escritura abortan la instalación. Los errores de escritura ordinarios restauran el binario, el archivo de inicio y la `PATH` anteriores; la reversión se niega a sobrescribir algo que cambió después de la instalación.

Los cambios de `PATH` se aplican **solo en terminales nuevas**; reinicia el shell o abre una terminal nueva para usar `zesh-skills`. `--help` y `--version` mantienen su precedencia y siguen sin escribir nada.

Límites de la evidencia: las pruebas unitarias y las pruebas de flujo de la CLI pública se ejecutan contra directorios personales temporales y una `PATH` de usuario falsa respaldada por un registro simulado, así que ninguna prueba lee ni escribe el estado real del usuario ni la `PATH` de la máquina. El adaptador del registro de Windows se compila para los dos destinos de Windows, pero se ejercita mediante ese almacén falso, y la instalación en macOS y Linux sigue siendo evidencia de compilación, en línea con los límites de plataforma existentes.

## Compila un ejecutable sin dependencias

El desarrollo requiere Go 1.25+; `golang.org/x/term` y su soporte de plataforma se descargan como módulos de Go fijados en la primera compilación. Se compilan dentro del ejecutable: el binario instalado no requiere Go, Node, Git ni el runtime de .NET. Se necesita acceso a la red para instalar o actualizar skills. El selector es compatible con las consolas de Windows y Windows Terminal y con terminales ANSI en Linux y macOS; cuando la entrada y la salida están redirigidas se usa el contrato de opciones no interactivas en lugar de cambiar los modos de la consola.

```bash
# Solo Go; el soporte de terminal se compila dentro del ejecutable.
CGO_ENABLED=0 go build -trimpath -o dist/zesh-skills ./cli

# Los scripts de conveniencia del repositorio usan Node solo durante el desarrollo.
npm run build:zesh-cli
npm run build:zesh-cli -- --all
# Preparación de release para el mantenedor; versión explícita, sin publicación:
npm run build:zesh-cli -- --release 2.8.0
```

El generador de conveniencia escribe `dist/zesh-skills/<os>-<arch>/zesh-skills[.exe]` y `SHA256SUMS`. `--all` compila de forma cruzada para Windows, Linux y macOS en amd64 y arm64 con CGO desactivado. La compilación cruzada no es una prueba de ejecución en esas plataformas. Las compilaciones predeterminada y `--all` conservan la versión de desarrollo `dev`. `--release <version>` acepta solo un `major.minor.patch` numérico simple (sin prefijo `v`, sin sufijos de preestreno ni de compilación); compila los seis destinos con `CGO_ENABLED=0`, `GOTOOLCHAIN=local` y `main.version` del enlazador fijado con ese único argumento. Para `2.8.0`, la salida es `dist/zesh-skills-release/v2.8.0/`: los seis nombres planos de arriba más `SHA256SUMS` (SHA-256 en minúsculas, dos espacios, el nombre real del archivo y un salto de línea). Solo esos siete archivos explícitos pertenecen a una carga; los archivos ajenos o antiguos no se incluyen ni se borran. Una compilación fallida no es un release publicable. Compila desde el commit final autorizado de `main` y verifica antes de publicar; este comando nunca crea etiquetas ni sobrescribe releases remotas. No existe un paso de publicación o instalación global automática.

Para tener `zesh-skills` disponible desde cualquier terminal, ejecuta `zesh-skills --install-cli` (ver arriba), o copia **el ejecutable de tu sistema operativo y arquitectura** a un directorio que ya esté en tu `PATH` y agrega ese directorio a la `PATH` del usuario. En Linux y macOS, `~/.local/bin` es una opción convencional (créalo, copia el binario, hazlo ejecutable y agrégalo a la `PATH` si hace falta). En Windows, elige una carpeta de herramientas del usuario y agrégala desde la configuración de variables de entorno del usuario; después abre una terminal nueva. No agregues un checkout del proyecto con permisos de escritura a una `PATH` de toda la máquina. Estas acciones manuales son tuyas; `--install-cli` realiza solo la versión de nivel de usuario descrita arriba.

## Archivos del proyecto y comportamiento de actualización

| Ubicación | Propósito |
| --- | --- |
| `.zesh-skills/content/skills/<domain>/<name>/` | Instrucciones de origen seleccionadas y sus recursos incluidos |
| `.zesh-skills/manifest.json` | Selecciones instaladas, último commit de origen y hashes de contenido |
| `.claude/skills/<domain>-<name>/SKILL.md` | Entrada de descubrimiento para Claude Code |
| `.agents/skills/<domain>-<name>/SKILL.md` | Entrada de descubrimiento para Codex |
| `.opencode/skills/<domain>-<name>/SKILL.md` | Entrada de descubrimiento para OpenCode |
| `.pi/skills/<domain>-<name>/SKILL.md` | Entrada de descubrimiento para Pi |
| `CLAUDE.md`, `AGENTS.md` | Índice administrado de skills instalados, descripciones y disparadores de uso |

Las entradas de descubrimiento enlazan a las instrucciones locales canónicas y conservan nombres calificados por dominio como `frontend-security` y `backend-security`. Elegir varios agentes crea sus entradas nativas; los agentes que además exploran directorios de compatibilidad pueden ver entradas duplicadas equivalentes. Los nombres y el contenido originales de los skills de origen permanecen en el almacenamiento canónico.

Las variantes de mayúsculas y minúsculas existentes, como `claude.md`, se reutilizan y no se duplican. El Markdown ajeno a `<!-- zesh-skills:start -->` / `<!-- zesh-skills:end -->` se conserva byte por byte, incluidos sus finales de línea. Los marcadores malformados o duplicados detienen la instalación. `GEMINI.md` y `copilot-instructions.md` en la raíz se detectan y se pueden elegir como archivo de instrucciones adicional; los archivos Markdown arbitrarios no se adivinan ni se reescriben. La configuración anidada personalizada de un agente no se modifica.

Ejecuta el asistente otra vez para agregar o actualizar una selección. Los skills y destinos de agente instalados antes se conservan; omitirlos no es una desinstalación. Los archivos sin modificar se actualizan automáticamente; las ediciones locales y los archivos en colisión sin propietario requieren confirmación explícita. Rechazar cualquier conflicto cancela toda la fase de escritura. Quitar skills o agentes no está implementado de forma intencional en esta primera versión.

Las descargas resuelven `main` **una sola vez** mediante la API pública de GitHub. El catálogo y cada archivo seleccionado provienen exactamente de ese commit, con validación SHA-256 proporcionada por el catálogo. Los archivos se obtienen antes de escribir. Las referencias a skills no seleccionados o a material exclusivo del repositorio se convierten en enlaces de GitHub fijados al commit, en lugar de instalaciones implícitas de skills. Las referencias de soporte compartidas que requiere un skill seleccionado permanecen locales. Los enlaces remotos requieren conexión al abrirlos; no se promete un cierre de dependencias totalmente sin conexión.

Las rutas se validan para ser portables, se rechazan los destinos que son enlaces simbólicos y las operaciones con rutas absolutas confinan las escrituras. Un bloqueo de proyecto evita ejecuciones simultáneas del instalador. Los fallos de escritura ordinarios activan la reversión; una pérdida de energía o la terminación del proceso no es una transacción atómica. Si se interrumpe, revisa `.zesh-skills.lock` y los archivos instalados antes de quitar manualmente un bloqueo obsoleto y volver a ejecutar. No quites un bloqueo que pertenece a un proceso activo. Durante la instalación no se ejecutan scripts del proyecto ni ejecutables descargados.

## Mantén el catálogo remoto

Las instalaciones leen el `skills-index.json` publicado en `main` del repositorio. Los cambios futuros de skills deben publicar el catálogo regenerado junto con sus cambios de origen. Las pruebas locales usan fixtures HTTP aislados; compilar un binario por sí solo no publica skills ni un release.

```bash
npm run catalog:generate -- --help  # Ayuda del comando de catálogo sin efectos secundarios
npm run catalog:generate
npm run catalog:check
npm run test:catalog
npm test                 # Suite existente, solo con Node, incluida la vigencia del catálogo
npm run test:zesh-cli     # Comprobaciones unitarias y de aceptación de Go, por separado
```

`skills-index.json` usa `schemaVersion: 1` e identificadores estables `<domain>/<directory>`. Las entradas contienen nombre, descripción, alcances, temas, tecnologías, disparadores, ruta de origen, versión y hashes de archivos. El campo opcional `descriptionEs` es una descripción en español de texto plano y tamaño acotado (hasta 1024 bytes UTF-8, sin caracteres de control); los valores ausentes o vacíos siguen siendo entrada heredada válida. El generador lee todo el frontmatter YAML de `SKILL.md`, incluidas las descripciones plegadas. Agrega un skill, regenera el catálogo y envía ambos juntos; no hace falta recompilar el binario. La canalización rechaza un catálogo desactualizado. Un esquema de catálogo futuro e incompatible requiere una actualización compatible de la CLI.

Las sobreescrituras opcionales del frontmatter permiten una clasificación precisa sin editar el generador:

```yaml
metadata:
  version: "1.0"
  catalog:
    scopes: [frontend]
    topics: [design]
    technologies: [react]
    triggers: ["Landing-page layout and visual composition"]
```

Para un skill nuevo, el campo opcional `metadata.catalog.descriptionEs` aporta la descripción en español de la interfaz sin recompilar el binario. En caso contrario, el generador usa `cli/skill_descriptions_es.json` solo para una coincidencia exacta de identificador y origen; ese archivo es opcional en los fixtures heredados y nunca determina el descubrimiento. Mantén sin cambios la `description` canónica en inglés. Las traducciones escritas inválidas se rechazan antes de generar el catálogo.

Sin sobreescrituras se usan los metadatos de dominio existentes y los temas predeterminados conservadores. Los alcances son `frontend`, `backend` y `app`; los temas son descriptivos y extensibles. Los skills compartidos pueden aplicar a varios alcances. La detección sugiere alcances a partir de manifiestos de paquetes conocidos y archivos de proyecto de .NET sin ejecutar herramientas; no oculta skills según una versión de framework inferida. Revisa la compatibilidad de versiones al elegir .NET 8 frente a .NET 10 o skills específicos de un framework.

## Referencias de integración y límites de la evidencia

Las rutas de descubrimiento siguen las convenciones documentadas de [Claude Code](https://code.claude.com/docs/en/skills), [Codex](https://developers.openai.com/codex/skills), [OpenCode](https://opencode.ai/docs/skills/) y [Pi](https://github.com/badlogic/pi-mono/blob/main/packages/coding-agent/docs/skills.md). Las estructuras generadas se comprueban localmente; la carga y el comportamiento reales del modelo y la interfaz de confianza de proyecto de cada host requieren verificación del lado del consumidor y no quedan establecidos por estas pruebas.
