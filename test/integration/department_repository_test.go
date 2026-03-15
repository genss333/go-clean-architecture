package repositories_integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	db_repositories "github.com/genss333/go-clean-architecture/infrastructure/database/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		panic("failed to start postgres container: " + err.Error())
	}
	defer pgContainer.Terminate(ctx) //nolint:errcheck

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("failed to get connection string: " + err.Error())
	}

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic("failed to parse pgx config: " + err.Error())
	}

	poolConfig.MaxConns = 50
	poolConfig.MinConns = 10

	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		panic("failed to create pgx pool: " + err.Error())
	}
	defer pool.Close()

	if err := runMigrations(ctx, pool); err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	os.Exit(m.Run())
}

func runMigrations(ctx context.Context, p *pgxpool.Pool) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS departments (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, sql := range migrations {
		if _, err := p.Exec(ctx, sql); err != nil {
			return err
		}
	}

	return nil
}

func seedDepartments(t testing.TB) {
	t.Helper()
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`INSERT INTO departments (id, name) VALUES ($1, $2), ($3, $4)
		 ON CONFLICT DO NOTHING`,
		1, "Engineering",
		2, "Human Resources",
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM departments WHERE id IN (1, 2)`)
	})
}

func TestDepartmentRepository_GetDepartmentByID(t *testing.T) {
	ctx := context.Background()
	repo := db_repositories.NewDepartmentSQLCRepository(pool)

	seedDepartments(t)

	tests := []struct {
		name     string
		id       int
		wantID   int32
		wantName string
		wantErr  bool
	}{
		{
			name:     "found department with id 1",
			id:       1,
			wantID:   1,
			wantName: "Engineering",
			wantErr:  false,
		},
		{
			name:     "found department with id 2",
			id:       2,
			wantID:   2,
			wantName: "Human Resources",
			wantErr:  false,
		},
		{
			name:    "not found returns error",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dept, err := repo.GetDepartmentByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantID, dept.DepartmentID)
			assert.Equal(t, tt.wantName, dept.Name)
		})
	}
}

// BenchmarkGetDepartmentByID_Found วัด throughput การ query department ที่มีอยู่จริง
func BenchmarkGetDepartmentByID_Found(b *testing.B) {
	seedDepartments(b)

	ctx := context.Background()
	repo := db_repositories.NewDepartmentSQLCRepository(pool)

	b.ResetTimer()
	for i := range b.N {
		_ = i
		_, err := repo.GetDepartmentByID(ctx, 1)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetDepartmentByID_NotFound วัด throughput การ query department ที่ไม่มีในฐานข้อมูล
func BenchmarkGetDepartmentByID_NotFound(b *testing.B) {
	ctx := context.Background()
	repo := db_repositories.NewDepartmentSQLCRepository(pool)

	b.ResetTimer()
	for i := range b.N {
		_ = i
		_, _ = repo.GetDepartmentByID(ctx, 999999)
	}
}

// BenchmarkGetDepartmentByID_Parallel วัด throughput เมื่อยิง query พร้อมกันหลาย goroutine
func BenchmarkGetDepartmentByID_Parallel(b *testing.B) {
	seedDepartments(b)

	ctx := context.Background()
	repo := db_repositories.NewDepartmentSQLCRepository(pool)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := repo.GetDepartmentByID(ctx, 1)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
