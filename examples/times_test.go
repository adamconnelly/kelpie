package examples

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/examples/mocks/alarmservice"
)

//go:generate go run ../cmd/kelpie generate --package github.com/adamconnelly/kelpie/examples --interfaces AlarmService
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
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Times(2).Return(errors.New("Cannot create alarm :(")))

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
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Once().Return(errors.New("Cannot create alarm :(")))

	// Act
	result1 := mock.Instance().CreateAlarm("Fire alarm")
	result2 := mock.Instance().CreateAlarm("Fire alarm")
	result3 := mock.Instance().CreateAlarm("Fire alarm")

	// Assert
	t.Error(result1)
	t.NoError(result2)
	t.NoError(result3)
}

func (t *TimesTests) Test_Mocking_CanMatchZeroTimes() {
	// Arrange
	mock := alarmservice.NewMock()
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Times(0).Return(errors.New("Cannot create alarm :(")))

	// Act
	result1 := mock.Instance().CreateAlarm("Fire alarm")

	// Assert
	t.NoError(result1)
}

func (t *TimesTests) Test_Mocking_CanMatchNever() {
	// Arrange
	mock := alarmservice.NewMock()
	mock.Setup(alarmservice.CreateAlarm(kelpie.Any[string]()).Never().Return(errors.New("Cannot create alarm :(")))

	// Act
	result1 := mock.Instance().CreateAlarm("Fire alarm")

	// Assert
	t.NoError(result1)
}

func (t *TimesTests) Test_Verification_CanMatchSpecificTimes() {
	// Arrange
	mock := alarmservice.NewMock()

	// Act
	mock.Instance().CreateAlarm("Alarm 1")
	mock.Instance().CreateAlarm("Alarm 2")

	// Assert
	t.True(mock.Called(alarmservice.CreateAlarm(kelpie.Any[string]()).Times(2)))
	t.True(mock.Called(alarmservice.CreateAlarm("Alarm 1").Once()))
	t.False(mock.Called(alarmservice.CreateAlarm("Alarm 2").Times(2)))
}

func (t *TimesTests) Test_Verification_Never() {
	// Arrange
	mock := alarmservice.NewMock()

	// Act
	mock.Instance().CreateAlarm("Alarm 1")
	mock.Instance().CreateAlarm("Alarm 2")

	// Assert
	t.False(mock.Called(alarmservice.CreateAlarm(kelpie.Any[string]()).Never()))
	t.True(mock.Called(alarmservice.CreateAlarm("Alarm 3").Never()))
}

func TestTimes(t *testing.T) {
	suite.Run(t, new(TimesTests))
}
