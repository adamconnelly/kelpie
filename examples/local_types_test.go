package examples

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adamconnelly/kelpie/examples/users"
	"github.com/adamconnelly/kelpie/examples/users/mocks/userrepository"
)

type LocalTypesTests struct {
	suite.Suite
}

func (t *LocalTypesTests) Test_CanReturnAValue() {
	// Arrange
	mock := userrepository.NewMock()
	mock.Setup(userrepository.FindUserByUsername("adam@kelpie.com").Return(&users.User{ID: 123}, nil))

	// Act
	user, err := mock.Instance().FindUserByUsername("adam@kelpie.com")

	// Assert
	t.NoError(err)
	t.NotNil(user)
	t.Equal(123, user.ID)
}

func (t *LocalTypesTests) Test_CanMatchOnALocalType() {
	// Arrange
	mock := userrepository.NewMock()
	mock.Setup(userrepository.GetAllUsersOfType(users.UserTypeAdmin).Return([]users.User{{ID: 1}, {ID: 2}}, nil))

	// Act
	admins, err := mock.Instance().GetAllUsersOfType(users.UserTypeAdmin)
	t.NoError(err)

	normalUsers, err := mock.Instance().GetAllUsersOfType(users.UserTypeUser)
	t.NoError(err)

	// Assert
	t.NotNil(admins)
	t.Len(admins, 2)
	t.Equal(1, admins[0].ID)
	t.Equal(2, admins[1].ID)

	t.Empty(normalUsers)
}

func TestLocalTypes(t *testing.T) {
	suite.Run(t, new(LocalTypesTests))
}
