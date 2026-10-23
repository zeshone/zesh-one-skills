# React SPA rendered-behavior trial

This small specimen validates an API-backed list and form through rendered accessible controls. It uses React19, Testing Library and jsdom; it is component acceptance with an injected controlled client, not production Query integration, an application scaffold or browser/E2E proof.

## Quick path

Default collection checks remain unchanged in prerequisites: `npm test` and `npm run verify` do not install or load these packages. The optional trial requires Node>=20 and its own local dependencies.

After explicit installation approval, from the collection root:

```sh
# Local project/cache/temp only. Disable lifecycle scripts and optional native peers.
# Set TEMP/TMP/TMPDIR to a project-local directory before npm on your platform.
npm ci --prefix tools/fixtures/frontend-dom --ignore-scripts --omit=optional --no-audit --no-fund --registry=https://registry.npmjs.org/ --cache .frontend-dom-npm/cache --userconfig .frontend-dom-npm/npmrc --globalconfig .frontend-dom-npm/global-npmrc
npm run test:frontend-dom
npm run test:frontend-dom -- --help
```

For the initial dependency setup, the maintainer used `npm install` with the same boundary flags to generate the committed lockfile. Later installs use `npm ci`. Do not run either without installation authorization; testing/help never installs. npm can contact the registry only during that separately approved installation. Downloaded packages/cache stay ignored under this repository; nothing is installed globally, and root package dependencies/CI are not changed.

Direct dependencies are fixed in [package.json](package.json); [package-lock.json](package-lock.json) fixes the resolved transitive graph/integrity. jsdom26.1.0 is deliberately compatible with the available Node23.11.0; jsdom30.1.3 requires a newer supported Node line. This trial does not authorize Node or consumer-framework upgrades. React and react-dom use the same19.3.1 version; Testing Library React16.3.3 supports React19 and DOM10.4.2. Upgrading them is a separate reviewed dependency change, not a floating runtime resolution.

## Contract and isolation

[contract.test.mjs](contract.test.mjs) exercises [NotesPanel.mjs](NotesPanel.mjs) through roles, accessible names, visible content and user events:

1. Loading: visible status, no premature empty/success/error, one list request and no writes.
2. List failure/retry: safe error, no leaked exception details, new request, cleared error and recovered records.
3. Empty: explicit empty result, not loading/failure/populated list.
4. Success: expected records as text, no stale states or HTML interpretation.
5. Invalid form: field error, no create call, unchanged API-visible records/storage and usable controls.
6. Valid form/recovery: exact trimmed payload, pending/disabled controls, no duplicate write, safe recoverable failure retaining input, then one stored success, cleared error/input and enabled controls.

Deferred promises settle async states without sleeps. The tests never call production APIs, fetch external resources, run page scripts or use real credentials. The tiny client adapter is in-memory synthetic test state. Each test unmounts roots and checks localStorage via its public interface; the shared jsdom is closed and temporary globals restored at suite end. The runner uses a unique project-local working/temp directory and removes only that owned directory, reporting cleanup independently. Node launches set `NODE_DISABLE_COMPILE_CACHE=1`.

Help exits0 on stdout; unknown runner options exit2 on stderr. Missing isolated packages exit2 with a prerequisite message and no install attempt. A failed/aborted/timed-out child exits nonzero. Test/cleanup failure is not success; interruption may prevent cleanup. Do not recursively delete unexpected leftovers without checking ownership.

## Evidence limits

- jsdom is a simulated DOM, not browser layout, focus/navigation fidelity, native form validation, visual accessibility, an accessibility audit or E2E.
- Using semantic queries does not itself establish accessibility compliance. An injected adapter does not prove real HTTP, authentication, .NET API or production cache/provider integration.
- This is executed specimen behavior, not model compliance or general security. npm packages and the inherited host environment remain trusted; ignored local directories are not a sandbox.
- The collection does not impose these dependencies, runner, form flow or sample state design on consumer projects. Existing React SPA client/router, Query/form choices and backend validation/authorization remain their owners.
