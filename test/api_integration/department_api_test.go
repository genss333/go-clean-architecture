package api_integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	db_repositories "github.com/genss333/go-clean-architecture/infrastructure/database/repositories"
	"github.com/genss333/go-clean-architecture/internal/delivery/http/handler"
	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pool   *pgxpool.Pool
	router *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
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

	// Wire up full stack: repo → usecase → handler → router
	repo := db_repositories.NewDepartmentSQLCRepository(pool)
	uc := depart_usecases.NewDepartmentUsecase(repo)
	h := handler.NewDepartmentHandler(uc)

	router = gin.New()
	router.Use(gin.Recovery())
	v1 := router.Group("/api/v1")
	v1.GET("/departments/:id", h.GetDepartmentByID)

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
		pool.Exec(ctx, `DELETE FROM departments WHERE id IN (1, 2)`) //nolint:errcheck
	})
}

func TestGetDepartmentByID_Success(t *testing.T) {
	seedDepartments(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &got)
	require.NoError(t, err)

	assert.Equal(t, float64(1), got["department_id"])
	assert.Equal(t, "Engineering", got["departmnet_name"])
}

func TestGetDepartmentByID_SecondDepartment(t *testing.T) {
	seedDepartments(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &got)
	require.NoError(t, err)

	assert.Equal(t, float64(2), got["department_id"])
	assert.Equal(t, "Human Resources", got["departmnet_name"])
}

func TestGetDepartmentByID_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var got map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &got)
	require.NoError(t, err)

	assert.Contains(t, got["error"], "no rows")
}

func TestGetDepartmentByID_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var got map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &got)
	require.NoError(t, err)

	assert.Equal(t, "invalid id", got["error"])
}

func TestGetDepartmentByID_NegativeID(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetDepartmentByID_ZeroID(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/0", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetDepartmentByID_ResponseContentType(t *testing.T) {
	seedDepartments(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestGetDepartmentByID_ResponseStructure(t *testing.T) {
	seedDepartments(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var got map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &got)
	require.NoError(t, err)

	// Verify response contains exactly the expected fields
	assert.Len(t, got, 2)
	assert.Contains(t, got, "department_id")
	assert.Contains(t, got, "departmnet_name")
}

// BenchmarkGetDepartmentByID_API measures full HTTP request → DB → response throughput.
func BenchmarkGetDepartmentByID_API(b *testing.B) {
	seedDepartments(b)

	b.ResetTimer()
	for i := range b.N {
		_ = i
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/1", nil)
		router.ServeHTTP(w, req)
	}
}

// BenchmarkGetDepartmentByID_API_NotFound measures full stack throughput for missing departments.
func BenchmarkGetDepartmentByID_API_NotFound(b *testing.B) {
	b.ResetTimer()
	for i := range b.N {
		_ = i
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/999", nil)
		router.ServeHTTP(w, req)
	}
}

// BenchmarkGetDepartmentByID_API_Parallel measures full stack throughput under concurrent load.
func BenchmarkGetDepartmentByID_API_Parallel(b *testing.B) {
	seedDepartments(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/1", nil)
			router.ServeHTTP(w, req)
		}
	})
}
