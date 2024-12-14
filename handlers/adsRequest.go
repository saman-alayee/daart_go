package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"net/http"
	"project/database"
	"project/enums"
	"project/models"
	"project/tools"
	"strconv"
	"strings"
	"time"
)
type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}
func SaveAdRequestToDB(viewID string, campaignID int, publisherID int, web bool, origin string) error {
	// Get the collection
	collection := database.GetCollection("ads_requests")

	// Create the AdsRequest object with the required fields
	adRequest := models.AdsRequest{
		Web:         web,
		CampaignID:  campaignID,
		PublisherID: publisherID,
		Created:     time.Now().Unix(), // Get the current Unix timestamp
		Origin:      origin,
	}

	// Insert the ad request into MongoDB
	_, err := collection.InsertOne(nil, adRequest)
	if err != nil {
		log.Printf("Failed to insert ad request: %v", err)
		return err
	}

	return nil
}

// Helper function to convert query parameter to boolean
func getBooleanQueryParam(c echo.Context, paramName string) (bool, bool) {
	value := c.QueryParam(paramName)
	if value == "" {
		return false, false // not provided
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, false // invalid boolean value
	}
	return boolValue, true // valid boolean and provided
}
func GetAdSize(c echo.Context) error {
	// Retrieve the token from the Authorization header
	tokenString := c.Request().Header.Get("Authorization")

	// Validate the token format
	token, err := parseBearerToken(tokenString)
	if err != nil {
		// Return error with status and message
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Decode the token to extract claims
	jwtHandler := tools.NewJwtHandler()
	claims, err := jwtHandler.DecodeToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalid or expired token",
		})
	}

	// Extract and validate the publisher ID from claims
	publisherID, ok := claims["publisher_id"].(string)
	if !ok || publisherID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Publisher ID is missing in the token",
		})
	}

	// Convert publisher ID to an integer
	publisherIDInt, err := strconv.Atoi(publisherID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Publisher ID",
		})
	}

	// Get the adsize parameter
	adsizeStr := c.QueryParam("Adsize")
	if adsizeStr == "" {
		adsizeStr = c.QueryParam("adsize") // Check for lowercase as fallback
	}
	if adsizeStr == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Adsize parameter is required",
		})
	}

	// Convert adsizeStr to an integer
	adsizeInt, err := strconv.Atoi(adsizeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid adsize value",
		})
	}

	// Retrieve the corresponding AdSize
	adsize, err := adsize.GetAdSizeByID(adsizeInt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid adsize value",
		})
	}

	// Get additional query parameters
	campaignType := c.QueryParam("type")
	activate, activateProvided := getBooleanQueryParam(c, "activate")
	forMobile, _ := getBooleanQueryParam(c, "forMobile")
	web, webProvided := getBooleanQueryParam(c, "web")
	install, installProvided := getBooleanQueryParam(c, "install")

	// Adjust web parameter based on forMobile
	if forMobile {
		web = false
	}

	// Build the MongoDB query filter
	filter := bson.M{"width": adsize.Width, "height": adsize.Height}
	if campaignType != "" {
		filter["type"] = campaignType
	}
	if activateProvided {
		filter["activate"] = activate
	}
	if webProvided {
		filter["web"] = web
	}
	if installProvided {
		filter["install"] = install
	}

	// Get the campaigns collection
	collection := database.GetCollection("campaigns")
	opts := options.FindOne().SetSort(bson.D{
		{Key: "bid", Value: -1},
		{Key: "dailyBudget", Value: -1},
	})

	// Retrieve the top matching campaign
	var campaign bson.M
	err = collection.FindOne(c.Request().Context(), filter, opts).Decode(&campaign)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "this size does not exist"})
	}

	// Generate a UUID for the view ID
	viewID := tools.UUID()

	// Prepare campaign fields
	campaignIDStr := fmt.Sprintf("%v", campaign["campaign_ID"])
	redirectURL := fmt.Sprintf("%v", campaign["redirect_url"])

	// Convert campaignID to integer
	campaignID, err := strconv.Atoi(campaignIDStr)
	if err != nil {
		log.Printf("Error parsing campaign_ID: %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid campaign ID",
		})
	}

	// Generate a token for redirection
	token, err = jwtHandler.EncodeToken(campaignIDStr, redirectURL, publisherID, viewID, 30) // 30-minute validity
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error generating token",
		})
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// Set redirect path based on forMobile
	redirectPath := baseURL + "/api/v1/Callback?token="
	if forMobile {
		redirectPath = baseURL + "/api/v1/MCallback?token="
	}

	// Prepare the filtered campaign response
	filteredCampaign := map[string]interface{}{
		"width":        campaign["width"],
		"height":       campaign["height"],
		"image_url":    campaign["image_url"],
		"redirect_url": redirectPath + token,
		"status":       200,
	}

	// Get and origin
	origin := tools.GetOrigin(c.Request())

	// Save the ad request to the database
	if err := SaveAdRequestToDB(viewID, campaignID, publisherIDInt, webProvided, origin); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to save ad request",
		})
	}

	// Return the campaign response
	return c.JSON(http.StatusOK, filteredCampaign)
}


// parseBearerToken validates and extracts the token from the Authorization header
func parseBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("Authorization token is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("Invalid token format")
	}

	return parts[1], nil
}
