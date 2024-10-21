package handlers

import (
    "net/http"
    "project/database"
    "project/models"
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(c echo.Context) error {
    user := models.User{
        ID: primitive.NewObjectID(), // Generate a new ObjectID
    }
    if err := c.Bind(&user); err != nil {
        return c.String(http.StatusBadRequest, "Invalid request")
    }

    collection := database.GetCollection("users")
    _, err := collection.InsertOne(c.Request().Context(), user)
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to create user")
    }

    return c.JSON(http.StatusCreated, user) // Return the created user
}

func GetUsers(c echo.Context) error {
    collection := database.GetCollection("users")

    cursor, err := collection.Find(c.Request().Context(), bson.D{})
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to fetch users")
    }
    defer cursor.Close(c.Request().Context())

    var users []models.User
    if err := cursor.All(c.Request().Context(), &users); err != nil {
        return c.String(http.StatusInternalServerError, "Failed to parse users")
    }

    return c.JSON(http.StatusOK, users)
}
