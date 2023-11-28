package examples

import (
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/accountservice"
	"github.com/stretchr/testify/suite"
)

//go:generate go run ../cmd/kelpie generate --source-file result_test.go --interfaces AccountService
type AccountService interface {
	SendActivationEmail(emailAddress string) bool
	DisableAccount(id uint)
	DisabledAccountIDs() []uint
}

type ResultTests struct {
	suite.Suite
}

func (t *ResultTests) Test_CanReturnAValue() {
	// Arrange
	mock := accountservice.NewMock()
	mock.Setup(accountservice.SendActivationEmail("a@b.com").Return(true))

	// Act
	result := mock.Instance().SendActivationEmail("a@b.com")

	// Assert
	t.True(result)
}

func (t *ResultTests) Test_CanPanic() {
	// Arrange
	mock := accountservice.NewMock()
	mock.Setup(accountservice.SendActivationEmail("a@b.com").Panic("Oh no - something really bad happened!"))

	// Act
	t.PanicsWithValue("Oh no - something really bad happened!", func() { mock.Instance().SendActivationEmail("a@b.com") })
}

func (t *ResultTests) Test_CustomAction() {
	// Arrange
	var recipientAddress string
	mock := accountservice.NewMock()
	mock.Setup(accountservice.SendActivationEmail("a@b.com").When(func(emailAddress string) bool {
		recipientAddress = emailAddress
		return true
	}))

	// Act
	mock.Instance().SendActivationEmail("a@b.com")

	// Assert
	t.Equal("a@b.com", recipientAddress)
}

func (t *ResultTests) Test_CanMockMethodsWithNoReturnArgs() {
	// Arrange
	var accountID uint
	mock := accountservice.NewMock()
	mock.Setup(accountservice.DisableAccount(kelpie.Any[uint]()).When(func(id uint) {
		accountID = id
	}))

	// Act
	mock.Instance().DisableAccount(uint(123))

	// Assert
	t.Equal(uint(123), accountID)
}

func (t *ResultTests) Test_CanMockMethodsWithNoArgs() {
	// Arrange
	mock := accountservice.NewMock()
	mock.Setup(accountservice.DisabledAccountIDs().Return([]uint{1, 2, 3}))

	// Act
	result := mock.Instance().DisabledAccountIDs()

	// Assert
	t.Equal([]uint{1, 2, 3}, result)
}

func TestResults(t *testing.T) {
	suite.Run(t, new(ResultTests))
}
