package examples

import (
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/maths"
	"github.com/stretchr/testify/suite"
)

//go:generate go run ../cmd/kelpie generate --interfaces Maths
type Maths interface {
	Add(a, b int) int
	ParseInt(input string) (int, error)
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

func TestArgumentMatching(t *testing.T) {
	suite.Run(t, new(ArgumentMatchingTests))
}
