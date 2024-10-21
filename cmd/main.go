package main

import (
    "log"
    "project/database"
    "project/router"

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

    // Start the server
    log.Fatal(e.Start(":8080"))  // Start server on port 8080
}
