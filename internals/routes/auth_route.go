package routes

import (
	"github.com/gin-gonic/gin"
	handlers "filmfolk/internals/handler"
)


func SetupAuthRoutes(router *gin.RouterGroup) {
	
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", handlers.RegisterUser)
		authGroup.POST("/login", handlers.LoginUser)
	}
}
