package main

import (
    "log"
    "project/database"
    "project/router"
    "github.com/labstack/echo/v4/middleware"
    "github.com/labstack/echo/v4"

    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to the database
    database.Connect()

    // Initialize the router
    e := router.Init()
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"}, // Allow all origins, or specify domains like "http://example.com"
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE}, // Allow specific HTTP methods
        AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization}, // Allow specific headers
    }))
    // Start the server
    log.Fatal(e.Start(":8080"))  // Start server on port 8080
}
