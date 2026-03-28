---
name: react-19
description: >
  React 19 patterns with React Compiler.
  Trigger: When writing React components - no useMemo/useCallback needed.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- Writing React components with React 19 and the React Compiler active.
- Replacing legacy `useMemo`/`useCallback` patterns — the compiler handles it.
- Using `use()` to read promises or context conditionally.
- Building forms with `useActionState` for pending and error states.
- Passing refs as plain props — no `forwardRef` needed in React 19.

## Critical Patterns

- NEVER add `useMemo` or `useCallback` — the React Compiler optimizes automatically (see `## No Manual Memoization`).
- ALWAYS use named imports from `"react"` — no default `React` import (see `## Imports`).
- Server Components are the default; add `"use client"` only for interactivity (see `## Server Components First`).
- Use `use()` for promise reading inside Suspense and for conditional context access.
- Use `useActionState` to track form submission state and pending indicators.

## No Manual Memoization (REQUIRED)

```typescript
// ✅ React Compiler handles optimization automatically
function Component({ items }) {
  const filtered = items.filter(x => x.active);
  const sorted = filtered.sort((a, b) => a.name.localeCompare(b.name));

  const handleClick = (id) => {
    console.log(id);
  };

  return <List items={sorted} onClick={handleClick} />;
}

// ❌ NEVER: Manual memoization
const filtered = useMemo(() => items.filter(x => x.active), [items]);
const handleClick = useCallback((id) => console.log(id), []);
```

## Imports (REQUIRED)

```typescript
// ✅ ALWAYS: Named imports
import { useState, useEffect, useRef } from "react";

// ❌ NEVER
import React from "react";
import * as React from "react";
```

## Server Components First

```typescript
// ✅ Server Component (default) - no directive
export default async function Page() {
  const data = await fetchData();
  return <ClientComponent data={data} />;
}

// ✅ Client Component - only when needed
"use client";
export function Interactive() {
  const [state, setState] = useState(false);
  return <button onClick={() => setState(!state)}>Toggle</button>;
}
```

## When to use "use client"

- useState, useEffect, useRef, useContext
- Event handlers (onClick, onChange)
- Browser APIs (window, localStorage)

## use() Hook

```typescript
import { use } from "react";

// Read promises (suspends until resolved)
function Comments({ promise }) {
  const comments = use(promise);
  return comments.map(c => <div key={c.id}>{c.text}</div>);
}

// Conditional context (not possible with useContext!)
function Theme({ showTheme }) {
  if (showTheme) {
    const theme = use(ThemeContext);
    return <div style={{ color: theme.primary }}>Themed</div>;
  }
  return <div>Plain</div>;
}
```

## Actions & useActionState

```typescript
"use server";
async function submitForm(formData: FormData) {
  await saveToDatabase(formData);
  revalidatePath("/");
}

// With pending state
import { useActionState } from "react";

function Form() {
  const [state, action, isPending] = useActionState(submitForm, null);
  return (
    <form action={action}>
      <button disabled={isPending}>
        {isPending ? "Saving..." : "Save"}
      </button>
    </form>
  );
}
```

## ref as Prop (No forwardRef)

```typescript
// ✅ React 19: ref is just a prop
function Input({ ref, ...props }) {
  return <input ref={ref} {...props} />;
}

// ❌ Old way (unnecessary now)
const Input = forwardRef((props, ref) => <input ref={ref} {...props} />);
```

## Resources

- [React 19 Release Notes](https://react.dev/blog/2024/12/05/react-19) — What's new.
- [React Compiler Docs](https://react.dev/learn/react-compiler) — Automatic optimization.
- [use() Hook Reference](https://react.dev/reference/react/use) — Reading promises and context.

## Keywords
react, react 19, compiler, useMemo, useCallback, server components, use hook
