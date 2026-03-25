#!/usr/bin/env node
// tools/test-lint-zero-skills.js
// Runtime proof: lint-skills.js MUST exit 1 when zero SKILL.md files are found.
// Uses LINT_SKILLS_ROOTS env-var override to point the linter at empty temp dirs.
// No new npm dependencies — relies solely on Node.js built-ins.

'use strict';

const { execFileSync } = require('child_process');
const fs   = require('fs');
const os   = require('os');
const path = require('path');

// ─── Helpers ────────────────────────────────────────────────────────────────

let passed = 0;
let failed = 0;

function assert(name, condition) {
  if (condition) {
    console.log(`  ✅ ${name}`);
    passed++;
  } else {
    console.error(`  ❌ ${name}`);
    failed++;
  }
}

// ─── Setup ──────────────────────────────────────────────────────────────────

// Create a unique temp base dir with three empty sub-dirs mirroring the
// default skills layout (backend / frontend / shared), but containing no
// SKILL.md files.
const tmpBase = fs.mkdtempSync(path.join(os.tmpdir(), 'lint-zero-'));
const roots = ['backend', 'frontend', 'shared'].map((d) => {
  const p = path.join(tmpBase, d);
  fs.mkdirSync(p);
  return p;
});

// ─── Scenario ───────────────────────────────────────────────────────────────

console.log('\n── Scenario: Linter finds no skills → must exit 1 ─────────────────────────\n');

let exitCode = 0;
try {
  // execFileSync throws when the child exits non-zero — that is exactly the
  // behaviour we want to assert, so we catch and record the exit status.
  execFileSync(
    process.execPath,
    [path.join(__dirname, 'lint-skills.js')],
    {
      env: { ...process.env, LINT_SKILLS_ROOTS: roots.join(path.delimiter) },
      stdio: 'pipe',
    },
  );
  // If we reach here the linter exited 0 — that is the failure case.
  exitCode = 0;
} catch (err) {
  exitCode = typeof err.status === 'number' ? err.status : 1;
}

assert('lint-skills exits with code 1 when zero SKILL.md files found', exitCode === 1);

// ─── Teardown ───────────────────────────────────────────────────────────────

fs.rmSync(tmpBase, { recursive: true, force: true });

// ─── Summary ────────────────────────────────────────────────────────────────

console.log(`\n─── Summary ─────────────────────────────────────────────`);
console.log(`  Passed: ${passed}`);
console.log(`  Failed: ${failed}`);

process.exit(failed > 0 ? 1 : 0);
