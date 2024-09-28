package main

import (

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var (
    googleOAuthConfig *oauth2.Config
    googleRedirectURL = "http://localhost:1337/callback"
)

func init() {
    googleOAuthConfig = &oauth2.Config{
        ClientID:     "1077157352206-a5vgdjg7vdbirtc6j7huup26m9qat1ia.apps.googleusercontent.com",
        ClientSecret: "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",
        RedirectURL:  googleRedirectURL,
        Scopes: []string{
            "https://www.googleapis.com/auth/calendar",
            // Add any additional scopes here
        },
        Endpoint: google.Endpoint,
    }
}

