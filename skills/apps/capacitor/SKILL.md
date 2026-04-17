---
name: capacitor
description: >
  ZeshOne Capacitor conventions for mobile native layer.
  Trigger: When using Capacitor plugins, platform detection, native APIs (camera, push, filesystem, preferences), capacitor.config.ts, or build/deploy to iOS/Android.
license: Apache-2.0
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
metadata:
  author: Zesh-One
  version: "1.0"
  inspired-by: capawesome-team/ionic-capacitor
---

## When to Use

Load this skill when working with the Capacitor native layer — plugins, platform detection, configuration, native APIs, or building and deploying to iOS/Android.

## Critical Patterns

- Run `ionic cap sync` after every Capacitor plugin install or web build that must reach the native projects.
- Use `Capacitor.isNativePlatform()` and `Capacitor.getPlatform()` for platform checks — never `navigator.userAgent`.
- Wrap native plugin calls in `try/catch` and treat user cancellation differently from real failures.
- Initialize push notifications after authentication so token registration and routing have user context.
- Clean up long-lived native listeners and watchers (`Network`, `Geolocation`, app lifecycle) when the view or feature stops using them.

## Capacitor Config

```typescript
// capacitor.config.ts
import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.yourcompany.appname',   // Reverse domain — must match App Store/Play Store
  appName: 'Your App Name',
  webDir: 'dist',                      // Build output — match your build tool's outDir
  server: {
    androidScheme: 'https',            // Always https in production
    // Dev live reload (remove in prod):
    // url: 'http://192.168.x.x:8100',
    // cleartext: true,
  },
  ios: {
    contentInset: 'automatic',         // Safe area handling
    preferredContentMode: 'mobile',
  },
  android: {
    allowMixedContent: false,          // Keep false in production
  },
  plugins: {
    SplashScreen: {
      launchAutoHide: true,
      androidScaleType: 'CENTER_CROP',
      splashFullScreen: true,
      splashImmersive: true,
    },
    StatusBar: {
      style: 'dark',                   // 'dark' | 'light'
    },
  },
};

export default config;
```

## Plugin Install Pattern

```bash
npm install @capacitor/<plugin>
ionic cap sync              # REQUIRED after every plugin install — syncs to native projects
```

**Never skip `ionic cap sync`.** `npm install` alone doesn't update native projects.

## Platform Detection

```typescript
import { Capacitor } from '@capacitor/core';

// Running on native device (iOS/Android) vs browser
if (Capacitor.isNativePlatform()) {
  // Use native Capacitor plugin
} else {
  // Web fallback (file input, etc.)
}

// Specific platform: 'ios' | 'android' | 'web'
const platform = Capacitor.getPlatform();
if (platform === 'ios') { /* iOS-specific */ }
if (platform === 'android') { /* Android-specific */ }
```

**Never use `navigator.userAgent` for platform detection.** Use `Capacitor.getPlatform()`.

## Error Handling for Native APIs

Always wrap Capacitor calls in try/catch. Distinguish user cancellation from real errors:

```typescript
const takePhoto = async () => {
  try {
    const image = await Camera.getPhoto({ quality: 90, resultType: CameraResultType.Uri });
    return image.webPath;
  } catch (error: unknown) {
    if (error instanceof Error && error.message?.includes('cancelled')) {
      return null; // User cancelled — not an error
    }
    throw error; // Real error — propagate
  }
};
```

## Core Native Plugins

### Camera

```typescript
import { Camera, CameraResultType, CameraSource } from '@capacitor/camera';

const photo = await Camera.getPhoto({
  quality: 90,
  allowEditing: false,
  resultType: CameraResultType.Uri,   // Uri → use webPath for <img src>
  source: CameraSource.Prompt,        // Camera | Photos | Prompt
  width: 1024,
  saveToGallery: false,
});
// photo.webPath — use directly in <img src>
// photo.base64String — if resultType: CameraResultType.Base64
```

### Preferences (Key-Value Storage)

```typescript
import { Preferences } from '@capacitor/preferences';

// Save
await Preferences.set({ key: 'token', value: authToken });
await Preferences.set({ key: 'user', value: JSON.stringify(user) });

// Read
const { value } = await Preferences.get({ key: 'token' });
const user = value ? JSON.parse(await Preferences.get({ key: 'user' }).then(r => r.value!)) : null;

// Remove
await Preferences.remove({ key: 'token' });

// Clear all
await Preferences.clear();
```

