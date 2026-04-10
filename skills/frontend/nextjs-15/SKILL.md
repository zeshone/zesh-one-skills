---
name: nextjs-15
description: >
  ZeshOne Next.js 15 conventions.
  Trigger: When working with Next.js — routing, Server Actions, data fetching, middleware.
license: Apache-2.0
metadata:
  author: zesh-one
  version: "1.0"
  inspired-by: gentleman-programming/nextjs-15
---

## When to Use

Load this skill when working with Next.js 15 — routing, Server Actions, data fetching, middleware, or project structure.

## Critical Patterns

- Use App Router exclusively — no `pages/` directory.
- Default to Server Components. Add `"use client"` only when the component needs interactivity.
- ALL external API calls go through Server Actions — never call an external API directly from the client. This hides endpoints from browser DevTools.
- Use Vertical Slices: each feature lives under `app/(features)/<feature>/`. Route groups keep URLs clean.
- `params` and `searchParams` are Promises in Next.js 15 — always `await` them.
- `cookies()` and `headers()` are async in Next.js 15 — always `await` them.
- Never call `redirect()` inside a `try/catch` — it throws internally and the catch will swallow it.

## Project Structure — Vertical Slices

```
app/
├── layout.tsx
├── page.tsx
├── (features)/
│   ├── users/
│   │   ├── page.tsx
│   │   ├── loading.tsx
│   │   ├── error.tsx
│   │   ├── actions.ts       # Server Actions for this feature
│   │   └── _components/     # Private components (not routed)
│   └── dashboard/
│       ├── page.tsx
│       └── actions.ts
└── api/                     # Route Handlers — webhooks only
```

## Server Actions — API Proxy Pattern (REQUIRED)

ALL calls to external APIs go through a Server Action. The client never knows the endpoint URL.

```typescript
// app/(features)/users/actions.ts
"use server";
import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

export async function createUser(formData: FormData) {
  const payload = {
    name: formData.get("name") as string,
    email: formData.get("email") as string,
  };

  await fetch(`${process.env.API_URL}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  revalidatePath("/users");
  redirect("/users"); // Outside try/catch — it throws internally
}
```

## Next.js 15 Breaking Changes — params and searchParams

`params` and `searchParams` are now Promises. Forgetting `await` produces stale or undefined values with no build error.

```typescript
// ✅ Next.js 15
export default async function Page({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ page?: string }>;
}) {
  const { id } = await params;
  const { page } = await searchParams;
  // ...
}

// ❌ Next.js 14 and earlier — breaks silently in v15
export default async function Page({ params }: { params: { id: string } }) {
  const { id } = params; // undefined in v15
}
```

## Next.js 15 Breaking Changes — cookies and headers

```typescript
import { cookies, headers } from "next/headers";

// ✅ Next.js 15 — async
const cookieStore = await cookies();
const token = cookieStore.get("token");

const headersList = await headers();
const auth = headersList.get("authorization");
```

## Data Fetching — Parallel by Default

```typescript
// ✅ Always fetch independent data in parallel
const [users, settings] = await Promise.all([
  fetchUsers(),
  fetchSettings(),
]);
```

## Route Handlers — Webhooks Only

Route Handlers (`route.ts`) are only for webhooks or third-party callbacks that cannot use Server Actions. All other mutations use Server Actions.

## Keywords
nextjs, next.js 15, app router, server actions, vertical slices, params, cookies, async
