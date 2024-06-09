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

func (t *ParserTests) Test_Parse_ReturnsPackageDirectory() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(recipient, message string) error
}`

	// Act
	result, packageDir, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Equal(*packageDir, result.PackageDirectory)
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 2)
	t.Equal("NotificationService", result.Mocks[0].Name)
	t.Equal("NotificationService", result.Mocks[0].FullName)
	t.Equal("UserService", result.Mocks[1].Name)
	t.Equal("UserService", result.Mocks[1].FullName)
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 1)
	t.Equal("NotificationService", result.Mocks[0].Name)
}

func (t *ParserTests) Test_Parse_PopulatesInterfaceDetails() {
	// Arrange
	input := `package test

type NotificationService interface {
	SendNotification(recipient, message string) error
	BroadcastNotification(message string) (recipients int, err error)
}`

	// Act
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 1)

	notificationService := slices.FirstOrPanic(result.Mocks, func(mock parser.MockedInterface) bool {
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 1)

	notificationService := slices.FirstOrPanic(result.Mocks, func(mock parser.MockedInterface) bool {
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 1)

	alarmService := slices.FirstOrPanic(result.Mocks, func(mock parser.MockedInterface) bool {
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 1)

	alarmService := slices.FirstOrPanic(result.Mocks, func(mock parser.MockedInterface) bool {
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	sendNotification := result.Mocks[0].Methods[0]
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	sendNotification := result.Mocks[0].Methods[0]
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result.Mocks[0]
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result.Mocks[0]
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result.Mocks[0]
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
	result, _, err := t.ParseInput("test", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	updater := slices.FirstOrPanic(result.Mocks, func(i parser.MockedInterface) bool { return i.Name == "Updater" })

	update := slices.FirstOrPanic(updater.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Update" })
	valuesParam := slices.FirstOrPanic(update.Parameters, func(p parser.ParameterDefinition) bool { return p.Name == "values" })
	t.Equal("map[string]io.Reader", valuesParam.Type)

	values := slices.FirstOrPanic(updater.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Values" })
	valuesResult := values.Results[0]
	t.Equal("map[string]io.Reader", valuesResult.Type)

	t.Len(updater.Imports, 1)
	t.Contains(updater.Imports, `"io"`)

	tracer := slices.FirstOrPanic(result.Mocks, func(i parser.MockedInterface) bool { return i.Name == "Tracer" })
	traces := slices.FirstOrPanic(tracer.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Traces" })
	tracesResult := traces.Results[0]
	t.Equal("map[time.Time]int", tracesResult.Type)

	t.Len(tracer.Imports, 1)
	t.Contains(tracer.Imports, `"time"`)
}

func (t *ParserTests) Test_Parse_SupportsTypesFromSamePackage() {
	// Arrange
	input := `package users

type UserType int

type User struct {
	ID   uint
	Name string
	Type UserType
}

type UserService interface {
	FindUser(id uint) (*User, error)
	GetAllUsersOfType(t UserType) ([]*User, error)
}`

	// Act
	result, _, err := t.ParseInput("users", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result.Mocks[0]

	findUser := slices.FirstOrPanic(requester.Methods, func(m parser.MethodDefinition) bool { return m.Name == "FindUser" })
	t.Equal("*users.User", findUser.Results[0].Type)

	getAllUsersOfType := slices.FirstOrPanic(requester.Methods, func(m parser.MethodDefinition) bool { return m.Name == "GetAllUsersOfType" })
	t.Equal("[]*users.User", getAllUsersOfType.Results[0].Type)

	t.Len(requester.Imports, 1)
	t.Contains(requester.Imports, `"github.com/adamconnelly/kelpie-test/users"`)
}

func (t *ParserTests) Test_Parse_SupportsFunctionsInParameters() {
	// Arrange
	input := `package users

type User struct {}

type UserService interface {
	UpdateUsers(callback func(id int, user User)) error
}`

	// Act
	result, _, err := t.ParseInput("users", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	userService := result.Mocks[0]

	updateUsers := slices.FirstOrPanic(userService.Methods, func(m parser.MethodDefinition) bool { return m.Name == "UpdateUsers" })
	t.Equal("callback", updateUsers.Parameters[0].Name)
	t.Equal("func(id int, user users.User)", updateUsers.Parameters[0].Type)

	t.Len(userService.Imports, 1)
	t.Contains(userService.Imports, `"github.com/adamconnelly/kelpie-test/users"`)
}

func (t *ParserTests) Test_Parse_SupportsNamelessParameters() {
	// Arrange
	input := `package users

type UserType int

type User struct {}

type UserService interface {
	FindUser(int, UserType) (*User, error)
}`

	// Act
	result, _, err := t.ParseInput("users", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	userService := result.Mocks[0]

	findUser := slices.FirstOrPanic(userService.Methods, func(m parser.MethodDefinition) bool { return m.Name == "FindUser" })
	t.Equal("_p0", findUser.Parameters[0].Name)
	t.Equal("int", findUser.Parameters[0].Type)

	t.Equal("_p1", findUser.Parameters[1].Name)
	t.Equal("users.UserType", findUser.Parameters[1].Type)
}

func (t *ParserTests) Test_Parse_SupportsFunctionsInResults() {
	// Arrange
	input := `package users

