---
name: ionic-angular
description: >
  ZeshOne Ionic + Angular conventions for mobile apps.
  Trigger: When building Ionic Angular apps — standalone components, navigation, lifecycle, forms, state with signals, services, performance.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.0"
  inspired-by: capawesome-team/ionic-angular, ionic-capacitor/angular.md
---

## When to Use

Load this skill when building mobile apps with Ionic + Angular — pages, navigation, lifecycle hooks, reactive forms, signals, services, or performance.

## Critical Patterns

- Detect whether the app is standalone or NgModule-based before adding Ionic imports or providers.
- Every routed page must keep the Ionic page shell: `<ion-page>` + header/content structure.
- Use Ionic lifecycle hooks such as `ionViewWillEnter` for refresh logic because cached pages do not re-run `ngOnInit` on revisit.
- Use `NavController` for animated navigation flows instead of Angular `Router` directly.
- Keep state local with signals when possible, and keep presentational components `OnPush` and input/output driven.

## Stack Base

| Package | Version |
|---------|---------|
| `@ionic/angular` | 8.x |
| `@angular/core` | 19.x |
| `ionicons` | 7.x |
| Architecture | **Standalone components** (no NgModules) |

## Auto-Detect Architecture

Before writing code, check `src/main.ts`:
- `bootstrapApplication` → **Standalone** (preferred, use `@ionic/angular/standalone` imports)
- `platformBrowserDynamic().bootstrapModule` → **NgModule** (use `IonicModule` globally, import from `@ionic/angular`)

| Aspect | Standalone | NgModule |
|--------|-----------|----------|
| Ionic setup | `provideIonicAngular({})` in `app.config.ts` | `IonicModule.forRoot()` in `app.module.ts` |
| Import source | `@ionic/angular/standalone` | `@ionic/angular` |
| Lazy loading | `loadComponent` | `loadChildren` |
| Icons | `addIcons()` required | Automatic |
| Tree-shaking | ✅ | ❌ |

## Project Structure

```
src/
├── app/
│   ├── features/              # Feature-based vertical slices
│   │   ├── auth/
│   │   │   ├── login.page.ts
│   │   │   └── login.page.html
│   │   ├── home/
│   │   │   ├── home.page.ts
│   │   │   └── components/   # Feature-specific components
│   │   └── settings/
│   ├── shared/
│   │   ├── components/        # Reusable UI components
│   │   └── pipes/
│   ├── core/
│   │   ├── guards/            # Auth guards
│   │   ├── interceptors/      # HTTP interceptors
│   │   └── services/          # App-wide services
│   ├── app.component.ts
│   ├── app.config.ts          # Providers (standalone)
│   └── app.routes.ts
├── environments/
│   ├── environment.ts         # Dev — apiUrl: 'http://localhost:3000/api'
│   └── environment.prod.ts    # Prod — apiUrl: 'https://api.example.com'
└── theme/
    └── variables.scss         # Ionic CSS custom properties
```

## Bootstrap (Standalone)

```typescript
// main.ts
import { bootstrapApplication } from '@angular/platform-browser';
import { provideIonicAngular } from '@ionic/angular/standalone';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { routes } from './app/app.routes';
import { AppComponent } from './app/app.component';
import { authInterceptor } from './app/core/interceptors/auth.interceptor';

bootstrapApplication(AppComponent, {
  providers: [
    provideIonicAngular({}),
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor])),
  ],
});

// app.component.ts
import { Component } from '@angular/core';
import { IonApp, IonRouterOutlet } from '@ionic/angular/standalone';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [IonApp, IonRouterOutlet],
  template: `<ion-app><ion-router-outlet></ion-router-outlet></ion-app>`,
})
export class AppComponent {}
```

## Page Structure — Never Skip

Every routed page must use this exact structure. Ionic's scroll, safe areas, and transitions depend on it.

```typescript
import { Component } from '@angular/core';
import {
  IonHeader, IonToolbar, IonTitle, IonContent,
  IonPage,
} from '@ionic/angular/standalone';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonPage],
  template: `
    <ion-page>
      <ion-header>
        <ion-toolbar>
          <ion-title>Home</ion-title>
        </ion-toolbar>
      </ion-header>
      <ion-content fullscreen>
        <!-- iOS large-title collapse effect -->
        <ion-header collapse="condense">
          <ion-toolbar>
            <ion-title size="large">Home</ion-title>
          </ion-toolbar>
        </ion-header>
        <!-- Page content -->
      </ion-content>
    </ion-page>
  `,
})
export class HomePage {}
```

