package examples

import (
	"errors"
	"testing"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/doublenestedservice"
	"github.com/adamconnelly/kelpie/examples/mocks/encrypter"
	"github.com/adamconnelly/kelpie/examples/mocks/storage"
	"github.com/stretchr/testify/suite"
)

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces ConfigService.Encrypter
//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces ConfigService.Storage
type ConfigService struct {
	Encrypter interface {
		Encrypt(value string) (string, error)
	}

	Storage interface {
		StoreConfigValue(key, value string) error
	}
}

func (c *ConfigService) StoreConfig(key, value string) error {
	encryptedValue, err := c.Encrypter.Encrypt(value)
	if err != nil {
		return err
	}

	return c.Storage.StoreConfigValue(key, encryptedValue)
}

type NestedInterfacesTests struct {
	suite.Suite
}

func (t *NestedInterfacesTests) Test_ConfigService_StoresEncryptedValue() {
	// Arrange
	encrypterMock := encrypter.NewMock()
	storageMock := storage.NewMock()

	encrypterMock.Setup(encrypter.Encrypt("unencrypted").Return("encrypted", nil))

	configService := &ConfigService{
		Encrypter: encrypterMock.Instance(),
		Storage:   storageMock.Instance(),
	}

	// Act
	err := configService.StoreConfig("kelpie.testSecret", "unencrypted")

	// Assert
	t.NoError(err)
	t.True(storageMock.Called(storage.StoreConfigValue("kelpie.testSecret", "encrypted")))
}

func (t *NestedInterfacesTests) Test_ConfigService_HandlesEncryptionFailures() {
	// Arrange
	encrypterMock := encrypter.NewMock()
	storageMock := storage.NewMock()

	encrypterMock.Setup(encrypter.Encrypt("unencrypted").Return("", errors.New("failure to encrypt")))

	configService := &ConfigService{
		Encrypter: encrypterMock.Instance(),
		Storage:   storageMock.Instance(),
	}

	// Act
	err := configService.StoreConfig("kelpie.testSecret", "unencrypted")

	// Assert
	t.ErrorContains(err, "failure to encrypt")
	t.False(storageMock.Called(storage.StoreConfigValue(kelpie.Any[string](), kelpie.Any[string]())))
}

func (t *NestedInterfacesTests) Test_ConfigService_HandlesStorageFailures() {
	// Arrange
	encrypterMock := encrypter.NewMock()
	storageMock := storage.NewMock()

	encrypterMock.Setup(encrypter.Encrypt("unencrypted").Return("encrypted", nil))
	storageMock.Setup(storage.StoreConfigValue(kelpie.Any[string](), kelpie.Any[string]()).Return(errors.New("could not store value")))

	configService := &ConfigService{
		Encrypter: encrypterMock.Instance(),
		Storage:   storageMock.Instance(),
	}

	// Act
	err := configService.StoreConfig("kelpie.testSecret", "unencrypted")

	// Assert
	t.ErrorContains(err, "could not store value")
}

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces DoubleNested.Internal.DoubleNestedService
type DoubleNested struct {
	Internal struct {
		DoubleNestedService interface {
			DoSomething()
		}
	}
}

func (d *DoubleNested) DoSomething() {
	d.Internal.DoubleNestedService.DoSomething()
}

func (t *NestedInterfacesTests) Test_SupportsDeeplyNestedInterfaces() {
	// Arrange
	mock := doublenestedservice.NewMock()

	doubleNested := &DoubleNested{
		Internal: struct{ DoubleNestedService interface{ DoSomething() } }{
			DoubleNestedService: mock.Instance(),
		},
	}

	// Act
	doubleNested.DoSomething()

	// Assert
	t.True(mock.Called(doublenestedservice.DoSomething()))
}

func TestNestedInterfaces(t *testing.T) {
	suite.Run(t, new(NestedInterfacesTests))
}
