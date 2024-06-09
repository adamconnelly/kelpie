package examples

import (
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/reader"
	"github.com/stretchr/testify/suite"
)

type ExternalTypesTests struct {
	suite.Suite
}

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

func TestExternalTypes(t *testing.T) {
	suite.Run(t, new(ExternalTypesTests))
}
