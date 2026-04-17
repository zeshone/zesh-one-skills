---
name: security
description: >
  ZeshOne Next.js security conventions.
  Trigger: When configuring security headers, CSP, auth guards, input validation, rate limiting, or cookie security.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.0"
  inspired-by: frontend-nextjs/security.md
---

## When to Use

Load this skill when implementing security: headers, CSP nonce, auth guards, input validation, cookie config, rate limiting, or CORS.

## Critical Patterns

- Set baseline security headers in `next.config.ts` for every route.
- Generate a per-request CSP nonce in `middleware.ts`, forward it through `x-nonce`, and apply it to inline or external scripts that need it.
- Enforce auth and authorization in middleware before rendering protected routes.
- Validate all mutation inputs with Zod and prefer `safeParse` in Server Actions and API handlers.
- Let React escape content by default; only use `dangerouslySetInnerHTML` after DOMPurify sanitization and URL validation.
- Keep sensitive cookies `httpOnly`, `secure`, and `sameSite`, and never expose secrets through `NEXT_PUBLIC_` variables.
- Add rate limiting to auth and other sensitive endpoints.

## Security Headers — next.config.ts

```typescript
// next.config.ts
const nextConfig: NextConfig = {
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "X-Frame-Options", value: "SAMEORIGIN" },
          { key: "X-XSS-Protection", value: "1; mode=block" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          { key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains; preload" },
          {
            key: "Permissions-Policy",
            value: "accelerometer=(), camera=(), geolocation=(self), microphone=(), payment=()",
          },
        ],
      },
    ];
  },
};
```

## CSP with Nonce — middleware.ts (or proxy.ts in v16)

Generates a unique nonce per request. Inline scripts and styles must reference it.

```typescript
// middleware.ts — single example with CSP nonce + auth guards
import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/auth";

export async function middleware(request: NextRequest) {
  const nonce = Buffer.from(crypto.randomUUID()).toString("base64");
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set("x-nonce", nonce);
  const pathname = request.nextUrl.pathname;
  const session = await auth();

  let response: NextResponse;

  if (pathname.startsWith("/dashboard") && !session?.user) {
    response = NextResponse.redirect(new URL("/login", request.url));
  } else if (pathname.startsWith("/admin") && session?.user?.role !== "admin") {
    response = NextResponse.redirect(new URL("/unauthorized", request.url));
  } else {
    response = NextResponse.next({ request: { headers: requestHeaders } });
  }

  const cspHeader = [
    `default-src 'self'`,
    `script-src 'self' 'nonce-${nonce}' 'strict-dynamic'`,
    `style-src 'self' 'nonce-${nonce}'`,
    `img-src 'self' data: https:`,
    `font-src 'self'`,
    `connect-src 'self'`,
    `frame-ancestors 'none'`,
    `base-uri 'self'`,
    `form-action 'self'`,
    `upgrade-insecure-requests`,
  ].join("; ");

  response.headers.set("x-nonce", nonce);
  response.headers.set("Content-Security-Policy", cspHeader);
  return response;
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico|public).*)"],
};
```

```tsx
// app/layout.tsx — use nonce for external scripts
import { headers } from "next/headers";
import Script from "next/script";

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const nonce = (await headers()).get("x-nonce") ?? undefined;
  return (
    <html lang="en">
      <head>
        <Script src="https://www.googletagmanager.com/gtag/js" strategy="afterInteractive" nonce={nonce} />
      </head>
      <body>{children}</body>
    </html>
  );
}
```

## Auth Guards in Middleware

**Rule:** Protect routes in `middleware.ts`, NOT only in Server Components. Middleware runs before rendering.

## Input Validation — Zod (Server + Client)

```typescript
// features/posts/schemas/post.schema.ts
// ⚠️ Requires Zod 4 — see skills/frontend/zod-4/SKILL.md for breaking changes (z.email() not z.string().email())
import { z } from "zod";

export const createPostSchema = z.object({
  title: z.string().min(3).max(200).trim(),
  content: z.string().min(10).max(50000).trim(),
  tags: z.array(z.string().max(50)).max(10).optional(),
  status: z.enum(["draft", "published"]).default("draft"),
});

export type CreatePostInput = z.infer<typeof createPostSchema>;
```

