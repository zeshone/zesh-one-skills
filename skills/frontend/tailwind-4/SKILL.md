---
name: tailwind-4
description: >
  ZeshOne Tailwind CSS 4 conventions.
  Trigger: When styling with Tailwind — cn(), theme variables, dynamic values, fluid typography, container queries.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "2.0"
  inspired-by: gentleman-programming/tailwind-4, frontend-nextjs
---

## When to Use

Load this skill when styling with Tailwind CSS 4 — `cn()`, theme variables, fluid typography, container queries, or third-party library integration.

## Critical Patterns

- NEVER use `var()` inside `className` — use Tailwind semantic classes (`bg-primary`, `text-slate-400`).
- NEVER use hex colors in `className` — always use Tailwind color classes.
- NEVER use `px` for font sizes or spacing — use `rem` or `clamp()`.
- `px` is only valid for borders, box-shadows, and outlines.
- Use `cn()` for conditional/merged classes. Plain `className` for static-only.
- `style` prop only for truly dynamic runtime values (percentages, progress bars, user-driven sizes).
- `var()` allowed in `style` props for third-party libs that cannot accept `className` (e.g., Recharts).

## Styling Decision Tree

```
Static classes only?              → className="..."
Conditional or mergeable classes? → cn("base", condition && "variant")
Truly dynamic runtime value?      → style={{ width: `${x}%` }}
Third-party lib prop (no class)?  → style prop with CHART_COLORS constants
Fluid spacing?                    → style={{ padding: 'var(--space-md)' }}
```

## cn() Utility — Implementation

```typescript
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

## Tailwind v4 — CSS-First Configuration

No `tailwind.config.js`. Configuration and theme tokens live in CSS via `@theme`.

```css
/* globals.css */
@import "tailwindcss";
@import "tw-animate-css";

@custom-variant dark (&:is(.dark *));

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-primary: var(--primary);
  --color-border: var(--border);
  --radius-lg: var(--radius);

  /* Fluid typography — text-lg, text-2xl, etc. use clamp() automatically */
  --font-size-xs:   var(--text-xs);
  --font-size-sm:   var(--text-sm);
  --font-size-base: var(--text-base);
  --font-size-lg:   var(--text-lg);
  --font-size-xl:   var(--text-xl);
  --font-size-2xl:  var(--text-2xl);
  --font-size-3xl:  var(--text-3xl);
  --font-size-4xl:  var(--text-4xl);
}

:root {
  /* Fluid type scale: min at 320px viewport, max at 1200px */
  --text-xs:   clamp(0.75rem,  0.7rem  + 0.25vw, 0.875rem);
  --text-sm:   clamp(0.875rem, 0.8rem  + 0.35vw, 1rem);
  --text-base: clamp(1rem,     0.9rem  + 0.5vw,  1.125rem);
  --text-lg:   clamp(1.125rem, 1rem    + 0.6vw,  1.25rem);
  --text-xl:   clamp(1.25rem,  1.1rem  + 0.75vw, 1.5rem);
  --text-2xl:  clamp(1.5rem,   1.2rem  + 1.5vw,  2rem);
  --text-3xl:  clamp(1.875rem, 1.5rem  + 1.875vw,2.5rem);
  --text-4xl:  clamp(2.25rem,  1.8rem  + 2.25vw, 3rem);

  /* Fluid spacing scale */
  --space-2xs: clamp(0.125rem, 0.1rem  + 0.125vw, 0.25rem);
  --space-xs:  clamp(0.25rem,  0.2rem  + 0.25vw,  0.5rem);
  --space-sm:  clamp(0.5rem,   0.4rem  + 0.5vw,   0.75rem);
  --space-md:  clamp(1rem,     0.8rem  + 1vw,     1.5rem);
  --space-lg:  clamp(1.5rem,   1.2rem  + 1.5vw,   2rem);
  --space-xl:  clamp(2rem,     1.5rem  + 2.5vw,   3rem);
  --space-2xl: clamp(3rem,     2rem    + 5vw,     5rem);

  /* Semantic color tokens */
  --background: oklch(1 0 0);
  --foreground: oklch(0.145 0 0);
  --primary: oklch(56% 0.2 250);
  --border: oklch(30% 0.02 250);
  --radius: 0.5rem;
}

