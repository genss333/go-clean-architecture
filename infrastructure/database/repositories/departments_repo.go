package db_repositories

import (
	"context"

	sqlcdb "github.com/genss333/go-clean-architecture/infrastructure/database/sqlc"
	"github.com/genss333/go-clean-architecture/internal/entity"
	"github.com/genss333/go-clean-architecture/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repositories.DepartmentRepository = (*DepartmentSQLCRepository)(nil)

type DepartmentSQLCRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewDepartmentSQLCRepository(pool *pgxpool.Pool) *DepartmentSQLCRepository {
	return &DepartmentSQLCRepository{
		pool:    pool,
		queries: sqlcdb.New(pool),
	}
}

func (d *DepartmentSQLCRepository) GetDepartmentByID(ctx context.Context, id int) (entity.Department, error) {
	row, err := d.queries.GetDepartmentByID(ctx, int32(id))
	if err != nil {
		return entity.Department{}, err
	}

	return entity.Department{
		DepartmentID: row.DepartmentID,
		Name:         row.DepartmentName,
	}, nil

}
