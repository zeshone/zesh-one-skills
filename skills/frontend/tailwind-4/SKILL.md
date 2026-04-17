---
name: tailwind-4
description: >
  ZeshOne Tailwind CSS 4 conventions.
  Trigger: When styling with Tailwind — cn(), theme variables, dynamic values.
license: Apache-2.0
metadata:
  author: zesh-one
  version: "1.0"
  inspired-by: gentleman-programming/tailwind-4
---

## When to Use

Load this skill when styling with Tailwind CSS 4 — `cn()`, theme variables, dynamic values, or third-party library integration.

## Critical Patterns

- NEVER use `var()` inside `className` — use Tailwind semantic classes (`bg-primary`, `text-slate-400`).
- NEVER use hex colors in `className` — always use Tailwind color classes.
- Use `cn()` for conditional or merged classes. Use plain `className` for static-only classes.
- Use the `style` prop only for truly dynamic runtime values (percentages, progress bars, user-driven sizes).
- `var()` is allowed ONLY inside `style` props for library integrations (e.g., Recharts) that cannot accept `className`.

## Styling Decision Tree

```
Static classes only?              → className="..."
Conditional or mergeable classes? → cn("base", condition && "variant")
Truly dynamic runtime value?      → style={{ width: `${x}%` }}
Third-party lib prop (no class)?  → style prop with CHART_COLORS constants
```

## cn() Utility — Implementation

```typescript
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

## Style Constants for Third-Party Libraries

When a library prop (like Recharts `fill`, `stroke`) cannot receive a `className`:

```typescript
// ✅ Define var() constants — ONLY for library props
const CHART_COLORS = {
  primary: "var(--color-primary)",
  secondary: "var(--color-secondary)",
  border: "var(--color-border)",
};

<XAxis tick={{ fill: CHART_COLORS.primary }} />
<CartesianGrid stroke={CHART_COLORS.border} />
```

## Tailwind v4 — CSS-First Configuration (Breaking Change)

Tailwind v4 removes `tailwind.config.js`. Configuration and theme tokens live in CSS via `@theme`.

```css
/* globals.css */
@import "tailwindcss";

@theme {
  --color-primary: oklch(56% 0.2 250);
  --color-secondary: oklch(70% 0.15 200);
  --color-border: oklch(30% 0.02 250);
}
```

This is why `bg-primary` works as a class — it maps to `--color-primary` defined in `@theme`. Without `@theme`, `bg-primary` does nothing.

## Keywords
tailwind, tailwind 4, css, cn, theme, var, styling
