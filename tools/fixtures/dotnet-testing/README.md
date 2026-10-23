# Isolated xUnit capacity-rule trial

Run `npm run test:dotnet-unit` after separately approved local prerequisite setup. This is a real .NET10/xUnit2.9 specimen, not console acceptance, a consumer scaffold or evidence for every testing recipe in the collection.

## Contract

The counter starts at8 with fixed capacity10. `TryAdd(int amount)` accepts a positive amount only when the new total fits; rejection returns false and preserves its public `Total`. No files, network, clock, cancellation, database, publication collaborator or security implementation belongs to this rule.

[CapacityCounterTests.cs](CapacityCounterTests.cs) contains two theories/six cases:

- `1` produces9 and `2` reaches exactly10; both return true.
- `0`, `-1`, `3` (one above remaining capacity) and `int.MaxValue` return false with unchanged total8. Checking remaining capacity before addition avoids an overflowing sum bypassing the limit.

Every case asserts the returned outcome and observable total through the public API. Tests use the .NET10 native Assert/naming/AAA style; they do not migrate the separate .NET8 profile. Time/cancellation changes in the skills are instruction evidence only, not executed time/cancellation behavior here.

## Prerequisites and reproduction

Node>=18 and an already installed stable .NET10 SDK/reference packs are required. Direct packages are xunit2.9.3, xunit.runner.visualstudio3.1.5 and Microsoft.NET.Test.Sdk17.14.1; [packages.lock.json](packages.lock.json) fixes the complete resolved graph/content hashes. No root dependencies or CI requirement changes.

Installation is separate from testing and needs explicit authorization. For a fresh checkout, create project-local package/home/temp/cache directories and set these environment variables before the approved restore:

| Variable | Absolute path under the collection root |
|---|---|
| DOTNET_CLI_HOME | `.dotnet-testing-nuget/home` |
| NUGET_PACKAGES | `.dotnet-testing-nuget/packages` |
| NUGET_HTTP_CACHE_PATH | `.dotnet-testing-nuget/http` |
| NUGET_PLUGINS_CACHE_PATH | `.dotnet-testing-nuget/plugins` |
| TEMP, TMP, TMPDIR | `.dotnet-testing-nuget/tmp` |

Set `DOTNET_CLI_TELEMETRY_OPTOUT=1`, `DOTNET_NOLOGO=1`, `DOTNET_GENERATE_ASPNET_CERTIFICATE=false`, `DOTNET_CLI_WORKLOAD_UPDATE_NOTIFY_DISABLE=true` and `NODE_DISABLE_COMPILE_CACHE=1`. Then use the following commands from the collection root, substituting the absolute project-local restore-obj path:

```sh
# Only after explicit prerequisite-download approval:
dotnet restore tools/fixtures/dotnet-testing/CapacityRule.Tests.csproj --locked-mode --configfile tools/fixtures/dotnet-testing/Install.NuGet.Config --packages "$NUGET_PACKAGES" -p:MSBuildProjectExtensionsPath="<absolute-root>/.dotnet-testing-nuget/restore-obj/" -p:ImportDirectoryBuildProps=false -p:ImportDirectoryBuildTargets=false -p:ImportDirectoryPackagesProps=false -p:NuGetAudit=false
# Ordinary execution/help never downloads prerequisites:
npm run test:dotnet-unit
npm run test:dotnet-unit -- --help
```

The initial approved restore generated the lockfile with `--use-lock-file`; subsequent setup uses locked mode. `Install.NuGet.Config` has only official nuget.org. Ordinary runner execution uses the separate [empty-source config](NuGet.Config), checked local packages, locked restore and `test --no-restore`; absent prerequisites exit2 without downloads. SDKs/tooling are not installed globally. NuGet packages can contain MSBuild/test-runtime code: pinning and local paths are not a sandbox or a dependency audit.

## Runner/evidence limits

Help exits0 on stdout; unknown options and missing local prerequisites exit2 with specific stderr. Tests execute in a unique project-local copy with redirected home/temp/caches and ancestor build imports/compiler servers disabled. Only the owned copy is removed; the approved local package cache remains a reusable prerequisite. TRX counters must show six passing cases and zero failures, so an empty discovery cannot silently succeed. Cleanup failure makes the runner fail; forced interruption/timeouts may prevent complete cleanup or descendant termination.

Default `npm test`/`verify` remain Node-only and do not load NuGet packages or invoke this trial. The standalone project's inline compiler settings are deliberate because its disposable copy does not inherit repository build props. Native Assert, xUnit theories, public results and state are verified; no HTTP/WAF, real-engine transactions/security, clock/cancellation, other platforms, code coverage or arbitrary consumer behavior is established.
