# Kelpie

## What is Kelpie

Kelpie is the most magical mock generator for Go. Kelpie aims to be easy to use, and generates fully type-safe mocks for Go interfaces.

## Project Status

At the moment Kelpie is very much in development, and there are missing features and some pretty rough edges. You're of course welcome to use Kelpie, but just be prepared to hit problems and raise issues or PRs!

The following is a list of known-outstanding features:

- [ ] Generating mocks for inline interfaces in structs.

## Quickstart

Install Kelpie:

```shell
go install github.com/adamconnelly/kelpie/cmd/kelpie@latest
```

Add a `go:generate` marker to the interface you want to mock:

```go
//go:generate kelpie generate --package github.com/someorg/some/package --interfaces EmailService
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

### Variable parameter lists (variadic functions)

You can mock methods that accept variable parameter lists, but there are some caveats to be aware of. Here's a simple example using exact matching:

```go
type Printer interface {
	Printf(formatString string, args ...interface{}) string
}

func (t *VariadicFunctionsTests) Test_Parameters_ExactMatch() {
	// Arrange
	mock := printer.NewMock()

	mock.Setup(printer.Printf("Hello %s. This is %s, %s.", "Dolly", "Louis", "Dolly").Return("Hello Dolly. This is Louis, Dolly."))

	// Act
	result := mock.Instance().Printf("Hello %s. This is %s, %s.", "Dolly", "Louis", "Dolly")

	// Assert
	t.Equal("Hello Dolly. This is Louis, Dolly.", result)
}
```

#### Mixing exact and custom matching

Because of the way generics work, you can't mix exact matching with custom matching. So for example the following will work:

```go
mock.Setup(printer.Printf("Hello %s. This is %s, %s.", kelpie.ExactMatch("Dolly"), kelpie.Any[string](), kelpie.ExactMatch("Dolly")).
		Return("Hello Dolly. This is Louis, Dolly."))
```

But the following will not compile:

```go
mock.Setup(printer.Printf("Hello %s. This is %s, %s.", "Dolly", kelpie.Any[string](), "Dolly").
		Return("Hello Dolly. This is Louis, Dolly."))
```

#### Mixing argument types

If your variadic parameter is `...any` or `...interface{}`, and you try to pass in multiple different types of argument, the Go compiler can't infer the types for you. Here's an example:

```go
// Fails with a "mismatched types untyped string and untyped int (cannot infer P1)" error
mock.Called(printer.Printf("Hello world!", "One", 2, 3.0))
```

To fix this, just specify the type parameters:

```go
mock.Called(printer.Printf[string, any]("Hello world!", "One", 2, 3.0))
```

#### Matching no arguments

If you want to match that a variadic function call is made with no arguments provided, you can use `kelpie.None[T]()`:

```go
mock.Setup(printer.Printf("Hello world", kelpie.None[any]()))
mock.Called(secrets.Get(kelpie.Any[context.Context](), kelpie.Any[string](), kelpie.None[any]()))
```

The reason for using `None` is that otherwise the Go compiler can't infer the type of the variadic parameter:

```go
// Fails with "cannot infer P1"
mock.Setup(printer.Printf("Nothing to say").Return("Nothing to say"))
```

Another option instead of using `None` is to specify the type arguments explicitly, but that can become very verbose, especially when using Kelpie's matching functions:

```go
secretsManagerMock.Called(
	secretsmanagerapi.PutSecretValue[mocking.Matcher[context.Context], mocking.Matcher[*secretsmanager.PutSecretValueInput], func(*secretsmanager.Options)](
		kelpie.Any[context.Context](), kelpie.Any[*secretsmanager.PutSecretValueInput]()))
```

#### Matching any arguments

Similar to the way that you can match against no parameters with `kelpie.None[T]()`, you can match that any amount of parameters are passed to a variadic function using `kelpie.AnyArgs[T]()`:

```go
mock.Setup(printer.Printf("Don't panic!", kelpie.AnyArgs[any]()).Panic("Ok!"))
```

### Interface parameters

Under the hood, Kelpie uses Go generics to allow either the actual parameter type or a Kelpie matcher to be passed in when setting up mocks or verifying expectations. For example, say we have the following method:

```go
Add(a, b int) int
```

Kelpie will generate the following method for configuring expectations on `Add`:

```go
func Add[P0 int | mocking.Matcher[int], P1 int | mocking.Matcher[int]](a P0, b P1) *addMethodMatcher {
```

This is neat, because it allows each parameter to either be an `int`, or a Kelpie matcher, allowing you to write simple setups like this:

```go
mock.Setup(maths.Add(10, 20).Return(30))
```

Unfortunately Go generics don't allow a union that contains a non-empty interface. Because of this if any of your parameters accept an interface, you need to use a Kelpie matcher. For example the following won't work:

```go
var ctx context.Context
mock.Setup(secrets.Get(ctx, "MySecret").Return("SuperSecret"))
```

But the following will:

```go
var ctx context.Context
mock.Setup(secrets.Get(kelpie.Any[context.Context](), "MySecret").Return("SuperSecret"))
```

### Mocking an interface from an external package

Kelpie can happily mock interfaces that aren't part of your own source. You don't need to do anything special to mock an "external" interface - just specify the package and interface name you want to mock:

```go
//go:generate go run ../cmd/kelpie generate --package io --interfaces Reader

func (t *ExternalTypesTests) Test_CanMockAnExternalType() {
	// Arrange
	var bytesRead []byte
	mock := reader.NewMock()
	mock.Setup(reader.Read(kelpie.Match(func(b []byte) bool {
		bytesRead = b
		return true
	})).Return(20, nil))

	// Act
	read, err := mock.Instance().Read([]byte("Hello World!"))

	// Assert
	t.NoError(err)
	t.Equal(20, read)
	t.Equal([]byte("Hello World!"), bytesRead)
}
```

## FAQ

### What makes Kelpie so magical

[Kelpies](https://en.wikipedia.org/wiki/Kelpie) are magical creatures from Scottish folk-lore that have shape-shifting abilities. This name seemed fitting for a mocking library, where generated mocks match the shape of interfaces that you want to simulate.

But other than that, there's nothing very magical about Kelpie.
