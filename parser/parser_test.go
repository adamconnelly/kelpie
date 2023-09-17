package parser_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie"
	"github.com/adamconnelly/kelpie/mocks/interfacefilter"
	"github.com/adamconnelly/kelpie/parser"
	"github.com/adamconnelly/kelpie/slices"
)

type ParserTest struct {
	suite.Suite
	interfaceFilter interfacefilter.Mock
}

func (t *ParserTest) SetupTest() {
	t.interfaceFilter.Setup(interfacefilter.Include(kelpie.Any[string]()).Return(true))
}

func (t *ParserTest) Test_Parse_ReturnsAllInterfaces() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(recipient, message string) error
}

type UserService interface {
	CreateUser(username string) (User, error)
}`

	// Act
	result, err := parser.Parse(strings.NewReader(input), "github.com/adamconnelly/kelpie/tests", &t.interfaceFilter)

	// Assert
	t.NoError(err)
	t.Len(result, 2)
	t.Equal("NotificationService", result[0].Name)
	t.Equal("UserService", result[1].Name)
}

func (t *ParserTest) Test_Parse_IgnoresInterfacesThatAreNotIncluded() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(message string) error
}

type UserService interface {
	CreateUser(username string) (User, error)
}`

	t.interfaceFilter.Setup(interfacefilter.Include("github.com/adamconnelly/kelpie/tests.UserService").Return(false))

	// Act
	result, err := parser.Parse(strings.NewReader(input), "github.com/adamconnelly/kelpie/tests", &t.interfaceFilter)

	// Assert
	t.NoError(err)
	t.Len(result, 1)
	t.Equal("NotificationService", result[0].Name)
}

func (t *ParserTest) Test_Parse_PopulatesInterfaceDetails() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(recipient, message string) error
	BroadcastNotification(message string) (recipients int, err error)
}`

	t.interfaceFilter.Setup(interfacefilter.Include("github.com/adamconnelly/kelpie/tests.UserService").Return(false))

	// Act
	result, err := parser.Parse(strings.NewReader(input), "github.com/adamconnelly/kelpie/tests", &t.interfaceFilter)

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

// TODO: what about empty interfaces? Return a warning?

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTest))
}
