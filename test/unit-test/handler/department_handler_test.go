package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	gin_handler "github.com/genss333/go-clean-architecture/internal/delivery/http/handler"
	"github.com/genss333/go-clean-architecture/internal/entity"
	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
	"github.com/genss333/go-clean-architecture/test/testhelper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(h *gin_handler.DepartmentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/v1/departments/:id", h.GetDepartmentByID)
	return r
}

func TestGetDepartmentByID(t *testing.T) {
	rows := testhelper.LoadCSV(t, "../../testdata/get_department_by_id.csv")

	for _, row := range rows {
		t.Run(row["name"], func(t *testing.T) {
			setupMock := row["setup_mock"] == "true"
			expectedStatus, _ := strconv.Atoi(row["expected_status"])

			repo := new(mockDepartmentRepo)

			if setupMock {
				mockID, _ := strconv.Atoi(row["mock_id"])
				mockDeptID, _ := strconv.Atoi(row["mock_dept_id"])

				var mockErr error
				if row["mock_error"] != "" {
					mockErr = errors.New(row["mock_error"])
				}

				repo.On("GetDepartmentByID", mock.Anything, mockID).
					Return(entity.Department{
						DepartmentID: int32(mockDeptID),
						Name:         row["mock_dept_name"],
					}, mockErr)
			}

			uc := depart_usecases.NewDepartmentUsecase(repo)
			h := gin_handler.NewDepartmentHandler(uc)
			r := setupRouter(h)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/departments/"+row["url_param"], nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, expectedStatus, w.Code)

			var body map[string]any
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

			if row["expected_error"] != "" {
				assert.Equal(t, row["expected_error"], body["error"])
			}
			if row["expected_dept_id"] != "" {
				expectedID, _ := strconv.Atoi(row["expected_dept_id"])
				assert.Equal(t, float64(expectedID), body["DepartmentID"])
			}
			if row["expected_dept_name"] != "" {
				assert.Equal(t, row["expected_dept_name"], body["Name"])
			}

			repo.AssertExpectations(t)
		})
	}
}
