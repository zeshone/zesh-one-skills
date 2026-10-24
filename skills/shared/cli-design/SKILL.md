---
name: cli-design
description: "Trigger: CLI design, commands, flags, help, stdin/stdout/stderr. Define compatible human and automation interfaces without choosing architecture."
license: MIT
allowed-tools: Read Write Edit Bash Grep Glob
metadata:
  author: Zesh-One
  version: "1.0"
---

# CLI interaction contracts

## When to Use

Load when designing or changing a user-facing command-line interface, not merely running existing commands. Own the interaction contract; leave project structure, language conventions and security controls with their existing owners.

## Critical Patterns

1. **Contract before code.** Preserve existing command names, flags, defaults, output schemas and exit codes that callers consume. Record inputs, modes, outputs and failure behavior; obtain explicit approval for breaking changes instead of silently redesigning automation.
2. **Help and status.** Provide side-effect-free `--help` on stdout with exit 0: usage, operands, flags, defaults and a minimal example. Send actionable usage/runtime errors to stderr with documented nonzero exit codes. Never return success for a failed operation or expose raw arguments, credentials or input in diagnostics.
3. **Streams compose.** Reserve stdout for requested results and stderr for diagnostics/progress. Document file versus stdin input, including `-` when supported; use `--` to disambiguate option-like operands. Pass subprocess arguments as an argv array, not interpolated shell text.
4. **Machine output is a contract.** When machine-readable output is required, keep its schema deterministic and free of banners, ANSI decoration and progress messages. Preserve existing formats; do not add JSON or output flags to a CLI that does not need them.
5. **Automation must not prompt.** In non-interactive mode, never wait for a prompt: use explicit arguments or fail with a useful error. Reading an explicitly requested stdin stream is not a prompt. Gate optional interactive UI on the actual terminal/mode; a `--yes` flag alone is not authorization for the agent to perform destructive work.
6. **Validate before effects.** Reject unknown flags, missing values and ambiguous operands before side effects or consuming unrelated input. Preserve existing files/state on rejected arguments; document partial-failure behavior for streaming or mutating operations rather than promising atomicity without evidence.

## Constraints & Tradeoffs

Reuse an appropriate installed parser; require no new framework, generic host, HTTP, database or authentication infrastructure just to satisfy this skill. Keep a small interface small. Apply the selected language's validation/testing conventions; this skill does not replace them.

## Anti-Patterns

Silent compatibility breaks; success exits after errors; log-contaminated JSON; prompts in automation; treating CLI flags or tool metadata as execution permission.

## Progressive Disclosure

Before implementation, state the bounded command contract and use synthetic inputs. Verify through the public entry point: assert exit code, exact stdout, exact stderr and unchanged stored input/state on rejection; exercise accepted file/pipe modes and relevant boundaries. Report observed results and untested modes, not model-compliance or application-security guarantees.

## Resources

- For new structure or boundary decisions, load [architecture selection](../architecture-selection.md); reuse an applicable recorded choice or wait for the user's selection.
- Always retain the [security baseline](../security/baseline.md); load its specialized checks only for affected capabilities.
- For .NET 10 testing details, use the existing [testing skill](../../dotnet-10/testing/SKILL.md); other stacks keep their own testing guidance.
- For an opt-in .NET 10 interaction trial, see the [executable acceptance specimen](assets/line-tally/README.md); it is not a project template or a security/model-compliance guarantee.
