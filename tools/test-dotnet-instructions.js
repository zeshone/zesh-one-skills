'use strict';

// Instruction/preservation checks only, not xUnit or application execution.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const root = path.resolve(__dirname, '..');
const read = file => fs.readFileSync(path.join(root, file), 'utf8');
const paths = ['skills/backend/testing-unit/SKILL.md', 'skills/dotnet-10/testing/SKILL.md'];
const guides = paths.map(read);
const includes = (text, parts) => parts.forEach(part => assert.ok(text.includes(part), `Missing: ${part}`));

test('.NET guidance selects exact and adjacent boundaries only for applicable rules', () => {
  for (const guide of guides) includes(guide, ['exact limit', 'adjacent', 'documented result',
    'parameterized', 'only when the contract']);
});

test('.NET rejection assertions include outcome and prohibited public effects', () => {
  for (const guide of guides) includes(guide, ['unchanged business state', 'no write', 'publication',
    'public', 'diagnostic logging', 'private implementation']);
});

test('.NET time and cancellation scenarios preserve actual commitment semantics', () => {
  for (const guide of guides) includes(guide, ['before, at and after', 'cutoff', 'no sleeps',
    'token forwarding', 'pre-cancelled', 'documented cancellation', 'already committed']);
  includes(guides[1], ['FakeTimeProvider', 'TimeProvider', 'without building a web host']);
});

test('.NET8 and .NET10 retain separate libraries, naming and existing route contracts', () => {
  includes(guides[0], ['xUnit + NSubstitute + FluentAssertions', '{Method}_When{Condition}_Should{Expected}',
    'single Act', 'GetByIdAsync_WhenOwnerDiffers_ShouldThrowForbiddenException']);
  includes(guides[1], ['native `Assert` only', 'MethodName_Condition_ExpectedResult',
    'manual hand-rolled spies', 'Never issue SQL from these tests', 'never mock a transaction seam',
    'hard-fail in CI when Docker is unavailable', 'No WAF or database', 'no reason message']);
  includes(read('AGENTS.backend.md'), ['testing-unit', 'NSubstitute']);
  includes(read('AGENTS.dotnet-10.md'), ['`testing`', 'WAF-tier']);
});

test('.NET instruction checks preserve default Node-only pipelines and root dependencies', () => {
  const { scripts, devDependencies } = JSON.parse(read('package.json'));
  assert.equal(scripts['test:dotnet-instructions'], 'node tools/test-dotnet-instructions.js');
  for (const script of [scripts.test, scripts.verify]) {
    includes(script, ['test-fixtures.js', 'spec-fixtures-logging.js', 'test-pilot-instructions.js',
      'test-cli-instructions.js', 'test-frontend-instructions.js', 'test-dotnet-instructions.js',
      'lint-skills.js', 'test-release-readiness.js', 'test-lint-zero-skills.js']);
    assert.doesNotMatch(script, /dotnet test|test-dotnet-unit|npm (ci|install)/);
  }
  assert.deepEqual(Object.keys(devDependencies), ['js-yaml']);
});

test('.NET touched guidance and adoption links resolve locally', () => {
  for (const file of [...paths, 'README.md', 'ADOPTING.md', 'CREDITS.md',
    'tools/fixtures/dotnet-testing/README.md']) {
    for (const [, target] of read(file).matchAll(/\]\(([^)]+)\)/g)) {
      if (/^(https?:|#)/.test(target)) continue;
      assert.ok(fs.existsSync(path.resolve(root, path.dirname(file), target.split('#')[0])), `${file}: ${target}`);
    }
  }
});

for (const spec of [
  { name: 'help', args: ['--help'], code: 0,
    stdout: 'Usage: npm run test:dotnet-unit [-- --help]\nRuns the optional .NET 10 xUnit capacity-rule trial using approved local packages only.\nSee tools/fixtures/dotnet-testing/README.md for isolated prerequisite setup and limits.\n', stderr: '' },
  { name: 'unknown option', args: ['--unknown'], code: 2, stdout: '',
    stderr: 'dotnet-unit: unknown option. Use --help.\n' },
  { name: 'missing prerequisites', args: [], code: 2, stdout: '',
    stderr: 'dotnet-unit: isolated prerequisites missing. See tools/fixtures/dotnet-testing/README.md; no download was attempted.\n' },
]) {
  test(`xUnit runner ${spec.name} avoids SDK/downloads and preserves stored state`, () => {
    const scratch = fs.mkdtempSync(path.join(root, '.dotnet-unit-trial-runner-'));
    try {
      fs.mkdirSync(path.join(scratch, 'tools'));
      fs.copyFileSync(path.join(root, 'tools/test-dotnet-unit.js'), path.join(scratch, 'tools/test-dotnet-unit.js'));
      fs.writeFileSync(path.join(scratch, 'untouched.txt'), 'synthetic-only');
      const env = Object.fromEntries(Object.entries(process.env)
        .filter(([key]) => key.toUpperCase() !== 'PATH'));
      Object.assign(env, { PATH: '', NODE_DISABLE_COMPILE_CACHE: '1', TEMP: scratch, TMP: scratch, TMPDIR: scratch });

      const result = spawnSync(process.execPath, ['tools/test-dotnet-unit.js', ...spec.args],
        { cwd: scratch, env, encoding: 'utf8', timeout: 10000, shell: false });

      assert.equal(result.error, undefined);
      assert.equal(result.signal, null);
      assert.equal(result.status, spec.code);
      assert.equal(result.stdout, spec.stdout);
      assert.equal(result.stderr, spec.stderr);
      assert.deepEqual(fs.readdirSync(scratch).sort(), ['tools', 'untouched.txt']);
      assert.deepEqual(fs.readdirSync(path.join(scratch, 'tools')), ['test-dotnet-unit.js']);
      assert.equal(fs.readFileSync(path.join(scratch, 'untouched.txt'), 'utf8'), 'synthetic-only');
    } finally { fs.rmSync(scratch, { recursive: true }); }
  });
}