type User struct {}

type UserService interface {
	GetUserFn() func (id int) (*User, error)
}`

	// Act
	result, _, err := t.ParseInput("users", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	userService := result.Mocks[0]

	updateUserFn := slices.FirstOrPanic(userService.Methods, func(m parser.MethodDefinition) bool { return m.Name == "GetUserFn" })
	t.Equal("", updateUserFn.Results[0].Name)
	t.Equal("func(id int) (*users.User, error)", updateUserFn.Results[0].Type)

	t.Len(userService.Imports, 1)
	t.Contains(userService.Imports, `"github.com/adamconnelly/kelpie-test/users"`)
}

func (t *ParserTests) Test_Parse_SupportsVariadicFunctions() {
	// Arrange
	input := `package users

type UserType int

type FindUsersOptions struct {
	Type *UserType
}

type User struct {
	ID   uint
	Name string
	Type UserType
}

type UserService interface {
	FindUsers(opts ...func(*FindUsersOptions)) ([]*User, error)
}`

	// Act
	result, _, err := t.ParseInput("users", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	requester := result.Mocks[0]

	findUsers := slices.FirstOrPanic(requester.Methods, func(m parser.MethodDefinition) bool { return m.Name == "FindUsers" })
	t.Equal("func(*users.FindUsersOptions)", findUsers.Parameters[0].Type)

	t.Len(requester.Imports, 1)
	t.Contains(requester.Imports, `"github.com/adamconnelly/kelpie-test/users"`)
}

func (t *ParserTests) Test_Parse_SupportsEmptyInterfaceAndAny() {
	// Arrange
	input := `package printing

type Printer interface {
	Printf(formatString string, args ...interface{}) string
	PrintfAny(formatString string, args ...any) string
}`

	// Act
	result, _, err := t.ParseInput("printing", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	printer := result.Mocks[0]

	printf := slices.FirstOrPanic(printer.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Printf" })
	t.Equal("interface{}", printf.Parameters[1].Type)
	t.True(printf.Parameters[1].IsVariadic)

	printfAny := slices.FirstOrPanic(printer.Methods, func(m parser.MethodDefinition) bool { return m.Name == "PrintfAny" })
	t.Equal("any", printfAny.Parameters[1].Type)
	t.True(printfAny.Parameters[1].IsVariadic)
}

func (t *ParserTests) Test_Parse_MarksNonEmptyInterfaceParameters() {
	// Arrange
	input := `package secrets

import "context"

type Encrypter interface {
	Encrypt(value string) string
}

type SecretsManager interface {
	GetSecret(ctx context.Context, name string) string
	PutSecret(name string, value string, encrypter Encrypter)
}`

	// Act
	result, _, err := t.ParseInput("secrets", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)

	secretsManager := slices.FirstOrPanic(result.Mocks, func(m parser.MockedInterface) bool { return m.Name == "SecretsManager" })

	getSecret := slices.FirstOrPanic(secretsManager.Methods, func(m parser.MethodDefinition) bool { return m.Name == "GetSecret" })
	t.Equal("ctx", getSecret.Parameters[0].Name)
	t.True(getSecret.Parameters[0].IsNonEmptyInterface)

	t.Equal("name", getSecret.Parameters[1].Name)
	t.False(getSecret.Parameters[1].IsNonEmptyInterface)

	putSecret := slices.FirstOrPanic(secretsManager.Methods, func(m parser.MethodDefinition) bool { return m.Name == "PutSecret" })
	t.Equal("name", putSecret.Parameters[0].Name)
	t.False(putSecret.Parameters[0].IsNonEmptyInterface)

	t.Equal("value", putSecret.Parameters[1].Name)
	t.False(putSecret.Parameters[1].IsNonEmptyInterface)

	t.Equal("encrypter", putSecret.Parameters[2].Name)
	t.True(putSecret.Parameters[2].IsNonEmptyInterface)
}

func (t *ParserTests) Test_Parse_SupportsNestedInterfaces() {
	// Arrange
	input := `package config

type IgnoredInterface interface {
	DoSomething() error
}

type ConfigService struct {
	Encrypter interface {
		Encrypt(value string) (string, error)
	}
}`

	t.interfaceFilter.Setup(interfacefilter.Include(kelpie.Any[string]()).Return(false))
	t.interfaceFilter.Setup(interfacefilter.Include("ConfigService.Encrypter").Return(true))

	// Act
	result, _, err := t.ParseInput("config", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 1)

	encrypter := result.Mocks[0]
	t.Equal("Encrypter", encrypter.Name)

	encrypt := slices.FirstOrPanic(encrypter.Methods, func(m parser.MethodDefinition) bool { return m.Name == "Encrypt" })
	t.Len(encrypt.Parameters, 1)
	t.Equal("value", encrypt.Parameters[0].Name)
	t.Equal("string", encrypt.Parameters[0].Type)

	t.Len(encrypt.Results, 2)
	t.Equal("string", encrypt.Results[0].Type)
	t.Equal("error", encrypt.Results[1].Type)
}

func (t *ParserTests) Test_Parse_SupportsDeeplyNestedInterfaces() {
	// Arrange
	input := `package doublenested

