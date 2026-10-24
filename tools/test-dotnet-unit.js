'use strict';

// Optional real xUnit trial. Package installation is a separate authorized action.
const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const root = path.resolve(__dirname, '..');
const fixture = path.join(root, 'tools/fixtures/dotnet-testing');
const packages = path.join(root, '.dotnet-testing-nuget/packages');
const args = process.argv.slice(2);
if (args.length === 1 && args[0] === '--help') {
  process.stdout.write('Usage: npm run test:dotnet-unit [-- --help]\nRuns the optional .NET 10 xUnit capacity-rule trial using approved local packages only.\nSee tools/fixtures/dotnet-testing/README.md for isolated prerequisite setup and limits.\n');
} else if (args.length) {
  process.stderr.write('dotnet-unit: unknown option. Use --help.\n');
  process.exitCode = 2;
} else run();

function run() {
  const lockPath = path.join(fixture, 'packages.lock.json');
  if (!fs.existsSync(lockPath)) return missing();
  const lock = JSON.parse(fs.readFileSync(lockPath, 'utf8'));
  const entries = Object.entries(lock.dependencies['net10.0'] ?? {});
  if (!entries.length || entries.some(([id, pkg]) =>
    !/^[a-z0-9.-]+$/i.test(id) || !/^[a-z0-9.+-]+$/i.test(pkg.resolved) ||
    !fs.existsSync(path.join(packages, id.toLowerCase(), pkg.resolved.toLowerCase(), '.nupkg.metadata')))) {
    return missing();
  }
  const scratch = fs.mkdtempSync(path.join(root, '.dotnet-unit-trial-'));
  try {
    for (const file of ['CapacityRule.Tests.csproj', 'CapacityCounter.cs', 'CapacityCounterTests.cs',
      'packages.lock.json', 'NuGet.Config']) fs.copyFileSync(path.join(fixture, file), path.join(scratch, file));
    for (const name of ['home', 'tmp', 'http', 'plugins', 'results']) fs.mkdirSync(path.join(scratch, name));
    fs.writeFileSync(path.join(scratch, 'global.json'), JSON.stringify({
      sdk: { version: '10.0.100', rollForward: 'latestFeature', allowPrerelease: false },
    }));
    const env = { ...process.env, NODE_DISABLE_COMPILE_CACHE: '1',
      DOTNET_CLI_HOME: path.join(scratch, 'home'), DOTNET_CLI_TELEMETRY_OPTOUT: '1',
      DOTNET_NOLOGO: '1', DOTNET_SKIP_FIRST_TIME_EXPERIENCE: '1',
      DOTNET_GENERATE_ASPNET_CERTIFICATE: 'false', DOTNET_CLI_WORKLOAD_UPDATE_NOTIFY_DISABLE: 'true',
      DOTNET_EnableDiagnostics: '0', MSBUILDDISABLENODEREUSE: '1',
      NUGET_PACKAGES: packages, NUGET_HTTP_CACHE_PATH: path.join(scratch, 'http'),
      NUGET_PLUGINS_CACHE_PATH: path.join(scratch, 'plugins'),
      TEMP: path.join(scratch, 'tmp'), TMP: path.join(scratch, 'tmp'), TMPDIR: path.join(scratch, 'tmp'),
    };
    const common = { cwd: scratch, env, shell: false, encoding: 'utf8', timeout: 120000, maxBuffer: 1024 * 1024 };
    const invoke = command => {
      const result = spawnSync('dotnet', command, common);
      process.stdout.write(result.stdout ?? '');
      process.stderr.write(result.stderr ?? '');
      if (result.error || result.signal) process.stderr.write('dotnet-unit: process failed or timed out.\n');
      return result.error || result.signal ? 1 : result.status ?? 1;
    };
    const props = ['-p:ImportDirectoryBuildProps=false', '-p:ImportDirectoryBuildTargets=false',
      '-p:ImportDirectoryPackagesProps=false', '-p:RestoreFallbackFolders=', '-p:NuGetAudit=false',
      '-p:UseSharedCompilation=false'];
    const restored = invoke(['restore', 'CapacityRule.Tests.csproj', '--locked-mode', '--configfile',
      path.join(scratch, 'NuGet.Config'), '--packages', packages, '--verbosity', 'quiet', ...props]);
    if (restored !== 0) { process.exitCode = restored; return; }
    process.exitCode = invoke(['test', 'CapacityRule.Tests.csproj', '--no-restore', '--configuration', 'Release',
      '--disable-build-servers', '--logger', 'trx;LogFileName=capacity.trx',
      '--results-directory', path.join(scratch, 'results'), '--verbosity', 'normal', ...props]);
    // A zero exit with no discovered/passed cases is not a successful trial.
    if (process.exitCode === 0) {
      const report = fs.readFileSync(path.join(scratch, 'results/capacity.trx'), 'utf8');
      if (!/<Counters\b[^>]*\btotal="6"/.test(report) || !/<Counters\b[^>]*\bpassed="6"/.test(report) ||
        !/<Counters\b[^>]*\bfailed="0"/.test(report)) {
        console.error('dotnet-unit: expected six passing xUnit cases, none failed.');
        process.exitCode = 1;
      }
    }
  } finally {
    try { fs.rmSync(scratch, { recursive: true }); console.log('XUNIT_TRIAL_CLEANUP=removed'); }
    catch (error) { console.error(`XUNIT_TRIAL_CLEANUP=failed: ${error.message}`); process.exitCode = 1; }
  }
}

function missing() {
  process.stderr.write('dotnet-unit: isolated prerequisites missing. See tools/fixtures/dotnet-testing/README.md; no download was attempted.\n');
  process.exitCode = 2;
}
