---
name: typescript
description: >
  ZeshOne TypeScript conventions.
  Trigger: When writing TypeScript types or interfaces.
license: Apache-2.0
metadata:
  author: zesh-one
  version: "1.0"
  inspired-by: gentleman-programming/typescript
---

## Const Types Pattern (REQUIRED)

```typescript
// ✅ ALWAYS: const object → extract type
const STATUS = {
  ACTIVE: "active",
  INACTIVE: "inactive",
  PENDING: "pending",
} as const;

type Status = (typeof STATUS)[keyof typeof STATUS];

// ❌ NEVER: raw union string literals
type Status = "active" | "inactive" | "pending";
```

## Flat Interfaces (REQUIRED)

```typescript
// ✅ ALWAYS: one level max — nested objects get their own interface
interface UserAddress {
  street: string;
  city: string;
}

interface User {
  id: string;
  name: string;
  address: UserAddress;
}

// ❌ NEVER: inline nested objects
interface User {
  address: { street: string; city: string };
}
```

## Never Use `any`

- Use `unknown` for truly dynamic input — narrow with a type guard before use.
- Use generics for flexible but typed functions.
- `any` is always wrong; there is always a better option.

## Critical Patterns

- ALWAYS use the Const Types Pattern for string literals — never raw unions.
- ALWAYS keep interfaces flat — extract nested shapes into their own interfaces.
- NEVER use `any`.
- Use `import type` for type-only imports.
- Use `z.infer<typeof schema>` from Zod schemas as the source of truth for form/API types — do not duplicate.

## Keywords
typescript, types, interfaces, const types, strict