type DoubleNested struct {
	FirstLevel struct {
		DoubleNestedService interface {
			DoSomething()
		}
		
		SecondLevel struct {
			TripleNestedService interface {
				DoSomethingElse(a, b int) (string, error)
			}
		}
	}
}`

	t.interfaceFilter.Setup(interfacefilter.Include(kelpie.Any[string]()).Return(false))
	t.interfaceFilter.Setup(interfacefilter.Include("DoubleNested.FirstLevel.DoubleNestedService").Return(true))
	t.interfaceFilter.Setup(interfacefilter.Include("DoubleNested.FirstLevel.SecondLevel.TripleNestedService").Return(true))

	// Act
	result, _, err := t.ParseInput("doublenested", input, t.interfaceFilter.Instance())

	// Assert
	t.NoError(err)
	t.Len(result.Mocks, 2)

	doubleNested := slices.FirstOrPanic(result.Mocks, func(i parser.MockedInterface) bool { return i.Name == "DoubleNestedService" })
	t.NotNil(doubleNested)
	t.Equal("DoubleNested.FirstLevel.DoubleNestedService", doubleNested.FullName)

	doSomething := slices.FirstOrPanic(doubleNested.Methods, func(m parser.MethodDefinition) bool { return m.Name == "DoSomething" })
	t.Equal("DoSomething", doSomething.Name)
	t.Len(doSomething.Parameters, 0)
	t.Len(doSomething.Results, 0)

	tripleNested := slices.FirstOrPanic(result.Mocks, func(i parser.MockedInterface) bool { return i.Name == "TripleNestedService" })
	t.NotNil(tripleNested)
	t.Equal("DoubleNested.FirstLevel.SecondLevel.TripleNestedService", tripleNested.FullName)

	DoSomethingElse := slices.FirstOrPanic(tripleNested.Methods, func(m parser.MethodDefinition) bool { return m.Name == "DoSomethingElse" })
	t.Equal("DoSomethingElse", DoSomethingElse.Name)
	t.Len(DoSomethingElse.Parameters, 2)
	t.Len(DoSomethingElse.Results, 2)
}

func (t *ParserTests) Test_MockedInterface_AnyMethodsHaveParameters_ReturnsFalseIfNoMethodsHaveParameters() {
	// Arrange
	mockedInterface := parser.MockedInterface{
		Methods: []parser.MethodDefinition{
			{
				Name: "MethodWithoutParameters",
			},
		},
	}

	// Act
	hasParameters := mockedInterface.AnyMethodsHaveParameters()

	// Assert
	t.False(hasParameters)
}

func (t *ParserTests) Test_MockedInterface_AnyMethodsHaveParameters_ReturnsTrueIfAtLeastOneMethodHasParameters() {
	// Arrange
	mockedInterface := parser.MockedInterface{
		Methods: []parser.MethodDefinition{
			{
				Name: "MethodWithoutParameters",
			},
			{
				Name: "MethodWithParameters",
				Parameters: []parser.ParameterDefinition{
					{
						Name: "value",
					},
				},
			},
		},
	}

	// Act
	hasParameters := mockedInterface.AnyMethodsHaveParameters()

	// Assert
	t.True(hasParameters)
}

// TODO: add a test for handling types that can't be resolved (e.g. because of a mistake in the code we're parsing)
// TODO: what about empty interfaces? Return a warning?

func (t *ParserTests) ParseInput(packageName, input string, filter parser.InterfaceFilter) (*parser.ParsedPackage, *string, error) {
	tmpDir, err := os.MkdirTemp("", "kelpie-parser-tests")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create temp dir for module")
	}
	defer os.RemoveAll(tmpDir)

	packageDir := filepath.Join(tmpDir, packageName)

	if err := os.Mkdir(packageDir, os.ModePerm); err != nil {
		return nil, nil, errors.Wrap(err, "could not create package directory")
	}

	goMod, err := os.Create(filepath.Join(tmpDir, "go.mod"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create file for go.mod")
	}

	if _, err := goMod.WriteString(`module github.com/adamconnelly/kelpie-test

go 1.22.1`); err != nil {
		return nil, nil, errors.Wrap(err, "could not write go.mod file")
	}

	testFile, err := os.Create(filepath.Join(packageDir, "test.go"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create temp file for source")
	}

	if _, err := testFile.WriteString(input); err != nil {
		return nil, nil, errors.Wrap(err, "could not write test case to file")
	}

	pkg, err := parser.Parse("github.com/adamconnelly/kelpie-test/"+packageName, tmpDir, filter)
	if err != nil {
		return pkg, nil, err
	}

	return pkg, &packageDir, nil
}

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTests))
}
