package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/parser"
	"github.com/adamconnelly/kelpie/parser/mocks/interfacefilter"
	"github.com/adamconnelly/kelpie/slices"
)

type ParserTests struct {
	suite.Suite
	interfaceFilter *interfacefilter.Mock
}

func (t *ParserTests) SetupTest() {
	t.interfaceFilter = interfacefilter.NewMock()
	t.interfaceFilter.Setup(interfacefilter.Include(kelpie.Any[string]()).Return(true))
}

func (t *ParserTests) Test_Parse_ReturnsAllInterfaces() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(recipient, message string) error
}

type UserService interface {
	CreateUser(username string) (User, error)
}`

	// Act
	result, err := parser.Parse(strings.NewReader(input), t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result, 2)
	t.Equal("NotificationService", result[0].Name)
	t.Equal("UserService", result[1].Name)
}

func (t *ParserTests) Test_Parse_IgnoresInterfacesThatAreNotIncluded() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(message string) error
}

type UserService interface {
	CreateUser(username string) (User, error)
}`

	t.interfaceFilter.Setup(interfacefilter.Include("UserService").Return(false))

	// Act
	result, err := parser.Parse(strings.NewReader(input), t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result, 1)
	t.Equal("NotificationService", result[0].Name)
}

func (t *ParserTests) Test_Parse_PopulatesInterfaceDetails() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(recipient, message string) error
	BroadcastNotification(message string) (recipients int, err error)
}`

	// Act
	result, err := parser.Parse(strings.NewReader(input), t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result, 1)

	notificationService := slices.FirstOrPanic(result, func(mock parser.MockedInterface) bool {
		return mock.Name == "NotificationService"
	})
	t.Equal("notificationservice", notificationService.PackageName)
	t.Len(notificationService.Methods, 2)

	sendNotification := slices.FirstOrPanic(notificationService.Methods, func(method parser.MethodDefinition) bool {
		return method.Name == "SendNotification"
	})

	t.Len(sendNotification.Parameters, 2)
	t.Equal("recipient", sendNotification.Parameters[0].Name)
	t.Equal("string", sendNotification.Parameters[0].Type)
	t.Equal("message", sendNotification.Parameters[1].Name)
	t.Equal("string", sendNotification.Parameters[1].Type)

	t.Len(sendNotification.Results, 1)
	t.Equal("", sendNotification.Results[0].Name)
	t.Equal("error", sendNotification.Results[0].Type)

	broadcastNotification := slices.FirstOrPanic(notificationService.Methods, func(method parser.MethodDefinition) bool {
		return method.Name == "BroadcastNotification"
	})

	t.Len(broadcastNotification.Parameters, 1)
	t.Equal("message", broadcastNotification.Parameters[0].Name)
	t.Equal("string", broadcastNotification.Parameters[0].Type)

	t.Len(broadcastNotification.Results, 2)
	t.Equal("recipients", broadcastNotification.Results[0].Name)
	t.Equal("int", broadcastNotification.Results[0].Type)
	t.Equal("err", broadcastNotification.Results[1].Name)
	t.Equal("error", broadcastNotification.Results[1].Type)
}

func (t *ParserTests) Test_Parse_SupportsMethodsWithNoResults() {
	// Arrange
	input := `package test

type NotificationService interface {
	Block(recipient string)
}`

	// Act
	result, err := parser.Parse(strings.NewReader(input), t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result, 1)

	notificationService := slices.FirstOrPanic(result, func(mock parser.MockedInterface) bool {
		return mock.Name == "NotificationService"
	})
	t.Equal("notificationservice", notificationService.PackageName)
	t.Len(notificationService.Methods, 1)

	block := slices.FirstOrPanic(notificationService.Methods, func(method parser.MethodDefinition) bool {
		return method.Name == "Block"
	})

	t.Len(block.Parameters, 1)
	t.Equal("recipient", block.Parameters[0].Name)
	t.Equal("string", block.Parameters[0].Type)

	t.Len(block.Results, 0)
}

func (t *ParserTests) Test_Parse_SupportsSlices() {
	// Arrange
	input := `package test

type AlarmService interface {
	AddAlarms(names []string) []int
}`

	// Act
	result, err := parser.Parse(strings.NewReader(input), t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result, 1)

	alarmService := slices.FirstOrPanic(result, func(mock parser.MockedInterface) bool {
		return mock.Name == "AlarmService"
	})
	t.Equal("alarmservice", alarmService.PackageName)
	t.Len(alarmService.Methods, 1)

	addAlarms := slices.FirstOrPanic(alarmService.Methods, func(method parser.MethodDefinition) bool {
		return method.Name == "AddAlarms"
	})
	t.Len(addAlarms.Parameters, 1)
	t.Equal("[]string", addAlarms.Parameters[0].Type)
	t.Len(addAlarms.Results, 1)
	t.Equal("[]int", addAlarms.Results[0].Type)
}

// TODO: what about empty interfaces? Return a warning?

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTests))
}
