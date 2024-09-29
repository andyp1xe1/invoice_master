package main

import(
	"context"
	"net/http"
	"fmt"
	"encoding/json"

	"golang.org/x/oauth2"
)



// Function to get the user's email using the access token
func getUserEmail(accessToken string) (string, error) {
    ctx := context.Background()
    client := googleOAuthConfig.Client(ctx, &oauth2.Token{AccessToken: accessToken})
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check if the response is successful
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to get user info: %s", resp.Status)
    }

    // Decode the JSON response
    var userInfo struct {
        Email string `json:"email"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        return "", err
    }

    return userInfo.Email, nil
}