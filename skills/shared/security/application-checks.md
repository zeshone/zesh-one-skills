# Application checks — load by changed boundary

Use only the relevant sections alongside the selected domain security skill. These are evidence prompts, not new architecture mandates. The [baseline](baseline.md) owns scope, secret-safe inspection, finding evidence, and reporting; existing domain rules own implementation defaults.

## Authentication

**Trigger:** identity establishment, credentials, tokens, cookies, recovery, MFA, or session lifecycle changes.
- Trace the actual scheme, token validation, and endpoint/group/global middleware before reporting missing authentication. Distinguish intentionally anonymous routes from protected operations.
- Check expiry, refresh/replay behavior, recovery errors, and revocation against the project's requirements. Inspect cryptography only where the changed flow uses it; preserve the domain's key and password rules.
- For cookie-authenticated state changes, check CSRF protections and cookie settings against the actual client/origin model. CORS is not authentication or a substitute for CSRF protection.
- Verify valid, invalid, expired, and revoked credentials only as applicable to the changed rule. Keep responses non-enumerating and test with synthetic credentials.

## Authorization

**Trigger:** ownership, policies, roles, resource queries, bulk operations, or new entry points.
- Trace effective enforcement through endpoint groups, policies, services, and scoped queries; a missing local annotation is not proof of missing authorization.
- Check that the authenticated identity, requested action, and target resource are bound together. Include list/export and alternate entry points when the change affects them, not only single-item reads.
- Treat unpredictable IDs as defense-in-depth, never authorization. Establish unauthorized access before claiming BOLA/IDOR.
- Test the allowed actor and a different authenticated actor, plus unauthenticated access where relevant; rejected mutations must leave data unchanged. Match session checks to the documented revocation policy.

## External inputs and sensitive outputs

**Trigger:** DTOs, validators, SQL, outbound URLs, file paths/uploads, process execution, redirects, or serialization changes.
- Follow attacker-controlled values to the sensitive operation; inspect parameterization, encoding, allowlists, canonicalization, and framework behavior where the value is used. Do not equate string interpolation everywhere with injection.
- For outbound requests or file access, consider redirects/resolved destinations or path containment when relevant; syntactically valid input can still cross an unauthorized boundary.
- For uploads, check allowed content, size limits, storage/execution exposure, and stream handling through the existing request/validation rules. Do not duplicate validators in a security layer.
- Check explicitly exposed fields and production error/log paths for sensitive data. Use the domain's mapping, validation, and observability rules rather than imposing one validator strategy on both .NET versions.
- Test the relevant rejected input and legitimate input, including observable error and state effects; report untested platform behavior separately.

## Secrets and sensitive configuration

**Trigger:** credential references, configuration sources, logging, key handling, or suspected exposure.
- Start with filenames, configuration schema, ignore/tracking metadata, and placeholder-only examples. Existence of a local `.env` does not prove a committed leak; `.gitignore` does not prove a tracked file was never exposed.
- Inspect values only through the baseline's pre-model sanitization boundary. Do not print credential-bearing commands or dump environment/configuration to diagnose failures.
- Distinguish a secret value from a path/reference, and a verified synthetic fixture from real material. If provenance is unknown, report a suspected exposure without reproducing the value or testing whether it is live.
- For confirmed exposure, recommend owner-managed revocation/rotation and exposure assessment; removing the visible string alone is not remediation. History rewriting and live rotation require separate authorization.
