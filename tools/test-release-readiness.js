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
  assert('ADOPTING.md exists (required for REQ-04)', exists);

  if (!exists) {
    // Skip content assertions when the file is missing — existence already failed.
  } else {
    const content = fs.readFileSync(adoptingPath, 'utf8');
    assert(
      'ADOPTING.md mentions copying AGENTS.md (cp or copy instruction)',
      content.includes('AGENTS.md')
    );
    assert(
      'ADOPTING.md references the cp/copy action for AGENTS distribution file',
      content.includes('cp AGENTS.') || content.includes('copy AGENTS.') || content.includes('Copiá')
    );
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

// CI bootstrap must work before the CLI module lands and on child PR bases.
{
  const yaml = require('js-yaml');
  const ci = yaml.load(fs.readFileSync(path.join(ROOT, '.github/workflows/ci.yml'), 'utf8'));
  const pr = yaml.load(fs.readFileSync(path.join(ROOT, '.github/workflows/pr-validation.yml'), 'utf8'));
  assert('CI accepts all PR bases', !ci.on.pull_request?.branches && Object.hasOwn(ci.on, 'pull_request'));
  assert('PR governance accepts all PR bases', !pr.on.pull_request?.branches && Object.hasOwn(pr.on, 'pull_request'));
  assert('no privileged PR target event', !ci.on.pull_request_target && !pr.on.pull_request_target);
  assert('read-only CI permission', ci.permissions.contents === 'read');
  assert('existing Node job preserved', ci.jobs.test.name === 'Test & lint skills' && ci.jobs.test.steps.some(s => s.run === 'npm test'));
  const go = ci.jobs.go;
  assert('Go three-host matrix', JSON.stringify(go?.strategy?.matrix?.os) === JSON.stringify(['ubuntu-latest', 'windows-latest', 'macos-latest']));
  const steps = go?.steps || [];
  const guard = steps.findIndex(s => s.id === 'module');
  assert('module detection after checkout', guard > 0 && steps[0].uses === 'actions/checkout@v4' && steps[guard].run.includes('test -f go.mod'));
  for (const name of ['setup', 'test', 'vet']) {
    const step = steps.find(s => s.id === name);
    assert(`Go ${name} guarded on module presence`, step?.if === "steps.module.outputs.present == 'true'");
  }
  const setup = steps.find(s => s.id === 'setup');
  assert('official Go setup uses module and no absent-sum cache', setup?.uses === 'actions/setup-go@v5' && setup.with['go-version-file'] === 'go.mod' && setup.with.cache === false);
  assert('Go test and vet commands', steps.some(s => s.run === 'go test ./cli/...') && steps.some(s => s.run === 'go vet ./cli/...'));
  const governance = pr.jobs.validate.steps[0].with.script;
  assert('approved issue and single type governance retained', governance.includes('status:approved') && governance.includes('typeLabels.length !== 1') && governance.includes('core.setFailed'));
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
