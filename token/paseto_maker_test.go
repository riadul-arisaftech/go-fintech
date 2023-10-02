package token

import (
	"testing"

	"codemead.com/go_fintech/fintech_backend/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := utils.RandomString(8)

	hashedPassword1, err := utils.HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = utils.CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := utils.RandomString(6)
	err = utils.CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := utils.HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
