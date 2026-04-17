---
name: zod-4
description: >
  ZeshOne Zod 4 conventions.
  Trigger: When creating or updating Zod schemas.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.0"
  inspired-by: gentleman-programming/zod-4
---

## When to Use

Load this skill when creating or updating Zod schemas — validation, type inference, or React Hook Form integration.

## Critical Patterns

- Use top-level validators: `z.email()`, `z.uuid()`, `z.url()` — NOT `z.string().email()` (Zod 4 breaking change).
- Use `error:` instead of `message:` for custom error strings (Zod 4 breaking change).
- Prefer `safeParse` over `parse` in all user-facing flows — never let validation throw to the UI.
- Derive TypeScript types from schemas with `z.infer<typeof schema>` — never duplicate type definitions.
- Use `superRefine` for cross-field validations (e.g., password confirmation).

## Zod 4 Breaking Changes — v3 vs v4

These are changes the agent may apply incorrectly if trained on v3:

```typescript
// ❌ Zod 3 (OLD — do not use)
z.string().email()
z.string().uuid()
z.string().url()
z.string().nonempty()
z.object({ name: z.string() }, { required_error: "Required" })

// ✅ Zod 4 (CORRECT)
z.email()
z.uuid()
z.url()
z.string().min(1)
z.object({ name: z.string() }, { error: "Required" })
```

```typescript
// ❌ Zod 3 — message param
z.string().min(1, { message: "Required" })

// ✅ Zod 4 — error param
z.string().min(1, { error: "Required" })
```

## Installation Note

As of 2026, `npm install zod` still installs Zod 3 (stable default).
Zod 4 requires an explicit install:

```bash
npm install zod@beta   # Verify this maps to v4: npm show zod dist-tags
```

Verify the beta tag maps to v4 before installing: `npm show zod dist-tags`. If Zod 4 is stable, use `npm install zod@^4` instead.

Do NOT assume Zod 4 is active unless `package.json` explicitly targets the beta tag or v4 stable has been released.

## Error Handling

```typescript
// Custom error string
z.email({ error: "Invalid email format" })

// Custom error function (full control)
z.string({
  error: (issue) => {
    if (issue.code === "too_small") return "Too short";
    return "Invalid value";
  },
});
```

## React Hook Form Integration

```typescript
import { zodResolver } from "@hookform/resolvers/zod";

const schema = z.object({
  email: z.email(),
  password: z.string().min(8),
});

type FormData = z.infer<typeof schema>;

const { register, handleSubmit } = useForm<FormData>({
  resolver: zodResolver(schema),
});
```

## Resources

- Zod Docs — https://zod.dev/
- Zod API Reference — https://zod.dev/api
- Zod Ecosystem — https://zod.dev/ecosystem
- React Hook Form Resolvers — https://github.com/react-hook-form/resolvers
