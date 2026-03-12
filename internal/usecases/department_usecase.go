package depart_usecases

import (
	"context"

	"github.com/genss333/go-clean-architecture/internal/entity"
	"github.com/genss333/go-clean-architecture/internal/repositories"
)

type DepartmentUsecase struct {
	repo repositories.DepartmentRepository
}

func NewDepartmentUsecase(repo repositories.DepartmentRepository) *DepartmentUsecase {
	return &DepartmentUsecase{
		repo: repo,
	}
}

func (uc *DepartmentUsecase) GetDepartmentByID(ctx context.Context, id int) (entity.Department, error) {
	result, err := uc.repo.GetDepartmentByID(ctx, id)

	if err != nil {
		return entity.Department{}, err
	}

	return result, nil
}
