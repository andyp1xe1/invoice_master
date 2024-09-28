package main

import (

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
	"github.com/joho/godotenv"
)

var (
    googleOAuthConfig *oauth2.Config
    googleRedirectURL = "http://localhost:1337/callback"
)

func init() {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    googleOAuthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  googleRedirectURL,
        Scopes: []string{
            "https://www.googleapis.com/auth/calendar",
            // Add any additional scopes here
        },
        Endpoint: google.Endpoint,
    }
}

