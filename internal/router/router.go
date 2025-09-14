// internal/router/router.go
package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db" // Prisma Client
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/handler"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/service"

	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/middleware"
)

func SetupRouter(dbClient *db.PrismaClient) *gin.Engine {
	// Inisialisasi Service dan Handler
	authService := service.NewAuthService(dbClient)
	authHandler := handler.NewAuthHandler(authService)
	userService := service.NewUserService(dbClient)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handler.GetHealthCheck)

		// Rute Otentikasi
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken) 
			auth.GET("/profile", middleware.Authenticate(dbClient), authHandler.GetProfile)
			auth.PUT("/change-password", middleware.Authenticate(dbClient), authHandler.ChangePassword)
		}
		users := v1.Group("/users")
		// Lindungi semua rute di grup ini dengan otentikasi DAN otorisasi admin
		users.Use(middleware.Authenticate(dbClient), middleware.Authorize("admin"))
		{
			users.GET("", userHandler.GetUsers) // URL: /api/v1/users
			users.POST("", userHandler.CreateUser) // URL: /api/v1/users
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	return router
}