---
name: net8-apirest-testing-unit
description: >
  Unit testing standards for ASP.NET Core 8 REST APIs using xUnit, FluentAssertions, and NSubstitute (canonical). Covers Services, Validators, Mappings, and Domain Exceptions with AAA pattern, naming conventions, Test Data Builders, and mocking guidance.
  Trigger: When writing unit tests, creating test projects, mocking dependencies, or reviewing test quality in a .NET 8 API. Scope: Services, Validators, Mappings, Domain Exceptions only — for integration or E2E testing load a different skill.
license: Apache-2.0
metadata:
  author: Zesh-One
  version: "1.0"
allowed-tools: Read, Edit, Write, Glob, Grep
---

## When to Use

- Creating a new unit test project or adding tests to an existing one
- Writing tests for Services, Validators, Mappings, or Domain Exceptions
- Deciding how to mock dependencies (NSubstitute vs Moq)
- Reviewing test quality, naming, or structure
- Generating Test Data Builders for complex entities

> **v1 scope note**: This skill covers unit tests only. Integration tests (WebApplicationFactory, TestContainers), repository tests (EF Core InMemory), and E2E tests are out of scope for this skill.

---

## Scope Boundaries

Define **what to test unitarily** and what to leave to other testing strategies.

### ✅ IN — Test with this skill

| Layer | Reason |
|---|---|
| **Services** | Core business logic, exception paths, repository interaction |
| **Validators** (FluentValidation) | Pure rule coverage without infrastructure |
| **Mappings** (extension methods / AutoMapper profiles) | Pure input → output, no side effects |
| **Domain Exceptions** | Verify message and type propagation |

### ❌ OUT — Do not unit test with this skill

| Layer | Reason |
|---|---|
| **Repositories** | Depend on EF Core DbContext → integration territory |
| **External HTTP clients** | Require network or `HttpMessageHandler` stubs → integration territory |
| **Background jobs / Hosted services** | Depend on IHostedService lifecycle → integration territory |

### ⚠️ Controller Unit Tests — Anti-Pattern by Default

Unit testing controllers couples the test to ASP.NET Core internals (`ModelState`, `ActionResult` casting, route resolution). Controllers are designed to be a **thin layer** (see `../general/SKILL.md`) and their behavior is best validated via integration tests.

**Only justify controller unit tests when ALL of these apply:**
1. The team has zero integration test coverage yet, AND
2. The controller contains non-trivial conditional logic that cannot be moved to a service

**Even then**: mark the test file with `// TODO: replace with integration test` and migrate as soon as WebApplicationFactory is available. See `../general/SKILL.md` for controller design rules.

---

## Critical Patterns

### Project Setup

```bash
# Create test project alongside the API project
dotnet new xunit -n MyApi.Tests --framework net8.0

# Add project reference
dotnet add MyApi.Tests/MyApi.Tests.csproj reference src/MyApi/MyApi.csproj

# Required packages
dotnet add MyApi.Tests package NSubstitute
dotnet add MyApi.Tests package FluentAssertions
dotnet add MyApi.Tests package FluentValidation
```

Recommended folder structure — mirrors the main project's Vertical Slice layout:

```
tests/
  MyApi.Tests/
    Features/
      Users/
        Services/
          UserServiceTests.cs
        Validators/
          CreateUserRequestValidatorTests.cs
        Mappings/
          UserMappingTests.cs
    Shared/
      Builders/
        UserBuilder.cs
        CreateUserRequestBuilder.cs
```

---

### Naming Convention

| Element | Pattern | Example |
|---|---|---|
| Test class | `{ClassUnderTest}Tests` | `UserServiceTests` |
| Test method | `{Method}_When{Condition}_Should{Expected}` | `GetByIdAsync_WhenUserNotFound_ShouldThrowNotFoundException` |

**Rules:**
- Names must be readable without opening the test body
- `When` describes the input state or scenario
- `Should` describes the observable outcome (not the internal step)
- Avoid vague names like `Test1`, `ShouldWork`, `CreateUser_Test`

