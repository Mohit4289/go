package router

import (
	"gin-quickstart/handler"
	"gin-quickstart/middleware"
	"gin-quickstart/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userService *service.UserService) {

	auth := r.Group("/auth")

	auth.POST("/user", handler.User(userService))
	auth.POST("/login", handler.Login(userService))
	auth.POST("/refresh_token", handler.RefreshToken(userService))

	auth.Use(middleware.TokenVerification())
	auth.GET("/checkuser", handler.GetUser(userService))
	auth.GET("/userdata", handler.MeUser(userService))
	auth.POST("/logout", handler.Logout(userService))

}

