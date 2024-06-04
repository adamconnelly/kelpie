package examples

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/printer"
	"github.com/adamconnelly/kelpie/examples/secretsmanager"
	sm "github.com/adamconnelly/kelpie/examples/secretsmanager/mocks/secretsmanager"
)

type VariadicFunctionsTests struct {
	suite.Suite
}

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces Printer
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

func (t *VariadicFunctionsTests) Test_Parameters_AnyMatch() {
	// Arrange
	mock := printer.NewMock()

	mock.Setup(printer.Printf("Hello %s. This is %s, %s.", kelpie.ExactMatch("Dolly"), kelpie.Any[string](), kelpie.ExactMatch("Dolly")).
		Return("Hello Dolly. This is Louis, Dolly."))

	// Act
	result := mock.Instance().Printf("Hello %s. This is %s, %s.", "Dolly", "Rab", "Dolly")

	// Assert
	t.Equal("Hello Dolly. This is Louis, Dolly.", result)
}

func (t *VariadicFunctionsTests) Test_Parameters_When() {
	// Arrange
	mock := printer.NewMock()

	var formatString string
	var args []interface{}

	mock.Setup(printer.Printf(kelpie.Any[string](), kelpie.Any[string](), kelpie.Any[string](), kelpie.Any[string]()).
		When(func(f string, a ...interface{}) string {
			formatString = f
			args = a

			return ""
		}))

	// Act
	mock.Instance().Printf("Hello %s. This is %s, %s.", "Dolly", "Rab", "Dolly")

	// Assert
	t.Equal("Hello %s. This is %s, %s.", formatString)
	t.Equal([]interface{}{"Dolly", "Rab", "Dolly"}, args)
}

func (t *VariadicFunctionsTests) Test_Parameters_Called() {
	// Arrange
	mock := printer.NewMock()

	mock.Setup(
		printer.Printf(kelpie.Any[string](), kelpie.Any[any](), kelpie.Any[any](), kelpie.Any[any]()).Return("abc"))

	// Act
	mock.Instance().Printf("Hello world!", "One", 2, 3.0)

	// Assert
	t.True(mock.Called(printer.Printf[string, any]("Hello world!", "One", 2, 3.0)))
	t.False(mock.Called(printer.Printf[string, any]("Hello world!", "Five", 2, 3.0)))
}

func (t *VariadicFunctionsTests) Test_OptionsStyleFunction() {
	// Arrange
	var secretName string
	options := secretsmanager.GetSecretOptions{}

	mock := sm.NewMock()
	mock.Setup(sm.GetSecret(kelpie.Any[context.Context](), kelpie.Any[string](), kelpie.Any[func(*secretsmanager.GetSecretOptions)]()).
		When(func(ctx context.Context, name string, opts ...func(*secretsmanager.GetSecretOptions)) (*secretsmanager.GetSecretResult, error) {
			secretName = name
			for _, opt := range opts {
				opt(&options)
			}

			return nil, nil
		}))

	// Act
	mock.Instance().GetSecret(context.Background(), "SuperSecret", secretsmanager.WithVersion("version123"))

	// Assert
	t.Equal("SuperSecret", secretName)
	t.Equal("version123", *options.Version)
}

func TestVariadicFunctions(t *testing.T) {
	suite.Run(t, new(VariadicFunctionsTests))
}
