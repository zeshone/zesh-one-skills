#!/usr/bin/env node
// tools/lint-skills.js — Skill documentation linter for Zesh-One-Skills
// Rules: FM-001..FM-008, SEC-001..SEC-004, LNK-001
// Profiles: strict (Zesh-One author) vs relaxed (vendor)
// Exit: 0 if clean or warnings-only, 1 if any error

'use strict';

const fs   = require('fs');
const path = require('path');
let yaml;
try {
  yaml = require('js-yaml');
} catch {
  console.error('ERROR: js-yaml not installed. Run: npm install');
  process.exit(2);
}

// ─── Constants ──────────────────────────────────────────────────────────────

const SKILLS_ROOTS = [
  path.resolve(__dirname, '..', 'backend', 'skills'),
  path.resolve(__dirname, '..', 'vendor', 'gentleman'),
  path.resolve(__dirname, '..', 'shared', 'skills'),
];
const ZESH_ONE_AUTHOR = 'Zesh-One';

// Rules table — level depends on profile at runtime
const RULE_META = {
  'FM-001': { desc: 'Frontmatter must be valid YAML (delimited by ---)' },
  'FM-002': { desc: 'Field `name` must be a non-empty string' },
  'FM-003': { desc: 'Field `description` must be a non-empty string' },
  'FM-004': { desc: 'Field `metadata.author` must be a non-empty string' },
  'FM-005': { desc: 'Field `metadata.version` must be a string (not a number)' },
  'FM-006': { desc: 'Field `license` recommended for Zesh-One skills' },
  'FM-007': { desc: 'Field `allowed-tools` recommended for Zesh-One skills' },
  'FM-008': { desc: '`description` should contain "Trigger:" (best practice)' },
  'SEC-001': { desc: 'Section `## When to Use` must be present' },
  'SEC-002': { desc: 'Section `## Critical Patterns` must be present' },
  'SEC-003': { desc: '`## Changelog` required when version >= 2.0 (Zesh-One)' },
  'SEC-004': { desc: 'Section `## Resources` recommended' },
  'LNK-001': { desc: 'Relative link in `## Resources` points to a non-existent file' },
};

// ─── Helpers ─────────────────────────────────────────────────────────────────

/**
 * Recursively discover all SKILL.md files under `rootDir`.
 * @param {string} rootDir
 * @returns {string[]} absolute paths
 */
function discoverSkillFiles(rootDir) {
  const results = [];
  if (!fs.existsSync(rootDir)) return results;

  function walk(dir) {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const full = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        walk(full);
      } else if (entry.isFile() && entry.name === 'SKILL.md') {
        results.push(full);
      }
    }
  }
  walk(rootDir);
  return results;
}

/**
 * Split a SKILL.md content into frontmatter string and body string.
 * Returns null for frontmatter if no valid `---` delimiters found.
 * @param {string} content
 * @returns {{ frontmatterRaw: string|null, body: string, frontmatterEndLine: number }}
 */
function splitFrontmatter(content) {
  // Normalize CRLF → LF to handle Windows line endings gracefully
  const normalized = content.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
  const lines = normalized.split('\n');
  if (lines[0] !== '---') {
    return { frontmatterRaw: null, body: normalized, frontmatterEndLine: 0 };
  }
  const closeIdx = lines.findIndex((l, i) => i > 0 && l === '---');
  if (closeIdx === -1) {
    return { frontmatterRaw: null, body: normalized, frontmatterEndLine: 0 };
  }
  const frontmatterRaw = lines.slice(1, closeIdx).join('\n');
  const body = lines.slice(closeIdx + 1).join('\n');
  return { frontmatterRaw, body, frontmatterEndLine: closeIdx + 1 };
}

/**
 * Parse a semantic version string and return major number.
 * Returns 0 if unparseable.
 * @param {string} versionStr
 * @returns {number}
 */
function parseMajorVersion(versionStr) {
  if (typeof versionStr !== 'string') return 0;
  const match = versionStr.match(/^(\d+)/);
  return match ? parseInt(match[1], 10) : 0;
}

/**
 * Extract relative Markdown links from a `## Resources` section.
 * Only returns links that don't start with http/https.
 * @param {string} body
 * @returns {{ link: string, line: number }[]}
 */
