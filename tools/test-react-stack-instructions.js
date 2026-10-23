'use strict';
// Instruction contracts only: no React rendering, SSR execution or upstream API probe.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
const load = file => fs.existsSync(path.join(root, file)) ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const skill = name => load(`skills/frontend/${name}/SKILL.md`);
const includes = (text, values) => values.forEach(value => assert.ok(text.includes(value), `Missing: ${value}`));

test('R01: forms preserve the established architecture with one submission owner', () => {
  for (const name of ['react-19', 'shadcn-ui']) includes(skill(name), [
    'existing form architecture', 'one owner for submission, pending state and errors',
    'React Actions', 'Query mutations',
  ]);
  assert.doesNotMatch(skill('react-19'), /Move forms to `useActionState`|forms migrated to action-based flow/);
  assert.doesNotMatch(skill('shadcn-ui'), /keep form contracts explicit \(`React Hook Form` \+ Zod/);
});

test('R02: useFormState is renamed and deprecated, not removed', () => {
  includes(skill('react-19'), ['`ReactDOM.useFormState` was renamed to `React.useActionState`', 'deprecated, not removed', 'prefer `useActionState`']);
  assert.doesNotMatch(skill('react-19'), /removed `(?:ReactDOM\.)?useFormState`/);
});

test('R03: Zod conventions distinguish supported aliases from removed parameters', () => {
  includes(skill('zod-4'), ['`message` is a deprecated alias of `error`', 'throws if both are given',
    '`required_error`, `invalid_type_error` and `errorMap` parameters are removed',
    '`.nonempty()` remains an alias of `.min(1)`', 'prefer `.min(1)`',
    'top-level format validators (`z.email()`, `z.uuid()`, `z.url()`)', '`z.string().email()` are deprecated but still work']);
  assert.doesNotMatch(skill('zod-4'), /removed `\.nonempty\(\)`|Mixing v3 and v4 APIs \(`message`, `\.nonempty\(\)`/);
});

test('R04: Zustand shallow paths are valid and equality APIs are scoped', () => {
  includes(skill('zustand-5'), ['`zustand/shallow` and `zustand/react/shallow` are both valid',
    'prefer `zustand/react/shallow`', '`create` no longer takes an equality function',
    '`createWithEqualityFn` from `zustand/traditional`', 'requires `use-sync-external-store`', 'do not add it without approval']);
  assert.doesNotMatch(skill('zustand-5'), /don't\*\* use v4 path `zustand\/shallow`|outdated v4 `useShallow` import path|`zustand\/react\/shallow` \(v5 path\)/);
});

test('R05: Query keys use deterministic content hashing rather than reference identity', () => {
  includes(skill('tanstack-query'), ['hashed deterministically', 'object key order does not matter',
    'JSON-serializable objects', 'every variable the `queryFn` uses', 'equivalent contents', 'functions', 'class instances']);
  assert.doesNotMatch(skill('tanstack-query'), /unstable references create hidden cache misses|unstable object\/array references|keys constructed inline with changing object literals|keys from stable, serializable primitives/i);
});

test('R06: Query and Zustand require selection and preserve established alternatives', () => {
  includes(skill('tanstack-query'), ['installed or explicitly selected', 'Preserve other established solutions',
    'remote state managed in its client cache', 'server-only RSC reads', 'approved migration']);
  includes(skill('zustand-5'), ['installed or explicitly selected', 'other established state solutions',
    'not applicable', 'approved architectural migration plan']);
  includes(load('AGENTS.frontend.md'), ['Query-managed remote state', 'installed or explicitly selected', 'preserve other established solutions']);
  assert.doesNotMatch(skill('tanstack-query'), /whenever state comes from an API|for ALL server\/API state/);
  assert.doesNotMatch(skill('zustand-5'), /Zustand is the standard only when Redux is not/);
});

test('R07: Tailwind requires theme mappings and statically detectable classes', () => {
  includes(skill('tailwind-4'), ['only if the project theme defines them', 'bg-(--brand)', 'bg-[var(--brand)]',
    'bg-${color}-500', 'not detected', 'lookup map of complete class strings', 'shadcn/ui only when opted in']);
  assert.doesNotMatch(skill('tailwind-4'), /DON'T place `var\(--token\)` directly in `className`/);
});

test('R08: shadcn examples follow the installed primitive without an implicit upgrade', () => {
  includes(skill('shadcn-ui'), ['documentation matching the installed setup', 'Radix-based example',
    '`asChild` comes from the Radix Slot primitive', 'follow the primitive actually installed', 'do not upgrade']);
  assert.doesNotMatch(skill('shadcn-ui'), /prefer the latest documented|prefer the most recent documented/);
});

test('R09: shadcn source is owned and shared fixes are intentional', () => {
  includes(skill('shadcn-ui'), ['owned source', 'feature-specific behavior', 'shared variants',
    'reusable primitive fixes', 'upstream/regeneration diff review']);
  assert.doesNotMatch(skill('shadcn-ui'), /vendor-like generated base|remain untouched unless explicitly justified/);
});

test('R10: SSR lifetimes isolate requests and preserve matching hydration state', () => {
  includes(skill('tanstack-query'), ['new `QueryClient` per request', 'one stable browser instance',
    'never share a server client across requests', 'matching initial snapshots', 'verify against the installed version']);
  includes(skill('zustand-5'), ['request-local stores', 'stable browser store', 'matching initial snapshots',
    'controlled persisted-state hydration', 'verify against the installed version']);
  assert.doesNotMatch(skill('tanstack-query'), /- Calling `new QueryClient\(\)` inside render paths\./);
});

test('R11: Next server caches and the Query browser cache have separate owners', () => {
  includes(skill('nextjs-15'), ['`fetch` is not cached by default', 'GET Route Handlers are not cached by default',
    'changed from Next.js 14', 'Next cached reads', 'Uncached reads', 'Query-managed reads',
    '`revalidatePath`/`revalidateTag` act on Next server caches only', 'do not invalidate a TanStack Query browser cache']);
  includes(skill('tanstack-query'), ['Next server caches', '`revalidatePath`/`revalidateTag`', 'browser cache', 'invalidateQueries']);
  assert.doesNotMatch(skill('nextjs-15'), /DO revalidate cache after successful mutations|Confirm revalidation happens after successful writes/);
});

test('R12: optimism cancels reads and bounds overlapping-write rollback', () => {
  includes(skill('tanstack-query'), ['await queryClient.cancelQueries({ queryKey })', 'getQueryData',
    'return context', '`onError`', 'restore from context', '`onSettled`',
    'one in-flight optimistic mutation per entity', 'shared snapshot/cache key',
    'older failure cannot roll back newer work', 'verify against the installed version']);
});

test('R13: schema-derived types distinguish raw input and parsed output and await async validation', () => {
  for (const name of ['zod-4', 'typescript']) includes(skill(name), ['z.input<typeof Schema>',
    'z.output<typeof Schema>', 'z.infer', 'form defaults', 'coercions/transforms',
    'safeParseAsync', 'async refinements/transforms']);
});

test('R14: compiler status checks do not authorize activation', () => {
  includes(skill('react-19'), ['check whether React Compiler is already enabled',
    'separate user authorization', 'If disabled, retain necessary memoization', 'verify against the installed version']);
  assert.doesNotMatch(skill('react-19'), /Enable and validate React Compiler/);
});

test('preservation: SPA and Next profile gates and all existing frontend routes remain', () => {
  const guide = load('AGENTS.frontend.md');
  includes(guide, ['React SPA', 'existing authenticated client', 'Do not introduce Next.js',
    'Client -> Server Action -> .NET API', '`params`/`searchParams`', 'backend authorization',
    'shared selection procedure', 'shared security baseline', 'preserve installed tooling and form/client flow']);
  for (const route of ['typescript', 'react-19', 'nextjs-15', 'tailwind-4', 'zod-4', 'zustand-5',
    'tanstack-query', 'shadcn-ui', 'security', 'github-pr']) includes(guide, [`\`${route}\``]);
  includes(guide, ['skills/shared/diagnosing-bugs/SKILL.md', 'skills/shared/accessibility/SKILL.md',
    'skills/frontend/testing/SKILL.md', 'React SPA only']);
  includes(skill('react-19'), ['no Server Components, Server Actions or `"use client"` scaffolding',
    'React 19 alone does not supply server infrastructure', 'existing client mutation path']);
  includes(skill('tanstack-query'), ['SPA without SSR', 'Do not add Next.js, Server Actions or an SSR host',
    'existing authenticated client']);
  includes(skill('nextjs-15'), ['Client -> fetch wrapper -> Server Action -> backend API']);
  includes(skill('shadcn-ui'), ['only when the developer explicitly requests shadcn/ui']);
  includes(skill('zustand-5'), ['Do **not** force migration, parallel adoption, or coexistence by default.']);
});
