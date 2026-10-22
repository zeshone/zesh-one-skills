#!/usr/bin/env node
'use strict';

// Public instruction-contract checks, not proof of model behavior or application security.
const fs = require('node:fs');
const path = require('node:path');
const assert = require('node:assert/strict');
const { test } = require('node:test');
const root = path.resolve(__dirname, '..');
const read = file => fs.existsSync(path.join(root, file))
  ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const selection = read('skills/shared/architecture-selection.md');
const architecture = read('skills/dotnet-10/architecture/SKILL.md');
const net8 = read('AGENTS.backend.md');
const net10 = read('AGENTS.dotnet-10.md');

function requireText(text, fragments) {
  for (const fragment of fragments) assert.ok(text.includes(fragment), `Missing instruction: ${fragment}`);
}

test('S2: contextual selection offers alternatives and waits for a user decision', () => {
  requireText(selection, ['application type', 'installed framework', 'existing architecture',
    'ownership', 'minimum sufficient', 'alternative', 'Wait for the user',
    'before scaffolding', 'already recorded', 'Clean/Hexagonal', 'DDD']);
  requireText(architecture, ['../../shared/architecture-selection.md', 'not size tiers']);
  requireText(read('skills/backend/general/SKILL.md'), ['../../shared/architecture-selection.md']);
});

test('S2: context count and speculative growth never force assembly promotion', () => {
  requireText(selection, ['A second bounded context alone', 'LOC', 'RPS', 'not automatic triggers']);
  for (const content of [architecture, net10]) {
    assert.doesNotMatch(content, /promote[^.\n]*when a second bounded context stabilizes/i);
    assert.doesNotMatch(content, /when a slice grows into a real bounded context, promote/i);
  }
  requireText(architecture, ['observed boundary problem', 'user-approved']);
});

test('S3: CLI stays minimal and separates product stdout from diagnostics', () => {
  requireText(selection, ['one application project plus tests', 'empty layers', 'CLI']);
  requireText(architecture, ['CLI stdout', 'diagnostics', 'Do not add HTTP', 'database', 'login']);
  requireText(net10, ['CLI', 'worker', 'only when', 'application type']);
  requireText(net8, ['CLI', 'HTTP-specific']);
});

test('S3: routed testing guidance does not add HTTP or databases to CLI file I/O', () => {
  const testing = read('skills/dotnet-10/testing/SKILL.md');
  requireText(testing, ['CLI/file I/O', 'temporary fixture files', 'stdout', 'stderr',
    'No WAF or database', 'only when', 'without building a web host']);
  assert.doesNotMatch(testing, /SQL or infrastructure seams/);
});

test('S5: applicable HTTP and real-database testing obligations remain intact', () => {
  requireText(read('skills/dotnet-10/testing/SKILL.md'), ['WebApplicationFactory<Program>',
    'Testcontainers with the production database engine', 'Never issue SQL from these tests',
    'hard-fail in CI when Docker is unavailable', 'never mock a transaction seam',
    'FakeTimeProvider', 'native `Assert` only', 'no reason message']);
});

test('S5: sound brownfield and existing .NET conventions remain intact', () => {
  requireText(architecture, ['work within the established structure', 'Domain/', 'IEntityTypeConfiguration<T>',
    'DbContext', 'MediatR', 'TypedResults', 'ProblemDetails', 'TimeProvider', 'never remove or bypass']);
  requireText(net8, ['FluentValidation', 'ResponseDTO<T>', 'AutoMapper/manual', 'NSubstitute']);
  requireText(net10, ['Result', 'FluentValidation / `AddValidation()` OFF', 'WAF-tier', 'InMemory']);
  requireText(read('skills/backend/general/SKILL.md'), ['app.UseAuthentication();', 'app.UseAuthorization();',
    'app.MapControllers();', '{SERVICE}_CONFIG_PATH']);
});

test('S4: integrated security retains sanitization, boundary scope and domain exceptions', () => {
  for (const guide of [net8, net10]) requireText(guide, ['skills/shared/security/baseline.md']);
  requireText(read('skills/shared/security/baseline.md'), ['before any output reaches the model',
    'Analysis-only requests authorize no fixes', 'Apply the selected domain', 'Do not weaken authentication']);
  requireText(read('skills/backend/security/SKILL.md'), ['Sequential IDs alone do not prove BOLA',
    'enforce ownership in the service layer', 'generic authentication failure']);
  requireText(read('skills/dotnet-10/security/SKILL.md'), ['when immediate revocation is required',
    'bound the window explicitly', '19456 KiB', 'may be owned at the edge',
    'FluentValidation validators are prohibited on authentication endpoints',
    'Do not add a second store lookup']);
});

test('S5: fixed scenario inputs and adopted local links remain usable', () => {
  const fixture = JSON.parse(read('tools/fixtures/architecture-adoption-pilot.json'));
  assert.equal(fixture.mode, 'response-only');
  assert.deepEqual(fixture.cases.map(c => c.id), ['cli-small', 'two-contexts', 'brownfield-clean',
    'react-spa', 'ionic-version', 'user-choice']);
  for (const c of fixture.cases) {
    assert.ok(c.prompt.length > 50);
    assert.ok(read(c.guide).length > 0, c.guide);
  }
  for (const file of ['skills/shared/architecture-selection.md', 'skills/dotnet-10/architecture/SKILL.md',
    'skills/backend/general/SKILL.md', 'skills/dotnet-10/testing/SKILL.md',
    'AGENTS.backend.md', 'AGENTS.dotnet-10.md', 'ADOPTING.md']) {
    const text = read(file);
    assert.ok(text.length > 0, file);
    for (const [, target] of text.matchAll(/\]\(([^)]+)\)/g)) {
      if (/^(https?:|#)/.test(target)) continue;
      assert.ok(fs.existsSync(path.resolve(root, path.dirname(file), target.split('#')[0])), `${file}: ${target}`);
    }
  }
});
