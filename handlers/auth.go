package handlers

import (
    "net/http"
    "project/database"
    "project/models"
    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/golang-jwt/jwt/v4"
	"time"
	"errors"
)
var jwtSecret = []byte("JIEAAIygJGWnE0y4G0EJnIyHIGOMryx1GKcIq056JGABZyWcGKcEAR5HDzcCE00kGIEJnScHFzkBI1xkGzcSZH5HoT1MZyS3Gz1nnR1HEGAnnzudGacerIcHFzcBryMeGyqSZ05dGKcnI0H1GxqSrR1gFKqAoIy5G1qSrR9KFKp=")

// Login - Authenticates the user and returns a JWT token
func Login(c echo.Context) error {
    // Define a struct for the request body
    var loginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Bind the request body to the loginRequest struct
    if err := c.Bind(&loginRequest); err != nil {
        return c.String(http.StatusBadRequest, "Invalid request format")
    }

    // Ensure the username and password are not empty
    if loginRequest.Username == "" || loginRequest.Password == "" {
        return c.String(http.StatusBadRequest, "Username and Password are required")
    }

    // Get the users collection from the database
    collection := database.GetCollection("users")

    // Find the user by username
    var user models.User
    err := collection.FindOne(c.Request().Context(), bson.M{"username": loginRequest.Username}).Decode(&user)
    if err != nil {
        return c.String(http.StatusUnauthorized, "Invalid username or password")
    }

    // Compare the provided password with the hashed password stored in the database
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
    if err != nil {
        return c.String(http.StatusUnauthorized, "Invalid username or password")
    }

    // Generate JWT token if authentication is successful
    token, err := generateJWT(user.Username, user.PublisherID) // Pass publisher ID here (assuming it's part of the user model)
    if err != nil {
        return c.String(http.StatusInternalServerError, "Failed to generate token")
    }

    // Return the JWT token in the response
    return c.JSON(http.StatusOK, echo.Map{
        "token": token,
    })
}

// Helper function to generate JWT token
func generateJWT(username, publisherID string) (string, error) {
    // Define the secret key and expiration time for the token
    secretKey := "JIEAAIygJGWnE0y4G0EJnIyHIGOMryx1GKcIq056JGABZyWcGKcEAR5HDzcCE00kGIEJnScHFzkBI1xkGzcSZH5HoT1MZyS3Gz1nnR1HEGAnnzudGacerIcHFzcBryMeGyqSZ05dGKcnI0H1GxqSrR1gFKqAoIy5G1qSrR9KFKp=" // Use environment variables in production
	now := time.Now()
    // Define the claims for the token
    claims := &jwt.MapClaims{
        "iss":          "daartads",               // Issuer (the app issuing the token)
        "sub":          username,               // Subject (the username of the user)
        "publisher_id": publisherID,            // Add publisher ID to the claims
        "iat":          now.Unix(),             // Issued At (time the token was generated)
        "nbf":          now.Unix(),             // Not Before (token will be valid immediately)
        "Domain":        "daartads.com",  
    }

    // Create a new token with the claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign the token with the secret key
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token and verify the signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}