---

### AAA Pattern

Every test follows **Arrange → Act → Assert**. Comments are mandatory when the test body exceeds 5 lines.

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

**Rules:**
- **One Act per test** — never chain two calls and assert on both
- Assert **only on the outcome of Act**, never on Arrange state
- Prefer `[Fact]` for single scenarios; use `[Theory]` + `[InlineData]` for parametric cases

---

### Test Data Builders

Use the Builder pattern to create entities and DTOs without repetition. Builders live in `tests/Shared/Builders/`.

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
    public UserBuilder WithFirstName(string name) { _firstName = name; return this; }
    public UserBuilder Inactive() { _isActive = false; return this; }

    public User Build() => new()
    {
        Id = _id,
        Email = _email,
        FirstName = _firstName,
        LastName = _lastName,
        IsActive = _isActive
    };
}
```

Usage:
```csharp
var user = new UserBuilder().WithEmail("admin@example.com").Inactive().Build();
var request = new CreateUserRequestBuilder().WithEmail("new@example.com").Build();
```

**Rules:**
- Default values must be valid — the zero-arg `Build()` must produce a passing entity
- Never create entities inline in tests — always use a Builder
- One Builder per entity or DTO type

---

### Mocking with NSubstitute (Default)

`NSubstitute` is the canonical mocking library for this ecosystem. Use it to substitute direct dependencies of the SUT.

**Setup in constructor — shared substitute across all tests in the class:**

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

**Stub a return value:**
```csharp
_repository.GetByIdAsync(Arg.Any<Guid>()).Returns(user);
_repository.GetByIdAsync(Arg.Any<Guid>()).Returns((User?)null);
```

**Verify an interaction:**
```csharp
await _repository.Received(1).GetByIdAsync(userId);
await _repository.DidNotReceive().SaveAsync(Arg.Any<User>());
```

**Stub a thrown exception:**
```csharp
_repository.GetByIdAsync(Arg.Any<Guid>())
    .Throws(new InvalidOperationException("DB error"));
```

**Rules:**
- Mock **only direct dependencies** of the SUT — never mock two levels deep
- Never mock the SUT itself
- Never mock value objects or simple DTOs — use builders instead
- Prefer `Arg.Any<T>()` for generic stubs; use exact values only when the test verifies routing logic

---

### Alternative: Moq

For teams already using Moq or with existing test suites built on it. The patterns are equivalent — choose one and stay consistent within a project.

| Operation | NSubstitute (canonical) | Moq (alternative) |
|---|---|---|
| Create substitute | `Substitute.For<IRepo>()` | `new Mock<IRepo>()` |
| Access mock object | *(substitute is already the object)* | `mock.Object` |
| Stub return | `repo.GetByIdAsync(Arg.Any<Guid>()).Returns(user)` | `mock.Setup(r => r.GetByIdAsync(It.IsAny<Guid>())).ReturnsAsync(user)` |
| Verify call | `await repo.Received(1).GetByIdAsync(id)` | `mock.Verify(r => r.GetByIdAsync(id), Times.Once)` |
| Stub exception | `.Throws(new Ex())` | `.ThrowsAsync(new Ex())` |

> **Note on Moq versions**: Moq introduced SponsorLink telemetry in v4.20. Pin to `4.18.x` or switch to NSubstitute to avoid the issue. See [Moq SponsorLink notice](https://github.com/devlooped/moq/issues/1372).

---

### Tests by Layer

#### Service Tests — Business Logic and Exception Paths

Test the logic and the exceptions. Always verify interactions with mocked dependencies.

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

    [Fact]
    public async Task GetByIdAsync_WhenUserNotFound_ShouldThrowNotFoundException()
    {
        // Arrange
        var userId = Guid.NewGuid();
        _repository.GetByIdAsync(userId).Returns((User?)null);

        // Act
        var act = async () => await _sut.GetByIdAsync(userId);

        // Assert
        await act.Should().ThrowAsync<NotFoundException>()
            .WithMessage($"*{userId}*");
    }

    [Fact]
    public async Task CreateAsync_WhenEmailAlreadyExists_ShouldThrowConflictException()
    {
        // Arrange
        var request = new CreateUserRequestBuilder().WithEmail("dupe@example.com").Build();
        _repository.ExistsByEmailAsync(request.Email).Returns(true);

        // Act
        var act = async () => await _sut.CreateAsync(request);

        // Assert
        await act.Should().ThrowAsync<ConflictException>();
        await _repository.DidNotReceive().SaveAsync(Arg.Any<User>());
    }

    [Fact]
    public async Task CreateAsync_WhenValid_ShouldCallSaveAndReturnDto()
    {
        // Arrange
        var request = new CreateUserRequestBuilder().Build();
        _repository.ExistsByEmailAsync(Arg.Any<string>()).Returns(false);

        // Act
        var result = await _sut.CreateAsync(request);

        // Assert
        await _repository.Received(1).SaveAsync(Arg.Any<User>());
        result.Should().NotBeNull();
        result.Email.Should().Be(request.Email);
    }
}
```

