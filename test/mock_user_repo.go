package test

import (
	"errors"

	"github.com/genss333/go-clean-architecture/infrastructure"
	"github.com/genss333/go-clean-architecture/internal"
)

var _ internal.UserRepository = (*MockUserRepository)(nil)

type MockUserRepository struct {
	mockData []infrastructure.UserModel
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		mockData: []infrastructure.UserModel{
			{ID: 1, UserName: "jdoe88", Firstname: "John", Lastname: "Doe"},
			{ID: 2, UserName: "marley_ghost", Firstname: "Jacob", Lastname: "Marley"},
		},
	}
}

func (r *MockUserRepository) GetUserByID(id int) (*internal.User, error) {
	for _, m := range r.mockData {
		if m.ID == id {
			userEntity := m.ToEntity()
			return &userEntity, nil
		}
	}
	return nil, errors.New("user not found")
}
