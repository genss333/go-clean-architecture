package router

import (
	"github.com/genss333/go-clean-architecture/internal/delivery/http/handler"
	"github.com/genss333/go-clean-architecture/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func New(departmentHandler *handler.DepartmentHandler) *gin.Engine {
	r := gin.New()

	r.Use(middleware.ZapLogger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		departments := v1.Group("/departments")
		{
			departments.GET("/:id", departmentHandler.GetDepartmentByID)
		}
	}

	return r
}
