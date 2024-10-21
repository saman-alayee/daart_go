package handlers

import (
    "net/http"
    "project/database"
    "project/enums" // Import the new package
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "strconv"
    "log"
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

// GetAdSize retrieves all campaigns based on the adsize and other optional parameters.
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

    // Slice to hold all matching campaigns
    var campaigns []bson.M

    // Iterate through the cursor and decode each document into the campaigns slice
    for cursor.Next(c.Request().Context()) {
        var campaign bson.M
        if err := cursor.Decode(&campaign); err != nil {
            log.Printf("Error decoding campaign: %v", err)
            return c.String(http.StatusInternalServerError, "Error decoding campaign")
        }
        campaigns = append(campaigns, campaign)
    }

    // Check for cursor errors
    if err := cursor.Err(); err != nil {
        log.Printf("Cursor error: %v", err)
        return c.String(http.StatusInternalServerError, "Error iterating through campaigns")
    }

    // // Log the retrieved campaign details
    // log.Printf("Found campaigns: %+v", campaigns)

    // Return the matching campaigns
    return c.JSON(http.StatusOK, campaigns)
}
