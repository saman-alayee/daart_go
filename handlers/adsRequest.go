package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"project/database"
	"project/enums" // Import the new package
	"project/tools"
	"strconv"
)

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
	// Get adsize parameter
	adsizeStr := c.QueryParam("adsize")
	if adsizeStr == "" {
		return c.String(http.StatusBadRequest, "Adsize parameter is required")
	}

	// Convert adsizeStr to an integer
	adsizeInt, err := strconv.Atoi(adsizeStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid adsize value")
	}

	// Get the corresponding AdSize from the adsize package
	adsize, err := adsize.GetAdSizeByID(adsizeInt)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid adsize value")
	}

	// Get 'type' parameter (string)
	campaignType := c.QueryParam("type")

	// Get boolean parameters
	activate, activateProvided := getBooleanQueryParam(c, "activate")
	mobile, mobileProvided := getBooleanQueryParam(c, "mobile")
	web, webProvided := getBooleanQueryParam(c, "web")
	install, installProvided := getBooleanQueryParam(c, "install")

	// Create the MongoDB query filter using width and height from the map
	filter := bson.M{
		"width":  adsize.Width,
		"height": adsize.Height,
	}

	// Add other query parameters to the filter if they are provided
	if campaignType != "" {
		filter["type"] = campaignType
	}
	if activateProvided {
		filter["activate"] = activate
	}
	if mobileProvided {
		filter["mobile"] = mobile
	}
	if webProvided {
		filter["web"] = web
	}
	if installProvided {
		filter["install"] = install
	}

	// Get the collection
	collection := database.GetCollection("campaigns")

	// Find all campaigns that match the filter
	cursor, err := collection.Find(c.Request().Context(), filter)
	if err != nil {
		log.Printf("Failed to fetch campaigns: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch campaigns")
	}
	defer cursor.Close(c.Request().Context())

	var filteredCampaigns []map[string]interface{}

	for cursor.Next(c.Request().Context()) {
		var campaign bson.M
		if err := cursor.Decode(&campaign); err != nil {
			log.Printf("Error decoding campaign: %v", err)
			return c.String(http.StatusInternalServerError, "Error decoding campaign")
		}

		// Create an instance of JwtHandler
		jwtHandler := tools.NewJwtHandler()

		// Directly convert and check if the value exists
		campaignID := fmt.Sprintf("%v", campaign["campaign_ID"])
		redirectURL := fmt.Sprintf("%v", campaign["redirect_url"])
		viewID := tools.UUID() // Assuming you meant to get view_id here

		// Check for empty values
		if campaignID == "" || redirectURL == "" || viewID == "" {
			log.Println("Missing values in campaign")
			return c.String(http.StatusBadRequest, "Missing required fields")
		}

		// Generate token with campaign_ID, redirect_url, and view_id
		token, err := jwtHandler.EncodeToken(campaignID, redirectURL, viewID, 30) // Token valid for 30 minutes
		if err != nil {
			log.Printf("Error generating token: %v", err)
			return c.String(http.StatusInternalServerError, "Error generating token")
		}

		// Replace redirect_url with the generated token
		filteredCampaign := map[string]interface{}{
			"image_url":    campaign["image_url"],
			"redirect_url": "/api/v1/Callback?data=" + token, // Use the encoded token as the redirect URL
			"status":       campaign["status"],
		}
		filteredCampaigns = append(filteredCampaigns, filteredCampaign)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return c.String(http.StatusInternalServerError, "Error iterating through campaigns")
	}

	// Return the filtered campaigns
	return c.JSON(http.StatusOK, filteredCampaigns)
}
