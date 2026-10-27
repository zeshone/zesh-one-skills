package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// hostCLIInstallService detects the real host. Tests replace it so the public
// CLI flow can run without reading or writing this machine's user state.
var hostCLIInstallService = newHostCLIInstallService

// runInstall copies this executable into a user-owned location and persists only
// user-level PATH state. It never downloads, never asks for project input and
// never writes into a project root.
func runInstall(input io.Reader, out, errOut io.Writer, interactive, confirmed bool, lang language) int {
	usage := func(err error) int {
		fmt.Fprintf(errOut, lang.text("zesh-skills: %s\nUse --help for usage.\n"), lang.errorText(err))
		return 2
	}
	failure := func(err error) int { fmt.Fprintf(errOut, "zesh-skills: %s\n", lang.errorText(err)); return 1 }
	if !interactive && !confirmed {
		return usage(fmt.Errorf("--install-cli requires --yes to confirm"))
	}
	service, err := hostCLIInstallService()
	if err != nil {
		return failure(err)
	}
	plan, err := service.Plan()
	if err != nil {
		return failure(err)
	}
	fmt.Fprintf(out, lang.text("\nInstall this executable for the current user\nSource: %s\nDestination: %s\n"), display(service.Executable), display(plan.Destination))
	if len(plan.ConfigFiles) > 0 {
		fmt.Fprintf(out, lang.text("Shell startup files: %s\n"), display(strings.Join(plan.ConfigFiles, ", ")))
	}
	if plan.PATHChanged {
		fmt.Fprintln(out, lang.text("The user PATH will be updated; new terminals only."))
	} else {
		fmt.Fprintln(out, lang.text("The user PATH already contains the destination."))
	}
	if plan.Replaces {
		fmt.Fprintf(out, lang.text("An existing different executable will be replaced: %s\n"), display(plan.Destination))
	}
	fmt.Fprintln(out, lang.text("Skills stay project-local; nothing is installed globally."))
	reader := bufio.NewReader(input)
	ask := func(label string) (string, error) {
		fmt.Fprintf(out, "%s [%s]: ", lang.text(label), "n")
		line, readErr := reader.ReadString('\n')
		if readErr != nil {
			return "", fmt.Errorf("input ended; no installation started")
		}
		line = strings.TrimSpace(line)
		if line == "" {
			line = "n"
		}
		return line, nil
	}
	if !confirmed {
		answer, askErr := ask("Proceed? (y/N)")
		if askErr != nil {
			return usage(askErr)
		}
		if !lang.affirmative(answer) {
			fmt.Fprintln(out, lang.text("Cancelled; no files changed."))
			return 0
		}
	}
	consent := CLIInstallConsent{Install: true}
	if plan.Replaces {
		if !interactive {
			return failure(errors.New("existing executable differs; remove it or run interactively to confirm replacement"))
		}
		answer, askErr := ask(fmt.Sprintf(lang.text("Replace the existing executable at %s? (y/N)"), display(plan.Destination)))
		if askErr != nil {
			return usage(askErr)
		}
		if !lang.affirmative(answer) {
			fmt.Fprintln(out, lang.text("Cancelled; no files changed."))
			return 0
		}
		consent.Replace = true
	}
	if err := service.Apply(plan, consent); err != nil {
		return failure(err)
	}
	fmt.Fprintln(out, lang.text("Installed. Open a new terminal to use zesh-skills."))
	return 0
}
