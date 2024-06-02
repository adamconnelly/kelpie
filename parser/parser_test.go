package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
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
	CreateUser(username string) (string, error)
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

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
	CreateUser(username string) (string, error)
}`

	t.interfaceFilter.Setup(interfacefilter.Include("UserService").Return(false))

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

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
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

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
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

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
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

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

func (t *ParserTests) Test_Parse_IncludesComments() {
	// Arrange
	input := `package test

// AlarmService can be used to create and manage various alarms.
type AlarmService interface {
	// AddAlarms adds new alarms, returning the alarm IDs.
	//
	// Here's some super-exciting information about this method.
	AddAlarms(names []string) []int
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result, 1)

	alarmService := slices.FirstOrPanic(result, func(mock parser.MockedInterface) bool {
		return mock.Name == "AlarmService"
	})

	addAlarms := slices.FirstOrPanic(alarmService.Methods, func(method parser.MethodDefinition) bool {
		return method.Name == "AddAlarms"
	})
	t.Equal(
		`AddAlarms adds new alarms, returning the alarm IDs.

Here's some super-exciting information about this method.`, addAlarms.Comment)
}

func (t *ParserTests) Test_Parse_HandlesPointers() {
	// Arrange
	input := `package test

type EmailSender interface {
	SendEmail(recipient string, title *string, message string) (*bool, error)
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	sendNotification := result[0].Methods[0]
	t.Equal("title", sendNotification.Parameters[1].Name)
	t.Equal("*string", sendNotification.Parameters[1].Type)

	t.Equal("*bool", sendNotification.Results[0].Type)
}

func (t *ParserTests) Test_ParsePackage_HandlesPointers() {
	// Arrange
	input := `package test

type EmailSender interface {
	SendEmail(recipient string, title *string, message string) (*bool, error)
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	sendNotification := result[0].Methods[0]
	t.Equal("title", sendNotification.Parameters[1].Name)
	t.Equal("*string", sendNotification.Parameters[1].Type)

	t.Equal("*bool", sendNotification.Results[0].Type)
}

func (t *ParserTests) Test_Parse_ReturnsImportInformation() {
	// Arrange
	input := `package test

import (
	"io"
	"net/http"
)

type Requester interface {
	MakeRequest(r *http.Request) io.Reader
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result[0]
	t.Len(requester.Imports, 2)
	t.Contains(requester.Imports, `"io"`)
	t.Contains(requester.Imports, `"net/http"`)
}

func (t *ParserTests) Test_Parse_SupportsDotImports() {
	// Arrange
	input := `package test

import . "net/http"

type Requester interface {
	MakeRequest(r *Request) *Response
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result[0]
	t.Len(requester.Imports, 1)
	t.Contains(requester.Imports, `. "net/http"`)
}

func (t *ParserTests) Test_Parse_SupportsNamedImports() {
	// Arrange
	input := `package test

import h "net/http"

type Requester interface {
	MakeRequest(r *h.Request) *h.Response
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result[0]
	t.Len(requester.Imports, 1)
	t.Contains(requester.Imports, `h "net/http"`)
}

func (t *ParserTests) Test_Parse_SupportsMaps() {
	// Arrange
	input := `package test

import (
	"io"
	"time"
)

type Updater interface {
	Update(values map[string]io.Reader)
	Values(name string) map[string]io.Reader
}

type Tracer interface {
	Traces() map[time.Time]int
}`

	// Act
	result, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	updater := slices.FirstOrPanic(result, func(i parser.MockedInterface) bool { return i.Name == "Updater" })

	update := slices.FirstOrPanic(updater.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Update" })
	valuesParam := slices.FirstOrPanic(update.Parameters, func(p parser.ParameterDefinition) bool { return p.Name == "values" })
	t.Equal("map[string]io.Reader", valuesParam.Type)

	values := slices.FirstOrPanic(updater.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Values" })
	valuesResult := values.Results[0]
	t.Equal("map[string]io.Reader", valuesResult.Type)

	t.Len(updater.Imports, 1)
	t.Contains(updater.Imports, `"io"`)

	tracer := slices.FirstOrPanic(result, func(i parser.MockedInterface) bool { return i.Name == "Tracer" })
	traces := slices.FirstOrPanic(tracer.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Traces" })
	tracesResult := traces.Results[0]
	t.Equal("map[time.Time]int", tracesResult.Type)

	t.Len(tracer.Imports, 1)
	t.Contains(tracer.Imports, `"time"`)
}

// TODO: add a test for handling types that can't be resolved (e.g. because of a mistake in the code we're parsing)
// TODO: add a test for types from the same package as the mock
// TODO: what about empty interfaces? Return a warning?

func (t *ParserTests) ParseInput(packageName, input string, filter parser.InterfaceFilter) ([]parser.MockedInterface, error) {
	tmpDir, err := os.MkdirTemp("", "kelpie-parser-tests")
	if err != nil {
		return nil, errors.Wrap(err, "could not create temp dir for module")
	}
	defer os.RemoveAll(tmpDir)

	packageDir := filepath.Join(tmpDir, packageName)

	if err := os.Mkdir(packageDir, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "could not create package directory")
	}

	goMod, err := os.Create(filepath.Join(tmpDir, "go.mod"))
	if err != nil {
		return nil, errors.Wrap(err, "could not create file for go.mod")
	}

	if _, err := goMod.WriteString(`module github.com/adamconnelly/kelpie-test

go 1.22.1`); err != nil {
		return nil, errors.Wrap(err, "could not write go.mod file")
	}

	testFile, err := os.Create(filepath.Join(packageDir, "test.go"))
	if err != nil {
		return nil, errors.Wrap(err, "could not create temp file for source")
	}

	if _, err := testFile.WriteString(input); err != nil {
		return nil, errors.Wrap(err, "could not write test case to file")
	}

	return parser.Parse("github.com/adamconnelly/kelpie-test/"+packageName, tmpDir, filter)
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTests))
}
