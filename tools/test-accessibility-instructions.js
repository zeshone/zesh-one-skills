'use strict';
// Instruction contracts only: no browser, screen reader or compliance audit is run.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
const load = file => fs.existsSync(path.join(root, file)) ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const skillPath = 'skills/shared/accessibility/SKILL.md';
const skill = load(skillPath);
const includes = (text, values) => values.forEach(value => assert.ok(text.includes(value), `Missing: ${value}`));

test('accessibility applies to changed web/hybrid UI and prefers semantic controls', () => {
  includes(skill, ['changed web or hybrid UI', 'native controls', 'accessible name', 'ARIA', 'not a full audit']);
});
test('keyboard guidance distinguishes modal containment and meaningful focus restoration', () => {
  includes(skill, ['keyboard', 'visible focus', 'logical order', 'modal', 'Escape', 'restore focus', 'not trap focus in ordinary pages']);
});
test('forms and dynamic feedback do not depend on color or steal focus', () => {
  includes(skill, ['associated labels', 'error text', 'invalid state', 'status announcements', 'not color alone', 'without stealing focus']);
});
test('visual checks cover contrast, zoom, reflow and motion without claiming conformance', () => {
  includes(skill, ['contrast', 'zoom', 'reflow', 'reduced motion', 'target standard', 'not a compliance guarantee']);
});
test('verification separates automated semantics from real browser and assistive technology', () => {
  includes(skill, ['existing tooling', 'actual browser', 'assistive technology', 'jsdom', 'unperformed', 'explicit approval', '../security/baseline.md']);
});
test('accessibility routes preserve React SPA and installed Angular Ionic flows', () => {
  includes(load('AGENTS.frontend.md'), [skillPath, 'UI accessibility, keyboard/focus, labels and feedback', 'React SPA', 'Client -> Server Action -> .NET API']);
  includes(load('AGENTS.apps.md'), [skillPath, 'UI accessibility, keyboard/focus, labels and feedback', 'Reactive Forms and HttpClient', 'valid page shell', 'NavController']);
  for (const guide of ['AGENTS.backend.md', 'AGENTS.dotnet-10.md']) assert.ok(!load(guide).includes(skillPath));
  includes(skill, ['React SPA', 'Angular/Ionic', 'existing UI library', 'not require shadcn', 'no native plugin']);
});
test('accessibility script is additive and preserves default Node-only checks and dependencies', () => {
  const { scripts, devDependencies } = JSON.parse(load('package.json'));
  assert.equal(scripts['test:accessibility-instructions'], 'node tools/test-accessibility-instructions.js');
  for (const name of ['test', 'verify']) {
    const commands = scripts[name].split(' && ');
    for (const file of ['test-fixtures', 'spec-fixtures-logging', 'test-pilot-instructions', 'test-cli-instructions',
      'test-frontend-instructions', 'test-dotnet-instructions', 'test-debugging-instructions', 'test-accessibility-instructions',
      'lint-skills', 'test-release-readiness', 'test-lint-zero-skills']) {
      assert.equal(commands.filter(command => command === `node tools/${file}.js`).length, 1, file);
    }
    assert.doesNotMatch(scripts[name], /npm (install|ci)|test-(cli-contract|frontend-dom|dotnet-unit)\.js/);
  }
  assert.deepEqual(devDependencies, { 'js-yaml': '^4.1.0' });
});
test('accessibility skill and touched documentation resolve local links', () => {
  for (const file of [skillPath, 'AGENTS.frontend.md', 'AGENTS.apps.md', 'README.md', 'ADOPTING.md', 'CREDITS.md']) {
    assert.ok(load(file), file);
    for (const [, target] of load(file).matchAll(/\]\(([^)]+)\)/g)) {
      if (/^(https?:|#)/.test(target)) continue;
      assert.ok(fs.existsSync(path.resolve(root, path.dirname(file), target.split('#')[0])), `${file}: ${target}`);
    }
  }
});
