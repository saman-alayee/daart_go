package tools

import (
    "errors"
    "time"
    "fmt"
    "github.com/golang-jwt/jwt/v4"
)

type JwtHandler struct {
    SecretKey string
    Domain    string
}

// Claims struct for custom claims
type Claims struct {
    CampaignID string `json:"campaign_id"`
    RedirectURL string `json:"redirect_url"`
    ViewID     string `json:"view_id"`
    PublisherID     string `json:"publisher_id"`
    jwt.RegisteredClaims
}

// NewJwtHandler initializes a new JwtHandler
func NewJwtHandler() *JwtHandler {
    return &JwtHandler{
        SecretKey: "JIEAAIygJGWnE0y4G0EJnIyHIGOMryx1GKcIq056JGABZyWcGKcEAR5HDzcCE00kGIEJnScHFzkBI1xkGzcSZH5HoT1MZyS3Gz1nnR1HEGAnnzudGacerIcHFzcBryMeGyqSZ05dGKcnI0H1GxqSrR1gFKqAoIy5G1qSrR9KFKp=",
        Domain:    "daartads.com", // Set the domain (use environment variable in production)
    }
}

// EncodeToken creates a JWT token with campaign_id, redirect_url, and view_id
func (j *JwtHandler) EncodeToken(campaignID, redirectURL,publisherID, viewID string, validFor int) (string, error) {
    expirationTime := time.Now().Add(time.Duration(validFor) * time.Minute)

    claims := &Claims{
        CampaignID: campaignID,
        RedirectURL: redirectURL,
        ViewID:     viewID,
        PublisherID:  publisherID,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    j.Domain,
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
    return token.SignedString([]byte(j.SecretKey))
}

// ValidateToken verifies and decodes the JWT token
func (j *JwtHandler) ValidateToken(tokenStr string) (bool, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(j.SecretKey), nil
    })

    if err != nil || !token.Valid {
        return false, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || claims.Issuer != j.Domain {
        return false, errors.New("invalid token or issuer")
    }

    return true, nil
}

// GetTokenData extracts the data from the JWT token
func (j *JwtHandler) GetTokenData(tokenStr string) (map[string]interface{}, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(j.SecretKey), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return map[string]interface{}{
        "campaign_id": claims.CampaignID,
        "redirect_url": claims.RedirectURL,
        "view_id":     claims.ViewID,
        "publisher_id":     claims.PublisherID,
    }, nil
}

// DecodeToken decodes a JWT token and returns the claims
func (h *JwtHandler) DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
