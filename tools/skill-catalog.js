'use strict';
// Runtime consumers read JSON; YAML is only a repository-maintenance dependency.
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');
const yaml = require('js-yaml');
const scopes = { backend: ['backend'], 'dotnet-10': ['backend'], frontend: ['frontend'], apps: ['app'], shared: ['frontend', 'backend', 'app'] };
const topics = {
  security: 'security', general: 'architecture', architecture: 'architecture',
  dataaccess: 'data-access', requests: 'data-access', responses: 'data-access',
  'tanstack-query': 'data-access', 'zod-4': 'validation', validations: 'validation',
  'tailwind-4': 'design', 'shadcn-ui': 'design', accessibility: 'accessibility',
  testing: 'testing', 'testing-unit': 'testing', performance: 'performance',
  'diagnosing-bugs': 'debugging', logging: 'observability', mapping: 'data-access',
  'github-pr': 'collaboration', 'cli-design': 'cli', 'zustand-5': 'state',
};
function filesUnder(dir) {
  if (!fs.existsSync(dir)) return [];
  return fs.readdirSync(dir, { withFileTypes: true }).sort((a, b) => a.name.localeCompare(b.name, 'en')).flatMap(entry => {
    if (entry.isSymbolicLink()) throw new Error(`Symlinks are not catalog inputs: ${entry.name}`);
    if (['node_modules', 'bin', 'obj', '.git'].includes(entry.name)) return [];
    const full = path.join(dir, entry.name);
    return entry.isDirectory() ? filesUnder(full) : [full];
  });
}
function buildCatalog(root) {
  const relative = p => path.relative(root, p).split(path.sep).join('/');
  const all = filesUnder(path.join(root, 'skills'));
  const skills = all.filter(p => path.basename(p) === 'SKILL.md').map(entry => {
    const text = fs.readFileSync(entry, 'utf8');
    const frontmatter = /^---\r?\n([\s\S]*?)\r?\n---/.exec(text);
    if (!frontmatter) throw new Error(`Missing metadata: ${entry}`);
    const meta = yaml.load(frontmatter[1]);
    const skillPath = relative(path.dirname(entry));
    const id = skillPath.replace(/^skills\//, '');
    const domain = id.split('/')[0];
    if (!scopes[domain] || !meta.name || !meta.description || !meta.metadata?.version) throw new Error(`Invalid metadata: ${id}`);
    // Optional catalog metadata allows a new skill to specialize without changing this generator.
    const custom = meta.metadata.catalog || {};
    const resources = new Set(filesUnder(path.dirname(entry)));
    // Preserve shared non-skill resources. Cross-skill links are not implicit installations.
    for (const resource of resources) {
      if (!resource.endsWith('.md')) continue;
      const body = fs.readFileSync(resource, 'utf8');
      for (const match of body.matchAll(/\]\(([^\s)#]+)(?:#[^)]*)?\)/g)) {
        if (/^[a-z]+:|^[/#]/i.test(match[1])) continue;
        const target = path.resolve(path.dirname(resource), decodeURIComponent(match[1]));
        const rel = relative(target);
        if (rel.startsWith('skills/shared/') && !rel.endsWith('/SKILL.md') && fs.existsSync(target) && fs.statSync(target).isFile()) resources.add(target);
      }
    }
    const description = String(meta.description).replace(/\s+/g, ' ').trim();
    return {
      id, name: String(meta.name), description, path: skillPath, version: String(meta.metadata.version),
      scopes: custom.scopes || scopes[domain],
      topics: custom.topics || [topics[meta.name] || 'framework'],
      technologies: custom.technologies || [domain === 'shared' ? 'any' : domain],
      triggers: custom.triggers || [description.replace(/^Trigger:\s*/i, '')],
      files: [...resources].sort().map(file => ({ path: relative(file), sha256: crypto.createHash('sha256').update(fs.readFileSync(file)).digest('hex') })),
    };
  }).sort((a, b) => a.id.localeCompare(b.id, 'en'));
  return { schemaVersion: 1, skills };
}
if (require.main === module) {
  const args = process.argv.slice(2);
  if (args.length > 1 || (args.length && !['--check', '--help'].includes(args[0]))) {
    console.error('skill-catalog: unknown arguments. Use --help.');
    process.exit(2);
  }
  if (args[0] === '--help') {
    console.log('Usage: node tools/skill-catalog.js [--check|--help]\nGenerate skills-index.json from local skill metadata, or validate freshness with --check.');
    process.exit(0);
  }
  const root = path.resolve(__dirname, '..');
  const output = JSON.stringify(buildCatalog(root), null, 2) + '\n';
  const target = path.join(root, 'skills-index.json');
  if (process.argv.includes('--check')) {
    if (!fs.existsSync(target) || fs.readFileSync(target, 'utf8') !== output) {
      console.error('Catalog is stale. Run npm run catalog:generate.');
      process.exitCode = 1;
    }
  } else {
    fs.writeFileSync(target, output);
    console.log('Generated skills-index.json');
  }
}
module.exports = { buildCatalog };
