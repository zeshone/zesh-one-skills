'use strict';
// Instruction contracts only: no consumer application, model trial or license audit is run.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
const load = file => fs.existsSync(path.join(root, file)) ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const guidePath = 'docs/MAINTAINING.md';
const guide = load(guidePath);
const includes = (text, values) => values.forEach(value => assert.ok(text.includes(value), `Missing: ${value}`));

test('maintainer guide exists, is linked and separates maintenance from adoption', () => {
  assert.ok(guide, guidePath);
  assert.match(load('README.md'), /\]\((?:\.\/)?docs\/MAINTAINING\.md\)/);
  includes(guide, ['maintainer workflow', 'skills-review-rubric.md', '../ADOPTING.md']);
});
test('owner selection prefers targeted edits and preserves contextual gates', () => {
  includes(guide, ['existing owner', 'targeted edit', 'installed versions', 'runtime/profile',
    'architecture', 'lifecycle', 'security baseline', 'testing', 'diagnosis', 'accessibility']);
});
test('first-party frontmatter and operational conventions remain explicit', () => {
  includes(guide, ['name:', 'description:', 'Trigger:', 'license:', 'allowed-tools: Read Write Edit Bash',
    'metadata:', 'author: Zesh-One', 'version: "1.0"', 'space-separated string',
    'English-only', 'skills and domain AGENTS files', 'advisory', '80–140', 'Changelog']);
});
test('conditional routes and relocatable references preserve domain ownership', () => {
  includes(guide, ['AGENTS.*.md', '../AGENTS.backend.md', '../AGENTS.dotnet-10.md',
    '../AGENTS.frontend.md', '../AGENTS.apps.md', 'conditional', 'relative',
    'shared/security/', 'not automatic loading']);
});
test('instruction tests use RED GREEN and positive negative preservation route checks', () => {
  includes(guide, ['tools/test-*-instructions.js', 'positive', 'negative', 'preservation', 'route',
    'RED', 'GREEN', 'node:test', 'test-accessibility-instructions.js']);
});
test('registration retains existing checks in both default pipelines', () => {
  includes(guide, ['package.json', '`test` and `verify`', 'exactly once', 'retain existing checks',
    'npm run verify', 'git diff --check', 'tools/test-cli-instructions.js', 'exact default pipeline list']);
});
test('default checks stay Node-only and isolated trials require opt-in authorization', () => {
  includes(guide, ['Node-only', 'opt-in', 'isolated', 'explicit approval', 'test:cli',
    'test:frontend-dom', 'test:dotnet-unit', 'not part of default verification']);
});
test('documentation skill version and provenance updates are conditional and bounded', () => {
  includes(guide, ['README', 'ADOPTING', 'CHANGELOG', 'CREDITS', 'metadata.version',
    'source URL', 'revision', 'license', 'not recorded', 'copied versus conceptual']);
});
test('lint coverage distinguishes blocking errors warnings and parser qualifications', () => {
  includes(guide, ['What lint checks / what it does not check', 'tools/lint-skills.js',
    'five skill roots', '.github/skills/', 'LINT_SKILLS_ROOTS',
    'nonempty strings', 'metadata.author', 'metadata.version', 'empty string',
    'When to Use', 'Critical Patterns', 'Resources', 'exact `Trigger:`',
    'Keywords/Tags', 'Zesh-One', 'major', 'warning', 'prefix',
    'levels 1–6', 'arrays', 'comma-separated', 'fragments', '##']);
});
test('lint limits do not imply editorial provenance or behavioral correctness', () => {
  includes(guide, ['name syntax/length/directory alignment', 'semantic-version validity',
    'file length', 'English-only', 'seven-section', 'links outside Resources',
    'bare paths/URLs', 'attribution completeness', 'instruction correctness']);
});
test('repository guidance takes precedence over verified vendor template differences', () => {
  includes(guide, ['takes precedence', '../.github/skills/skill-creator/SKILL.md',
    'gentleman-programming', 'Apache-2.0', 'MIT', 'YAML list', 'Commands', 'domain']);
});
test('evidence checklist distinguishes static checks proposals and executed trials', () => {
  includes(guide, ['PR checklist', 'exit code', 'counts', 'unperformed', 'model compliance',
    'application correctness', 'security guarantee', 'proposal', 'executed', 'odd/evidence/']);
});
for (const file of [guidePath, 'README.md', 'CREDITS.md']) {
  test(`${file} local Markdown links resolve (fragments are not validated)`, () => {
    const text = load(file);
    assert.ok(text, file);
    for (const [, target] of text.matchAll(/\]\(([^)]+)\)/g)) {
      if (/^(https?:|mailto:|#)/.test(target)) continue;
      assert.ok(fs.existsSync(path.resolve(root, path.dirname(file), target.split('#')[0])), `${file}: ${target}`);
    }
  });
}
// Discover actual files, not a fixed list, including this script.
const instructionFiles = fs.readdirSync(__dirname)
  .filter(file => /^test-.*-instructions\.js$/.test(file)).sort();
for (const name of ['test', 'verify']) {
  test(`every instruction script appears exactly once in scripts.${name}`, () => {
    const { scripts } = JSON.parse(load('package.json'));
    assert.equal(typeof scripts[name], 'string');
    const commands = scripts[name].split('&&').map(command => command.trim());
    const mismatches = instructionFiles.map(file => {
      const count = commands.filter(command => command === `node tools/${file}`).length;
      return count === 1 ? null : `tools/${file}: expected 1, found ${count}`;
    }).filter(Boolean);
    assert.deepEqual(mismatches, [], `scripts.${name} parity:\n${mismatches.join('\n')}`);
  });
}
