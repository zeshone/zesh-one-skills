#!/usr/bin/env node
// tools/test-fixtures.js — In-memory fixture runner for skills-doc-lint
// Covers: FM-001, FM-004, SEC-003 (v2), SEC-003 negativo (v1), skill limpia, LNK-001 (broken xref), validations limpia
// Exit: 0 if all pass, 1 if any fail
// No external dependencies — uses only Node.js builtins + lint-skills.js exports

'use strict';

const { lintSkillFile } = require('./lint-skills.js');

let passed = 0;
let failed = 0;

/**
 * Assert a condition and print pass/fail.
 * @param {string} name  Human-readable label
 * @param {boolean} condition
 */
function assert(name, condition) {
  if (condition) {
    console.log(`  ✅ ${name}`);
    passed++;
  } else {
    console.error(`  ❌ ${name}`);
    failed++;
  }
}

// ─── Fixtures ─────────────────────────────────────────────────────────────────

// FM-001: YAML frontmatter inválido (YAML malformado)
const FIXTURE_FM001_INVALID_YAML = `---
name: test-skill
metadata: {author: Zesh-One, version: "1.0"
---

## When to Use
Use this skill.

## Critical Patterns
Pattern here.
`;

// FM-004: frontmatter válido pero sin metadata.author
const FIXTURE_FM004_MISSING_AUTHOR = `---
name: test-skill
description: "A test skill. Trigger: when testing."
license: MIT
metadata:
  version: "1.0"
---

## When to Use
Use this skill.

## Critical Patterns
Pattern here.
`;

// SEC-003 positivo: skill v2 sin sección ## Changelog
const FIXTURE_SEC003_V2_NO_CHANGELOG = `---
name: test-skill
description: "A test skill. Trigger: when testing."
license: MIT
metadata:
  author: Zesh-One
  version: "2.0"
---

## When to Use
Use this skill.

## Critical Patterns
Pattern here.
`;

// SEC-003 negativo: skill v1 sin ## Changelog — NO debe reportar SEC-003
const FIXTURE_SEC003_V1_NO_CHANGELOG = `---
name: test-skill
description: "A test skill. Trigger: when testing."
license: MIT
metadata:
  author: Zesh-One
  version: "1.0"
---

## When to Use
Use this skill.

## Critical Patterns
Pattern here.
`;

// Skill limpia: v2 con todos los campos, todas las secciones y changelog
const FIXTURE_VALID_SKILL = `---
name: test-skill
description: "A test skill. Trigger: when testing."
license: MIT
allowed-tools:
  - bash
metadata:
  author: Zesh-One
  version: "2.0"
---

## When to Use
Use this skill when testing.

## Critical Patterns
Pattern here.

## Resources
No external resources.

## Changelog

- 2.0: Initial release.
`;

// LNK-001: cross-reference roto estilo pre-fix de validations (../../net8-testing/unit/SKILL.md)
const FIXTURE_VALIDATIONS_BROKEN_XREF = `---
name: net8-apirest-validations
description: "Skill de validaciones para .NET 8. Trigger: when validating."
license: MIT
metadata:
  author: Zesh-One
  version: "1.0"
---

## When to Use
Use this skill when adding FluentValidation to .NET 8 API endpoints.

## Critical Patterns
Always register validators in DI container.

## Resources
- **Testing validators**: See [testing](../../net8-testing/unit/SKILL.md)
`;

// Validations limpia: misma estructura, pero sin links relativos rotos
const FIXTURE_VALIDATIONS_CLEAN = `---
name: net8-apirest-validations
description: "Skill de validaciones para .NET 8. Trigger: when validating."
license: MIT
metadata:
  author: Zesh-One
  version: "1.0"
---

## When to Use
Use this skill when adding FluentValidation to .NET 8 API endpoints.

## Critical Patterns
Always register validators in DI container.

## Resources
- **FluentValidation docs**: https://docs.fluentvalidation.net/en/latest/aspnet.html
`;

// ─── Scenarios ────────────────────────────────────────────────────────────────

console.log('\n── Scenario 1: FM-001 — Invalid YAML frontmatter ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_FM001_INVALID_YAML);
  assert(
    'reports ruleId FM-001',
    results.some(r => r.ruleId === 'FM-001')
  );
}

console.log('\n── Scenario 2: FM-004 — Missing metadata.author ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_FM004_MISSING_AUTHOR);
  assert(
    'reports ruleId FM-004',
    results.some(r => r.ruleId === 'FM-004')
  );
}

console.log('\n── Scenario 3: SEC-003 — v2 skill without Changelog ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_SEC003_V2_NO_CHANGELOG);
  assert(
    'reports ruleId SEC-003',
    results.some(r => r.ruleId === 'SEC-003')
  );
}

console.log('\n── Scenario 4: SEC-003 negative — v1 skill without Changelog ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_SEC003_V1_NO_CHANGELOG);
  assert(
    'does NOT report SEC-003 for v1 skill',
    !results.some(r => r.ruleId === 'SEC-003')
  );
}

console.log('\n── Scenario 5: Valid skill — zero errors ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_VALID_SKILL);
  const errors = results.filter(r => r.level === 'error');
  assert(
    'zero errors for fully valid skill',
    errors.length === 0
  );
}

console.log('\n── Scenario 6: LNK-001 — broken cross-reference (pre-fix validations path) ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_VALIDATIONS_BROKEN_XREF);
  assert(
    'reports ruleId LNK-001 for broken relative path',
    results.some(r => r.ruleId === 'LNK-001')
  );
}

console.log('\n── Scenario 7: validations clean — zero errors ──');
{
  const results = lintSkillFile('/fake/SKILL.md', FIXTURE_VALIDATIONS_CLEAN);
  assert(
    'zero errors for clean validations skill',
    results.filter(r => r.level === 'error').length === 0
  );
}

// ─── Summary ─────────────────────────────────────────────────────────────────

console.log('\n' + '─'.repeat(50));
console.log(`Results: ${passed} passed, ${failed} failed`);

if (failed === 0) {
  console.log('✅ All fixture assertions passed.');
} else {
  console.error(`❌ ${failed} assertion(s) failed.`);
}

process.exit(failed > 0 ? 1 : 0);
