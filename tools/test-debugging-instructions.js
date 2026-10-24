'use strict';
// Instruction contracts, not executable bug reproduction or model-compliance proof.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
const load = file => fs.existsSync(path.join(root, file)) ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const skillPath = 'skills/shared/diagnosing-bugs/SKILL.md';
const skill = load(skillPath);
const includes = (text, values) => values.forEach(value => assert.ok(text.includes(value), `Missing: ${value}`));

test('diagnosis starts from a bounded reproducible observation, not a guessed cause', () => {
  includes(skill, ['expected versus actual', 'minimal reproduction', 'observations from hypotheses', 'cannot reproduce']);
});
test('diagnosis uses one discriminating experiment before the smallest correction', () => {
  includes(skill, ['failing boundary', 'one hypothesis', 'discriminating experiment', 'smallest correction', 'stop speculative edits']);
});
test('regression checks pin rejection effects and preserve existing test profiles', () => {
  includes(skill, ['RED', 'GREEN', 'public interface', 'unchanged state', 'existing testing profile', 'not introduce infrastructure']);
});
test('diagnostic evidence preserves safety and reports missing proof honestly', () => {
  includes(skill, ['synthetic', 'sanitized', 'no production replay', 'unknown', 'unperformed', '../security/baseline.md']);
});
test('debugging routes are conditional and retain each existing runtime owner', () => {
  for (const [file, preserved] of [
    ['AGENTS.backend.md', 'NSubstitute, FluentAssertions'], ['AGENTS.dotnet-10.md', '`FakeTimeProvider`'],
    ['AGENTS.frontend.md', 'Client -> Server Action -> .NET API'], ['AGENTS.apps.md', 'Reactive Forms and HttpClient'],
  ]) includes(load(file), [skillPath, 'Bug reproduction, diagnosis, regression investigation', preserved, 'skills/shared/security/baseline.md']);
});
test('debugging check is additive in both Node-only pipelines', () => {
  const { scripts, devDependencies } = JSON.parse(load('package.json'));
  assert.equal(scripts['test:debugging-instructions'], 'node tools/test-debugging-instructions.js');
  for (const name of ['test', 'verify']) {
    includes(scripts[name], ['test-debugging-instructions.js', 'test-dotnet-instructions.js', 'test-frontend-instructions.js', 'test-cli-instructions.js']);
    assert.doesNotMatch(scripts[name], /npm (install|ci)|test-(cli-contract|frontend-dom|dotnet-unit)\.js/);
  }
  assert.deepEqual(devDependencies, { 'js-yaml': '^4.1.0' });
});
test('debugging documentation resolves local links and frozen scenarios retain their scope', () => {
  for (const file of [skillPath, 'README.md', 'ADOPTING.md', 'CREDITS.md']) {
    assert.ok(load(file), file);
    for (const [, target] of load(file).matchAll(/\]\(([^)]+)\)/g)) {
      if (/^(https?:|#)/.test(target)) continue;
      assert.ok(fs.existsSync(path.resolve(root, path.dirname(file), target.split('#')[0])), `${file}: ${target}`);
    }
  }
  const fixture = JSON.parse(load('tools/fixtures/debugging-accessibility-proposals.json'));
  assert.equal(fixture.baseline, 'e50ea92');
  assert.deepEqual(fixture.cases.map(c => c.id), ['rejected-capacity', 'react-dialog', 'ionic-focus']);
});
