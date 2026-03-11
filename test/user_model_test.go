package test

import (
	"testing"

	"github.com/genss333/go-clean-architecture/infrastructure"
	"github.com/genss333/go-clean-architecture/internal"
	"github.com/stretchr/testify/assert"
)

func TestUserModelToEntity(t *testing.T) {
	var entity internal.User

	m := infrastructure.UserModel{
		ID:        1,
		UserName:  "gg",
		Firstname: "test",
		Lastname:  "test",
	}

	entity = m.ToEntity()

	assert.Equal(t, entity.ID, 1)
	assert.Equal(t, entity.UserName, "gg")
	assert.Equal(t, entity.Firstname, "test")
	assert.Equal(t, entity.Lastname, "test")

}
