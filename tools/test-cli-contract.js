'use strict';

// Opt-in black-box acceptance, not xUnit coverage or an OS sandbox.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const { createHash } = require('node:crypto');
const { execFile, execFileSync } = require('node:child_process');
const root = path.resolve(__dirname, '..');
const assets = path.join(root, 'skills/shared/cli-design/assets/line-tally');

function validateContract(contract) {
  assert.equal(contract.schemaVersion, 1);
  assert.ok(Array.isArray(contract.cases) && contract.cases.length > 0);
  const ids = new Set();
  for (const spec of contract.cases) {
    assert.match(spec.id, /^[a-z0-9-]+$/);
    assert.ok(!ids.has(spec.id), 'Duplicate case id');
    ids.add(spec.id);
    assert.ok(Array.isArray(spec.args) && spec.args.every((arg) => typeof arg === 'string'));
    assert.ok(spec.stdin === null || typeof spec.stdin === 'string');
    assert.ok(Number.isInteger(spec.code));
    assert.equal(typeof spec.stdout, 'string');
    assert.equal(typeof spec.stderr, 'string');
    assert.ok(spec.files && typeof spec.files === 'object' && !Array.isArray(spec.files));
    for (const [name, contents] of Object.entries(spec.files)) {
      assert.ok(name && name !== '.' && name !== '..' && path.basename(name) === name);
      assert.equal(typeof contents, 'string');
    }
  }
}

function snapshot(directory) {
  return fs.readdirSync(directory, { withFileTypes: true }).map((entry) => ({
    name: entry.name,
    kind: entry.isFile() ? 'file' : entry.isDirectory() ? 'directory' : 'other',
    hash: entry.isFile() ? createHash('sha256')
      .update(fs.readFileSync(path.join(directory, entry.name))).digest('hex') : null,
  })).sort((a, b) => a.name.localeCompare(b.name));
}

function invoke(dll, spec, cwd, env) {
  return new Promise((resolve) => {
    const child = execFile('dotnet', [dll, ...spec.args], {
      cwd, env, encoding: 'utf8', shell: false, windowsHide: true,
      timeout: 3000, maxBuffer: 16 * 1024,
    }, (error, stdout, stderr) => resolve({
      code: error ? error.code : 0, signal: error?.signal ?? null, stdout, stderr,
    }));
    // Early rejection can close the pipe; status and exact output are still asserted.
    child.stdin.on('error', () => {});
    if (spec.stdin !== null) child.stdin.end(spec.stdin);
    // null deliberately keeps stdin open: help/file/errors must not prompt or read it.
  });
}

async function main() {
  const { test } = require('node:test');
  const contract = JSON.parse(fs.readFileSync(path.join(assets, 'contract.json'), 'utf8'));
  validateContract(contract); // Validate paths before creating or writing any trial data.
  const scratch = fs.mkdtempSync(path.join(root, '.cli-design-trial-'));
  try {
    for (const directory of ['app', 'out', 'home', 'tmp', 'packages', 'http-cache', 'plugin-cache']) {
      fs.mkdirSync(path.join(scratch, directory));
    }
    for (const file of ['LineTally.csproj', 'Program.cs']) {
      fs.copyFileSync(path.join(assets, file), path.join(scratch, 'app', file));
    }
    fs.writeFileSync(path.join(scratch, 'global.json'), JSON.stringify({
      sdk: { version: '10.0.100', rollForward: 'latestFeature', allowPrerelease: false },
    }));
    const config = path.join(scratch, 'NuGet.Config');
    fs.writeFileSync(config, '<configuration><packageSources><clear /></packageSources><auditSources><clear /></auditSources></configuration>');
    const env = { ...process.env,
      DOTNET_CLI_HOME: path.join(scratch, 'home'), DOTNET_CLI_TELEMETRY_OPTOUT: '1',
      DOTNET_NOLOGO: '1', DOTNET_SKIP_FIRST_TIME_EXPERIENCE: '1',
      DOTNET_GENERATE_ASPNET_CERTIFICATE: 'false', DOTNET_CLI_WORKLOAD_UPDATE_NOTIFY_DISABLE: 'true',
      DOTNET_EnableDiagnostics: '0', MSBUILDDISABLENODEREUSE: '1',
      NUGET_PACKAGES: path.join(scratch, 'packages'),
      NUGET_HTTP_CACHE_PATH: path.join(scratch, 'http-cache'),
      NUGET_PLUGINS_CACHE_PATH: path.join(scratch, 'plugin-cache'),
      TEMP: path.join(scratch, 'tmp'), TMP: path.join(scratch, 'tmp'), TMPDIR: path.join(scratch, 'tmp'),
    };
    try {
      execFileSync('dotnet', ['build', 'app/LineTally.csproj', '--configuration', 'Release',
        '--output', path.join(scratch, 'out'), '--disable-build-servers', '--verbosity', 'quiet',
        `-p:RestoreConfigFile=${config}`, '-p:ImportDirectoryBuildProps=false',
        '-p:ImportDirectoryBuildTargets=false', '-p:ImportDirectoryPackagesProps=false',
        '-p:UseSharedCompilation=false'], {
        cwd: scratch, env, encoding: 'utf8', timeout: 120000, maxBuffer: 2 * 1024 * 1024,
      });
    } catch (error) {
      throw new Error(`.NET 10 build failed; SDK/reference packs must already be installed. No feeds are configured.\n${error.message}\n${String(error.stdout ?? '').slice(-6000)}`);
    }
    console.log('Built the isolated .NET 10 acceptance specimen.');
    const dll = path.join(scratch, 'out', 'LineTally.dll');
    for (const spec of contract.cases) {
      await test(`CLI public contract: ${spec.id}`, async () => {
        const cwd = path.join(scratch, spec.id);
        fs.mkdirSync(cwd);
        for (const [name, contents] of Object.entries({ 'untouched.txt': 'synthetic-state\n', ...spec.files })) {
          fs.writeFileSync(path.join(cwd, name), contents);
        }
        const before = snapshot(cwd);

        const result = await invoke(dll, spec, cwd, env);

        assert.deepEqual(snapshot(cwd), before, 'Fixture state changed');
        assert.equal(result.signal, null, 'Unexpected termination or timeout');
        assert.equal(result.code, spec.code, 'Exit code');
        assert.equal(result.stdout, spec.stdout, 'Exact stdout');
        assert.equal(result.stderr, spec.stderr, 'Exact stderr');
      });
    }
  } finally {
    // Only remove the unique directory created by this invocation, never shared caches.
    try {
      fs.rmSync(scratch, { recursive: true });
      console.log('TRIAL_CLEANUP=removed');
    } catch (error) {
      console.error(`TRIAL_CLEANUP=failed: ${error.message}`);
      process.exitCode = 1;
    }
  }
}

const args = process.argv.slice(2);
if (args.length === 0) {
  main().catch((error) => { console.error(error.message); process.exitCode = 1; });
} else if (args.length === 1 && args[0] === '--help') {
  process.stdout.write('Usage: npm run test:cli [-- --help]\n' +
    'Runs the opt-in .NET 10 CLI acceptance specimen in disposable project-local state.\n' +
    'Requires an installed .NET 10 SDK/reference packs. No SDK/package installation or external feeds.\n');
} else {
  process.stderr.write('cli-contract: unknown option. Use --help.\n');
  process.exitCode = 2;
}
