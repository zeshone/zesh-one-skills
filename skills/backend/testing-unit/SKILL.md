---
name: net8-apirest-testing-unit
description: >
  Unit testing standards for ASP.NET Core 8 REST APIs using xUnit, FluentAssertions, and NSubstitute (canonical). Covers Services, Validators, Mappings, and Domain Exceptions with AAA pattern, naming conventions, Test Data Builders, and mocking guidance.
  Trigger: When writing unit tests, creating test projects, mocking dependencies, or reviewing test quality in a .NET 8 API. Scope: Services, Validators, Mappings, Domain Exceptions only — for integration or E2E testing load a different skill.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.2"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Creating a new test project or adding tests to an existing one
- Writing tests for Services, Validators, Mappings, or Domain Exceptions
- Deciding how to mock dependencies
- Reviewing test quality, naming, or structure

---

## Scope Boundaries

### ✅ IN — Test with this skill

| Layer | Reason |
|---|---|
| **Services** | Core business logic, exception paths, interaction with repositories |
| **Validators** (FluentValidation) | Rule coverage without infrastructure |
| **Mappings** (extension methods / AutoMapper profiles) | Pure input → output, no side effects |
| **Domain Exceptions** | Verify message and type propagation |

### ❌ OUT — Do not unit test with this skill

| Layer | Reason |
|---|---|
| **Repositories** | Depend on EF Core DbContext → territory of integration tests |
| **External HTTP clients** | Require network or `HttpMessageHandler` stubs |
| **Background jobs / Hosted services** | Depend on `IHostedService` lifecycle |

### ⚠️ Controller Unit Tests — Anti-Pattern by default

Unit testing controllers couples tests to ASP.NET Core internals (`ModelState`, `ActionResult` casting, route resolution). Controllers are thin by design — their behavior is validated in integration tests.

**Only justified if:**
1. The team has zero integration test coverage, AND
2. The controller has non-trivial conditional logic that cannot be moved to the service

Even then: mark with `// TODO: replace with integration test`.

---

## Critical Patterns

### Naming Convention

| Element | Pattern | Example |
|---|---|---|
| Test class | `{ClassUnderTest}Tests` | `UserServiceTests` |
| Test method | `{Method}_When{Condition}_Should{Expected}` | `GetByIdAsync_WhenUserNotFound_ShouldThrowNotFoundException` |

Names must be readable without opening the test body.

### AAA — Arrange / Act / Assert

Comments are mandatory when the body exceeds 5 lines. **One single Act per test.**

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

### Test Data Builders — Always, never inline

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

Rules:
- `Build()` with no args must produce a valid entity
- Never create entities inline in tests — always use Builder
- One Builder per entity or DTO type

### Mocking — NSubstitute (canonical) vs Moq

NSubstitute is the canonical choice. Moq is an alternative for teams already using it — pick one and stay consistent throughout the project.

| Operation | NSubstitute (canonical) | Moq (alternative) |
|---|---|---|
| Create substitute | `Substitute.For<IRepo>()` | `new Mock<IRepo>()` |
| Access the object | *(the substitute is already the object)* | `mock.Object` |
| Stub return | `repo.Method(Arg.Any<Guid>()).Returns(user)` | `mock.Setup(r => r.Method(It.IsAny<Guid>())).ReturnsAsync(user)` |
| Verify call | `await repo.Received(1).Method(id)` | `mock.Verify(r => r.Method(id), Times.Once)` |
| Stub exception | `.ThrowsAsync(new Ex())` — async methods; `.Throws(new Ex())` — sync methods | `.ThrowsAsync(new Ex())` |

> **Moq SponsorLink**: Moq v4.20+ has SponsorLink telemetry. Pin to `4.18.x` or use NSubstitute.

Setup in constructor — substitute shared across all tests in the class:

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

| Anti-pattern | Problem |
|---|---|
| Testing implementation details instead of behavior | Tests break on every refactor even when behavior does not change |
| Over-mocking — mocking value objects and DTOs | The test does not exercise real logic; false confidence |
| Multiple Acts in one test | Impossible to tell which Act caused the failure |
| Name without context (`Test1`, `ShouldWork`) | Failure reports are useless |
| Repository unit tests with EF Core | DbContext requires a database — territory of integration |
| Controller unit tests as default | Couples tests to ASP.NET internals; prefer integration tests |
| `Assert.Equal` instead of FluentAssertions | Lower readability and diagnostics on failures |

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

### v1.2 — 2026-03-28
- **Fixed (W-12)**: NSubstitute stub exception column now distinguishes `.ThrowsAsync(new Ex())` for async methods and `.Throws(new Ex())` for sync methods. The previous entry only showed `.Throws()`, which silently passes for sync but does NOT stub async Task-returning methods correctly.

### v1.1 — 2026-03-28
- **Removed**: Full test examples by layer (Service, Validator, Mapping, Exception) — the agent knows how to write tests once it knows the conventions
- **Removed**: Full CRUD service test class example of 60 lines
- **Removed**: Project setup section with detailed folder structure
- **Kept**: Scope boundaries (decision on what to test), naming convention, AAA, Test Data Builders, NSubstitute vs Moq comparison, anti-patterns
