package examples

import (
	"errors"
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/maths"
	"github.com/adamconnelly/kelpie/examples/mocks/sender"
	"github.com/stretchr/testify/suite"
)

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces Maths
type Maths interface {
	// Add adds a and b together and returns the result.
	Add(a, b int) int

	// ParseInt interprets a string s in the given base (0, 2 to 36) and bit size (0 to 64)
	// and returns the corresponding value i.
	//
	// The string may begin with a leading sign: "+" or "-".
	//
	// If the base argument is 0, the true base is implied by the string's prefix following
	// the sign (if present): 2 for "0b", 8 for "0" or "0o", 16 for "0x", and 10 otherwise.
	// Also, for argument base 0 only, underscore characters are permitted as defined by the
	// Go syntax for integer literals.
	//
	// The bitSize argument specifies the integer type that the result must fit into. Bit
	// sizes 0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64. If bitSize
	// is below 0 or above 64, an error is returned.
	//
	// The errors that ParseInt returns have concrete type *NumError and include err.Num = s.
	// If s is empty or contains invalid digits, err.Err = ErrSyntax and the returned value is
	// 0; if the value corresponding to s cannot be represented by a signed integer of the given
	// size, err.Err = ErrRange and the returned value is the maximum magnitude integer of the
	// appropriate bitSize and sign.
	ParseInt(input string) (int, error)
}

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces Sender
type Sender interface {
	SendMessage(title *string, message string) error
}

type ArgumentMatchingTests struct {
	suite.Suite
}

func (t *ArgumentMatchingTests) Test_ReturnsDefaultValuesForNoMatch() {
	// Arrange
	mock := maths.NewMock()

	// Act
	result, err := mock.Instance().ParseInt("10")

	// Assert
	t.Equal(0, result)
	t.Nil(err)
}

func (t *ArgumentMatchingTests) Test_CanPerformAnExactMatch() {
	// Arrange
	mock := maths.NewMock()
	mock.Setup(maths.Add(10, 20).Return(30))

	// Act
	result := mock.Instance().Add(10, 20)

	// Assert
	t.Equal(30, result)
}

func (t *ArgumentMatchingTests) Test_CanMatchAnyValue() {
	// Arrange
	mock := maths.NewMock()
	mock.Setup(maths.Add(kelpie.Any[int](), 20).Return(30))

	// Act
	result := mock.Instance().Add(123, 20)

	// Assert
	t.Equal(30, result)
}

func (t *ArgumentMatchingTests) Test_CanUseCustomMatchingLogic() {
	// Arrange
	mock := maths.NewMock()
	mock.Setup(maths.Add(kelpie.Match[int](func(arg int) bool {
		return arg > 0
	}), 20).Return(30))

	// Act
	result1 := mock.Instance().Add(123, 20)
	result2 := mock.Instance().Add(-1, 20)

	// Assert
	t.Equal(30, result1)
	t.Equal(0, result2)
}

func (t *ArgumentMatchingTests) Test_CanMatchNil() {
	// Arrange
	mock := sender.NewMock()
	mock.Setup(sender.SendMessage((*string)(nil), "Open the gates").Return(errors.New("the way is shut!")))

	// Act
	err := mock.Instance().SendMessage(nil, "Open the gates")

	// Assert
	t.ErrorContains(err, "the way is shut")
}

func TestArgumentMatching(t *testing.T) {
	suite.Run(t, new(ArgumentMatchingTests))
}
