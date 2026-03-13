package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/genss333/go-clean-architecture/internal/entity"
	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
)

type mockDepartmentRepo struct {
	dept entity.Department
	err  error
}

func (m *mockDepartmentRepo) GetDepartmentByID(_ context.Context, _ int) (entity.Department, error) {
	return m.dept, m.err
}

// BenchmarkUsecase_GetDepartmentByID_Found วัด throughput เมื่อ repository คืนค่าสำเร็จ
func BenchmarkUsecase_GetDepartmentByID_Found(b *testing.B) {
	repo := &mockDepartmentRepo{
		dept: entity.Department{DepartmentID: 1, Name: "Engineering"},
	}
	uc := depart_usecases.NewDepartmentUsecase(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		_ = i
		_, err := uc.GetDepartmentByID(ctx, 1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUsecase_GetDepartmentByID_NotFound วัด throughput เมื่อ repository คืน error
func BenchmarkUsecase_GetDepartmentByID_NotFound(b *testing.B) {
	repo := &mockDepartmentRepo{
		err: errors.New("department not found"),
	}
	uc := depart_usecases.NewDepartmentUsecase(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := range b.N {
		_ = i
		_, _ = uc.GetDepartmentByID(ctx, 999)
	}
}
