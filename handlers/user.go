package handlers

import (
    "net/http"
    "project/database"
    "project/models"
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/mongo"

)

func CreateUser(c echo.Context) error {
    user := models.User{
        ID: primitive.NewObjectID(), // Generate a new ObjectID
    }
    if err := c.Bind(&user); err != nil {
        return c.String(http.StatusBadRequest, "Invalid request format")
    }

    if user.Username == "" || user.Password == "" || user.PublisherID == "" {
        return c.String(http.StatusBadRequest, "Username and Password and PublisherId are required")
    }

    // Hash the password before saving it to the database
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to hash password")
    }
    user.Password = string(hashedPassword)

    collection := database.GetCollection("users")

    // Check if a user with the same username already exists
    filter := bson.M{"username": user.Username}
    var existingUser models.User
    err = collection.FindOne(c.Request().Context(), filter).Decode(&existingUser)
    if err == nil {
        return c.String(http.StatusConflict, "User with this username already exists")
    } else if err != mongo.ErrNoDocuments {
        return c.String(http.StatusInternalServerError, "Failed to check if user exists")
    }

    // Insert the new user with hashed password
    _, err = collection.InsertOne(c.Request().Context(), user)
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to create user")
    }

    return c.JSON(http.StatusCreated, user)
}

func GetUsers(c echo.Context) error {
    // Get the users collection from the database
    collection := database.GetCollection("users")

    // Find all users in the collection
    cursor, err := collection.Find(c.Request().Context(), bson.D{})
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to fetch users")
    }
    defer cursor.Close(c.Request().Context())

    // Create a slice to hold all the users
    var users []models.User
    if err := cursor.All(c.Request().Context(), &users); err != nil {
        return c.String(http.StatusInternalServerError, "Failed to parse users")
    }

    // Return the list of users
    return c.JSON(http.StatusOK, users)
}
