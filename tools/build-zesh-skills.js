'use strict';
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');
const { spawnSync } = require('node:child_process');
const usage = 'Usage: node tools/build-zesh-skills.js [--all|--release <version>|--help]';
const targets = ['windows', 'linux', 'darwin'].flatMap(os => ['amd64', 'arm64'].map(arch => [os, arch]));

function parseArgs(args) {
  if (!args.length) return { mode: 'host' };
  if (args.length === 1 && args[0] === '--help') return { mode: 'help' };
  if (args.length === 1 && args[0] === '--all') return { mode: 'all' };
  // A plain numeric release triplet only; no path components, shell text or aliases.
  if (args.length === 2 && args[0] === '--release' && /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$/.test(args[1])) {
    return { mode: 'release', version: args[1] };
  }
  throw new Error(usage);
}
function assetName(os, arch) {
  return `zesh-skills-${os}-${arch}${os === 'windows' ? '.exe' : ''}`;
}
function checksum(destination, relative) {
  return `${crypto.createHash('sha256').update(fs.readFileSync(destination)).digest('hex')}  ${relative}`;
}
// Read exactly the six known inputs before replacing the manifest. Unknown files
// remain untouched and are never included in the upload/checksum contract.
function packageRelease(directory) {
  const lines = targets.map(([os, arch]) => {
    const name = assetName(os, arch);
    return checksum(path.join(directory, name), name);
  });
  fs.writeFileSync(path.join(directory, 'SHA256SUMS'), lines.join('\n') + '\n');
}
function build(args, { root = path.resolve(__dirname, '..'), run = spawnSync } = {}) {
  const options = parseArgs(args); // Validate before any writes or Go invocation.
  if (options.mode === 'help') {
    console.log(`${usage}\nBuild the host binary, or six Windows/Linux/macOS amd64/arm64 binaries.\nDevelopment: dist/zesh-skills/. Release: dist/zesh-skills-release/v<version>/.\nRequires Go 1.25+. No publishing or global installation.`);
    return;
  }
  const hostOS = { win32: 'windows', linux: 'linux', darwin: 'darwin' }[process.platform];
  const hostArch = { x64: 'amd64', arm64: 'arm64' }[process.arch];
  const selected = options.mode === 'host' ? [[hostOS, hostArch]] : targets;
  if (selected.some(target => target.some(value => !value))) {
    throw new Error('Unsupported build host. Use --all on a supported Go installation.');
  }
  const release = options.mode === 'release';
  const directory = path.join(root, release ? `dist/zesh-skills-release/v${options.version}` : 'dist/zesh-skills');
  const checksums = [];
  for (const [os, arch] of selected) {
    const relative = release ? assetName(os, arch) : `${os}-${arch}/zesh-skills${os === 'windows' ? '.exe' : ''}`;
    const destination = path.join(directory, relative);
    fs.mkdirSync(path.dirname(destination), { recursive: true });
    const flags = `-ldflags=-s -w${release ? ` -X main.version=${options.version}` : ''}`;
    const result = run('go', ['build', '-trimpath', flags, '-o', destination, './cli'], {
      cwd: root, env: { ...process.env, GOOS: os, GOARCH: arch, CGO_ENABLED: '0', GOTOOLCHAIN: 'local' }, stdio: 'inherit', shell: false,
    });
    if (result.error || result.status !== 0) {
      const error = new Error(`Build failed for ${os}/${arch}${result.error ? `: ${result.error.message}` : ''}`);
      error.exitCode = result.status || 1;
      throw error;
    }
    checksums.push(checksum(destination, relative));
    console.log(`Built ${relative}`);
  }
  if (release) packageRelease(directory);
  else fs.writeFileSync(path.join(directory, 'SHA256SUMS'), checksums.join('\n') + '\n');
}
if (require.main === module) {
  try { build(process.argv.slice(2)); }
  catch (error) {
    console.error(error.message);
    process.exitCode = error.message === usage || error.message.startsWith('Unsupported build host') ? 2 : (error.exitCode || 1);
  }
}
module.exports = { parseArgs, build, packageRelease, targets };
