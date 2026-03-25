#!/usr/bin/env node
// tools/test-release-readiness.js — Semantic runner for release-readiness
// Validates REQ-01..REQ-05: suite is in a releasable state.
// Exit: 0 if all pass, 1 if any fail
// No external dependencies — uses only Node.js builtins (fs, path, child_process)

'use strict';

const fs   = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const ROOT = path.join(__dirname, '..');

let passed = 0;
let failed = 0;

/**
 * Assert a condition and print pass/fail.
 * @param {string} name       Human-readable label
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

// ─── REQ-01: .git directory exists AND baseline commit evidence is present ────

console.log('\n── REQ-01: Git repo initialized with baseline commit ──');
{
  const gitDir = path.join(ROOT, '.git');
  const gitExists = fs.existsSync(gitDir);
  assert('.git directory exists', gitExists);

  if (gitExists) {
    let commitCount = 0;
    try {
      const output = execSync('git rev-list --count HEAD', {
        cwd: ROOT,
        stdio: ['pipe', 'pipe', 'pipe'],
      }).toString().trim();
      commitCount = parseInt(output, 10);
    } catch (_) {
      commitCount = 0;
    }
    assert('at least one baseline commit exists (git rev-list --count HEAD >= 1)', commitCount >= 1);
  }
}

// ─── REQ-02: Minimum project documentation covers repo identity + adoption ────

console.log('\n── REQ-02: Project documentation covers repo identity and adoption guidance ──');
{
  // Repo identity: README.md exists and identifies the repo/project
  const readmePath = path.join(ROOT, 'README.md');
  const readmeExists = fs.existsSync(readmePath);
  assert('README.md exists (repo identity)', readmeExists);

  if (readmeExists) {
    const content = fs.readFileSync(readmePath, 'utf8');
    assert(
      'README.md identifies the repo name (Zesh-One-Skills)',
      content.toLowerCase().includes('zesh-one-skills') || content.toLowerCase().includes('zesh-one skills')
    );
    assert(
      'README.md references AGENTS distribution files or skills/ (structure coverage)',
      content.includes('AGENTS.backend.md') || content.includes('skills/')
    );
  }

  // Adoption guidance: ADOPTING.md exists and covers onboarding
  const adoptingPath = path.join(ROOT, 'ADOPTING.md');
  const adoptingExists = fs.existsSync(adoptingPath);
  assert('ADOPTING.md exists (adoption guidance)', adoptingExists);

  if (adoptingExists) {
    const content = fs.readFileSync(adoptingPath, 'utf8');
    assert(
      'ADOPTING.md references npm run verify (verification step)',
      content.includes('npm run verify')
    );
  }
}

// ─── REQ-03: Vendor attribution includes BOTH github-pr AND skill-creator ─────

console.log('\n── REQ-03: CREDITS.md attributes both github-pr and skill-creator ──');
{
  const creditsPath = path.join(ROOT, 'CREDITS.md');
  const exists = fs.existsSync(creditsPath);
  assert('CREDITS.md exists', exists);

  if (exists) {
    const content = fs.readFileSync(creditsPath, 'utf8');
    assert(
      'CREDITS.md attributes github-pr vendor asset',
      content.includes('github-pr')
    );
    assert(
      'CREDITS.md attributes skill-creator vendor asset',
      content.includes('skill-creator')
    );
  }
}

// ─── REQ-04: Adoption criterion explicitly includes copying/using AGENTS.md ───

console.log('\n── REQ-04: ADOPTING.md explicitly guides copying AGENTS.md ──');
{
  const adoptingPath = path.join(ROOT, 'ADOPTING.md');
  const exists = fs.existsSync(adoptingPath);

  if (exists) {
    const content = fs.readFileSync(adoptingPath, 'utf8');
    assert(
      'ADOPTING.md mentions copying AGENTS.md (cp or copy instruction)',
      content.includes('AGENTS.md')
    );
    assert(
      'ADOPTING.md references the cp/copy action for AGENTS distribution file',
      content.includes('cp AGENTS.') || content.includes('copy AGENTS.') || content.includes('Copiá')
    );
  } else {
    assert('ADOPTING.md exists (required for REQ-04)', false);
    assert('ADOPTING.md copy instruction present (required for REQ-04)', false);
  }
}

// ─── REQ-05: Root distribution files AGENTS.backend.md and AGENTS.frontend.md exist ─

console.log('\n── REQ-05: Root distribution files AGENTS.backend.md and AGENTS.frontend.md exist ──');
{
  assert(
    'AGENTS.backend.md exists at repo root',
    fs.existsSync(path.join(ROOT, 'AGENTS.backend.md'))
  );
  assert(
    'AGENTS.frontend.md exists at repo root',
    fs.existsSync(path.join(ROOT, 'AGENTS.frontend.md'))
  );
}

// ─── Summary ─────────────────────────────────────────────────────────────────

console.log('\n' + '─'.repeat(50));
console.log(`Results: ${passed} passed, ${failed} failed`);

if (failed === 0) {
  console.log('✅ All release-readiness assertions passed.');
} else {
  console.error(`❌ ${failed} assertion(s) failed.`);
}

process.exit(failed > 0 ? 1 : 0);
