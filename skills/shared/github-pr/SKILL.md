---
name: github-pr
description: >
  ZeshOne Pull Request conventions.
  Trigger: When creating or reviewing a Pull Request.
license: Apache-2.0
metadata:
  author: zesh-one
  version: "1.0"
  inspired-by: gentleman-programming/github-pr
---

## When to Use

Load this skill when creating or reviewing a Pull Request — title, branch naming, description structure, or `gh` CLI usage.

## Critical Patterns

- PR title follows Conventional Commits format: `<type>(<scope>): <description>`.
- Branch naming matches the commit type: `feat/`, `fix/`, `refactor/`, `chore/`, `docs/`.
- Base branch is `alfa` — NOT `main`. `alfa` is the concentration branch before promoting to `main`.
- Merge strategy depends on context — do not assume squash always.
- PR description always includes Summary, Changes, and a closing issue reference.
- One logical concern per PR — if it touches more than one feature, split it.

## Branch Naming

```
feat/<short-description>     New feature
fix/<short-description>      Bug fix
refactor/<short-description> Code restructuring without behavior change
chore/<short-description>    Tooling, deps, config
docs/<short-description>     Documentation only
```

## PR Title — Conventional Commits

```
feat(users): add profile picture upload
fix(auth): prevent session expiry on active use
refactor(dashboard): extract chart to standalone component
chore(deps): upgrade zod to 4.0 stable
```

## PR Description Structure

```markdown
## Summary
- What was done and why (1–3 bullets)

## Changes
- List of concrete changes

## Testing
- [ ] Tests added or updated
- [ ] Manual testing done

Closes #<issue-number>
```

## PR Creation — HEREDOC (Preferred for Complex Descriptions)

```bash
gh pr create \
  --title "feat(users): add profile upload" \
  --base alfa \
  --body "$(cat <<'EOF'
## Summary
- Add profile picture upload to user settings

## Changes
- Added UploadAvatar server action
- Created AvatarCropper client component

## Testing
- [x] Upload tested with jpg, png, webp
- [x] Error state for oversized files

Closes #87
EOF
)"
```

## Keywords
github, pull request, conventional commits, branch, alfa, gh cli
