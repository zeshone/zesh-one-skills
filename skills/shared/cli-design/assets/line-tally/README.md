# Executable CLI acceptance specimen

`line-tally` is a small read-only .NET 10 specimen for the CLI interaction contract, not a scaffolding template or a replacement for a project's architecture/testing decisions.

## Run the trial

From the collection root, with Node.js >=18 and a .NET 10 SDK/reference packs already installed:

```sh
npm run test:cli
npm run test:cli -- --help
```

Help exits0 without a build; unknown runner options exit2 on stderr without creating trial state. This is opt-in: `npm test` and `npm run verify` remain Node-only (including runner help/error checks with no SDK on PATH). The runner copies the project into a unique project-local scratch directory, selects an installed stable .NET 10 SDK, configures empty NuGet sources, redirects CLI caches/temp and disables inherited directory build imports. It does not install SDKs/packages or configure global settings. It removes only the directory it created, and reports cleanup failure separately.

The runner invokes `dotnet LineTally.dll` directly with argument arrays, not a shell command string. `line-tally` below is the logical command name; nothing is installed on PATH.

## Public contract

- `line-tally [--json] [--] <file|->`: exactly one input; `-` reads stdin to EOF. Count logical lines, including blank/final unterminated lines; empty input has zero lines.
- `--help` alone: usage/options/default output/example/statuses on stdout, exit0, no input consumption.
- Default output: invariant decimal count plus LF. `--json`: exactly `{"lines":N}` plus LF, without decoration or progress. stderr is empty on success.
- Unknown flags, duplicate `--json`, missing/extra operands: exit2, empty stdout, exactly `line-tally: invalid arguments. See --help.` plus LF on stderr.
- Input read failure: exit1, empty stdout, exactly `line-tally: unable to read input.` plus LF on stderr; no filename, input or exception details.
- No prompts, output files or stored counters. All input/state remains unchanged. `--` permits an option-like filename.

[`contract.json`](contract.json) contains13 independent cases. Every case asserts status, exact stdout/stderr, termination without a signal/timeout, and the unchanged fixture directory. Cases not using stdin leave its pipe open to detect accidental prompting/consumption. File, pipe, JSON, CRLF/LF, empty/blank/final lines, option-like/spaced filenames and rejected arguments are covered.

## Evidence boundary

This is a Node black-box acceptance harness, not xUnit/unit coverage. It demonstrates the checked specimen's public behavior, not arbitrary generated CLIs, automatic skill discovery, model compliance, destructive-command safeguards or application security. Test-state isolation is not an OS sandbox. Interrupt/signal mapping, broken pipes, huge inputs and other operating systems are not established by one local run; real products retain their own applicable tests and controls.
