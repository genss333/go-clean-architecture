package repositories

import "github.com/genss333/go-clean-architecture/internal/entity"

type DepartmentRepository interface {
	GetDepartmentByID(id int) (entity.Department, error)
}
