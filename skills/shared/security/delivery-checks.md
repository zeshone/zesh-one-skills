# Delivery checks — load by changed boundary

Use the matching section only. Follow the [baseline](baseline.md) for evidence, execution authorization, and sanitization. Do not install tools, restore dependencies, execute an unfamiliar script, contact a registry, or submit private manifests merely because a checker is available; establish its behavior and authorization first.

## Dependencies and restore sources

**Trigger:** `.csproj`, `Directory.Packages.props`, `packages.lock.json`, `NuGet.Config`, package sources, or restore/build configuration changes.
- Inspect declared and resolved versions, including transitive packages when resolution evidence is available. Do not invent resolved versions or assume the absence of a lock file proves compromise; check the project's restore/reproducibility policy.
- For an advisory, record the package/version, authoritative advisory source and date checked, affected range, fixed version if known, and relevant deployment conditions. Distinguish an affected installed version from a demonstrated reachable exploit; unknown reachability does not dismiss the advisory.
- Inspect source mapping, trusted feeds, and package identity where a change creates dependency-confusion or provenance risk. Treat feed credentials as secrets; inspect sanitized configuration only.
- Use existing approved SCA tooling only within the authorized boundary. If unavailable, report dependency scanning as not performed; manual manifest review is not an equivalent CVE scan. Do not introduce automatic package upgrades or fix commands.

## Infrastructure and containers

**Trigger:** Dockerfiles, image references, Compose, IaC, reverse proxy/TLS, exposed services, or deployment configuration changes.
- Check the effective runtime user, capabilities, mounts, network exposure, and secret delivery. Inspect inherited image/deployment settings before claiming a missing local directive is a vulnerability.
- Evaluate image provenance and pinning against the deployment policy. A floating tag is a reproducibility/hardening concern, not proof of malicious contents.
- Establish which layer owns TLS and security headers before reporting their absence. Preserve the domain's documented proxy exceptions; do not require duplicate controls blindly.
- Do not apply infrastructure, start containers, access production, or fetch images during analysis-only work. When runtime evidence is needed, state what remains unverified and propose an authorized isolated check.

## CI/CD and build trust

**Trigger:** workflow events, checkout refs, permissions, actions, build scripts, artifacts, publishing, or deployment jobs.
- Trace untrusted event/PR data into shell execution; distinguish data passed safely as arguments from executable interpolation. Establish attacker control and execution context before claiming injection.
- Check effective token permissions, secret availability, and trust boundaries around forked PR code, privileged workflow events, reusable workflows, self-hosted runners, and downloaded artifacts where touched.
- Inspect `pull_request_target` together with checkout/execution provenance and privileges; the event name alone is not proof of an exploit.
- Check third-party action pinning/provenance and artifact origin before privileged reuse. Classify policy/hardening gaps separately from demonstrated compromise paths.
- Verify fixes with sanitized fixtures or isolated checks where applicable. Do not trigger live workflows, publishing, deployments, or secret-bearing jobs merely to prove the change.
