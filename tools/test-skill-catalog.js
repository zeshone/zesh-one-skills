'use strict';
const assert = require('node:assert/strict');
const fs = require('node:fs');
const os = require('node:os');
const path = require('node:path');
const { buildCatalog } = require('./skill-catalog');
const root = path.resolve(__dirname, '..');
const catalog = buildCatalog(root);
assert.equal(catalog.schemaVersion, 1);
assert.equal(new Set(catalog.skills.map(s => s.id)).size, catalog.skills.length);
const discovered = fs.readdirSync(path.join(root, 'skills'), { recursive: true })
  .filter(name => path.basename(name) === 'SKILL.md');
assert.ok(catalog.skills.length > 0);
assert.equal(catalog.skills.length, discovered.length);
for (const skill of catalog.skills) {
  assert.ok(skill.description && skill.triggers.length && skill.topics.length && skill.scopes.length);
  assert.ok(skill.files.some(f => f.path === `${skill.path}/SKILL.md`));
  for (const file of skill.files) {
    assert.match(file.sha256, /^[a-f0-9]{64}$/);
    assert.ok(!file.path.includes('..'));
  }
}
assert.deepEqual(catalog, buildCatalog(root), 'generation must be deterministic');
assert.ok(catalog.skills.some(s => s.id === 'frontend/security'));
assert.ok(catalog.skills.some(s => s.id === 'backend/security'));
assert.ok(catalog.skills.find(s => s.id === 'shared/cli-design').files.some(f => f.path.endsWith('baseline.md')));
assert.deepEqual(JSON.parse(fs.readFileSync(path.join(root, 'skills-index.json'), 'utf8')), catalog, 'catalog is stale: npm run catalog:generate');
const tmp = fs.mkdtempSync(path.join(os.tmpdir(), 'zesh-catalog-'));
try {
  fs.mkdirSync(path.join(tmp, 'skills/frontend/new-design'), { recursive: true });
  fs.writeFileSync(path.join(tmp, 'skills/frontend/new-design/SKILL.md'), '---\nname: new-design\ndescription: >\n  Trigger: design.\n  New UI guidance.\nmetadata:\n  version: "1.0"\n---\n# Design\n');
  const added = buildCatalog(tmp).skills[0];
  assert.equal(added.id, 'frontend/new-design');
  assert.match(added.description, /New UI guidance/);
  assert.deepEqual(added.scopes, ['frontend']);
  assert.ok(added.topics.length);
} finally { fs.rmSync(tmp, { recursive: true, force: true }); }
// Exercise the actual maintenance command against isolated state, never the real catalog.
const { spawnSync } = require('node:child_process');
fs.mkdirSync(path.join(root, 'dist'), { recursive: true });
const scratch = fs.mkdtempSync(path.join(root, 'dist/catalog-command-'));
try {
  fs.mkdirSync(path.join(scratch, 'tools'));
  fs.mkdirSync(path.join(scratch, 'skills/frontend/design'), { recursive: true });
  fs.copyFileSync(path.join(__dirname, 'skill-catalog.js'), path.join(scratch, 'tools/skill-catalog.js'));
  fs.writeFileSync(path.join(scratch, 'skills/frontend/design/SKILL.md'), '---\nname: design\ndescription: Design\nmetadata:\n  version: "1"\n---\n');
  const target = path.join(scratch, 'skills-index.json');
  const invoke = args => spawnSync(process.execPath, [path.join(scratch, 'tools/skill-catalog.js'), ...args], { encoding: 'utf8', shell: false });
  const help = 'Usage: node tools/skill-catalog.js [--check|--help]\nGenerate skills-index.json from local skill metadata, or validate freshness with --check.\n';
  for (const [args, status, stdout, stderr] of [
    [['--help'], 0, help, ''],
    [['--unknown'], 2, '', 'skill-catalog: unknown arguments. Use --help.\n'],
    [['--check', '--help'], 2, '', 'skill-catalog: unknown arguments. Use --help.\n'],
    [['--check'], 1, '', 'Catalog is stale. Run npm run catalog:generate.\n'],
  ]) {
    fs.writeFileSync(target, 'local catalog edits\n');
    const result = invoke(args);
    assert.equal(result.status, status, args.join(' '));
    assert.equal(result.stdout, stdout);
    assert.equal(result.stderr, stderr);
    assert.equal(fs.readFileSync(target, 'utf8'), 'local catalog edits\n');
  }
  const generated = invoke([]);
  assert.equal(generated.status, 0);
  assert.equal(generated.stdout, 'Generated skills-index.json\n');
  assert.equal(generated.stderr, '');
  const bytes = fs.readFileSync(target, 'utf8');
  assert.deepEqual(JSON.parse(bytes), buildCatalog(scratch));
  const checked = invoke(['--check']);
  assert.equal(checked.status, 0);
  assert.equal(checked.stdout, '');
  assert.equal(checked.stderr, '');
  assert.equal(fs.readFileSync(target, 'utf8'), bytes);
} finally { fs.rmSync(scratch, { recursive: true, force: true }); }
console.log(`Catalog checks passed (${catalog.skills.length} skills; discovery, metadata, hashes, resources, freshness and command safety).`);
