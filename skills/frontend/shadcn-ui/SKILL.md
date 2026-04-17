---
name: shadcn-ui
description: >
  ZeshOne shadcn/ui + React Hook Form conventions.
  Trigger: When building UI components, forms, dark mode, or using shadcn/ui primitives.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.0"
  inspired-by: frontend-nextjs/components.md
---

## When to Use

Load this skill when building UI components with shadcn/ui, React Hook Form + Zod, dark mode, or animated components.

## Critical Patterns

- Install and add components with the `shadcn` CLI, not with `npm install` for individual UI primitives.
- Treat generated files in `components/ui/` as vendor-like building blocks — customize through composition and CSS variables instead of editing them directly.
- Build forms with the shadcn `Form` primitives on top of React Hook Form and Zod.
- Use `next-themes` class-based dark mode through a shared `ThemeProvider`.
- Keep heavy animated libraries in client-only boundaries and use `dynamic()` when SSR should be disabled.

## shadcn/ui Setup

```bash
# Initialize in Next.js project
npx shadcn@latest init

# Add components
npx shadcn@latest add button card dialog form input label
npx shadcn@latest add dropdown-menu navigation-menu sheet
npx shadcn@latest add table tabs sonner skeleton sidebar
```

```json
// components.json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": true,
  "tsx": true,
  "tailwind": {
    "config": "",
    "css": "app/globals.css",
    "baseColor": "neutral",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "hooks": "@/hooks"
  }
}
```

**Key facts:**
- Uses unified `radix-ui` package — not individual `@radix-ui/react-*`
- All theme colors use oklch color model
- Files land in `components/ui/` — never edit them directly
- Customize via CSS variables, not component props
- Install via `npx shadcn@latest add [component]`, never `npm install`

## Component Composition Pattern

```tsx
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

function ProductCard({ product }: { product: Product }) {
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg">{product.name}</CardTitle>
          <Badge variant={product.inStock ? "default" : "secondary"}>
            {product.inStock ? "In Stock" : "Out of Stock"}
          </Badge>
        </div>
      </CardHeader>
      <CardContent>
        <p className="text-2xl font-bold">${product.price}</p>
      </CardContent>
      <CardFooter>
        <Button className="w-full" disabled={!product.inStock}>Add to Cart</Button>
      </CardFooter>
    </Card>
  );
}
```

## Form Pattern — React Hook Form + Zod + shadcn

```tsx
"use client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

// ⚠️ Requires Zod 4 (npm install zod@beta). See zod-4 skill for API differences.
const formSchema = z.object({
  username: z.string().min(2).max(50),
  email: z.email(),
});

export function ProfileForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { username: "", email: "" },
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    // Call Server Action or mutation
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl><Input placeholder="johndoe" {...field} /></FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Save</Button>
      </form>
    </Form>
  );
}
```

## Dark Mode — next-themes

```tsx
// providers/theme-provider.tsx
"use client";
import { ThemeProvider as NextThemesProvider } from "next-themes";

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  return (
    <NextThemesProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
      {children}
    </NextThemesProvider>
  );
}
```

```tsx
// components/shared/theme-toggle.tsx
"use client";
import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

export function ThemeToggle() {
  const { setTheme } = useTheme();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="icon">
          <Sun className="h-[1.2rem] w-[1.2rem] scale-100 dark:scale-0" />
          <Moon className="absolute h-[1.2rem] w-[1.2rem] scale-0 dark:scale-100" />
          <span className="sr-only">Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme("light")}>Light</DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("dark")}>Dark</DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme("system")}>System</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
```

## Usage Decision — shadcn vs Animation Libraries

| Need | Use |
|------|-----|
| Standard UI (buttons, forms, dialogs) | **shadcn/ui** |
| Data display (tables, charts) | **shadcn/ui** |
| Hero sections, landing pages | **React Bits / Framer Motion** |
| Animated text effects | **React Bits** |
| Interactive backgrounds | **React Bits** |
| Navigation & layout | **shadcn/ui** |

**Animated components (React Bits):** Load only in Client Components. Use `dynamic()` import for heavy ones to avoid SSR overhead.

```tsx
import dynamic from "next/dynamic";
const Particles = dynamic(() => import("@/components/animated/particles"), { ssr: false });
```

## Resources

- shadcn/ui Docs — https://ui.shadcn.com/docs
- shadcn/ui Next.js Installation — https://ui.shadcn.com/docs/installation/next
- shadcn/ui Components — https://ui.shadcn.com/docs/components
- shadcn/ui Forms — https://ui.shadcn.com/docs/components/form
- shadcn/ui Dark Mode — https://ui.shadcn.com/docs/dark-mode/next
