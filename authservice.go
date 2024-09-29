package main

import (
	"log"
	"os"



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
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
            // Add any additional scopes here
        },
        Endpoint: google.Endpoint,
    }
}


func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-id")
    if err == nil && session.Values["accessToken"] != nil {
        // If a valid session is found, redirect 
        log.Println("User already logged in, redirecting...")
        http.Redirect(w, r, "http://localhost:1337", http.StatusSeeOther)
        return
    }

    // If the session is invalid initiate Google login
    url := googleOAuthConfig.AuthCodeURL("state")
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    code := r.URL.Query().Get("code")

    // Exchange the authorization code for an access token
    token, err := googleOAuthConfig.Exchange(ctx, code)
    if err != nil {
        http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
        return
    }

	// Retrieve user email using the access token
	email, err := getUserEmail(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user email: "+err.Error(), http.StatusInternalServerError)
		return
	}


    // Store the access token in the session
    session, _ := store.Get(r, "session-id") 
    session.Values["accessToken"] = token.AccessToken 
	session.Values["email"] = email // Store the email in the session
    session.Options.MaxAge = 3600 
    session.Save(r, w)

    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Login successful!")) 
}

