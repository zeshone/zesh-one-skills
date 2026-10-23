'use strict';

// Instruction contracts only: these do not render React or establish model compliance.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const root = path.resolve(__dirname, '..');
const load = file => fs.existsSync(path.join(root, file)) ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const skillPath = 'skills/frontend/testing/SKILL.md';
const skill = load(skillPath);
const guide = load('AGENTS.frontend.md');
const includes = (text, values) => values.forEach(value => assert.ok(text.includes(value), `Missing: ${value}`));

test('frontend testing selects the smallest useful layer without imposing tools', () => {
  includes(skill, ['installed versions', 'Unit', 'Component', 'Integration', 'E2E',
    'existing test runner', 'explicit approval', 'not a fixed pyramid']);
});

test('frontend checks visible asynchronous states and preserves the API boundary', () => {
  includes(skill, ['loading', 'error', 'empty', 'success', 'retry', 'role', 'accessible name',
    'existing authenticated client', 'backend authorization', 'deferred', 'cleanup']);
});

test('form tests preserve the installed flow and assert rejected-input effects', () => {
  includes(skill, ['invalid', 'valid', 'pending', 'disabled', 'no API call', 'unchanged stored',
    'recoverable', 'not migrate forms', 'internal hook', 'sleep']);
});

test('React SPA routing retains Next, Ionic and security owners', () => {
  includes(guide, ['skills/frontend/testing/SKILL.md', 'React SPA only']);
  includes(skill, ['Next', 'Ionic', 'not browser', 'not model', '../react-19/SKILL.md',
    '../tanstack-query/SKILL.md', '../../shared/security/baseline.md']);
  includes(guide, ['Client -> Server Action -> .NET API', 'existing authenticated client',
    '`params`/`searchParams`', 'backend authorization', 'skills/shared/security/baseline.md']);
});

test('collection test pipelines keep current checks and do not install frontend prerequisites', () => {
  const scripts = JSON.parse(load('package.json')).scripts;
  assert.equal(scripts['test:frontend-instructions'], 'node tools/test-frontend-instructions.js');
  for (const command of [scripts.test, scripts.verify]) {
    includes(command, ['test-fixtures.js', 'spec-fixtures-logging.js', 'test-pilot-instructions.js',
      'test-cli-instructions.js', 'test-frontend-instructions.js', 'lint-skills.js',
      'test-release-readiness.js', 'test-lint-zero-skills.js']);
    assert.doesNotMatch(command, /npm (install|ci)|test:frontend-dom|test-frontend-dom/);
  }
  assert.deepEqual(Object.keys(JSON.parse(load('package.json')).devDependencies), ['js-yaml']);
});

test('frontend skill and touched documents resolve their local links', () => {
  for (const file of [skillPath, 'AGENTS.frontend.md', 'README.md', 'ADOPTING.md', 'CREDITS.md',
    'tools/fixtures/frontend-dom/README.md']) {
    const text = load(file);
    assert.ok(text.length, file);
    for (const [, target] of text.matchAll(/\]\(([^)]+)\)/g)) {
      if (/^(https?:|#)/.test(target)) continue;
      assert.ok(fs.existsSync(path.resolve(root, path.dirname(file), target.split('#')[0])), `${file}: ${target}`);
    }
  }
});

for (const spec of [
  { name: 'help', args: ['--help'], code: 0,
    stdout: 'Usage: npm run test:frontend-dom [-- --help]\nRuns the optional React SPA component trial in jsdom, not a browser.\nInstall its isolated prerequisites only with approval; see tools/fixtures/frontend-dom/README.md.\n', stderr: '' },
  { name: 'unknown option', args: ['--unknown'], code: 2, stdout: '',
    stderr: 'frontend-dom: unknown option. Use --help.\n' },
  { name: 'missing prerequisites', args: [], code: 2, stdout: '',
    stderr: 'frontend-dom: isolated prerequisites missing. See tools/fixtures/frontend-dom/README.md; no installation was attempted.\n' },
]) {
  test(`DOM runner ${spec.name} does not install or mutate stored state`, () => {
    const scratch = fs.mkdtempSync(path.join(root, '.frontend-dom-trial-runner-'));
    try {
      fs.mkdirSync(path.join(scratch, 'tools'));
      fs.copyFileSync(path.join(root, 'tools/test-frontend-dom.js'), path.join(scratch, 'tools/test-frontend-dom.js'));
      fs.writeFileSync(path.join(scratch, 'untouched.txt'), 'synthetic-only');
      const env = Object.fromEntries(Object.entries(process.env)
        .filter(([key]) => key.toUpperCase() !== 'PATH'));
      Object.assign(env, { PATH: '', NODE_DISABLE_COMPILE_CACHE: '1', TEMP: scratch, TMP: scratch, TMPDIR: scratch });

      const result = spawnSync(process.execPath, ['tools/test-frontend-dom.js', ...spec.args],
        { cwd: scratch, env, encoding: 'utf8', timeout: 10000, shell: false });

      assert.equal(result.error, undefined);
      assert.equal(result.signal, null);
      assert.equal(result.status, spec.code);
      assert.equal(result.stdout, spec.stdout);
      assert.equal(result.stderr, spec.stderr);
      assert.deepEqual(fs.readdirSync(scratch).sort(), ['tools', 'untouched.txt']);
      assert.deepEqual(fs.readdirSync(path.join(scratch, 'tools')), ['test-frontend-dom.js']);
      assert.equal(fs.readFileSync(path.join(scratch, 'untouched.txt'), 'utf8'), 'synthetic-only');
    } finally { fs.rmSync(scratch, { recursive: true }); }
  });
}