.dark {
  --background: oklch(0.145 0 0);
  --foreground: oklch(0.985 0 0);
}
```

**Why this works:** `text-lg` maps to `--font-size-lg` which maps to `var(--text-lg)` (a clamp value). Fluid without manual `style` props.

**WCAG zoom warning:** Never use bare `vw` — always combine with `rem`:
```css
/* WRONG — fails WCAG 1.4.4 */
font-size: clamp(1rem, 2vw, 2rem);
/* CORRECT — rem component scales with browser zoom */
font-size: clamp(1rem, 0.9rem + 0.5vw, 1.125rem);
```

## Container Queries — Component-Level Responsive

```tsx
/* Tailwind v4 — first-class @container support */
<div className="@container">
  <div className="flex flex-col @sm:flex-row @lg:grid @lg:grid-cols-3 gap-4" />
</div>
```

Three layers of responsive:
- **Media queries** (`sm:`, `md:`, `lg:`) → page-level layout changes
- **Container queries** (`@container`, `@sm:`, `@md:`) → component-level adaptation
- **`clamp()`** → fluid typography/spacing, no breakpoints needed

## Mobile-First Rules

```tsx
{/* Base = mobile, prefix = larger */}
<div className="flex flex-col gap-4 md:flex-row md:gap-6" />
```

Critical mobile patterns:
- Touch targets: `min-h-11 min-w-11` (44px) on all interactive elements
- Form inputs: always `text-base` (prevents iOS auto-zoom)
- Full height: `min-h-dvh` (not `min-h-screen`)
- Motion: always add `motion-reduce:` variant to animations

## Style Constants for Third-Party Libraries

```typescript
// Only for library props (Recharts fill, stroke, etc.) that cannot receive className
const CHART_COLORS = {
  primary: "var(--color-primary)",
  secondary: "var(--color-secondary)",
  border: "var(--color-border)",
};

<XAxis tick={{ fill: CHART_COLORS.primary }} />
<CartesianGrid stroke={CHART_COLORS.border} />
```

## Dark Mode Pattern (class-based via next-themes)

```css
@custom-variant dark (&:is(.dark *));
/* next-themes adds .dark to <html> — Tailwind dark: variants work automatically */
```

## CSS Modern Features (in Next.js + Tailwind stack)

**CSS :has() — style parent based on child:**
```css
.card:has(img) { grid-template-columns: 12.5rem 1fr; }
.form-group:has(input:invalid) .label { color: var(--destructive); }
```

**CSS Nesting — max 3 levels deep:**
```css
.card {
  padding: var(--space-md);
  &:hover { box-shadow: 0 0.25rem 1rem rgba(0,0,0,0.1); }
  @media (width >= 48rem) { flex-direction: row; }
  @container card (width >= 25rem) { flex-direction: row; }
}
```

**CSS @layer — declare order upfront:**
```css
@layer reset, base, tokens, components, utilities, overrides;
/* Unlayered styles override all layered styles */
```

**Logical properties — use instead of physical left/right:**
```css
/* Physical (avoid) → Logical (use) */
margin-left/right  → margin-inline
padding-left/right → padding-inline
text-align: left   → text-align: start
/* Tailwind: ms-, me-, ps-, pe-, start-, end- */
```

## Resources

- Tailwind CSS Docs — https://tailwindcss.com/docs
- Tailwind CSS Installation (Next.js) — https://tailwindcss.com/docs/installation/framework-guides/nextjs
- Tailwind CSS Theme Variables — https://tailwindcss.com/docs/theme
- Tailwind CSS Dark Mode — https://tailwindcss.com/docs/dark-mode
- Tailwind CSS Responsive Design — https://tailwindcss.com/docs/responsive-design

## Changelog

### v2.0 — 2026-04-16
- Added Tailwind v4 CSS-first configuration guidance, fluid typography and spacing tokens, container queries, and modern CSS integration rules.
