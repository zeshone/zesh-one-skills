---
name: net8-apirest-performance
description: >
  Performance, concurrency and resilience rules for ASP.NET Core 8 REST APIs covering EF Core optimization, async patterns, Polly resilience pipelines, and parallelism best practices.
  Trigger: When optimizing API performance, configuring EF Core, implementing resilience with Polly, or working with concurrent operations in a .NET 8 API.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Configurando EF Core para escenarios de alta concurrencia
- Implementando resilience pipelines con Polly (retry, circuit breaker, timeout)
- Paralerizando operaciones async independientes
- Optimizando queries de base de datos

---

## Critical Patterns

### Circuit Breaker — Per-Request, No Global

El circuit breaker debe ser **per-request scoped**, no global. Un circuit breaker global genera false positives bajo concurrencia — un request lento abre el circuito para todos los demás.

```csharp
builder.Services.AddResiliencePipeline("external-api", pipeline =>
{
    pipeline
        .AddRetry(new RetryStrategyOptions
        {
            ShouldHandle = new PredicateBuilder().Handle<HttpRequestException>(),
            MaxRetryAttempts = 3,
            Delay = TimeSpan.FromMilliseconds(500),
            BackoffType = DelayBackoffType.Exponential
        })
        .AddCircuitBreaker(new CircuitBreakerStrategyOptions
        {
            ShouldHandle = new PredicateBuilder().Handle<HttpRequestException>(),
            FailureRatio = 0.5,
            MinimumThroughput = 10,          // previene false positives
            SamplingDuration = TimeSpan.FromSeconds(30),
            BreakDuration = TimeSpan.FromSeconds(15)
        })
        .AddTimeout(TimeSpan.FromSeconds(10));
});
```

Uso en service via inyección:
```csharp
public class ExternalPaymentService
{
    private readonly ResiliencePipeline _pipeline;
    private readonly HttpClient _httpClient;

    public ExternalPaymentService(ResiliencePipelineProvider<string> provider, HttpClient httpClient)
    {
        _pipeline = provider.GetPipeline("external-api");
        _httpClient = httpClient;
    }

    public async Task<PaymentResult> ProcessAsync(PaymentRequest request) =>
        await _pipeline.ExecuteAsync(async ct =>
        {
            var response = await _httpClient.PostAsJsonAsync("/payments", request, ct);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<PaymentResult>(ct);
        });
}
```

### Kestrel — Production Configuration

```csharp
builder.WebHost.ConfigureKestrel(options =>
{
    options.Limits.MaxConcurrentConnections = 1000;
    options.Limits.MaxRequestBodySize = 10 * 1024 * 1024; // 10 MB
    options.Limits.RequestHeadersTimeout = TimeSpan.FromSeconds(15);
});
```

### Response Caching — Para datos estables

```csharp
// Output caching (ASP.NET Core 7+) — preferido sobre ResponseCache
builder.Services.AddOutputCache(options =>
{
    options.AddPolicy("products", b => b.Expire(TimeSpan.FromMinutes(5)));
});
app.UseOutputCache();

[HttpGet]
[OutputCache(PolicyName = "products")]
public async Task<IActionResult> GetProducts() { ... }
```

### `Task.WhenAll` — Operaciones independientes en paralelo

```csharp
// CORRECT — paralelo
var userTask = _userRepository.GetByIdAsync(userId);
var ordersTask = _orderRepository.GetRecentByUserAsync(userId, limit: 10);
await Task.WhenAll(userTask, ordersTask);

// WRONG — secuencial sin necesidad
var user = await _userRepository.GetByIdAsync(userId);
var orders = await _orderRepository.GetRecentByUserAsync(userId, limit: 10);
```

---

## Anti-Patterns

| Anti-pattern | Problema |
|---|---|
| Circuit Breaker global | False positives bajo concurrencia — un request lento afecta a todos |
| `AspNetCoreRateLimit` | Conflictos de concurrencia con MemoryCache — usar built-in rate limiter |
| `MemoryCache` para estado compartido | Interferencia de concurrencia — usar distributed cache o scoping correcto |

---

## Commands

```bash
dotnet add package Polly
dotnet add package Microsoft.Extensions.Http.Resilience
dotnet add package Microsoft.AspNetCore.OutputCaching
```

---

## Resources

- **General conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **Security / Rate Limiting**: See [../security/SKILL.md](../security/SKILL.md)
- **Data access / DbContextPool / AsNoTracking / Compiled Queries**: See [`../dataaccess/SKILL.md`](../dataaccess/SKILL.md)

---

## Changelog

### v1.2 — 2026-03-28
- **Removed**: `AddDbContextPool` y `AsNoTracking` (ya cubiertos en `dataaccess/SKILL.md`)
- **Removed**: Rate limiting (ya cubierto en `security/SKILL.md`)
- **Removed**: GUIDs para IDs (ya en `security/SKILL.md`)
- **Removed**: Compiled Queries example (ya en `dataaccess/SKILL.md`)
- **Kept**: Circuit breaker per-request (decisión de diseño, no doc estándar), Kestrel config, response caching, `Task.WhenAll` reminder
- **Added**: Énfasis en la razón del circuit breaker per-request (false positives bajo concurrencia)

### v1.1 — 2026-03-24
- Fixed: `GetConnectionString("Default")` → `"DefaultConnection"`
- Added: Cross-refs a `security`, `dataaccess`, `requests`
