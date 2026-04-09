---
name: typescript
description: >
  TypeScript strict patterns and best practices.
  Trigger: When writing TypeScript code - types, interfaces, generics.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Const Types Pattern (REQUIRED)

```typescript
// ✅ ALWAYS: Create const object first, then extract type
const STATUS = {
  ACTIVE: "active",
  INACTIVE: "inactive",
  PENDING: "pending",
} as const;

type Status = (typeof STATUS)[keyof typeof STATUS];

// ❌ NEVER: Direct union types
type Status = "active" | "inactive" | "pending";
```

**Why?** Single source of truth, runtime values, autocomplete, easier refactoring.

## Flat Interfaces (REQUIRED)

```typescript
// ✅ ALWAYS: One level depth, nested objects → dedicated interface
interface UserAddress {
  street: string;
  city: string;
}

interface User {
  id: string;
  name: string;
  address: UserAddress;  // Reference, not inline
}

interface Admin extends User {
  permissions: string[];
}

// ❌ NEVER: Inline nested objects
interface User {
  address: { street: string; city: string };  // NO!
}
```

## Never Use `any`

```typescript
// ✅ Use unknown for truly unknown types
function parse(input: unknown): User {
  if (isUser(input)) return input;
  throw new Error("Invalid input");
}

// ✅ Use generics for flexible types
function first<T>(arr: T[]): T | undefined {
  return arr[0];
}

// ❌ NEVER
function parse(input: any): any { }
```

## Utility Types

```typescript
Pick<User, "id" | "name">     // Select fields
Omit<User, "id">              // Exclude fields
Partial<User>                 // All optional
Required<User>                // All required
Readonly<User>                // All readonly
Record<string, User>          // Object type
Extract<Union, "a" | "b">     // Extract from union
Exclude<Union, "a">           // Exclude from union
NonNullable<T | null>         // Remove null/undefined
ReturnType<typeof fn>         // Function return type
Parameters<typeof fn>         // Function params tuple
```

## Type Guards

```typescript
function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    "id" in value &&
    "name" in value
  );
}
```

## Import Types

```typescript
import type { User } from "./types";
import { createUser, type Config } from "./utils";
```

## When to Use

- Writing any TypeScript file where strict type safety is required.
- Defining shared types, interfaces, or enums across modules.
- Replacing `any` with `unknown` + type guards for safer dynamic data.
- Using utility types (`Pick`, `Omit`, `Partial`, etc.) to derive types from existing ones.
- Importing types with `import type` to keep runtime bundles clean.

## Critical Patterns

- ALWAYS use the Const Types Pattern for string literals — never raw union types (see `## Const Types Pattern`).
- ALWAYS keep interfaces flat — extract nested objects into dedicated interfaces (see `## Flat Interfaces`).
- NEVER use `any` — use `unknown` for truly dynamic types, generics for flexible ones (see `## Never Use any`).
- Use `import type` for type-only imports to avoid circular dependency issues.
- Use type guards (`value is T`) to safely narrow `unknown` to concrete types.

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html) — Official guide.
- [TypeScript Utility Types](https://www.typescriptlang.org/docs/handbook/utility-types.html) — Full reference.
- [TypeScript Strict Mode](https://www.typescriptlang.org/tsconfig#strict) — Compiler options.

## Keywords
typescript, ts, types, interfaces, generics, strict mode, utility types
