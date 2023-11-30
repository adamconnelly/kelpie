package examples

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/alarmservice"
)

//go:generate go run ../cmd/kelpie generate --interfaces AlarmService
type AlarmService interface {
	CreateAlarm(name string) error
}

type TimesTests struct {
	suite.Suite
}

func (t *TimesTests) Test_Mocking_DefaultsToMatchingUnlimitedTimes() {
	// Arrange
	mock := alarmservice.NewMock()
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Return(errors.New("Cannot create alarm :(")))

	// Act
	result1 := mock.Instance().CreateAlarm("Fire alarm")
	result2 := mock.Instance().CreateAlarm("Fire alarm")
	result3 := mock.Instance().CreateAlarm("Fire alarm")

	// Assert
	t.Error(result1)
	t.Error(result2)
	t.Error(result3)
}

func (t *TimesTests) Test_Mocking_CanMatchSpecificTimes() {
	// Arrange
	mock := alarmservice.NewMock()
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Return(errors.New("Cannot create alarm :(")).Times(2))

	// Act
	result1 := mock.Instance().CreateAlarm("Fire alarm")
	result2 := mock.Instance().CreateAlarm("Fire alarm")
	result3 := mock.Instance().CreateAlarm("Fire alarm")

	// Assert
	t.Error(result1)
	t.Error(result2)
	t.NoError(result3)
}

func (t *TimesTests) Test_Mocking_CanMatchOnce() {
	// Arrange
	mock := alarmservice.NewMock()
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Return(errors.New("Cannot create alarm :(")).Once())

	// Act
	result1 := mock.Instance().CreateAlarm("Fire alarm")
	result2 := mock.Instance().CreateAlarm("Fire alarm")
	result3 := mock.Instance().CreateAlarm("Fire alarm")

	// Assert
	t.Error(result1)
	t.NoError(result2)
	t.NoError(result3)
}

// TODO: add tests for verification after refactoring setup vs verifying

func TestTimes(t *testing.T) {
	suite.Run(t, new(TimesTests))
}
