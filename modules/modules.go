package modules

import (
	"gin-gonic/middlewares"
	"gin-gonic/modules/books"
	"gin-gonic/modules/loans"
	"gin-gonic/modules/users"

	_ "gin-gonic/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute(app *gin.Engine) {
	// Swagger (public)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := app.Group("/api/v1")

	// Protected routes group
	protected := api.Group("")
	protected.Use(middlewares.JWTMiddleware())

	// Admin routes group
	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	// Register module routes
	users.RegisterRoutes(api, protected, admin)
	books.RegisterRoutes(api, admin)
	loans.RegisterRoutes(protected, admin)
}