## Lifecycle Hooks — Critical Difference from Angular

Ionic caches pages in the DOM after navigation. `ngOnInit` fires **once**; Ionic hooks fire on **every** visit.

| Hook | Fires | Use for |
|------|-------|---------|
| `ngOnInit` | Once (first creation) | One-time setup, subscriptions |
| `ionViewWillEnter` | Every visit (before animation) | **Refresh data** |
| `ionViewDidEnter` | Every visit (after animation) | Heavy work that would block animation |
| `ionViewWillLeave` | Before leaving | Save drafts, pause media |
| `ionViewDidLeave` | After leaving | Cleanup |
| `ngOnDestroy` | Page popped from stack | Final cleanup, unsubscribe |

```typescript
import { ViewWillEnter, ViewWillLeave, ViewDidLeave } from '@ionic/angular/standalone';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

export class DashboardPage implements ViewWillEnter, ViewWillLeave {
  private leave$ = new Subject<void>();

  ionViewWillEnter() {
    this.loadData(); // Always refresh on enter

    // Restart subscriptions with takeUntil(leave$)
    this.dataService.stream$
      .pipe(takeUntil(this.leave$))
      .subscribe(data => this.data = data);
  }

  ionViewWillLeave() {
    this.leave$.next(); // Kill subscriptions when leaving
  }
}
```

## Navigation

### Route Config (Standalone + Tabs + Auth Guard)

```typescript
// app.routes.ts
export const routes: Routes = [
  { path: '', redirectTo: 'tabs', pathMatch: 'full' },
  {
    path: 'login',
    loadComponent: () => import('./features/auth/login.page').then(m => m.LoginPage),
  },
  {
    path: 'tabs',
    loadComponent: () => import('./features/tabs/tabs.page').then(m => m.TabsPage),
    canActivate: [authGuard],
    children: [
      { path: 'home', loadComponent: () => import('./features/home/home.page').then(m => m.HomePage) },
      { path: 'search', loadComponent: () => import('./features/search/search.page').then(m => m.SearchPage) },
      { path: 'profile', loadComponent: () => import('./features/profile/profile.page').then(m => m.ProfilePage) },
      { path: '', redirectTo: 'home', pathMatch: 'full' },
    ],
  },
];
```

### Tabs Page

```typescript
import { IonTabs, IonTabBar, IonTabButton, IonIcon, IonLabel } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { home, search, person } from 'ionicons/icons';

@Component({
  standalone: true,
  imports: [IonTabs, IonTabBar, IonTabButton, IonIcon, IonLabel],
  template: `
    <ion-tabs>
      <ion-tab-bar slot="bottom">
        <ion-tab-button tab="home"><ion-icon name="home"></ion-icon><ion-label>Home</ion-label></ion-tab-button>
        <ion-tab-button tab="search"><ion-icon name="search"></ion-icon><ion-label>Search</ion-label></ion-tab-button>
        <ion-tab-button tab="profile"><ion-icon name="person"></ion-icon><ion-label>Profile</ion-label></ion-tab-button>
      </ion-tab-bar>
    </ion-tabs>
  `,
})
export class TabsPage {
  constructor() { addIcons({ home, search, person }); }
}
```

**Tab rules:** Each tab has its own independent navigation stack. Never navigate between tabs programmatically — users switch via the tab bar only. The `tab` attribute on `<ion-tab-button>` must match the child route `path`.

### Programmatic Navigation

```typescript
import { NavController } from '@ionic/angular/standalone';

private navCtrl = inject(NavController);

goToDetail(id: string) { this.navCtrl.navigateForward(`/detail/${id}`); }
goBack()               { this.navCtrl.back(); }
goToRoot()             { this.navCtrl.navigateRoot('/tabs/home'); }
```

Use `NavController` (not `Router`) for page transitions — it handles Ionic animations. Add `routerDirection="forward|back|root"` on template `routerLink` elements.

### Functional Route Guard

```typescript
// core/guards/auth.guard.ts
import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);
  return authService.isAuthenticated() ? true : router.createUrlTree(['/login']);
};
```

### Modal

```typescript
import { ModalController } from '@ionic/angular/standalone';

private modalCtrl = inject(ModalController);

async openEdit(itemId: string) {
  const modal = await this.modalCtrl.create({
    component: EditModalComponent,
    componentProps: { itemId },
  });
  await modal.present();
  const { data, role } = await modal.onDidDismiss();
  if (role === 'confirm') { /* handle result */ }
}

// Inside the modal component:
confirm() { this.modalCtrl.dismiss({ saved: true }, 'confirm'); }
cancel()  { this.modalCtrl.dismiss(null, 'cancel'); }
```