Use `Preferences` for session tokens, user preferences, and lightweight app state. Not for large data — use `Filesystem` for files.

### Push Notifications

```typescript
import { PushNotifications } from '@capacitor/push-notifications';

async function initPush({
  registerToken,
  navigateForward,
}: {
  registerToken: (token: string) => Promise<void>;
  navigateForward: (route: string) => Promise<boolean> | void;
}) {
  const permStatus = await PushNotifications.requestPermissions();
  if (permStatus.receive !== 'granted') return;

  await PushNotifications.register();

  PushNotifications.addListener('registration', async ({ value: token }) => {
    // Send token to your .NET API
    await registerToken(token);
  });

  PushNotifications.addListener('pushNotificationReceived', notification => {
    // App is in foreground — show in-app alert or toast
  });

  PushNotifications.addListener('pushNotificationActionPerformed', async action => {
    // User tapped notification — navigate based on action.notification.data
    const route = action.notification.data?.route;
    if (route) await navigateForward(route);
  });
}
```

Call `initPush()` in the app root component after authentication, not at app start (permissions need user context).

### Filesystem

```typescript
import { Filesystem, Directory, Encoding } from '@capacitor/filesystem';

// Write
await Filesystem.writeFile({
  path: 'data/config.json',
  data: JSON.stringify(config),
  directory: Directory.Data,
  encoding: Encoding.UTF8,
  recursive: true,               // Create dirs as needed
});

// Read
const result = await Filesystem.readFile({
  path: 'data/config.json',
  directory: Directory.Data,
  encoding: Encoding.UTF8,
});
const config = JSON.parse(result.data as string);

// Delete
await Filesystem.deleteFile({ path: 'data/config.json', directory: Directory.Data });
```

**Directories:** `Directory.Data` (app-private, persists), `Directory.Cache` (may be cleared by OS), `Directory.Documents` (user-accessible on iOS).

### Geolocation

```typescript
import { Geolocation } from '@capacitor/geolocation';

const getPosition = async () => {
  const coords = await Geolocation.getCurrentPosition({ enableHighAccuracy: true, timeout: 10000 });
  return { lat: coords.coords.latitude, lng: coords.coords.longitude };
};

// Watch — returns watchId for cleanup
const watchId = await Geolocation.watchPosition({ enableHighAccuracy: true }, (position, err) => {
  if (err) return;
  updateMap(position.coords);
});

// Cleanup
await Geolocation.clearWatch({ id: watchId });
```

### Network Status

```typescript
import { Network } from '@capacitor/network';

// One-time check
const status = await Network.getStatus();
if (!status.connected) showOfflineBanner();

// Listen for changes
const listener = await Network.addListener('networkStatusChange', status => {
  if (status.connected) hideOfflineBanner();
  else showOfflineBanner();
});

// Cleanup (call in ngOnDestroy or ionViewWillLeave)
await listener.remove();
```

### Haptics

```typescript
import { Haptics, ImpactStyle } from '@capacitor/haptics';

// Feedback patterns
await Haptics.impact({ style: ImpactStyle.Medium }); // Light | Medium | Heavy
await Haptics.notification({ type: 'success' });       // success | warning | error
await Haptics.vibrate();                               // Default vibration
```

Use haptics for confirmations, errors, and interactive feedback. Do NOT overuse.

### Other Core Plugins

| Plugin | Install | Key Methods |
|--------|---------|-------------|
| `@capacitor/app` | `npm i @capacitor/app` | `App.addListener('appStateChange')`, `App.getLaunchUrl()` |
| `@capacitor/browser` | `npm i @capacitor/browser` | `Browser.open({ url })` |
| `@capacitor/clipboard` | `npm i @capacitor/clipboard` | `Clipboard.write()`, `Clipboard.read()` |
| `@capacitor/device` | `npm i @capacitor/device` | `Device.getInfo()`, `Device.getId()` |
| `@capacitor/keyboard` | `npm i @capacitor/keyboard` | `Keyboard.show()`, `Keyboard.hide()` |
| `@capacitor/share` | `npm i @capacitor/share` | `Share.share({ title, text, url })` |
| `@capacitor/splash-screen` | `npm i @capacitor/splash-screen` | `SplashScreen.hide()` |
| `@capacitor/status-bar` | `npm i @capacitor/status-bar` | `StatusBar.setStyle()`, `StatusBar.hide()` |
| `@capacitor/local-notifications` | `npm i @capacitor/local-notifications` | `LocalNotifications.schedule()` |

