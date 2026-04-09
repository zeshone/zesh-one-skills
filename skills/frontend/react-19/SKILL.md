---
name: react-19
description: >
  ZeshOne React 19 conventions.
  Trigger: When writing React components.
license: Apache-2.0
metadata:
  author: zesh-one
  version: "1.0"
  inspired-by: gentleman-programming/react-19
---

## Critical Patterns

- NEVER add `useMemo` or `useCallback` — the React Compiler handles memoization automatically.
- ALWAYS use named imports from `"react"` — never `import React from "react"`.
- Server Components are the default — add `"use client"` only when the component needs state, events, or browser APIs.
- Use `useActionState` (NOT `useFormState` — removed in React 19) for form submission state and pending indicators.
- `ref` is a plain prop in React 19 — do NOT use `forwardRef`.

## React Compiler — Prerequisite (Non-Obvious)

The compiler is NOT active by default. Without it, removing `useMemo`/`useCallback` causes re-render regressions.

- **Next.js 15**: add `experimental: { reactCompiler: true }` to `next.config.ts` — compiler is bundled.
- **Vite / other bundlers**: install `babel-plugin-react-compiler` and configure Babel explicitly.

Verify the compiler is active before removing any manual memoization from existing code.

## use() Hook — New React 19 API

`use()` does two things `useContext` cannot:

1. Reads a Promise inside a Suspense boundary (suspends until resolved).
2. Called conditionally — unlike `useContext`, it is not subject to the Rules of Hooks position constraint.

```typescript
import { use } from "react";

// Reads promise — suspends until resolved
function Comments({ promise }: { promise: Promise<Comment[]> }) {
  const comments = use(promise);
  return comments.map(c => <div key={c.id}>{c.text}</div>);
}

// Conditional context — not possible with useContext
function Theme({ showTheme }: { showTheme: boolean }) {
  if (showTheme) {
    const theme = use(ThemeContext);
    return <div style={{ color: theme.primary }}>Themed</div>;
  }
  return <div>Plain</div>;
}
```

## useActionState — Replaces useFormState (Breaking Change)

`useFormState` was removed in React 19. The replacement is `useActionState` from `"react"` (not from `"react-dom"`).

```typescript
import { useActionState } from "react";

function Form() {
  const [state, action, isPending] = useActionState(submitAction, null);
  return (
    <form action={action}>
      <button disabled={isPending}>
        {isPending ? "Saving..." : "Save"}
      </button>
    </form>
  );
}
```

## ref as Plain Prop — forwardRef Removed (Breaking Change)

```typescript
// ✅ React 19
function Input({ ref, ...props }: React.ComponentProps<"input">) {
  return <input ref={ref} {...props} />;
}

// ❌ forwardRef no longer needed — do not use it
const Input = forwardRef((props, ref) => <input ref={ref} {...props} />);
```

## Keywords
react, react 19, compiler, use hook, useActionState, forwardRef, server components
