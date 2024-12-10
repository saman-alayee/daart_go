package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// CheckTokenMiddleware validates the JWT token in the Authorization header
func CheckTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")

		// Check if the Authorization header is missing or malformed
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "Missing or invalid Authorization header",
			})
		}

		// Extract the token part from the Authorization header
		tokenString := authHeader[len("Bearer "):]

		// Define the secret key (use environment variables in production)
		secretKey := "JIEAAIygJGWnE0y4G0EJnIyHIGOMryx1GKcIq056JGABZyWcGKcEAR5HDzcCE00kGIEJnScHFzkBI1xkGzcSZH5HoT1MZyS3Gz1nnR1HEGAnnzudGacerIcHFzcBryMeGyqSZ05dGKcnI0H1GxqSrR1gFKqAoIy5G1qSrR9KFKp="

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token is signed with the correct method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, map[string]interface{}{
					"status":  http.StatusUnauthorized,
					"message": "Unexpected signing method",
				})
			}
			return []byte(secretKey), nil
		})

		// Handle parsing errors or invalid tokens
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "Invalid or expired token",
			})
		}

		// Extract and validate claims from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Retrieve specific claims (e.g., username and publisher ID)
			username, _ := claims["sub"].(string)
			publisherID, _ := claims["publisher_id"].(string)

			// Store claims in the context for downstream use
			c.Set("username", username)
			c.Set("publisher_id", publisherID)

			// Proceed to the next handler
			return next(c)
		}

		// If token validation fails
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"message": "Invalid token claims",
		})
	}
}
