package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"project/database"
	"project/enums"
	"project/models"
	"project/tools"
	"strconv"
)

func SaveAdRequestToDB(viewID string, campaignID int, publisherID int, ip string, web bool, origin string) error {
	// Get the collection
	collection := database.GetCollection("ads_requests")

	// Create the AdsRequest object with the required fields
	adRequest := models.AdsRequest{
		Web:         web,
		CampaignID:  campaignID,
		PublisherID: publisherID,
		Created:     time.Now().Unix(), // Get the current Unix timestamp
		Origin:      origin,
		IP:          ip,
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
	forMobile, _ := getBooleanQueryParam(c, "forMobile") // forMobile used only to determine redirection
	web, webProvided := getBooleanQueryParam(c, "web")
	install, installProvided := getBooleanQueryParam(c, "install")

	// If forMobile is true, set web to false
	if forMobile {
		web = false
	}

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

		// Convert fields to strings
		campaignIDStr := fmt.Sprintf("%v", campaign["campaign_ID"])
		redirectURL := fmt.Sprintf("%v", campaign["redirect_url"])
		viewID := tools.UUID() // Generate a UUID for the view_id
		publisherID := "55"

		// Parse the campaignID into an integer
		campaignID, err := strconv.Atoi(campaignIDStr)
		if err != nil {
			log.Printf("Error parsing campaign_ID: %v", err)
			return c.String(http.StatusBadRequest, "Invalid campaign ID")
		}

		// Check for missing values
		if campaignIDStr == "" || redirectURL == "" || viewID == "" {
			log.Println("Missing values in campaign")
			return c.String(http.StatusBadRequest, "Missing required fields")
		}

		// Generate token with campaign_ID, redirect_url, and view_id
		token, err := jwtHandler.EncodeToken(campaignIDStr, redirectURL, publisherID, viewID, 30) // Token valid for 30 minutes
		if err != nil {
			log.Printf("Error generating token: %v", err)
			return c.String(http.StatusInternalServerError, "Error generating token")
		}

		// Set the redirect URL based on the "forMobile" parameter
		redirectPath := "/Callback?token="
		if forMobile {
			redirectPath = "/MCallback?token="
		}

		// Replace redirect_url with the generated token
		filteredCampaign := map[string]interface{}{
			"image_url":    campaign["image_url"],
			"redirect_url": redirectPath + token, // Use the encoded token as the redirect URL
			"status":       200,
			"width":        campaign["width"],
			"height":       campaign["height"],
		}
		filteredCampaigns = append(filteredCampaigns, filteredCampaign)

		// Get the request's IP and origin (referer header)
		ip := tools.GetClientIP(c.Request())
		origin := tools.GetOrigin(c.Request())

		// Save the ad request data into the database
		if err := SaveAdRequestToDB(viewID, campaignID, 55, ip, webProvided, origin); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to save ad request")
		}
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return c.String(http.StatusInternalServerError, "Error iterating through campaigns")
	}

	// Return the filtered campaigns
	return c.JSON(http.StatusOK, filteredCampaigns)
}



