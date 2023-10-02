package db_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	db "codemead.com/go_fintech/fintech_backend/db/sqlc"
	"codemead.com/go_fintech/fintech_backend/utils"
	"github.com/stretchr/testify/require"
)

func clean_up() {
	err := testQuery.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func CreateRandomUser(t *testing.T) db.User {
	hashedPassword, err := utils.HashedPassword(utils.RandomString(8))
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQuery.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	defer clean_up()

	user1 := CreateRandomUser(t)

	user2, err := testQuery.CreateUser(context.Background(), db.CreateUserParams{
		Email:          user1.Email,
		HashedPassword: user1.HashedPassword,
	})
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestUpdateUser(t *testing.T) {
	defer clean_up()

	user := CreateRandomUser(t)

	newPassword, err := utils.HashedPassword(utils.RandomString(8))
	require.NoError(t, err)

	arg := db.UpdateUserPasswordParams{
		HashedPassword: newPassword,
		ID:             user.ID,
		UpdatedAt:      time.Now(),
	}

	newUser, err := testQuery.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, newUser.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.Email, newUser.Email)
	require.NotZero(t, user.UpdatedAt)
}

func TestGetUserByID(t *testing.T) {
	defer clean_up()

	user := CreateRandomUser(t)

	newUser, err := testQuery.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.HashedPassword, user.HashedPassword)
	require.Equal(t, newUser.Email, user.Email)
}

func TestGetUserByEmail(t *testing.T) {
	defer clean_up()

	user := CreateRandomUser(t)

	newUser, err := testQuery.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, newUser.HashedPassword, user.HashedPassword)
	require.Equal(t, newUser.Email, user.Email)
}

func TestDeleteUser(t *testing.T) {
	defer clean_up()

	user := CreateRandomUser(t)

	err := testQuery.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	newUser, err := testQuery.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, newUser)
}

func TestListUser(t *testing.T) {
	defer clean_up()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			CreateRandomUser(t)
		}()
	}
	wg.Wait()

	arg := db.ListUsersParams{
		Offset: 0,
		Limit:  10,
	}

	users, err := testQuery.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, len(users), 10)
}