## State — Angular Signals (Preferred)

```typescript
import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Preferences } from '@capacitor/preferences';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private http = inject(HttpClient);
  private userSignal = signal<User | null>(null);
  private tokenSignal = signal<string | null>(null);

  user = this.userSignal.asReadonly();
  isAuthenticated = computed(() => this.tokenSignal() !== null);

  async login(credentials: Credentials): Promise<void> {
    const res = await firstValueFrom(
      this.http.post<AuthResponse>('/api/auth/login', credentials)
    );
    this.userSignal.set(res.user);
    this.tokenSignal.set(res.token);
    await Preferences.set({ key: 'token', value: res.token });
  }

  async logout(): Promise<void> {
    this.userSignal.set(null);
    this.tokenSignal.set(null);
    await Preferences.remove({ key: 'token' });
  }

  async restoreSession(): Promise<void> {
    const { value } = await Preferences.get({ key: 'token' });
    if (value) this.tokenSignal.set(value);
  }
}
```

**Component-level signals:**
```typescript
items = signal<Item[]>([]);
searchTerm = signal('');
isLoading = signal(false);
filteredItems = computed(() =>
  this.items().filter(i => i.name.toLowerCase().includes(this.searchTerm().toLowerCase()))
);
```

**Resource API (Angular 19+ — server data):**
```typescript
import { resource } from '@angular/core';

itemsResource = resource({
  request: () => this.categoryId(),
  loader: async ({ request: category }) => {
    const res = await fetch(`/api/items?category=${category}`);
    return res.json() as Promise<Item[]>;
  },
});
// Template: itemsResource.isLoading(), itemsResource.value()
```

## HTTP Service + Interceptor Pattern

```typescript
// core/services/items.service.ts
@Injectable({ providedIn: 'root' })
export class ItemsService {
  private http = inject(HttpClient);
  private baseUrl = environment.apiUrl + '/items';

  getAll()                            { return this.http.get<Item[]>(this.baseUrl); }
  getById(id: string)                 { return this.http.get<Item>(`${this.baseUrl}/${id}`); }
  create(item: Partial<Item>)         { return this.http.post<Item>(this.baseUrl, item); }
  update(id: string, data: Partial<Item>) { return this.http.put<Item>(`${this.baseUrl}/${id}`, data); }
  delete(id: string)                  { return this.http.delete<void>(`${this.baseUrl}/${id}`); }
}

// core/interceptors/auth.interceptor.ts
export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const token = inject(AuthService).token();
  if (token) req = req.clone({ setHeaders: { Authorization: `Bearer ${token}` } });
  return next(req);
};
// Register: provideHttpClient(withInterceptors([authInterceptor]))
```

**Use `ionViewWillEnter` (not constructor/ngOnInit) for HTTP calls** — Ionic caches pages and ngOnInit won't re-fire on back-navigation.

## Forms — Reactive Forms + Ionic

```typescript
import { ReactiveFormsModule, FormBuilder, Validators, inject } from '@angular/forms';
import { IonList, IonItem, IonInput, IonSelect, IonSelectOption, IonButton } from '@ionic/angular/standalone';

@Component({
  standalone: true,
  imports: [ReactiveFormsModule, IonList, IonItem, IonInput, IonSelect, IonSelectOption, IonButton],
  template: `
    <form [formGroup]="form" (ngSubmit)="onSubmit()">
      <ion-list>
        <ion-item>
          <ion-input
            formControlName="name"
            label="Name" labelPlacement="stacked" fill="outline"
            [errorText]="form.controls.name.errors?.['required'] ? 'Required' : ''"
            [class.ion-invalid]="form.controls.name.touched && form.controls.name.invalid"
            [class.ion-touched]="form.controls.name.touched"
          ></ion-input>
        </ion-item>
        <ion-item>
          <ion-input formControlName="email" label="Email" labelPlacement="stacked" fill="outline" type="email"></ion-input>
        </ion-item>
      </ion-list>
      <ion-button expand="block" type="submit" [disabled]="form.invalid">Save</ion-button>
    </form>
  `,
})
export class MyFormPage {
  private fb = inject(FormBuilder);
  form = this.fb.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
    email: ['', [Validators.required, Validators.email]],
  });
  onSubmit() { if (this.form.valid) { /* call service */ } }
}
```

Use `labelPlacement="stacked"` + `fill="outline"` for modern Ionic form styling.

