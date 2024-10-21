// tools/responses.go
package tools

import "net/http"

// APIResponse represents the structure of the API response.
type APIResponse struct {
	Header struct {
		URI     string `json:"uri"`
		Builder string `json:"builder"`
	} `json:"header"`
	Status int         `json:"status"`
	Result interface{} `json:"result,omitempty"`
}

// API generates a standard API response.
func API(uri string, statusCode int, data interface{}, errMsg string) APIResponse {
	response := APIResponse{
		Status: statusCode,
	}
	response.Header.URI = uri
	response.Header.Builder = "AliAPI/2.0"

	if errMsg != "" {
		response.Result = map[string]string{"Error-Msg": errMsg}
	} else {
		response.Result = data
	}

	return response
}

// API_Default returns predefined error responses based on the status code.
func API_Default(uri string, statusCode int) APIResponse {
	response := APIResponse{
		Status: statusCode,
	}
	response.Header.URI = uri
	response.Header.Builder = "AliAPI/2.0"

	switch statusCode {
	case http.StatusUnauthorized:
		response.Result = map[string]string{"Error-Msg": "Unauthorized access"}
	case http.StatusForbidden:
		response.Result = map[string]string{"Error-Msg": "Access Prohibited"}
	case http.StatusNotFound:
		response.Result = map[string]string{"Error-Msg": "Requested API not found"}
	case http.StatusNotAcceptable:
		response.Result = map[string]string{"Error-Msg": "Request not acceptable"}
	case http.StatusUnsupportedMediaType:
		response.Result = map[string]string{"Error-Msg": "Required field not found"}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
		response.Result = map[string]string{"Error-Msg": "We have a problem Here :("}
	default:
		response.Result = map[string]string{"Error-Msg": "Unknown error"}
	}

	return response
}
