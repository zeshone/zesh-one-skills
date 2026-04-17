---
name: zustand-5
description: >
  ZeshOne Zustand 5 conventions.
  Trigger: When creating or configuring Zustand stores.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.2"
  inspired-by: gentleman-programming/zustand-5
---

## When to Use

Load this skill when creating or configuring Zustand stores — store structure, selectors, async actions, or persistence.

## Critical Patterns

- One store per feature/domain — never one global store with everything.
- ALWAYS select specific fields with a selector — never consume the entire store object.
- Use `useShallow` when selecting multiple fields simultaneously.
- Use local state (`useState`) for data that lives in a single component or small subtree.
- Zustand → global state shared across features, with persistence or devtools.
- Use `immer` middleware when mutations would otherwise require manual spreading of deeply nested state.

## Store Organization — Per Feature

```typescript
// stores/user.store.ts  ← one file per domain
import { create } from "zustand";

interface UserStore {
  user: User | null;
  setUser: (user: User) => void;
  clearUser: () => void;
}

export const useUserStore = create<UserStore>((set) => ({
  user: null,
  setUser: (user) => set({ user }),
  clearUser: () => set({ user: null }),
}));
```

## Selectors — Always Specific (REQUIRED)

```typescript
// ✅ Select only what the component needs
const name = useUserStore((state) => state.name);

// ✅ Multiple fields — use useShallow to avoid extra re-renders
import { useShallow } from "zustand/react/shallow"; // v5 import path

const { name, email } = useUserStore(
  useShallow((state) => ({ name: state.name, email: state.email }))
);

// ❌ NEVER — re-renders on any state change
const store = useUserStore();
```

> **Zustand v5 breaking change**: `useShallow` import path changed from `zustand/shallow` (v4) to `zustand/react/shallow` (v5). Use the v5 path.

## Async Actions — Standard Shape

Actions live in `@/features/{feature}/actions` — import them into the store, never call services directly.

```typescript
import { fetchProductsAction } from "@/features/products/actions";

// Inside create():
fetchProducts: async () => {
  set({ loading: true, error: null });
  try {
    const products = await fetchProductsAction();
    set({ products, loading: false });
  } catch {
    set({ error: "Failed to load products", loading: false });
  }
},
```

## Persist Middleware

```typescript
import { persist } from "zustand/middleware";

export const useSettingsStore = create<SettingsStore>()(
  persist(
    (set) => ({
      theme: "light",
      setTheme: (theme) => set({ theme }),
    }),
    { name: "zesh-settings" } // localStorage key — use project prefix
  )
);
```

## Resources

- Zustand Docs — https://zustand.docs.pmnd.rs/
- Zustand Getting Started — https://zustand.docs.pmnd.rs/getting-started/introduction
- Zustand Prevent Rerenders with useShallow — https://zustand.docs.pmnd.rs/guides/prevent-rerenders-with-use-shallow
- Zustand Persist Middleware — https://zustand.docs.pmnd.rs/integrations/persisting-store-data
