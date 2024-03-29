# Kelpie

## What is Kelpie

Kelpie is the most magical mock generator for Go. Kelpie aims to be easy to use, and generates fully type-safe mocks for Go interfaces.

## Project Status

At the moment Kelpie is very much in development, and there are missing features and some pretty rough edges. You're of course welcome to use Kelpie, but just be prepared to hit problems and raise issues or PRs!

The following is a list of known-outstanding features:

- [ ] Generating mocks for external (i.e. not in your source tree) interfaces by name.
- [ ] Generating mocks for inline interfaces in structs.
- [ ] Support for types that require imports.

## Quickstart

Install Kelpie:

```shell
go install github.com/adamconnelly/kelpie/cmd/kelpie@latest
```

Add a `go:generate` marker to the interface you want to mock:

```go
//go:generate kelpie generate --interfaces EmailService
type EmailService interface {
	Send(sender, recipient, body string) (cost float64, err error)
}
```

Use the mock:

```go
emailServiceMock := emailservice.NewMock()
emailServiceMock.Setup(
	emailservice.Send("sender@sender.com", "someone@receiver.com", kelpie.Any[string]()).
	Return(100.54, nil)
)
emailServiceMock.Setup(
	emailservice.Send("sender@sender.com", "someone@forbidden.com", kelpie.Any[string]()).
	Return(0, errors.New("that domain is forbidden!"))
)

service := emailServiceMock.Instance()

service.Send("sender@sender.com", "someone@receiver.com", "Amazing message")
// Returns 100.54, nil

service.Send("sender@sender.com", "someone@forbidden.com", "Hello")
// Returns 0, errors.New("that domain is forbidden!)
```

## Using Kelpie

### Default Behaviour

No setup, no big deal. Kelpie returns the default values for method calls instead of panicking:

```go
mock := emailservice.NewMock()
mock.Instance().Send("sender@sender.com", "someone@receiver.com", "Hello world")
// Returns 0, nil
```

### Overriding an Expectation

Kelpie always uses the most recent expectation when trying to match a method call. That way you can easily override behaviour. This is really useful if you want to for example specify a default behaviour, and then later test an error condition:

```go
mock := emailservice.NewMock()

// Setup an initial behaviour
mock.Setup(
	emailservice.Send(kelpie.Any[string](), kelpie.Any[string](), kelpie.Any[string]()).
	Return(200, nil)
)

service := mock.Instance()

cost, err := service.Send("sender@sender.com", "someone@receiver.com", "Hello world")
t.Equal(200, cost)
t.NoError(err)

// We override the mock, to allow us to test an error condition
mock.Setup(
	emailservice.Send(kelpie.Any[string](), kelpie.Any[string](), kelpie.Any[string]()).
	Return(0, errors.New("no way!"))
)

cost, err := service.Send("sender@sender.com", "someone@receiver.com", "Hello world")
t.Equal(0, cost)
t.ErrorEqual(err, "no way!")
```

### Argument Matching

#### Exact Matching

By default Kelpie uses exact matching, and any parameters in a method call need to exactly match those specified in the setup:

```go
emailServiceMock.Setup(
	emailservice.Send("sender@sender.com", "someone@receiver.com", "Hello world").
	Return(100.54, nil)
)
```

#### Any

You can match against any possible values of a particular parameter using `kelpie.Any[T]()`:

```go
emailServiceMock.Setup(
	emailservice.Send(kelpie.Any[string](), "someone@receiver.com", "Hello world").
	Return(100.54, nil)
)
```

#### Custom Matching

You can add custom argument matching functionality using `kelpie.Match[T](isMatch)`:

```go
emailServiceMock.Setup(
	emailservice.Send(
		kelpie.Match(func(sender string) bool {
			return strings.HasSuffix(sender, "@discounted-sender.com")
		}),
		"someone@receiver.com",
		"Hello world!").
		Return(50, nil))
```

### Setting Behaviour

#### Returning a Value

To return a specific value from a method call, use `Return()`:

```go
emailServiceMock.Setup(
	emailservice.Send("sender@sender.com", "someone@receiver.com", "Hello world").
	Return(100.54, nil)
)
```

#### Panic

To panic, use `Panic()`:

```go
emailServiceMock.Setup(
	emailservice.Send("panic@thedisco.com", kelpie.Any[string](), "testing").
	Panic("Something has gone badly wrong!")
)
```

#### Custom Action

To perform a custom action, use `When()`:

```go
emailServiceMock.Setup(
	emailservice.Send(kelpie.Any[string](), kelpie.Any[string](), kelpie.Any[string]()).
		When(func(sender, recipient, body string) (float64, error) {
			// Do something
			return 0, nil
		}))
```

### Verifying Method Calls

You can verify that a method has been called using the `mock.Called()` method:

```go
// Arrange
mock := registrationservice.NewMock()

// Act
mock.Instance().Register("Mark")
mock.Instance().Register("Jim")

// Assert
t.True(mock.Called(registrationservice.Register("Mark")))
t.True(mock.Called(registrationservice.Register(kelpie.Any[string]()).Times(2)))
t.False(mock.Called(registrationservice.Register("Wendy")))
```

### Times

You can configure a method call to only match a certain number of times, or verify a method has been called a certain number of times using the `Times()`, `Once()` and `Never()` helpers:

```go
// Arrange
mock := registrationservice.NewMock()

// Act
mock.Instance().Register("Mark")
mock.Instance().Register("Jim")

// Assert
t.True(mock.Called(registrationservice.Register("Mark").Once()))
t.True(mock.Called(registrationservice.Register(kelpie.Any[string]()).Times(2)))
t.True(mock.Called(registrationservice.Register("Wendy").Never()))
```

## FAQ

### What makes Kelpie so magical

[Kelpies](https://en.wikipedia.org/wiki/Kelpie) are magical creatures from Scottish folk-lore that have shape-shifting abilities. This name seemed fitting for a mocking library, where generated mocks match the shape of interfaces that you want to simulate.

But other than that, there's nothing very magical about Kelpie.
