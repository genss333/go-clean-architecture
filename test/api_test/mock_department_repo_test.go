package handler_test

import (
	"context"

	"github.com/genss333/go-clean-architecture/internal/entity"
)

type mockDepartmentRepo struct {
	dept entity.Department
	err  error
}

func (m *mockDepartmentRepo) GetDepartmentByID(_ context.Context, _ int) (entity.Department, error) {
	return m.dept, m.err
}
