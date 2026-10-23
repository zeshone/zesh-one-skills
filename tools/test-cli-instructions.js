'use strict';

// Instruction contracts, not evidence of model compliance or CLI execution.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const root = path.resolve(__dirname, '..');
const skillPath = 'skills/shared/cli-design/SKILL.md';
const read = (file) => fs.existsSync(path.join(root, file))
  ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const guides = ['AGENTS.backend.md', 'AGENTS.dotnet-10.md'];

test('CLI activation owns interaction contracts, not architecture or permissions', () => {
  const skill = read(skillPath);
  assert.match(skill, /^name: cli-design$/m);
  assert.match(skill, /description: "Trigger: CLI design/);
  assert.match(skill, /not merely running existing commands/i);
  assert.ok(skill.includes('../architecture-selection.md'));
  assert.ok(skill.includes('../security/baseline.md'));
});

test('CLI guidance covers help, arguments, streams, statuses and pipe input', () => {
  const skill = read(skillPath);
  for (const term of ['--help', 'stdout', 'stderr', 'exit', 'unknown flags',
    'missing values', 'stdin', '`--`', 'machine-readable']) {
    assert.ok(skill.includes(term), `Missing CLI contract: ${term}`);
  }
  assert.match(skill, /before.*side effects/i);
  assert.match(skill, /partial.*failure/i);
});

test('CLI automation preserves compatibility and does not infer destructive approval', () => {
  const skill = read(skillPath);
  assert.match(skill, /Preserve existing command names/);
  assert.match(skill, /never wait for a prompt/i);
  assert.match(skill, /--yes.*not authorization/i);
  assert.match(skill, /no new framework/i);
  assert.match(skill, /exit code.*stdout.*stderr.*unchanged/i);
});

test('both .NET guides route CLI interface work explicitly', () => {
  for (const guide of guides) {
    assert.match(read(guide), /\| CLI commands, flags, help, exit codes, stdin\/stdout\/stderr, non-interactive use \|/);
    assert.ok(read(guide).includes('skills/shared/cli-design/SKILL.md'), guide);
  }
});

test('existing .NET architecture and testing routes remain intact', () => {
  assert.match(read(guides[0]), /\| Architecture, project structure, feature scaffolding \| `net8-apirest-general` \|/);
  assert.match(read(guides[0]), /\| Unit tests \(xUnit, NSubstitute, FluentAssertions\) \| `net8-apirest-testing-unit` \|/);
  assert.match(read(guides[1]), /\| Project structure, contextual architecture choice, CLI\/worker applicability, layering, boundary\/seam placement \| `architecture` \|/);
  assert.match(read(guides[1]), /\| Tests, RealPg\/Testcontainers, `FakeTimeProvider`, xUnit, test naming, coverage strategy \| `testing` \|/);
});

test('instruction checks join both default pipelines without requiring .NET', () => {
  const { scripts } = JSON.parse(read('package.json'));
  const added = 'node tools/test-cli-instructions.js';
  const previous = ['node tools/test-fixtures.js', 'node tools/spec-fixtures-logging.js',
    'node tools/test-pilot-instructions.js', 'node tools/test-frontend-instructions.js',
    'node tools/test-dotnet-instructions.js', 'node tools/test-debugging-instructions.js',
    'node tools/test-accessibility-instructions.js', 'node tools/test-apps-instructions.js',
    'node tools/lint-skills.js',
    'node tools/test-release-readiness.js', 'node tools/test-lint-zero-skills.js'];
  assert.equal(scripts['test:cli-instructions'], added);
  for (const name of ['test', 'verify']) {
    const commands = scripts[name].split(' && ');
    assert.equal(commands.filter((command) => command === added).length, 1);
    assert.deepEqual(commands.filter((command) => command !== added), previous);
  }
});

test('CLI skill and touched guide/document links resolve locally', () => {
  for (const file of [skillPath, ...guides, 'README.md', 'ADOPTING.md', 'CREDITS.md',
    'skills/shared/cli-design/assets/line-tally/README.md']) {
    assert.ok(fs.existsSync(path.join(root, file)), `Missing ${file}`);
    for (const [, target] of read(file).matchAll(/\[[^\]]*\]\(([^\s)]+)\)/g)) {
      if (/^(https?:|mailto:|#)/.test(target)) continue;
      const resolved = path.resolve(root, path.dirname(file), target.split('#')[0]);
      assert.ok(fs.existsSync(resolved), `${file}: ${target}`);
    }
  }
});

for (const spec of [
  { name: 'help', args: ['--help'], code: 0,
    stdout: 'Usage: npm run test:cli [-- --help]\nRuns the opt-in .NET 10 CLI acceptance specimen in disposable project-local state.\nRequires an installed .NET 10 SDK/reference packs. No SDK/package installation or external feeds.\n', stderr: '' },
  { name: 'unknown option', args: ['--unknown'], code: 2, stdout: '',
    stderr: 'cli-contract: unknown option. Use --help.\n' },
]) {
  test(`CLI runner ${spec.name} requires no SDK or trial state`, () => {
    const env = Object.fromEntries(Object.entries(process.env)
      .filter(([key]) => key.toUpperCase() !== 'PATH'));
    env.PATH = '';
    env.NODE_DISABLE_COMPILE_CACHE = '1';
    const scratch = () => fs.readdirSync(root).filter((name) => name.startsWith('.cli-design-trial-')).sort();
    const before = scratch();

    const result = spawnSync(process.execPath, [path.join(root, 'tools/test-cli-contract.js'), ...spec.args],
      { cwd: root, env, encoding: 'utf8', timeout: 10000, shell: false });

    assert.deepEqual(scratch(), before, 'Trial state left behind');
    assert.equal(result.error, undefined);
    assert.equal(result.signal, null);
    assert.equal(result.status, spec.code);
    assert.equal(result.stdout, spec.stdout);
    assert.equal(result.stderr, spec.stderr);
  });
}
