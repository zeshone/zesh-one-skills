---
name: security
description: >
  Operational security defaults for backend APIs with OWASP API focus: ownership checks (BOLA), JWT/policy boundaries, CORS/rate limits, and secret handling.
  Trigger: When implementing or reviewing authn/authz, resource access controls, API hardening, user-owned data, external inputs, secrets, dependencies, infrastructure, or CI/CD.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "2.1"
allowed-tools: Read Edit Write Bash Glob Grep
---

## When to Use

- Designing or reviewing endpoint-level access control
- Implementing ownership checks for user-scoped resources (BOLA/IDOR prevention)
- Defining JWT claim boundaries between controller, service, and policy layers
- Hardening API surface with CORS and rate limiting
- Validating security behavior in service-level tests
- Changing external inputs, sensitive outputs, secret handling, dependencies, infrastructure, or CI/CD

Load the [shared baseline](../../shared/security/baseline.md) first. Apply the domain rules below to the affected flow; they are not a mandate to audit every endpoint or migrate an existing application.

| Changed boundary | Load on demand |
|---|---|
| Authentication, authorization, inputs, or secrets | Matching section of [application checks](../../shared/security/application-checks.md), then the relevant local rules |
| Packages, restore sources, infrastructure, or CI/CD | Matching section of [delivery checks](../../shared/security/delivery-checks.md) |
| Validation, error mapping, or tests | Existing neighboring skill in Resources; do not duplicate its rules |

---

## Critical Patterns

1) **DO** use non-enumerable public IDs (GUID/ULID) for externally addressable resources; **DON'T** expose sequential `int` IDs.  
Why: unpredictable IDs reduce enumeration opportunities; this is a project convention and defense-in-depth, not authorization. Sequential IDs alone do not prove BOLA, and GUIDs do not prevent it. Rule 2 remains the enforcement boundary.

2) **DO** enforce ownership in the service layer for every `Get/Update/Delete by Id`; **DON'T** leave it only in controllers.  
Why: services are the real policy boundary reused across entry points.

3) **DO** treat missing/invalid identity claim as unauthenticated (`401`); **DON'T** map it to domain-forbidden (`403`).  
Why: `401` = identity not established; `403` = identity established but not allowed.

4) **DO** use explicit policy/role checks for non-owner access paths; **DON'T** rely on implicit assumptions from UI flows.  
Why: API is callable without UI constraints.

5) **DO** return generic authentication failure messages; **DON'T** reveal if user/account exists.  
Why: prevents user enumeration during credential attacks.

6) **DO** segment rate-limit policies by risk profile (`auth` stricter than `read`); **DON'T** apply one global limiter blindly.  
Why: brute-force surfaces need tighter controls than normal traffic.

7) **DO** define allowlist CORS origins from configuration; **DON'T** use `AllowAnyOrigin` in production.  
Why: wide origins increase CSRF/data exfiltration exposure.

8) **DO** keep secrets outside committed config and inject at runtime; **DON'T** store API/JWT secrets in repo files.  
Why: git history is long-lived and difficult to fully sanitize.

9) **DO** centralize exception-to-HTTP mapping in middleware; **DON'T** duplicate response-shaping logic per controller action.  
Why: consistent contracts and fewer drift bugs.

10) **DO** write unit tests for ownership/policy branches in services; **DON'T** assume auth middleware alone proves authorization correctness.  
Why: BOLA defects usually live in business paths, not token parsing.

```csharp
if (!Guid.TryParse(User.FindFirstValue(ClaimTypes.NameIdentifier), out var userId))
    return Unauthorized(new ProblemDetails { Status = 401, Title = "Unauthorized", Detail = "Missing or invalid identity claim." });
return Ok(await _orders.GetByIdAsync(id, userId)); // service enforces ownership check
```

---

## Constraints & Tradeoffs

- Strong ownership checks add boilerplate in service methods, but they are NON-NEGOTIABLE for BOLA prevention.
- Generic auth errors improve security but reduce diagnostics for clients; use internal logs/telemetry for detail.
- Strict rate limits reduce abuse but can impact legitimate bursts; tune per route criticality.
- Fine-grained policies improve least-privilege control but increase maintenance complexity.
- CORS allowlists improve safety but require environment-aware config discipline.
- Centralized policy logic improves consistency but demands clear naming and documentation for policy intent.

---

## Anti-Patterns

- Ownership check implemented only in controller actions.
- `401`/`403` semantics inverted for missing identity claims.
- Treating unpredictable public IDs as authorization, or labeling sequential IDs alone as a confirmed BOLA vulnerability (rule 1 remains the public-ID convention).
- CORS `AllowAnyOrigin` in production APIs.
- Secrets in `appsettings.json`, `.env` committed, or docs examples with real values.
- Security controls explained as tutorial/setup instead of repository defaults and constraints.
- Assuming JWT validation alone enforces per-resource authorization.

---

## Progressive Disclosure

1. **Start here:** shared baseline and the changed-boundary table; apply rules 1, 2, 3, 5, 7, 8 where relevant.
2. **Then harden abuse controls:** rule 6 + per-route limiter strategy.
3. **Then tighten policy boundaries:** rule 4 + reusable policy requirements.
4. **Then validate behavior:** rule 10 with focused unit tests on service branches; follow the baseline's regression-test and evidence contract for authorized fixes.
5. **Close security-relevant work:** distinguish checked scope, pending hypotheses, areas not reviewed, and any justified deeper-audit need as defined in the baseline. No security score or blanket assurance.

---

## Resources

- Always-loaded development guardrails and evidence/output contract: [shared baseline](../../shared/security/baseline.md)
- On authentication, authorization, inputs, or secret changes: [application checks](../../shared/security/application-checks.md)
- On dependency, infrastructure, or CI/CD changes: [delivery checks](../../shared/security/delivery-checks.md)
- API response contract and exception mapping: [`../responses/SKILL.md`](../responses/SKILL.md)
- Middleware order and base backend conventions: [`../general/SKILL.md`](../general/SKILL.md)
- Validation boundaries: [`../validations/SKILL.md`](../validations/SKILL.md)
- Unit-testing ownership and policy logic: [`../testing-unit/SKILL.md`](../testing-unit/SKILL.md)
- OWASP API Security Top 10 (2023): https://owasp.org/API-Security/editions/2023/en/0x00-introduction/

---

## Changelog

### v2.1
- Added shared everyday guardrails and change-triggered references without replacing .NET 8 controls or validation/testing conventions.
- Distinguished public-ID hardening from proven BOLA; adopted contextual evidence, pre-model secret sanitization, regression verification, and bounded coverage reporting.

### v2.0 — 2026-04-21
- P2 trim: reduced encyclopedic setup content and kept an operational default set.
- Standardized structure to: When to Use, Critical Patterns, Constraints & Tradeoffs, Anti-Patterns, Progressive Disclosure, Resources, Changelog.
- Consolidated to 10 atomic Do/Don't security rules with concise rationale.
- Preserved high-value cross-references and one short canonical snippet.
