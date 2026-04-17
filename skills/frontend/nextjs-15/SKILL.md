---
name: nextjs-15
description: >
  ZeshOne Next.js 15 conventions.
  Trigger: When working with Next.js — routing, Server Actions, data fetching, middleware, project structure, performance.
license: Apache-2.0
allowed-tools:
  - Read
  - Edit
  - Write
  - Glob
  - Grep
metadata:
  author: Zesh-One
  version: "2.0"
  inspired-by: gentleman-programming/nextjs-15, frontend-nextjs
---

## When to Use

Load this skill when working with Next.js 15 — routing, Server Actions, data fetching, middleware, project structure, or performance.

## Critical Patterns

- App Router exclusively — no `pages/` directory.
- Default to Server Components. Add `"use client"` only for state, events, or browser APIs.
- ALL external API calls go through Server Actions — client never knows the endpoint URL (endpoints hidden from DevTools).
- Vertical Slices: each feature under `features/<feature>/`. Route groups keep URLs clean.
- `params` and `searchParams` are Promises in Next.js 15 — always `await` them.
- `cookies()` and `headers()` are async in Next.js 15 — always `await` them.
- Never call `redirect()` inside `try/catch` — it throws internally and catch will swallow it.
- Use `loading.tsx` for Suspense boundaries on slow data fetches. Use `error.tsx` for error boundaries.
- Push `"use client"` boundary as deep as possible in the component tree.

## Project Structure — Vertical Slices

```
project-root/
├── app/
│   ├── (auth)/             # Route group: login, register — minimal layout
│   │   ├── login/page.tsx
│   │   └── layout.tsx
│   ├── (dashboard)/        # Route group: authed pages — sidebar layout
│   │   ├── dashboard/page.tsx
│   │   ├── settings/page.tsx
│   │   └── layout.tsx
│   ├── (marketing)/        # Route group: public pages
│   │   ├── page.tsx
│   │   └── layout.tsx
│   ├── api/                # Route Handlers — webhooks ONLY
│   ├── layout.tsx          # Root layout (html, body, providers)
│   └── globals.css
├── features/               # Feature-based modules (self-contained)
│   └── <feature>/
│       ├── components/     # Feature UI components
│       ├── hooks/          # Feature-specific hooks
│       ├── actions/        # Server Actions ('use server')
│       ├── queries/        # TanStack Query keys + options
│       ├── schemas/        # Zod validation schemas
│       └── types/          # TypeScript types
├── components/
│   ├── ui/                 # shadcn/ui primitives (never edit directly)
│   ├── layout/             # header, footer, sidebar, mobile-nav
│   ├── forms/              # Reusable form elements
│   └── shared/             # Cross-feature shared components
├── lib/
│   ├── api/                # Fetch wrapper + endpoint constants
│   ├── query-client.ts     # TanStack Query client setup
│   ├── utils.ts            # cn(), formatDate, etc.
│   └── constants.ts
├── hooks/                  # Shared custom hooks
├── providers/              # QueryProvider, ThemeProvider, composed Providers
├── types/                  # Global TS types, env.d.ts
└── middleware.ts            # Auth guards + i18n (or proxy.ts in v16)
```

**Naming conventions:**

| Type | Convention |
|------|-----------|
| Files | kebab-case (`post-card.tsx`) |
| Components | PascalCase export |
| Hooks | `use-` prefix (`use-posts.ts`) |
| Actions | verb-first (`create-post.ts`) |
| Schemas | `.schema.ts` suffix |
| Types | `.types.ts` suffix |

**Barrel exports:** Only at feature boundary (`features/posts/index.ts`). Never inside `components/ui/` — direct imports only (tree-shaking).

## Server Actions — API Proxy Pattern (REQUIRED)

ALL calls to external .NET APIs go through a Server Action. Client never receives endpoint URL.
Pattern: Client request → Server Action → fetch to .NET API → return to client.