function extractResourceLinks(body) {
  const lines = body.split('\n');
  const links = [];
  let inResources = false;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (/^##\s+Resources?/i.test(line)) {
      inResources = true;
      continue;
    }
    // Exit resources section on next ## heading
    if (inResources && /^##\s/.test(line)) {
      inResources = false;
    }
    if (!inResources) continue;

    // Match all [text](url) patterns
    const linkRe = /\[([^\]]*)\]\(([^)]+)\)/g;
    let m;
    while ((m = linkRe.exec(line)) !== null) {
      const url = m[2].trim();
      // Skip external URLs and anchors
      if (/^https?:\/\//i.test(url) || url.startsWith('#')) continue;
      links.push({ link: url, line: i + 1 }); // body line numbers are relative; caller adjusts
    }
  }
  return links;
}

// ─── Core Linter ─────────────────────────────────────────────────────────────

/**
 * @typedef {{ file: string, line: number, level: 'error'|'warning', ruleId: string, message: string }} LintResult
 */

/**
 * Lint a single SKILL.md file.
 * @param {string} filePath  Absolute path
 * @param {string} content   Raw file content
 * @returns {LintResult[]}
 */
function lintSkillFile(filePath, content) {
  const results = [];
  const relPath  = path.relative(process.cwd(), filePath).replace(/\\/g, '/');

  function push(level, ruleId, line, message) {
    results.push({ file: relPath, line, level, ruleId, message });
  }

  // ── Frontmatter ────────────────────────────────────────────────────────────
  const { frontmatterRaw, body, frontmatterEndLine } = splitFrontmatter(content);

  if (frontmatterRaw === null) {
    push('error', 'FM-001', 1, 'Missing or unclosed frontmatter delimiters (---)');
    // Cannot continue without frontmatter
    return results;
  }

  let fm;
  try {
    fm = yaml.load(frontmatterRaw) ?? {};
  } catch (e) {
    push('error', 'FM-001', 1, `Invalid YAML frontmatter: ${e.message}`);
    return results;
  }

  // Determine authorship profile
  const author     = fm?.metadata?.author ?? '';
  const isZeshOne  = author === ZESH_ONE_AUTHOR;

  // FM-002: name
  if (typeof fm.name !== 'string' || fm.name.trim() === '') {
    push('error', 'FM-002', 1, 'Missing or invalid `name` field in frontmatter');
  }

  // FM-003: description
  if (typeof fm.description !== 'string' || fm.description.trim() === '') {
    push('error', 'FM-003', 1, 'Missing or invalid `description` field in frontmatter');
  } else {
    // FM-008: description should contain "Trigger:"
    if (!fm.description.includes('Trigger:')) {
      push('warning', 'FM-008', 1, '`description` does not contain "Trigger:" (best practice)');
    }
  }

  // FM-004: metadata.author
  if (typeof author !== 'string' || author.trim() === '') {
    push('error', 'FM-004', 1, 'Missing or invalid `metadata.author` field');
  }

  // FM-005: metadata.version must be a string
  const version = fm?.metadata?.version;
  if (version === undefined || version === null) {
    push('error', 'FM-005', 1, 'Missing `metadata.version` field');
  } else if (typeof version !== 'string') {
    push('error', 'FM-005', 1, `\`metadata.version\` must be a string, got ${typeof version}`);
  }

  // FM-006: license (Zesh-One warning)
  if (isZeshOne && (typeof fm.license !== 'string' || fm.license.trim() === '')) {
    push('warning', 'FM-006', 1, '`license` field recommended for Zesh-One skills');
  }

  // FM-007: allowed-tools (Zesh-One warning)
  if (isZeshOne && !fm['allowed-tools']) {
    push('warning', 'FM-007', 1, '`allowed-tools` field recommended for Zesh-One skills');
  }

  // ── Section checks ─────────────────────────────────────────────────────────
  const headings = [];
  for (const line of body.split('\n')) {
    const m = line.match(/^(#{1,6})\s+(.+)/);
    if (m) headings.push(m[2].trim());
  }

  const hasHeading = (pattern) =>
    headings.some((h) => pattern.test(h));

  // SEC-001: ## When to Use
  if (!hasHeading(/^When\s+to\s+Use/i)) {
    push('error', 'SEC-001', frontmatterEndLine + 1, 'Missing required section `## When to Use`');
  }

  // SEC-002: ## Critical Patterns
  if (!hasHeading(/^Critical\s+Patterns/i)) {
    push('error', 'SEC-002', frontmatterEndLine + 1, 'Missing required section `## Critical Patterns`');
  }

  // SEC-003: ## Changelog required if version >= 2.0 (Zesh-One only)
  if (isZeshOne && typeof version === 'string') {
    const major = parseMajorVersion(version);
    if (major >= 2 && !hasHeading(/^Changelog/i)) {
      push('warning', 'SEC-003', 1, '`## Changelog` section required when version >= 2.0');
    }
  }

  // SEC-004: ## Resources recommended
  if (!hasHeading(/^Resources?/i)) {
    push('warning', 'SEC-004', frontmatterEndLine + 1, 'Section `## Resources` not found (recommended)');
  }

  // ── Link validation ─────────────────────────────────────────────────────────
  const resourceLinks = extractResourceLinks(body);
  const fileDir       = path.dirname(filePath);

  for (const { link, line } of resourceLinks) {
    // Strip fragment (#section)
    const cleanLink = link.split('#')[0];
    if (!cleanLink) continue; // pure anchor link

    const resolved = path.resolve(fileDir, cleanLink);
    if (!fs.existsSync(resolved)) {
      // Absolute line in file = frontmatter lines + body line
      const absoluteLine = frontmatterEndLine + line;
      push('error', 'LNK-001', absoluteLine, `Broken relative link: ${link}`);
    }
  }

  return results;
}

// ─── Main ────────────────────────────────────────────────────────────────────

function main() {
  const files = SKILLS_ROOTS.flatMap((root) => discoverSkillFiles(root));

  if (files.length === 0) {
    console.log('No SKILL.md files found. Nothing to lint.');
    process.exit(0);
  }

  console.log(`Linting ${files.length} SKILL.md file(s)...\n`);

  /** @type {LintResult[]} */
  const allResults = [];

  for (const filePath of files) {
    let content;
    try {
      content = fs.readFileSync(filePath, 'utf8');
    } catch (e) {
      const relPath = path.relative(process.cwd(), filePath).replace(/\\/g, '/');
      allResults.push({
        file: relPath,
        line: 0,
        level: 'error',
        ruleId: 'IO-001',
        message: `Cannot read file: ${e.message}`,
      });
      continue;
    }
    allResults.push(...lintSkillFile(filePath, content));
  }

  // ── Report ────────────────────────────────────────────────────────────────
  const errors   = allResults.filter((r) => r.level === 'error');
  const warnings = allResults.filter((r) => r.level === 'warning');

  // Sort: errors first, then warnings; within each group by file+line
  const sorted = [
    ...errors.sort((a, b) => a.file.localeCompare(b.file) || a.line - b.line),
    ...warnings.sort((a, b) => a.file.localeCompare(b.file) || a.line - b.line),
  ];

  for (const r of sorted) {
    console.log(`${r.file}:${r.line} ${r.level.toUpperCase().padEnd(7)} ${r.ruleId.padEnd(8)} ${r.message}`);
  }

  // ── Summary ───────────────────────────────────────────────────────────────
  console.log('');
  console.log(`─── Summary ─────────────────────────────────────────────`);
  console.log(`  Files checked : ${files.length}`);
  console.log(`  Errors        : ${errors.length}`);
  console.log(`  Warnings      : ${warnings.length}`);

  if (errors.length === 0 && warnings.length === 0) {
    console.log('\n✅ All skills are clean!');
  } else if (errors.length === 0) {
    console.log('\n⚠️  Warnings only — exit 0 (non-blocking)');
  } else {
    console.log(`\n❌ ${errors.length} error(s) found — exit 1 (blocking)`);
  }

  // Exit code: 1 if any errors, 0 otherwise
  process.exit(errors.length > 0 ? 1 : 0);
}

if (require.main === module) {
  main();
} else {
  module.exports = { lintSkillFile, splitFrontmatter, parseMajorVersion };
}
