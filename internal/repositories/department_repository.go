package repositories

import (
	"context"

	"github.com/genss333/go-clean-architecture/internal/entity"
)

type DepartmentRepository interface {
	GetDepartmentByID(ctx context.Context, id int) (entity.Department, error)
}