#### Validator Tests — Pure Rule Coverage

Validators are value-pure: no mocks needed for sync rules. Use `[Theory]` for boundary conditions.

```csharp
public class CreateUserRequestValidatorTests
{
    private readonly CreateUserRequestValidator _validator = new();

    [Theory]
    [InlineData("")]
    [InlineData(null)]
    [InlineData("not-an-email")]
    public void Email_WhenInvalid_ShouldHaveValidationError(string? email)
    {
        // Arrange
        var request = new CreateUserRequestBuilder().WithEmail(email!).Build();

        // Act
        var result = _validator.Validate(request);

        // Assert
        result.IsValid.Should().BeFalse();
        result.Errors.Should().Contain(e => e.PropertyName == "Email");
    }

    [Fact]
    public void Validate_WhenAllFieldsValid_ShouldPass()
    {
        // Arrange
        var request = new CreateUserRequestBuilder().Build();

        // Act
        var result = _validator.Validate(request);

        // Assert
        result.IsValid.Should().BeTrue();
    }
}
```

**For `MustAsync` validators** — mock only the async dependency:

```csharp
public class CreateUserRequestValidatorTests
{
    private readonly IUserRepository _repository;
    private readonly CreateUserRequestValidator _validator;

    public CreateUserRequestValidatorTests()
    {
        _repository = Substitute.For<IUserRepository>();
        _validator = new CreateUserRequestValidator(_repository);
    }

    [Fact]
    public async Task Email_WhenAlreadyRegistered_ShouldHaveValidationError()
    {
        // Arrange
        var request = new CreateUserRequestBuilder().WithEmail("taken@example.com").Build();
        _repository.ExistsByEmailAsync("taken@example.com").Returns(true);

        // Act
        var result = await _validator.ValidateAsync(request);

        // Assert
        result.IsValid.Should().BeFalse();
        result.Errors.Should().Contain(e => e.PropertyName == "Email");
    }
}
```

#### Mapping Tests — Pure Input/Output

Mappings are pure functions. No mocks, no async, just verify property transformation.

```csharp
public class UserMappingTests
{
    [Fact]
    public void ToDto_WhenValidUser_ShouldMapAllProperties()
    {
        // Arrange
        var user = new UserBuilder()
            .WithId(Guid.NewGuid())
            .WithEmail("map@example.com")
            .WithFirstName("Map")
            .Build();

        // Act
        var dto = user.ToDto();  // extension method

        // Assert
        dto.Id.Should().Be(user.Id);
        dto.Email.Should().Be(user.Email);
        dto.FirstName.Should().Be(user.FirstName);
    }
}
```

#### Domain Exception Tests — Type and Message

```csharp
public class NotFoundExceptionTests
{
    [Fact]
    public void Constructor_WhenCalled_ShouldIncludeResourceAndId()
    {
        // Arrange
        var id = Guid.NewGuid();

        // Act
        var ex = new NotFoundException("User", id);

        // Assert
        ex.Message.Should().Contain("User");
        ex.Message.Should().Contain(id.ToString());
    }
}
```

