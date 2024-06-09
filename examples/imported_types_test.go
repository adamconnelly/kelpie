package examples

import (
	"errors"
	"io"
	. "net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/requester"
)

type Requester interface {
	MakeRequest(r *Request) (io.Reader, error)
}

type ImportedTypesTest struct {
	suite.Suite
}

func (t *ImportedTypesTest) Test_CanUseAMockWithImportedTypes() {
	// Arrange
	mock := requester.NewMock()
	mock.Setup(requester.MakeRequest(kelpie.Any[*Request]()).Return(nil, errors.New("error making request")))

	// Act
	result, err := mock.Instance().MakeRequest(&Request{})

	// Assert
	t.Nil(result)
	t.ErrorContains(err, "error making request")
}

func TestImportedTypes(t *testing.T) {
	suite.Run(t, new(ImportedTypesTest))
}
