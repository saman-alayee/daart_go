package router

import (
    "project/handlers"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "fmt"     // New: for forcing new line printing
)

func Init() *echo.Echo {
    e := echo.New()

    // Custom logger to force new line with Printf formatting
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: `${method} ${uri} ${status} ${latency_human} ${time_rfc3339}`, // Ensure newline formatting
        Output: middleware.DefaultLoggerConfig.Output, // Use default log output
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

    // Routes
    e.POST("/users", handlers.CreateUser)  // Create a new user
    e.GET("/users", handlers.GetUsers)     // Get all users
    e.GET("/GetAd", handlers.GetAdSize)    // Get ad based on size
    e.GET("/Callback", handlers.CallbackHandler)
    e.GET("/MCallback", handlers.CallbackHandlerMobile)


    return e
}
