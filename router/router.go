package router

import (
	"fmt"
	"project/handlers"
	"project/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	// Custom logger to force new line with Printf formatting
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${method} ${uri} ${status} ${latency_human} ${time_rfc3339}\n`,
		Output: middleware.DefaultLoggerConfig.Output,
	}))

	e.Use(middleware.Recover())

	// Print an extra blank line after each request log for debugging
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			fmt.Println() // Force a new line after each request for debugging
			return err
		}
	})

	// Group all routes under "/api"
	api := e.Group("/api/v1")

	// Public Routes (No middleware applied)
	api.POST("/users", handlers.CreateUser)        // Create a new user
	api.POST("/auth/login", handlers.Login)        // Login and get a token
	api.GET("/Callback", handlers.CallbackHandler) // Public callback
	api.GET("/MCallback", handlers.CallbackHandlerMobile)

	// Protected Routes (Require JWT token)
	apiProtected := api.Group("") // Nested under "/api"
	apiProtected.Use(middlewares.CheckTokenMiddleware)// Apply token-checking middleware

	apiProtected.GET("/users", handlers.GetUsers)     // Get all users (protected)
	apiProtected.GET("/GetAd", handlers.GetAdSize)    // Get ad based on size (protected)

	return e
}
