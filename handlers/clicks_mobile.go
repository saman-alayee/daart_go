package handlers

import (
    "log"
    "net/http"
    "time"
    "strconv"
    "github.com/labstack/echo/v4"
    "project/database"
    "project/models"
    "project/tools"
)

// ErrorResponseMClick struct for structured error responses
type ErrorResponseMClick struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}

// CallbackHandlerMobile handles the callback request, decodes the token, saves the click information, and redirects to the URL
func CallbackHandlerMobile(c echo.Context) error {
    // Get the token from the query parameters
    token := c.QueryParam("token")
    if token == "" {
        return handleError(c, http.StatusBadRequest, "Token is required")
    }

    // Create an instance of JwtHandler
    jwtHandler := tools.NewJwtHandler()

    // Decode the token
    claims, err := jwtHandler.DecodeToken(token)
    if err != nil {
        log.Printf("Failed to decode token: %v", err)
        return handleError(c, http.StatusBadRequest, "Invalid token")
    }

    // Extract fields from the token (assuming claims contains campaign_ID, view_ID, publisher_ID, and redirect_url)
    campaignIDStr, ok := claims["campaign_id"].(string)
    if !ok {
        return handleError(c, http.StatusBadRequest, "Missing or invalid campaign_ID in token")
    }
    publisherIDStr, ok := claims["publisher_id"].(string)
    if !ok {
        return handleError(c, http.StatusBadRequest, "Missing or invalid publisher_ID in token")
    }
    redirectURL, ok := claims["redirect_url"].(string)
    if !ok {
        return handleError(c, http.StatusBadRequest, "Missing or invalid redirect_url in token")
    }

    // Convert campaignID and publisherID to int
    campaignID, err := strconv.Atoi(campaignIDStr)
    if err != nil {
        log.Printf("Failed to convert campaignID to int: %v", err)
        return handleError(c, http.StatusBadRequest, "Invalid campaign_ID")
    }
    publisherID, err := strconv.Atoi(publisherIDStr)
    if err != nil {
        log.Printf("Failed to convert publisherID to int: %v", err)
        return handleError(c, http.StatusBadRequest, "Invalid publisher_ID")
    }

    // Get the IP of the client
    ip := tools.GetClientIP(c.Request())

    // Create the Click object
    click := models.Click_mobile{
        PublisherID: publisherIDStr,
        CampaignID:  campaignIDStr,
        Requested:   time.Now().Unix(),
        Counted:     true, // Set to 1 for initial count
        IP:          ip,
    }

    // Get the MongoDB collection for clicks
    collection := database.GetCollection("clicks_mobile")

    // Save the click information in the collection
    _, err = collection.InsertOne(c.Request().Context(), click)
    if err != nil {
        log.Printf("Failed to save click data: %v", err)
        return handleError(c, http.StatusInternalServerError, "Failed to save click data")
    }

    // Call WakefulCP with campaignID and publisherID
    success, err := tools.WakefulCP(c.Request(), campaignID, publisherID)
    if err != nil || !success {
        log.Printf("Failed to call WakefulCP: %v", err)
        return handleError(c, http.StatusInternalServerError, "Failed to process campaign request")
    }

    // Redirect to the redirect_url from the token
    return c.Redirect(http.StatusFound, redirectURL)
}

// handleError is a helper function to structure error responses with status and message
func handleError(c echo.Context, statusCode int, message string) error {
    return c.JSON(statusCode, ErrorResponseMClick{
        Status:  statusCode,
        Message: message,
    })
}
