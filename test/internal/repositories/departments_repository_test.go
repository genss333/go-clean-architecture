package department_repositories_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/genss333/go-clean-architecture/internal/entity"
	"github.com/genss333/go-clean-architecture/internal/repositories"
	depart_usecases "github.com/genss333/go-clean-architecture/internal/usecases"
	"github.com/genss333/go-clean-architecture/test/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ repositories.DepartmentRepository = (*MockDepartmentRepository)(nil)

type MockDepartmentRepository struct {
	mockData []entity.Department
}

func NewMockDepartmentRepository(rows []map[string]string) *MockDepartmentRepository {
	result := make([]entity.Department, len(rows))

	for idx, row := range rows {
		dpID, _ := strconv.Atoi(row["department_id"])
		dp := entity.Department{
			DepartmentID: dpID,
			Name:         row["department_name"],
		}
		result[idx] = dp
	}
	return &MockDepartmentRepository{
		mockData: result,
	}
}

func (m *MockDepartmentRepository) GetDepartmentByID(id int) (entity.Department, error) {
	for _, dp := range m.mockData {
		if int(dp.DepartmentID) == id {
			return dp, nil
		}
	}
	return entity.Department{}, errors.New("not found departments")
}

func TestDepartmentUescases(t *testing.T) {
	rows := testhelper.LoadCSV(t, "../testdata/departments.csv")
	mockRepo := NewMockDepartmentRepository(rows)
	uc := depart_usecases.NewDepartmentUsecase(mockRepo)

	tests := []struct {
		name    string
		id      int
		want    entity.Department
		wantErr bool
	}{
		{
			name:    "Success - QA Department",
			id:      1,
			want:    entity.Department{DepartmentID: 1, Name: "QA"},
			wantErr: false,
		},
		{
			name:    "Failed - QA Department",
			id:      99,
			want:    entity.Department{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := uc.GetDepartmentByID(tt.id)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, actual)
			}
		})
	}
}