---

## Anti-Patterns

| Anti-pattern | Problem |
|---|---|
| Testing implementation details instead of behavior | Tests break on every refactor, even when behavior is unchanged |
| Over-mocking — mocking everything including value objects | The test doesn't exercise real logic; false confidence |
| Multiple Acts in one test | Impossible to know which Act caused the failure |
| No descriptive method name (`Test1`, `ShouldWork`) | Failure reports are useless without context |
| Testing getters and setters | Zero business value; noise in coverage metrics |
| Unit testing repositories with EF Core | DbContext requires a database — this is integration territory |
| Controller unit tests as the default | Couples tests to ASP.NET internals; prefer integration tests (see Scope Boundaries) |
| Coverage gaming — writing tests to hit percentages | Produces low-value tests that satisfy a number, not the behavior |
| Using `Assert.Equal` instead of FluentAssertions | Reduced readability; harder to diagnose failures |
| Shared mutable state between tests | Tests pollute each other; create substitutes per test or in constructor |

---

## Code Examples

### Full Service Test Class Structure

```csharp
// tests/Features/Users/Services/UserServiceTests.cs
public class UserServiceTests
{
    private readonly IUserRepository _repository;
    private readonly UserService _sut;

    public UserServiceTests()
    {
        _repository = Substitute.For<IUserRepository>();
        _sut = new UserService(_repository);
    }

    [Fact]
    public async Task GetByIdAsync_WhenUserExists_ShouldReturnMappedDto()
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
        await _repository.Received(1).GetByIdAsync(userId);
    }

    [Fact]
    public async Task GetByIdAsync_WhenUserNotFound_ShouldThrowNotFoundException()
    {
        // Arrange
        var userId = Guid.NewGuid();
        _repository.GetByIdAsync(userId).Returns((User?)null);

        // Act
        var act = async () => await _sut.GetByIdAsync(userId);

        // Assert
        await act.Should().ThrowAsync<NotFoundException>();
    }

    [Theory]
    [InlineData(true)]
    [InlineData(false)]
    public async Task SetActiveAsync_WhenCalled_ShouldUpdateAndSave(bool activeState)
    {
        // Arrange
        var user = new UserBuilder().Build();
        _repository.GetByIdAsync(Arg.Any<Guid>()).Returns(user);

        // Act
        await _sut.SetActiveAsync(user.Id, activeState);

        // Assert
        user.IsActive.Should().Be(activeState);
        await _repository.Received(1).SaveAsync(user);
    }
}
```

---

## Commands

```bash
# Create xUnit test project
dotnet new xunit -n MyApi.Tests --framework net8.0

# Add project reference
dotnet add MyApi.Tests/MyApi.Tests.csproj reference src/MyApi/MyApi.csproj

# Add testing packages
dotnet add MyApi.Tests package NSubstitute
dotnet add MyApi.Tests package FluentAssertions

# Run all tests
dotnet test --verbosity normal

# Run tests with coverage report
dotnet test --collect:"XPlat Code Coverage"

# Run only a specific test class
dotnet test --filter "FullyQualifiedName~UserServiceTests"

# Run tests matching a name pattern
dotnet test --filter "Name~WhenUserNotFound"
```

---

## Resources

- **Architecture and conventions**: See [../general/SKILL.md](../general/SKILL.md)
- **Data access patterns** (what repositories look like): See [../dataaccess/SKILL.md](../dataaccess/SKILL.md)
- **Validation patterns** (what validators look like): See [../validations/SKILL.md](../validations/SKILL.md)
- **Mapping patterns** (what mappings look like): See [../mapping/SKILL.md](../mapping/SKILL.md)
- **NSubstitute docs**: https://nsubstitute.github.io/
- **FluentAssertions docs**: https://fluentassertions.com/introduction
- **xUnit docs**: https://xunit.net/docs/getting-started/netcore/cmdline
