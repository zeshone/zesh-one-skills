package main

import (
	"strings"
)

// Language affects presentation only, never identifiers, user values or stored formats.
type language bool

func (l language) text(english string) string {
	if l {
		if value, ok := spanishText[english]; ok {
			return value
		}
	}
	return english
}
func (m *pickerModel) inLanguage(l language) *pickerModel { m.lang = l; return m }
func (l language) affirmative(s string) bool {
	return affirmative(s) || (bool(l) && (strings.EqualFold(s, "s") || strings.EqualFold(s, "si") || strings.EqualFold(s, "sí")))
}
func (l language) reason(s string) string {
	const prefix = "Matches project dependency "
	if strings.HasPrefix(s, prefix) {
		return l.text(prefix) + strings.TrimPrefix(s, prefix)
	}
	return l.text(s)
}
func (l language) topics(values []string) string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = l.text(value)
	}
	return strings.Join(out, ",")
}
func (l language) confirmationName(name string) string {
	const suffix = " (modified managed index)"
	if strings.HasSuffix(name, suffix) {
		return strings.TrimSuffix(name, suffix) + l.text(suffix)
	}
	return name
}
func (l language) errorText(err error) string {
	raw := err.Error()
	if !l {
		return raw
	}
	if value, ok := spanishText[raw]; ok {
		return value
	}
	for _, pair := range [][2]string{
		{"flag provided but not defined: ", "opción no definida: "},
		{"incomplete skill: ", "skill incompleta: "},
		{"unknown scope: ", "alcance desconocido: "}, {"unknown agent: ", "agente desconocido: "}, {"unknown skill: ", "skill desconocida: "}, {"invalid skill number: ", "número de skill no válido: "},
		{"cancelled: local changes preserved in ", "cancelado: se conservaron los cambios locales en "}, {"checksum mismatch: ", "verificación SHA-256 fallida: "},
		{"cannot prepare interactive terminal: ", "no se pudo preparar la terminal interactiva: "}, {"cannot restore terminal: ", "no se pudo restaurar la terminal: "},
		{"refusing symbolic link: ", "se rechazó el enlace simbólico: "}, {"not a directory: ", "no es un directorio: "},
	} {
		if strings.HasPrefix(raw, pair[0]) {
			return pair[1] + strings.TrimPrefix(raw, pair[0])
		}
	}
	const drift = "instruction changed during installation: "
	if strings.HasPrefix(raw, drift) {
		return "el archivo de instrucciones cambió durante la instalación: " + strings.TrimSuffix(strings.TrimPrefix(raw, drift), "; retry") + "; vuelve a intentarlo"
	}
	return "Detalle técnico: " + raw
}
func (l language) help() string {
	if l {
		return helpTextES
	}
	return helpText
}

