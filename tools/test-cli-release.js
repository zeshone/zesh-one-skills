'use strict';
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const crypto = require('node:crypto');
const { spawnSync } = require('node:child_process');
const script = path.join(__dirname, 'build-zesh-skills.js');
const scratch = fs.mkdtempSync(path.join(os.tmpdir(), 'zesh-release-'));
// Intentionally retained isolated fixtures; never delete unknown files.
function cli(args) {
  return spawnSync(process.execPath, [script, ...args], {
    cwd: scratch, encoding: 'utf8', timeout: 5000, env: { ...process.env, PATH: '' },
  });
}
const help = cli(['--help']);
assert.equal(help.status, 0);
assert.equal(help.stderr, '');
assert.match(help.stdout, /--release <version>/);
assert.match(help.stdout, /--all/);
for (const args of [['--release'], ['--release', 'v2.8.0'], ['--release', '../2.8.0'],
  ['--release', '02.8.0'], ['--release', '2.8'], ['--release', '2.8.0;echo'],
  ['--release', '2.8.0', '--all'], ['--help', '--release', '2.8.0'], ['--wat']]) {
  const result = cli(args);
  assert.equal(result.status, 2, JSON.stringify(args));
  assert.equal(result.stdout, '');
  assert.match(result.stderr, /Usage:/);
}
assert.deepEqual(fs.readdirSync(scratch), []);
// Only require after the public help contract: old builder must not build on import.
const { parseArgs, build, packageRelease, targets } = require(script);
assert.deepEqual(targets, [['windows', 'amd64'], ['windows', 'arm64'], ['linux', 'amd64'], ['linux', 'arm64'], ['darwin', 'amd64'], ['darwin', 'arm64']]);
assert.deepEqual(parseArgs([]), { mode: 'host' });
assert.deepEqual(parseArgs(['--all']), { mode: 'all' });
assert.deepEqual(parseArgs(['--release', '2.8.0']), { mode: 'release', version: '2.8.0' });
let calls = [];
const root = path.join(scratch, 'fixture');
const run = (command, args, options) => {
  calls.push({ command, args, options });
  fs.writeFileSync(args[args.indexOf('-o') + 1], Buffer.from(`${options.env.GOOS}/${options.env.GOARCH}`));
  return { status: 0 };
};
assert.throws(() => build(['--release', 'bad'], { root, run }), /Usage:/);
assert.equal(calls.length, 0);
assert.equal(fs.existsSync(root), false);
build(['--release', '2.8.0'], { root, run });
assert.equal(calls.length, 6);
const out = path.join(root, 'dist/zesh-skills-release/v2.8.0');
const names = targets.map(([os, arch]) => `zesh-skills-${os}-${arch}${os === 'windows' ? '.exe' : ''}`);
assert.deepEqual(fs.readdirSync(out).sort(), [...names, 'SHA256SUMS'].sort());
for (const call of calls) {
  assert.equal(call.command, 'go');
  assert.equal(call.options.shell, false);
  assert.equal(call.options.env.CGO_ENABLED, '0');
  assert.equal(call.options.env.GOTOOLCHAIN, 'local');
  assert.ok(call.args.includes('-ldflags=-s -w -X main.version=2.8.0'));
}
const expected = names.map(name => `${crypto.createHash('sha256').update(fs.readFileSync(path.join(out, name))).digest('hex')}  ${name}\n`).join('');
assert.equal(fs.readFileSync(path.join(out, 'SHA256SUMS'), 'utf8'), expected);
fs.writeFileSync(path.join(out, 'unknown-old.zip'), 'keep');
packageRelease(out);
assert.equal(fs.readFileSync(path.join(out, 'unknown-old.zip'), 'utf8'), 'keep');
assert.equal(fs.readFileSync(path.join(out, 'SHA256SUMS'), 'utf8'), expected);
const missing = path.join(scratch, 'missing');
fs.mkdirSync(missing);
fs.writeFileSync(path.join(missing, 'SHA256SUMS'), 'prior');
assert.throws(() => packageRelease(missing), /ENOENT/);
assert.equal(fs.readFileSync(path.join(missing, 'SHA256SUMS'), 'utf8'), 'prior');
calls = [];
build(['--all'], { root, run });
assert.equal(calls.length, 6);
assert.ok(calls.every(call => call.args.includes('-ldflags=-s -w')));
assert.match(fs.readFileSync(path.join(root, 'dist/zesh-skills/SHA256SUMS'), 'utf8'), /  windows-amd64\/zesh-skills.exe\n/);
calls = [];
build([], { root, run });
assert.equal(calls.length, 1);
const failedRoot = path.join(scratch, 'failed-build');
let failures = 0;
assert.throws(() => build(['--release', '2.8.0'], {
  root: failedRoot,
  run: () => { failures++; return { status: 7 }; },
}), error => error.exitCode === 7 && /Build failed for windows\/amd64/.test(error.message));
assert.equal(failures, 1);
assert.equal(fs.existsSync(path.join(failedRoot, 'dist/zesh-skills-release/v2.8.0/SHA256SUMS')), false);
calls = [];
build(['--release', '3.0.1'], { root, run });
assert.equal(calls.length, 6);
assert.ok(calls.every(call => call.args.includes('-ldflags=-s -w -X main.version=3.0.1')));
assert.ok(fs.existsSync(path.join(root, 'dist/zesh-skills-release/v3.0.1/SHA256SUMS')));
assert.equal(fs.readFileSync(path.join(out, 'SHA256SUMS'), 'utf8'), expected);
console.log('CLI release parser, six assets, checksums, stale/missing/failure boundaries and development compatibility passed.');
