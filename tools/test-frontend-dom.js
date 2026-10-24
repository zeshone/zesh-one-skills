'use strict';

// Optional rendered-component trial. Never install packages from test/help paths.
const fs = require('node:fs');
const path = require('node:path');
const { spawnSync } = require('node:child_process');
const root = path.resolve(__dirname, '..');
const fixture = path.join(root, 'tools/fixtures/frontend-dom');
const args = process.argv.slice(2);
if (args.length === 1 && args[0] === '--help') {
  process.stdout.write('Usage: npm run test:frontend-dom [-- --help]\nRuns the optional React SPA component trial in jsdom, not a browser.\nInstall its isolated prerequisites only with approval; see tools/fixtures/frontend-dom/README.md.\n');
} else if (args.length) {
  process.stderr.write('frontend-dom: unknown option. Use --help.\n');
  process.exitCode = 2;
} else {
  run();
}

function run() {
  const packages = ['react', 'react-dom', '@testing-library/react', '@testing-library/dom', 'jsdom'];
  if (packages.some(name => !fs.existsSync(path.join(fixture, 'node_modules', name, 'package.json')))) {
    process.stderr.write('frontend-dom: isolated prerequisites missing. See tools/fixtures/frontend-dom/README.md; no installation was attempted.\n');
    process.exitCode = 2;
    return;
  }
  const scratch = fs.mkdtempSync(path.join(root, '.frontend-dom-trial-'));
  try {
    const result = spawnSync(process.execPath, ['--test', '--test-concurrency=1',
      path.join(fixture, 'contract.test.mjs')], {
      cwd: scratch, shell: false, encoding: 'utf8', timeout: 15000, maxBuffer: 1024 * 1024,
      env: { ...process.env, NODE_DISABLE_COMPILE_CACHE: '1',
        TEMP: scratch, TMP: scratch, TMPDIR: scratch },
    });
    process.stdout.write(result.stdout ?? '');
    process.stderr.write(result.stderr ?? '');
    if (result.error || result.signal) {
      process.stderr.write('frontend-dom: test process failed or timed out.\n');
      process.exitCode = 1;
    } else process.exitCode = result.status ?? 1;
  } finally {
    try { fs.rmSync(scratch, { recursive: true }); console.log('DOM_TRIAL_CLEANUP=removed'); }
    catch (error) { console.error(`DOM_TRIAL_CLEANUP=failed: ${error.message}`); process.exitCode = 1; }
  }
}
