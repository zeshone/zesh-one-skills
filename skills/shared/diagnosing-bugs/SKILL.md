---
name: diagnosing-bugs
description: "Trigger: bug, regression, unexpected behavior, diagnosis. Reproduce, isolate and verify the smallest evidence-based correction."
license: MIT
allowed-tools: Read Write Edit Bash Grep Glob
metadata:
  author: Zesh-One
  version: "1.0"
---

# Evidence-led bug diagnosis

## When to Use

Load for a reported defect or regression, not routine feature work or an automatic repository audit. Keep architecture, security and the existing testing profile with their current owners.

## Critical Patterns

1. **Reproduce first.** Record expected versus actual behavior, inputs, runtime/version and relevant state. Use a minimal reproduction through the public interface with synthetic data. Separate observations from hypotheses; if you cannot reproduce, report the missing evidence instead of declaring success.
2. **Locate the boundary.** Trace the failing boundary from observable input to output/state using bounded source reads and sanitized diagnostics. Compare the nearest passing case or last known good behavior; do not rewrite neighboring code merely because it looks suspicious.
3. **Test the explanation.** Choose one hypothesis and a discriminating experiment with a predicted outcome. Change one relevant variable, compare the observation and reject unsupported explanations. If evidence contradicts the hypothesis, stop speculative edits and reassess; do not stack guessed fixes.
4. **Correct and pin.** Apply the smallest correction supported by evidence. With an applicable deterministic runner, observe a RED regression before the fix and GREEN afterward; assert the result and all relevant observable effects, including unchanged state on rejection. Cover named boundaries and retain existing behavior in shared code. Do not introduce infrastructure just to test a pure rule.
5. **Keep diagnostics safe.** Follow the security baseline: no production replay, secret dumps, destructive probes or new dependencies without authorization. Reuse installed tools, sanitized examples and the selected runtime's testing/logging conventions.

## Constraints & Tradeoffs

A diagnosis may remain unknown. For timing, cancellation or intermittent failures, capture the relevant ordering and control supported clocks/signals rather than adding sleeps or claiming a flaky pass proves the cause. Do not widen scope into architecture, dependency upgrades or unrelated cleanup.

## Anti-Patterns

Patch before reproduction; treating correlation as cause; speculative refactors; weakening validation to make tests pass; leaking real payloads; claiming proposed checks were executed.

## Progressive Disclosure

Report: reproduction and observations; supported or unknown cause; bounded correction; exact checks/results and unperformed checks. Run the focused regression and applicable existing suite; disclose unavailable proof. A response-only proposal is not execution evidence or model-compliance proof.

## Resources

- Retain the [security baseline](../security/baseline.md) and its capability-specific routes.
- Reuse [.NET 8 testing](../../backend/testing-unit/SKILL.md), [.NET 10 testing](../../dotnet-10/testing/SKILL.md) or [React SPA testing](../../frontend/testing/SKILL.md) only within their applicability; other runtimes keep their installed test flow.
- For process/stream defects, retain the [CLI contract](../cli-design/SKILL.md).
