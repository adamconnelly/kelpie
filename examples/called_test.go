package examples

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/registrationservice"
)

type RegistrationService interface {
	// Register registers the item with the specified name.
	Register(name string) error
}

type CalledTests struct {
	suite.Suite
}

func (t *CalledTests) Test_Called_VerifiesWhetherMethodHasBeenCalled() {
	// Arrange
	mock := registrationservice.NewMock()

	// Act
	mock.Instance().Register("Mark")
	mock.Instance().Register("Jim")

	// Assert
	t.True(mock.Called(registrationservice.Register("Mark")))
	t.True(mock.Called(registrationservice.Register(kelpie.Any[string]()).Times(2)))
	t.False(mock.Called(registrationservice.Register("Wendy")))
}

func TestCalled(t *testing.T) {
	suite.Run(t, new(CalledTests))
}
