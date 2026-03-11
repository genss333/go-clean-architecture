package test

import (
	"testing"

	"github.com/genss333/go-clean-architecture/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	repo := NewMockUserRepository()

	uc := internal.NewUserUsecase(repo)

	result, err := uc.GetUserByID(1)

	require.NoError(t, err, "not have any errors")

	assert.Equal(t, result.ID, 1)
	assert.Equal(t, result.UserName, "jdoe88")
	assert.Equal(t, result.Firstname, "John")
	assert.Equal(t, result.Lastname, "Doe")
}
