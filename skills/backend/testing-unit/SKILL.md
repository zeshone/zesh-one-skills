---
name: net8-apirest-testing-unit
description: >
  Unit testing standards for ASP.NET Core 8 REST APIs using xUnit, FluentAssertions, and NSubstitute (canonical). Covers Services, Validators, Mappings, and Domain Exceptions with AAA pattern, naming conventions, Test Data Builders, and mocking guidance.
  Trigger: When writing unit tests, creating test projects, mocking dependencies, or reviewing test quality in a .NET 8 API. Scope: Services, Validators, Mappings, Domain Exceptions only — for integration or E2E testing load a different skill.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.1"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Creando un nuevo test project o agregando tests a uno existente
- Escribiendo tests para Services, Validators, Mappings, o Domain Exceptions
- Decidiendo cómo mockear dependencias
- Revisando calidad, naming o estructura de tests

---

## Scope Boundaries

### ✅ IN — Testear con esta skill

| Layer | Razón |
|---|---|
| **Services** | Core business logic, exception paths, interacción con repositories |
| **Validators** (FluentValidation) | Cobertura de reglas sin infraestructura |
| **Mappings** (extension methods / AutoMapper profiles) | Pure input → output, sin side effects |
| **Domain Exceptions** | Verificar message y type propagation |

### ❌ OUT — No unit test con esta skill

| Layer | Razón |
|---|---|
| **Repositories** | Dependen de EF Core DbContext → territorio de integration tests |
| **External HTTP clients** | Requieren network o `HttpMessageHandler` stubs |
| **Background jobs / Hosted services** | Dependen del lifecycle de `IHostedService` |

### ⚠️ Controller Unit Tests — Anti-Pattern por default

Unit testing controllers acopla los tests a internals de ASP.NET Core (`ModelState`, `ActionResult` casting, route resolution). Los controllers son thin por diseño — su comportamiento se valida en integration tests.

**Solo justificar si:**
1. El equipo tiene cero cobertura de integration tests, Y
2. El controller tiene lógica condicional no trivial que no puede moverse al service

Incluso entonces: marcar con `// TODO: replace with integration test`.

---

## Critical Patterns

### Naming Convention

| Elemento | Patrón | Ejemplo |
|---|---|---|
| Test class | `{ClassUnderTest}Tests` | `UserServiceTests` |
| Test method | `{Method}_When{Condition}_Should{Expected}` | `GetByIdAsync_WhenUserNotFound_ShouldThrowNotFoundException` |

Nombres deben ser legibles sin abrir el body del test.

### AAA — Arrange / Act / Assert

Comentarios obligatorios cuando el body supera 5 líneas. **Un solo Act por test.**

```csharp
[Fact]
public async Task GetByIdAsync_WhenUserExists_ShouldReturnUserDto()
{
    // Arrange
    var userId = Guid.NewGuid();
    var user = new UserBuilder().WithId(userId).Build();
    _repository.GetByIdAsync(userId).Returns(user);

    // Act
    var result = await _sut.GetByIdAsync(userId);

    // Assert
    result.Should().NotBeNull();
    result.Id.Should().Be(userId);
}
```

### Test Data Builders — Siempre, nunca inline

```csharp
// tests/Shared/Builders/UserBuilder.cs
public class UserBuilder
{
    private Guid _id = Guid.NewGuid();
    private string _email = "test@example.com";
    private string _firstName = "Test";
    private string _lastName = "User";
    private bool _isActive = true;

    public UserBuilder WithId(Guid id) { _id = id; return this; }
    public UserBuilder WithEmail(string email) { _email = email; return this; }
    public UserBuilder Inactive() { _isActive = false; return this; }

    public User Build() => new() { Id = _id, Email = _email, FirstName = _firstName, LastName = _lastName, IsActive = _isActive };
}
```

Reglas:
- El `Build()` sin args debe producir una entidad válida
- Nunca crear entidades inline en los tests — siempre Builder
- Un Builder por tipo de entidad o DTO

### Mocking — NSubstitute (canónico) vs Moq

NSubstitute es el canonical. Moq es alternativa para equipos que ya lo usan — elegir uno y ser consistente en el proyecto.

| Operación | NSubstitute (canónico) | Moq (alternativa) |
|---|---|---|
| Crear substitute | `Substitute.For<IRepo>()` | `new Mock<IRepo>()` |
| Acceder al objeto | *(el substitute ya es el objeto)* | `mock.Object` |
| Stub return | `repo.Method(Arg.Any<Guid>()).Returns(user)` | `mock.Setup(r => r.Method(It.IsAny<Guid>())).ReturnsAsync(user)` |
| Verificar llamada | `await repo.Received(1).Method(id)` | `mock.Verify(r => r.Method(id), Times.Once)` |
| Stub exception | `.Throws(new Ex())` | `.ThrowsAsync(new Ex())` |

> **Moq SponsorLink**: Moq v4.20+ tiene telemetría SponsorLink. Pinear a `4.18.x` o usar NSubstitute.

Setup en constructor — substitute compartido entre todos los tests de la clase:

```csharp
public class UserServiceTests
{
    private readonly IUserRepository _repository;
    private readonly UserService _sut;

    public UserServiceTests()
    {
        _repository = Substitute.For<IUserRepository>();
        _sut = new UserService(_repository);
    }
}
```

---

## Anti-Patterns

| Anti-pattern | Problema |
|---|---|
| Testear implementation details en vez de behavior | Los tests rompen en cada refactor aunque el comportamiento no cambie |
| Over-mocking — mockear value objects y DTOs | El test no ejercita lógica real; falsa confianza |
| Múltiples Acts en un test | Imposible saber qué Act causó el fallo |
| Nombre sin contexto (`Test1`, `ShouldWork`) | Los failure reports son inútiles |
| Unit tests de repositories con EF Core | DbContext requiere base de datos — territorio de integration |
| Controller unit tests como default | Acopla tests a internals de ASP.NET; preferir integration tests |
| `Assert.Equal` en vez de FluentAssertions | Menor legibilidad y diagnóstico en fallos |

---

## Commands

```bash
dotnet new xunit -n MyApi.Tests --framework net8.0
dotnet add MyApi.Tests/MyApi.Tests.csproj reference src/MyApi/MyApi.csproj
dotnet add MyApi.Tests package NSubstitute
dotnet add MyApi.Tests package FluentAssertions
dotnet test --verbosity normal
dotnet test --filter "FullyQualifiedName~UserServiceTests"
```

---

## Resources

- **Architecture and conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **NSubstitute docs**: https://nsubstitute.github.io/
- **FluentAssertions docs**: https://fluentassertions.com/introduction
- **xUnit docs**: https://xunit.net/docs/getting-started/netcore/cmdline

---

## Changelog

### v1.1 — 2026-03-28
- **Removed**: Ejemplos completos de tests por layer (Service, Validator, Mapping, Exception) — el agente sabe escribir tests una vez conoce las convenciones
- **Removed**: Full CRUD service test class example de 60 líneas
- **Removed**: Sección de project setup con folder structure detallada
- **Kept**: Scope boundaries (decisión de qué testear), naming convention, AAA, Test Data Builders, NSubstitute vs Moq comparison, anti-patterns
