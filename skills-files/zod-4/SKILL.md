---
name: zod-4
description: >
  Zod 4 schema validation patterns.
  Trigger: When using Zod for validation - breaking changes from v3.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Breaking Changes from Zod 3

```typescript
// ❌ Zod 3 (OLD)
z.string().email()
z.string().uuid()
z.string().url()
z.string().nonempty()
z.object({ name: z.string() }).required_error("Required")

// ✅ Zod 4 (NEW)
z.email()
z.uuid()
z.url()
z.string().min(1)
z.object({ name: z.string() }, { error: "Required" })
```

## Basic Schemas

```typescript
import { z } from "zod";

// Primitives
const stringSchema = z.string();
const numberSchema = z.number();
const booleanSchema = z.boolean();
const dateSchema = z.date();

// Top-level validators (Zod 4)
const emailSchema = z.email();
const uuidSchema = z.uuid();
const urlSchema = z.url();

// With constraints
const nameSchema = z.string().min(1).max(100);
const ageSchema = z.number().int().positive().max(150);
const priceSchema = z.number().min(0).multipleOf(0.01);
```

## Object Schemas

```typescript
const userSchema = z.object({
  id: z.uuid(),
  email: z.email({ error: "Invalid email address" }),
  name: z.string().min(1, { error: "Name is required" }),
  age: z.number().int().positive().optional(),
  role: z.enum(["admin", "user", "guest"]),
  metadata: z.record(z.string(), z.unknown()).optional(),
});

type User = z.infer<typeof userSchema>;

// Parsing
const user = userSchema.parse(data);  // Throws on error
const result = userSchema.safeParse(data);  // Returns { success, data/error }

if (result.success) {
  console.log(result.data);
} else {
  console.log(result.error.issues);
}
```

## Arrays and Records

```typescript
// Arrays
const tagsSchema = z.array(z.string()).min(1).max(10);
const numbersSchema = z.array(z.number()).nonempty();

// Records (objects with dynamic keys)
const scoresSchema = z.record(z.string(), z.number());
// { [key: string]: number }

// Tuples
const coordinatesSchema = z.tuple([z.number(), z.number()]);
// [number, number]
```

## Unions and Discriminated Unions

```typescript
// Simple union
const stringOrNumber = z.union([z.string(), z.number()]);

// Discriminated union (more efficient)
const resultSchema = z.discriminatedUnion("status", [
  z.object({ status: z.literal("success"), data: z.unknown() }),
  z.object({ status: z.literal("error"), error: z.string() }),
]);
```

## Transformations

```typescript
// Transform during parsing
const lowercaseEmail = z.email().transform(email => email.toLowerCase());

// Coercion (convert types)
const numberFromString = z.coerce.number();  // "42" → 42
const dateFromString = z.coerce.date();      // "2024-01-01" → Date

// Preprocessing
const trimmedString = z.preprocess(
  val => typeof val === "string" ? val.trim() : val,
  z.string()
);
```

## Refinements

```typescript
const passwordSchema = z.string()
  .min(8)
  .refine(val => /[A-Z]/.test(val), {
    message: "Must contain uppercase letter",
  })
  .refine(val => /[0-9]/.test(val), {
    message: "Must contain number",
  });

// With superRefine for multiple errors
const formSchema = z.object({
  password: z.string(),
  confirmPassword: z.string(),
}).superRefine((data, ctx) => {
  if (data.password !== data.confirmPassword) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: "Passwords don't match",
      path: ["confirmPassword"],
    });
  }
});
```

## Optional and Nullable

```typescript
// Optional (T | undefined)
z.string().optional()

// Nullable (T | null)
z.string().nullable()

// Both (T | null | undefined)
z.string().nullish()

// Default values
z.string().default("unknown")
z.number().default(() => Math.random())
```

## Error Handling

```typescript
// Zod 4: Use 'error' param instead of 'message'
const schema = z.object({
  name: z.string({ error: "Name must be a string" }),
  email: z.email({ error: "Invalid email format" }),
  age: z.number().min(18, { error: "Must be 18 or older" }),
});

// Custom error map
const customSchema = z.string({
  error: (issue) => {
    if (issue.code === "too_small") {
      return "String is too short";
    }
    return "Invalid string";
  },
});
```

## React Hook Form Integration

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

const schema = z.object({
  email: z.email(),
  password: z.string().min(8),
});

type FormData = z.infer<typeof schema>;

function Form() {
  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input {...register("email")} />
      {errors.email && <span>{errors.email.message}</span>}
    </form>
  );
}
```

## When to Use

- Validating form data, API responses, or any external input at runtime.
- Migrating from Zod 3 — top-level validators (`z.email()`, `z.uuid()`) replace `z.string().*`.
- Inferring TypeScript types from schemas with `z.infer<typeof schema>`.
- Integrating with React Hook Form via `zodResolver`.
- Parsing untrusted data safely with `safeParse` to avoid thrown exceptions.

## Critical Patterns

- Use top-level validators in Zod 4: `z.email()`, `z.uuid()`, `z.url()` — NOT `z.string().email()` (see `## Breaking Changes from Zod 3`).
- Use `error` param instead of `message` for custom error strings in Zod 4.
- Prefer `safeParse` over `parse` in user-facing flows to handle errors gracefully.
- Use `z.infer<typeof schema>` to derive types — avoids duplicating type definitions.
- Use `superRefine` when a single object needs multiple cross-field validations.

## Resources

- [Zod 4 Docs](https://zod.dev) — Official documentation with v4 changes.
- [Zod GitHub](https://github.com/colinhacks/zod) — Changelog and migration guide.
- [React Hook Form + Zod](https://react-hook-form.com/get-started#SchemaValidation) — Integration guide.

## Keywords
zod, validation, schema, typescript, forms, parsing
