'use strict';
// Instruction contracts only: no Angular app, browser, native runtime or device is executed.
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const root = path.resolve(__dirname, '..');
const load = file => fs.existsSync(path.join(root, file)) ? fs.readFileSync(path.join(root, file), 'utf8') : '';
const ionic = load('skills/apps/ionic-angular/SKILL.md');
const capacitor = load('skills/apps/capacitor/SKILL.md');
const guide = load('AGENTS.apps.md');
const includes = (text, values) => values.forEach(value => assert.ok(text.includes(value), `Missing: ${value}`));

test('A1: standalone and signals defaults preserve existing NgModule and RxJS ownership', () => {
  includes(ionic, ['Angular 19 makes components standalone by default', 'NgModule-based apps remain supported',
    'loadComponent', 'loadChildren', 'Preserve valid NgModule and RxJS ownership', 'no migration is required']);
  includes(guide, ['Preserve valid NgModule and RxJS ownership', 'signals where appropriate', 'standalone routes']);
  for (const text of [ionic, guide]) {
    assert.doesNotMatch(text, /Standalone-first is mandatory|Signals[^.\n]*for new code|Don't[^\n]*introduce new NgModule-only patterns/i);
  }
});

test('A2: entry refresh is conditional and cached-page work stops before destruction', () => {
  includes(ionic, ['refreshable or stale data', 'ngOnInit', 'once per page instance', 'ngOnDestroy',
    'removed from the stack', 'ionViewWillEnter', 'ionViewDidEnter', 'ionViewWillLeave', 'ionViewDidLeave',
    'ViewWillEnter', 'ViewDidEnter', 'ViewWillLeave', 'ViewDidLeave',
    'Stop work that must not run while hidden on leave', 'destruction-only cleanup is insufficient',
    'Do not fetch on every entry when cached data is still valid']);
  assert.doesNotMatch(ionic, /refresh is mandatory on routed pages|rather than skipping refresh/);
});

test('A3: Ionic navigation supports Router integration and independent tab stacks', () => {
  for (const text of [ionic, guide]) {
    includes(text, ['NavController', 'navigateForward', 'navigateBack', 'navigateRoot',
      'Router', 'routerLink', 'routerDirection', 'Each tab keeps its own stack',
      "another tab's child pages", 'notification or deep link']);
    assert.doesNotMatch(text, /Do not programmatically jump between tab roots|Don't[^\n]*default to `Router.navigate`/);
  }
});

test('A4: Capacitor checks plugin capability, permissions and native configuration without excluding web', () => {
  includes(capacitor, ['installed plugin', 'Capacitor.isNativePlatform()', 'Capacitor.getPlatform()',
    "Capacitor.isPluginAvailable('<Name>')", 'web implementations', 'Do not blanket-disable web',
    'checkPermissions()', 'requestPermissions()', 'granted', 'denied', 'prompt', 'prompt-with-rationale',
    'Info.plist', 'Android manifest', "plugin's own docs", 'denied/restricted', 'unsupported',
    'cancellation', 'operational failures', 'verify against the installed version']);
  assert.doesNotMatch(capacitor, /gate plugin execution by runtime capability \(`native` vs `web`\)|only run plugins on native/i);
});

test('A5: Angular template and DI idioms are version-gated and preserve working conventions', () => {
  includes(ionic, ['@if', '@for', '@switch', 'stable since Angular 17', '@for` requires `track`',
    'stable entity identity', 'structural directives', 'constructor injection',
    'inject()', 'injection context', 'constructor', 'field initializer', 'provider factory',
    'runInInjectionContext', 'verify against the installed version']);
  assert.doesNotMatch(ionic, /always replace structural directives|inject\(\) anywhere|constructor injection is prohibited/i);
});

test('A6: routine behavioral verification covers changed Ionic routes using installed tooling', () => {
  includes(ionic, ['Routine behavior checks', 'installed testing tooling', 'initial entry and revisit',
    'loading/empty/error/retry', 'stale request races', 'back/tab behavior', 'retained form state',
    'Keep native/device checks in the Capacitor skill', 'unperformed checks']);
  assert.doesNotMatch(ionic, /install a new test runner|browser checks prove native behavior/i);
});

test('preservation: installed context, app ownership, lifecycle and shared routes remain intact', () => {
  includes(ionic, ['installed versions', 'Do not upgrade', 'Reactive Forms', 'HttpClient',
    'valid page shell', 'presentational', 'OnPush', 'ion-router-outlet', 'markForCheck',
    '@ionic/angular/standalone', 'idempotent',
    '../../shared/security/baseline.md', '../../shared/security/application-checks.md']);
  includes(capacitor, ['installed Capacitor', 'Do not install a plugin or run native sync for a UI-only task',
    'ionic cap sync', 'try/catch', 'cleartext/mixed content', 'real devices',
    '../../shared/security/baseline.md', '../../shared/security/delivery-checks.md']);
  includes(guide, ['installed versions', 'Reactive Forms and HttpClient', 'component -> service (`HttpClient`)',
    'valid page shell', 'OnPush', 'Do not generalize this to every routed page',
    'skills/shared/security/baseline.md', 'skills/shared/security/application-checks.md',
    'skills/shared/security/delivery-checks.md', 'skills/shared/diagnosing-bugs/SKILL.md',
    'skills/shared/accessibility/SKILL.md', 'UI accessibility, keyboard/focus, labels and feedback',
    'server-side authorization', 'Client bundles and Angular environment files cannot keep secrets private']);
});
