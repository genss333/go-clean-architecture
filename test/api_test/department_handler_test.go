package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/genss333/go-clean-architecture/internal/delivery/http/handler"
	"github.com/genss333/go-clean-architecture/internal/entity"
	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newHandler(dept entity.Department, err error) (*handler.DepartmentHandler, *gin.Engine) {
	repo := &mockDepartmentRepo{dept: dept, err: err}
	uc := depart_usecases.NewDepartmentUsecase(repo)
	h := handler.NewDepartmentHandler(uc)
	return h, setupRouter(h)
}

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(h *handler.DepartmentHandler) *gin.Engine {
	r := gin.New()
	r.GET("/api/v1/departments/:id", h.GetDepartmentByID)
	return r
}

func TestGetDepartmentByID(t *testing.T) {
	tests := []struct {
		name       string
		paramID    string
		mockDept   entity.Department
		mockErr    error
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name:    "success",
			paramID: "1",
			mockDept: entity.Department{
				DepartmentID: 1,
				Name:         "Engineering",
			},
			mockErr:    nil,
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"department_id": float64(1), "departmnet_name": "Engineering"},
		},
		{
			name:       "invalid id - not a number",
			paramID:    "abc",
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "invalid id"},
		},
		{
			name:       "not found",
			paramID:    "99",
			mockErr:    errors.New("department not found"),
			wantStatus: http.StatusNotFound,
			wantBody:   map[string]any{"error": "department not found"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, router := newHandler(tc.mockDept, tc.mockErr)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/"+tc.paramID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			var got map[string]any
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, got)
		})
	}
}

// BenchmarkGetDepartmentByID_Success measures throughput for a successful department lookup.
func BenchmarkGetDepartmentByID_Success(b *testing.B) {
	dept := entity.Department{DepartmentID: 1, Name: "Engineering"}
	_, router := newHandler(dept, nil)

	b.ResetTimer()
	for i := range b.N {
		_ = i
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/1", nil)
		router.ServeHTTP(w, req)
	}
}

// BenchmarkGetDepartmentByID_NotFound measures throughput when the department is missing.
func BenchmarkGetDepartmentByID_NotFound(b *testing.B) {
	_, router := newHandler(entity.Department{}, errors.New("department not found"))

	b.ResetTimer()
	for i := range b.N {
		_ = i
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/99", nil)
		router.ServeHTTP(w, req)
	}
}

// BenchmarkGetDepartmentByID_InvalidID measures throughput for malformed ID parsing.
func BenchmarkGetDepartmentByID_InvalidID(b *testing.B) {
	_, router := newHandler(entity.Department{}, nil)

	b.ResetTimer()
	for i := range b.N {
		_ = i
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/abc", nil)
		router.ServeHTTP(w, req)
	}
}
