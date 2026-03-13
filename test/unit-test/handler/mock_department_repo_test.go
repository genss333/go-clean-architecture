package handler_test

import (
	"context"

	"github.com/genss333/go-clean-architecture/internal/entity"
	"github.com/stretchr/testify/mock"
)

type mockDepartmentRepo struct {
	mock.Mock
}

func (m *mockDepartmentRepo) GetDepartmentByID(ctx context.Context, id int) (entity.Department, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Department), args.Error(1)
}