## App Lifecycle (Foreground/Background)

```typescript
import { App } from '@capacitor/app';

// In AppComponent or root service
App.addListener('appStateChange', ({ isActive }) => {
  if (isActive) {
    // App came to foreground — refresh data, restart timers
  } else {
    // App went to background — save state, pause media
  }
});

// Deep link handling
App.addListener('appUrlOpen', ({ url }) => {
  const slug = new URL(url).pathname;
  if (slug) this.navCtrl.navigateForward(slug);
});
```

## Theming

```css
/* src/theme/variables.scss */
:root {
  /* Primary brand color — generates shade/tint automatically */
  --ion-color-primary: #3880ff;
  --ion-color-primary-rgb: 56, 128, 255;
  --ion-color-primary-contrast: #ffffff;
  --ion-color-primary-shade: #3171e0;
  --ion-color-primary-tint: #4c8dff;

  --ion-background-color: #ffffff;
  --ion-text-color: #000000;
  --ion-font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
```

**Dark mode** — pick one approach:
```css
@import '@ionic/angular/css/palettes/dark.system.css'; /* Follow system (recommended) */
@import '@ionic/angular/css/palettes/dark.class.css';  /* Toggle via .ion-palette-dark on <html> */
```

```html
<!-- Required for native status bar / OS UI elements -->
<meta name="color-scheme" content="light dark" />
```

**Per-component styling** via CSS Shadow Parts and CSS variables:
```css
ion-button::part(native) { border-radius: 20px; }
ion-item { --background: transparent; --padding-start: 16px; }
```

## Development Workflow

```bash
ionic serve                           # Browser preview + hot reload
ionic cap run ios -l --external       # Live reload on iOS device (same WiFi)
ionic cap run android -l --external   # Live reload on Android device

ionic build && ionic cap sync         # Build web + sync to native
ionic cap open ios                    # Open Xcode
ionic cap open android                # Open Android Studio
```

**Debugging:**
- Browser: `ionic serve` → Chrome DevTools
- iOS: Safari > Develop > [Device Name] > [App]
- Android: `chrome://inspect` → Remote Targets

## Asset Generation (Icons + Splash)

```bash
npm install @capacitor/assets --save-dev
# Place logo at assets/logo.png (1024x1024px min)
npx @capacitor/assets generate \
  --iconBackgroundColor '#ffffff' \
  --iconBackgroundColorDark '#111111' \
  --splashBackgroundColor '#ffffff' \
  --splashBackgroundColorDark '#111111'
```

## Deployment

**iOS:** Apple Developer ($99/yr) → Xcode signing → `Product > Archive` → upload via Xcode Organizer.

**Android:** Google Play ($25 one-time) → `Build > Generate Signed Bundle (AAB)` → upload to Play Console.

## Common Pitfalls

- **Not running `ionic cap sync`** after plugin install — native side won't have the plugin.
- **Hardcoding platform checks** via user-agent — use `Capacitor.getPlatform()`.
- **Ignoring safe areas** — use `<ion-content fullscreen>` and Ionic's built-in safe area handling. Don't manually pad for notches.
- **Using browser APIs directly** — `navigator.camera` doesn't exist natively. Always use `@capacitor/camera`.
- **Not testing on real devices** — emulators miss touch performance, native quirks, and plugin behavior.
- **Not handling `allowMixedContent: false`** — all API calls in production must use HTTPS.

## Resources

- Capacitor Docs — https://capacitorjs.com/docs
- Capacitor Config Reference — https://capacitorjs.com/docs/config
- Capacitor Push Notifications — https://capacitorjs.com/docs/apis/push-notifications
- Capacitor iOS Deployment — https://capacitorjs.com/docs/ios/deploying-to-app-store
- Capacitor Android Deployment — https://capacitorjs.com/docs/android/deploying-to-google-play
