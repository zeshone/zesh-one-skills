#!/usr/bin/env node
// tools/spec-fixtures-logging.js — Semantic fixture runner for net8-apirest-logging
// Validates 5 spec scenarios against the REAL content of backend/skills/net8-apirest/logging/SKILL.md
// Exit: 0 if all pass, 1 if any fail
// No external dependencies — uses only Node.js builtins

'use strict';

const fs = require('fs');
const path = require('path');

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

// ─── Read real skill file ─────────────────────────────────────────────────────

const skillPath = path.join(__dirname, '..', 'backend', 'skills', 'net8-apirest', 'logging', 'SKILL.md');
const content = fs.readFileSync(skillPath, 'utf8');

// ─── Scenarios ────────────────────────────────────────────────────────────────

console.log('\n── Scenario 1: Minimum Traceability Fields ──');
{
  // SKILL.md must document the 9 mandatory traceability fields
  assert(
    'contains CorrelationId field',
    content.includes('CorrelationId')
  );
  assert(
    'contains RequestPath or Endpoint field',
    content.includes('RequestPath') || content.includes('Endpoint')
  );
  assert(
    'contains StatusCode field',
    content.includes('StatusCode')
  );
  assert(
    'contains DurationMs or Duration field',
    content.includes('DurationMs') || content.includes('Duration')
  );
  assert(
    'contains UserId field',
    content.includes('UserId')
  );
  assert(
    'contains ExceptionType field',
    content.includes('ExceptionType')
  );
  assert(
    'contains Method or RequestMethod field',
    content.includes('RequestMethod') || content.includes('Method')
  );
}

console.log('\n── Scenario 2: CorrelationId in ProblemDetails Extensions ──');
{
  // SKILL.md must show correlationId exposed via ProblemDetails.Extensions
  assert(
    'contains Extensions["correlationId"]',
    content.includes('Extensions["correlationId"]')
  );
}

console.log('\n── Scenario 3: Canonical Exception Handling Middleware ──');
{
  // SKILL.md must present ExceptionHandlingMiddleware as the canonical mechanism
  assert(
    'contains ExceptionHandlingMiddleware',
    content.includes('ExceptionHandlingMiddleware')
  );
}

console.log('\n── Scenario 4: Cross-Skill Alignment References ──');
{
  // SKILL.md must reference sibling skills using relative paths
  assert(
    'references ../responses/SKILL.md',
    content.includes('../responses/SKILL.md')
  );
  assert(
    'references ../security/SKILL.md',
    content.includes('../security/SKILL.md')
  );
  assert(
    'references ../general/SKILL.md',
    content.includes('../general/SKILL.md')
  );
}

console.log('\n── Scenario 5: Exception Mapping Contract ──');
{
  // SKILL.md must document the canonical exception-to-HTTP mapping
  assert(
    'contains ValidationException',
    content.includes('ValidationException')
  );
  assert(
    'contains ForbiddenException',
    content.includes('ForbiddenException')
  );
  assert(
    'contains NotFoundException',
    content.includes('NotFoundException')
  );
  assert(
    'contains ConflictException',
    content.includes('ConflictException')
  );
}

// ─── Summary ─────────────────────────────────────────────────────────────────

console.log('\n' + '─'.repeat(50));
console.log(`Results: ${passed} passed, ${failed} failed`);

if (failed === 0) {
  console.log('✅ All semantic assertions passed.');
} else {
  console.error(`❌ ${failed} assertion(s) failed.`);
}

process.exit(failed > 0 ? 1 : 0);
