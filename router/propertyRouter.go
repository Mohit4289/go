package router

import (
	propertyHandler "gin-quickstart/handler/property"
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func SetupPropertyRoutes(r *gin.Engine, propertyService *service.PropertyService) {
	propertyGroup := r.Group("/property")
	propertyGroup.POST("/add", propertyHandler.AddProperty(propertyService))
}

