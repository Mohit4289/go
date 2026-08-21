package router

import (
	"gin-quickstart/handler"
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userService *service.UserService) {

	auth := r.Group("/auth")

	auth.POST("/user", handler.User(userService))
	auth.POST("/login", handler.Login)
}