var spanishText = map[string]string{
	"zesh-skills: %s\nUse --help for usage.\n":                                   "zesh-skills: %s\nUsa --help para ver el uso.\n",
	"unexpected positional arguments":                                            "no se esperaban argumentos posicionales",
	"non-interactive mode requires --skills, --agents, --architecture and --yes": "el modo no interactivo requiere --skills, --agents, --architecture y --yes",
	"invalid architecture text":                                                  "texto de arquitectura no válido",
	"--instructions requires a root-level .md filename":                          "--instructions requiere un archivo .md en la raíz",
	"instruction file must be a root-level .md filename":                         "el archivo de instrucciones debe ser un .md en la raíz",
	"invalid project root":                                                       "raíz del proyecto no válida", "choose at least one scope": "selecciona al menos un alcance",
	"scope is required": "se requiere un alcance", "agents and architecture are required": "se requieren agentes y arquitectura",
	"skills, agents and architecture are required":                 "se requieren skills, agentes y arquitectura",
	"input ended; no installation started":                         "terminó la entrada; no se inició la instalación",
	"project is locked or not writable; inspect .zesh-skills.lock": "el proyecto está bloqueado o no permite escritura; revisa .zesh-skills.lock",
	"download failed: check network connectivity":                  "falló la descarga; revisa la conexión de red",
	"download interrupted":                                         "se interrumpió la descarga", "invalid catalog JSON": "el JSON del catálogo no es válido",
	"invalid main commit response": "la respuesta del commit de main no es válida", "unsafe download path": "ruta de descarga no segura",
	"Cancelled; no files changed.": "Cancelado; no se modificaron archivos.",
	"Project root":                 "Raíz del proyecto",
	"Project: %s\nDetected technologies: %s\nExisting instructions: %s\n": "Proyecto: %s\nTecnologías detectadas: %s\nInstrucciones existentes: %s\n",
	"Choose project scope": "Selecciona el alcance del proyecto", "Detected project scope": "Alcance detectado en el proyecto",
	"Choose architecture": "Elige la arquitectura", "Custom architecture": "Arquitectura personalizada",
	"Agents (claude,codex,opencode,pi)":                "Agentes (claude,codex,opencode,pi)",
	"Additional root instruction file (.md; optional)": "Archivo adicional de instrucciones en la raíz (.md; opcional)",
	"\nInstall %d selected skill(s) into %s\nSource: main @ %s\nArchitecture: %s\nAgents: %s\nSkills: %s\nContent: .zesh-skills/content/\n": "\nSkills a instalar: %d\nDestino: %s\nFuente: main @ %s\nArquitectura: %s\nAgentes: %s\nSkills: %s\nContenido: .zesh-skills/content/\n",
	"Instruction files: ":                                             "Archivos de instrucciones: ",
	"\n  %s (managed index and triggers)\n":                           "\n  %s (índice administrado y triggers)\n",
	"Proceed? (y/N)":                                                  "¿Continuar? (s/N)",
	"Replace local changes in %s? (y/N)":                              "¿Reemplazar los cambios locales en %s? (s/N)",
	" (modified managed index)":                                       " (índice administrado modificado)",
	"Installed %d selected skill(s). Existing selections retained.\n": "Instalación completada. Skills seleccionadas: %d. Se conservan las anteriores.\n",
	"Select project skills":                                           "Selecciona las skills del proyecto",
	"%d selected | %d suggested | %d available":                       "%d seleccionadas | %d sugeridas | %d disponibles",
	"Suggestions are optional; every choice is editable.":             "Las sugerencias son opcionales; puedes cambiar cada elección.",
	"Choose one profile | %d available":                               "Elige un perfil | %d disponibles",
	"This describes the project; it does not scaffold it.":            "Describe el proyecto; no crea su estructura ni infraestructura.",
	"* suggested": "* sugerida", "* suggested: %s": "* sugerida: %s",
	"Only one architecture profile is selected.":           "Solo se selecciona un perfil de arquitectura.",
	"Up/Down or j/k: move | Space: toggle":                 "Arriba/Abajo o j/k: mover | Espacio: marcar",
	"Enter: confirm | a: all | n: none | r: reset":         "Enter: confirmar | a: todas | n: ninguna | r: restablecer",
	"q or Ctrl-C: cancel (no files changed)":               "q o Ctrl-C: cancelar (sin modificar archivos)",
	"Enter: choose highlighted profile | q/Ctrl-C: cancel": "Enter: elegir perfil resaltado | q/Ctrl-C: cancelar",
	"Framework-neutral debugging baseline":                 "Base de debugging independiente del framework",
	"UI keyboard, focus and feedback baseline":             "Base de teclado, focus y feedback para UI",
	"Matches project dependency ":                          "Coincide con la dependencia del proyecto ",
	"accessibility":                                        "Accesibilidad (Accessibility)", "architecture": "Arquitectura (Architecture)", "cli": "CLI", "collaboration": "Colaboración (Collaboration)",
	"data-access": "Acceso a datos (Data Access)", "debugging": "Depuración (Debugging)", "design": "Diseño (Design)", "framework": "Framework",
	"observability": "Observabilidad (Observability)", "performance": "Rendimiento (Performance)", "security": "Seguridad (Security)", "state": "Estado (State)",
	"testing": "Pruebas (Testing)", "validation": "Validación (Validation)",
}

const helpTextES = `Uso: zesh-skills [opciones]
Instala skills locales seleccionadas desde main de zeshone/zesh-one-skills.
Sin argumentos: asistente con selectores de alcance, arquitectura y skills.
Flechas/j/k: mover; Espacio: marcar; Enter: confirmar; a/n/r: todas/ninguna/restablecer.
q/Ctrl-C cancela. Las sugerencias son opcionales; la automatización usa elecciones explícitas.

  --root RUTA             Raíz del proyecto (predeterminada: carpeta actual)
  --scope ALCANCES        frontend,backend,app separados por comas; permite mixed
  --agents AGENTES        claude,codex,opencode,pi separados por comas
  --architecture TEXTO    Perfil o descripción personalizada de arquitectura
  --skills IDS            IDs del catálogo separados por comas, para automatización
  --instructions ARCH.md  Archivo adicional de instrucciones en la raíz
  --list                  Lista skills por tema sin instalar
  --yes                   Confirma la instalación, NO reemplaza cambios locales
  --non-interactive       No pregunta; exige skills, agentes, arquitectura y --yes
  --es                    Interfaz en español (sin este parámetro: inglés)
  --version               Muestra la versión sin red ni escrituras
  --help                  Muestra esta ayuda sin red ni escrituras

Ejemplos:
  zesh-skills --es
  zesh-skills --es --scope frontend --list
  zesh-skills --es --non-interactive --agents pi --architecture mockup --skills frontend/tailwind-4 --yes

Los nombres técnicos se conservan en inglés, solos o entre paréntesis:
Arquitectura Limpia (Clean Architecture).
En el catálogo solo se traducen las descripciones de skills disponibles;
los nombres, IDs y temas conservan el inglés original. Se prefieren descripciones
remotas en español; hay 35 traducciones integradas para descripciones originales
conocidas. Sin traducción o si cambia el original, se muestra el inglés.
Los valores personalizados y documentos descargados no se traducen.
Se conservan las selecciones instaladas; los cambios locales requieren confirmación.
Códigos de salida: 0 éxito/cancelación, 1 fallo de ejecución, 2 uso/entrada no válidos.
`
