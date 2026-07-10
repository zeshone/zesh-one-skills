<!--
  Every PR must link an approved issue and carry exactly one `type:*` label.
  The PR Validation workflow enforces both.
-->

## Linked Issue
<!-- REQUIRED: the linked issue must have the `status:approved` label. -->
Closes #

## Type
<!-- Check exactly ONE and add the matching `type:*` label to the PR. -->
- [ ] Bug fix (`type:bug`)
- [ ] New feature (`type:feature`)
- [ ] Documentation only (`type:docs`)
- [ ] Code/content refactoring (`type:refactor`)
- [ ] Maintenance / tooling (`type:chore`)
- [ ] Breaking change (`type:breaking-change`)

## Summary
<!-- 1-3 bullets: what this PR does and why. -->
-

## Changes
| File | Change |
|------|--------|
|      |        |

## Test Plan
- [ ] `npm test` passes
- [ ] Skills lint clean (`npm run lint:skills`)
- [ ] Manually verified the affected behavior

## Checklist
- [ ] Linked an approved issue (`status:approved`)
- [ ] Exactly one `type:*` label
- [ ] Conventional commit messages
- [ ] No `Co-Authored-By` trailers
- [ ] Docs updated if behavior changed