## Component Patterns

### Smart / Presentational Split

```typescript
// Smart — fetches data, handles logic
@Component({ /* IonPage wrapper, fullscreen content */ })
export class ItemListPage implements ViewWillEnter {
  private svc = inject(ItemsService);
  items = signal<Item[]>([]);
  isLoading = signal(true);

  ionViewWillEnter() { this.load(); }

  private async load() {
    this.isLoading.set(true);
    this.items.set(await firstValueFrom(this.svc.getAll()));
    this.isLoading.set(false);
  }
}

// Presentational — pure rendering, ChangeDetection.OnPush
@Component({
  changeDetection: ChangeDetectionStrategy.OnPush,
  // ...
})
export class ItemListComponent {
  @Input() items: Item[] = [];
  @Input() loading = false;
  @Output() itemSelected = new EventEmitter<Item>();
  @Output() itemDeleted = new EventEmitter<string>();
}
```

**All presentational components use `ChangeDetectionStrategy.OnPush`.**

### Skeleton Loading

```html
@if (loading) {
  <ion-list>
    @for (i of [1,2,3,4,5]; track i) {
      <ion-item>
        <ion-label>
          <h2><ion-skeleton-text [animated]="true" style="width: 60%"></ion-skeleton-text></h2>
          <p><ion-skeleton-text [animated]="true" style="width: 80%"></ion-skeleton-text></p>
        </ion-label>
      </ion-item>
    }
  </ion-list>
} @else {
  <!-- actual list -->
}
```

### Pull-to-Refresh + Infinite Scroll

```html
<ion-content>
  <ion-refresher slot="fixed" (ionRefresh)="handleRefresh($event)">
    <ion-refresher-content></ion-refresher-content>
  </ion-refresher>

  <!-- list content -->

  <ion-infinite-scroll (ionInfinite)="loadMore($event)" [disabled]="noMoreData">
    <ion-infinite-scroll-content loadingSpinner="bubbles"></ion-infinite-scroll-content>
  </ion-infinite-scroll>
</ion-content>
```

## Toast + Loading Controllers

```typescript
import { ToastController, LoadingController } from '@ionic/angular/standalone';

private toastCtrl = inject(ToastController);
private loadingCtrl = inject(LoadingController);

async showSuccess(msg: string) {
  const t = await this.toastCtrl.create({ message: msg, duration: 2000, color: 'success', position: 'bottom' });
  await t.present();
}

async showError(msg: string) {
  const t = await this.toastCtrl.create({ message: msg, duration: 3000, color: 'danger', position: 'bottom' });
  await t.present();
}

async withLoading<T>(fn: () => Promise<T>, message = 'Loading...'): Promise<T> {
  const loading = await this.loadingCtrl.create({ message });
  await loading.present();
  try { return await fn(); } finally { await loading.dismiss(); }
}
```

## Performance Rules

- **Lazy loading:** `loadComponent` on all routes — never import pages eagerly.
- **`@for` with `track`:** always track by `item.id` (not index).
- **Virtual scroll:** for lists 100+ items use CDK `cdk-virtual-scroll-viewport`.
- **OnPush:** all presentational components.
- **Icons:** import individually — `import { heart } from 'ionicons/icons'`, never the full library.
- **Bundle:** `import { IonButton } from '@ionic/angular/standalone'` — never `import * as Ionic`.
- **Preloading:** `PreloadAllModules` in NgModule apps; standalone uses `loadComponent` which already defers.

## Common Pitfalls

- **Missing `<ion-page>` wrapper** — every routed page must be wrapped. Without it, transitions break.
- **Using `ngOnInit` for data fetching** — use `ionViewWillEnter`; `ngOnInit` won't re-fire on back-navigation.
- **Icons not showing (standalone)** — call `addIcons()` in constructor + import `IonIcon` from `@ionic/angular/standalone`.
- **`routerLink` not working** — import `RouterLink` from `@angular/router` in standalone component's `imports`.
- **Tab navigation broken** — `tab` attribute on `<ion-tab-button>` must match the child route `path`.
- **Page transitions not animating** — use `NavController`, not Angular `Router`, for animated navigation.

## Resources

- Ionic Angular Overview — https://ionicframework.com/docs/angular/overview
- Ionic Navigation — https://ionicframework.com/docs/angular/navigation
- Ionic Page Life Cycle — https://ionicframework.com/docs/angular/lifecycle
- Ionic Theming — https://ionicframework.com/docs/theming
- Angular Standalone Components — https://angular.dev/guide/components/importing