```typescript
// Server Action validation pattern
"use server";
export async function createPost(formData: FormData) {
  const session = await auth();
  if (!session?.user?.id) throw new Error("Unauthorized");

  const result = createPostSchema.safeParse({
    title: formData.get("title"),
    content: formData.get("content"),
  });
  if (!result.success) return { error: result.error.flatten().fieldErrors };

  await fetch(`${process.env.API_URL}/posts`, {
    method: "POST",
    headers: { "Content-Type": "application/json", Authorization: `Bearer ${session.token}` },
    body: JSON.stringify(result.data),
  });

  revalidatePath("/dashboard/posts");
  return { success: true };
}
```

## XSS Prevention

```tsx
// React auto-escapes — always safe:
<div>{userContent}</div>

// NEVER without sanitization:
// <div dangerouslySetInnerHTML={{ __html: userContent }} />

// If HTML rendering is unavoidable:
import DOMPurify from "isomorphic-dompurify";
function SafeHTML({ html }: { html: string }) {
  const clean = DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ["b", "i", "em", "strong", "a", "p", "br", "ul", "ol", "li"],
    ALLOWED_ATTR: ["href", "target", "rel"],
  });
  return <div dangerouslySetInnerHTML={{ __html: clean }} />;
}

// Validate URLs before rendering
function safeUrl(url: string): string {
  try {
    const parsed = new URL(url);
    return ["http:", "https:", "mailto:"].includes(parsed.protocol) ? url : "";
  } catch { return ""; }
}
```

## Cookie Security

```typescript
// All sensitive cookies must set:
{
  httpOnly: true,                                         // no JS access
  secure: process.env.NODE_ENV === "production",          // HTTPS only
  sameSite: "lax",                                        // CSRF protection
  path: "/",
}
```

## Environment Variable Rules

```bash
# Server-only — NEVER prefix with NEXT_PUBLIC_
API_URL=
AUTH_SECRET=
STRIPE_SECRET_KEY=

# Client-safe — intentionally public
NEXT_PUBLIC_APP_URL=
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=
```

- Server vars: accessible in Server Components, Server Actions, API routes, middleware
- `NEXT_PUBLIC_` vars: bundled into client JS, visible to everyone
- Never put API keys, DB URLs, or auth secrets in `NEXT_PUBLIC_`
- Commit `.env.example` with empty values; `.env.local` is git-ignored

## Server Component Security

```tsx
interface AnalyticsRow {
  date: string;
  count: number;
}

export default async function AdminDashboard() {
  const session = await auth();
  if (session?.user?.role !== "admin") redirect("/unauthorized");

  const sensitiveData: AnalyticsRow[] = await fetch(`${process.env.API_URL}/analytics`).then(r => r.json());

  return (
    // Only serialize what the client needs — never pass raw sensitive fields
    <AnalyticsChart data={sensitiveData.map((d) => ({ date: d.date, count: d.count }))} />
  );
}
```

## Rate Limiting (Upstash Redis)

```typescript
// lib/rate-limit.ts
import { Ratelimit } from "@upstash/ratelimit";
import { Redis } from "@upstash/redis";

const redis = new Redis({ url: process.env.UPSTASH_REDIS_REST_URL!, token: process.env.UPSTASH_REDIS_REST_TOKEN! });

export const authRateLimit = new Ratelimit({ redis, limiter: Ratelimit.slidingWindow(5, "15 m"), prefix: "rl:auth" });
export const apiRateLimit = new Ratelimit({ redis, limiter: Ratelimit.slidingWindow(100, "1 h"), prefix: "rl:api" });

// Usage in API route:
const ip = request.headers.get("x-forwarded-for") ?? "anonymous";
const { success } = await authRateLimit.limit(ip);
if (!success) return Response.json({ error: "Too many attempts" }, { status: 429 });
```

## Security Checklist

- [ ] CSP nonce in middleware
- [ ] Security headers in next.config.ts (HSTS, X-Frame-Options, etc.)
- [ ] Zod validation on all Server Actions and API routes
- [ ] Auth check before every mutation
- [ ] Route guards in middleware (not only in components)
- [ ] No secrets in `NEXT_PUBLIC_` variables
- [ ] `httpOnly + secure + sameSite` on all sensitive cookies
- [ ] No `dangerouslySetInnerHTML` without DOMPurify
- [ ] URL validation before rendering links
- [ ] Rate limiting on auth and sensitive endpoints

## Resources

- Next.js Middleware — https://nextjs.org/docs/app/building-your-application/routing/middleware
- Next.js Content Security Policy Guide — https://nextjs.org/docs/app/guides/content-security-policy
- OWASP Cheat Sheet Series — https://cheatsheetseries.owasp.org/
- OWASP Cross Site Scripting Prevention Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html
- OWASP Session Management Cheat Sheet — https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
