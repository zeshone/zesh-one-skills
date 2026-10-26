'use strict';
const fs = require('node:fs');
const path = require('node:path');
const crypto = require('node:crypto');
const { spawnSync } = require('node:child_process');
const args = process.argv.slice(2);
if (args.length > 1 || (args.length && !['--all', '--help'].includes(args[0]))) {
  console.error('Usage: node tools/build-zesh-skills.js [--all|--help]');
  process.exit(2);
}
if (args[0] === '--help') {
  console.log('Usage: node tools/build-zesh-skills.js [--all|--help]\nBuild the host binary, or six Windows/Linux/macOS amd64/arm64 binaries, in dist/zesh-skills/.\nRequires Go 1.25+. No publishing or global installation.');
  process.exit(0);
}
const root = path.resolve(__dirname, '..');
const hostOS = { win32: 'windows', linux: 'linux', darwin: 'darwin' }[process.platform];
const hostArch = { x64: 'amd64', arm64: 'arm64' }[process.arch];
const targets = args[0] === '--all'
  ? ['windows', 'linux', 'darwin'].flatMap(os => ['amd64', 'arm64'].map(arch => [os, arch]))
  : [[hostOS, hostArch]];
if (targets.some(target => target.some(value => !value))) {
  console.error('Unsupported build host. Use --all on a supported Go installation.');
  process.exit(2);
}
const checksums = [];
for (const [os, arch] of targets) {
  const relative = `${os}-${arch}/zesh-skills${os === 'windows' ? '.exe' : ''}`;
  const destination = path.join(root, 'dist/zesh-skills', relative);
  fs.mkdirSync(path.dirname(destination), { recursive: true });
  const result = spawnSync('go', ['build', '-trimpath', '-ldflags=-s -w', '-o', destination, './cli'], {
    cwd: root, env: { ...process.env, GOOS: os, GOARCH: arch, CGO_ENABLED: '0', GOTOOLCHAIN: 'local' }, stdio: 'inherit', shell: false,
  });
  if (result.error || result.status !== 0) {
    console.error(`Build failed for ${os}/${arch}${result.error ? `: ${result.error.message}` : ''}`);
    process.exit(result.status || 1);
  }
  checksums.push(`${crypto.createHash('sha256').update(fs.readFileSync(destination)).digest('hex')}  ${relative}`);
  console.log(`Built ${relative}`);
}
fs.writeFileSync(path.join(root, 'dist/zesh-skills/SHA256SUMS'), checksums.join('\n') + '\n');