```typescript
// features/users/actions/create-user.ts
"use server";
import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { createUserSchema } from "../schemas/user.schema";

export async function createUser(formData: FormData) {
  const result = createUserSchema.safeParse({
    name: formData.get("name"),
    email: formData.get("email"),
  });

  if (!result.success) return { error: result.error.flatten().fieldErrors };

  await fetch(`${process.env.API_URL}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(result.data),
  });

  revalidatePath("/users");
  redirect("/users"); // Outside try/catch — throws internally
}
```

**Server Action rules:**
- Always validate with Zod before fetching
- Always authenticate before mutating
- `revalidatePath` / `revalidateTag` after mutations that affect cached pages
- Return `{ error }` on validation failure, never throw client-visible errors

## Next.js 15 Breaking Changes

```typescript
// params + searchParams — now Promises
export default async function Page({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ page?: string }>;
}) {
  const { id } = await params;
  const { page } = await searchParams;
}

// cookies + headers — async
import { cookies, headers } from "next/headers";
const cookieStore = await cookies();
const token = cookieStore.get("token");
```

## Data Fetching Patterns

```typescript
// Parallel fetch — always for independent data
const [users, settings] = await Promise.all([fetchUsers(), fetchSettings()]);

// Server prefetch + TanStack Query hydration (see tanstack-query skill)
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { getQueryClient } from "@/lib/query-client";

export default async function PostsPage() {
  const queryClient = getQueryClient();
  await queryClient.prefetchQuery(postListOptions({ status: "published" }));
  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <PostList />
    </HydrationBoundary>
  );
}
```

## Provider Composition

```tsx
// providers/providers.tsx
"use client";
import { QueryProvider } from "./query-provider";
import { ThemeProvider } from "./theme-provider";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <QueryProvider>
      <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
        {children}
      </ThemeProvider>
    </QueryProvider>
  );
}
```

## Performance Rules (Vercel Patterns)

- **Eliminate waterfalls:** `Promise.all()` for independent fetches. Suspense to stream slow parts.
- **Bundle:** Import directly — avoid barrel files in `components/ui/`. Use `next/dynamic` for heavy components.
- **Server cache:** `React.cache()` for per-request dedup. LRU cache for cross-request.
- **Client re-renders:** Use `startTransition` for non-urgent updates. Don't subscribe to state only used in callbacks.
- **Conditional render:** Use ternary (`a ? b : c`), not `&&` — avoids rendering `0` as text.
- **Static JSX:** Extract static JSX outside component functions to prevent re-creation.

```typescript
// Dynamic import for heavy components
import dynamic from "next/dynamic";
const HeavyChart = dynamic(() => import("@/components/heavy-chart"), { ssr: false });

// React.cache — deduplicates within a single request
import { cache } from "react";
export const getUser = cache(async (id: string) => {
  return fetch(`${process.env.API_URL}/users/${id}`).then(r => r.json());
});
```

## Route Handlers — Webhooks Only

`route.ts` files only for webhooks or third-party callbacks that cannot use Server Actions.

## Environment Variables

```bash
# Server-only — never prefix with NEXT_PUBLIC_
API_URL=https://api.internal.com
AUTH_SECRET=
# Client-safe — intentionally public
NEXT_PUBLIC_APP_URL=https://app.com
```

Type-safe env via `types/env.d.ts`:
```typescript
declare namespace NodeJS {
  interface ProcessEnv {
    API_URL: string;
    AUTH_SECRET: string;
    NEXT_PUBLIC_APP_URL: string;
  }
}
```

## Resources

- [Next.js App Router Documentation](https://nextjs.org/docs/app)
- [Next.js Server Actions](https://nextjs.org/docs/app/building-your-application/data-fetching/server-actions-and-mutations)
- [Next.js 15 Release Notes](https://nextjs.org/blog/next-15)

## Changelog

### v2.0 — 2026-04-16
- Expanded the skill with Next.js 15 async request APIs, Server Action proxy guidance, performance rules, and vertical-slice project structure.
