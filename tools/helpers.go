package tools

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// UUID generates a version 4 UUID
func UUID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04x%04x-%04x-%04x-%04x-%04x%04x%04x",
		rand.Intn(0x10000), rand.Intn(0x10000), rand.Intn(0x10000),
		rand.Intn(0x1000)|0x4000, rand.Intn(0x4000)|0x8000,
		rand.Intn(0x10000), rand.Intn(0x10000), rand.Intn(0x10000),
	)
}

// GetOrigin retrieves the origin of the request. It checks for the 'Origin' Host, then 'Referer', and finally falls back to the request's IP address.
func GetOrigin(r *http.Request) string {
	// Check for the 'Origin' header
	if origin := r.Host; origin != "" {

		return origin
	}

	// Check for the 'Referer' Host
	if referer := r.Host; referer != "" {
		return referer
	}

	// If neither is found, return the requester's IP address
	return r.RemoteAddr
}
func GetClientIP(r *http.Request) string {
	// Check if the HTTP_CLIENT_IP header exists
	clientIP := r.Header.Get("HTTP_CLIENT_IP")
	if clientIP != "" {
		return clientIP
	}

	// Check if the X-Forwarded-For header exists
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fallback to the remote address
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx] // Remove the port number
	}

	// Handle IPv6 loopback and normalize it to "127.0.0.1"
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	return ip
}

func WakefulCP(r *http.Request, cid int, source int) (bool, error) {
    clientIP := GetClientIP(r)
    url := fmt.Sprintf("https://ir.daartads.com/CP.php?Cid=%d&Source=%d&clientIP=%s&noRedirect=1&origin=api.daartads.com", cid, source, clientIP)
    resp, err := http.Get(url)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    // Optional: you can read the response if needed
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return false, err
    }

    fmt.Println("Response Body:", string(body))

    return resp.StatusCode == http.StatusOK, nil
}

