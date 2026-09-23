package router

import (
	propertyHandler "gin-quickstart/handler/property"
	"gin-quickstart/middleware"
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func SetupPropertyRoutes(r *gin.Engine, propertyService *service.PropertyService) {
	propertyGroup := r.Group("/property")
	propertyGroup.Use(middleware.TokenVerification())
	propertyGroup.POST("/add", propertyHandler.AddProperty(propertyService))
	propertyGroup.DELETE("/delete", propertyHandler.DeleteProperty(propertyService))
}
