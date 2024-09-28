package main

import (
	//"encoding/json"
	"fmt"
	"log/slog"
	//"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"log"
	"context"
	"strings"


	"golang.org/x/oauth2"
	"github.com/otiai10/gosseract/v2"
	"github.com/gorilla/sessions"

)

var (
    store = sessions.NewCookieStore([]byte("zel!@N7-U$NUjw9BQj+S%8DMS1XA?z%1cgJp-sE0IVY2G6P9Fq?TDImfbqnX")) // Replace with a secure key
)

type Response map[string]interface{}

type Server struct {
	listenAddr string
	tess       gosseract.Client
}

func NewServer(addr string) *Server {
	//err := godotenv.Load()
	//if err != nil {
	//	slog.Error(err.Error())
	//}

	s := &Server{}
	s.listenAddr = addr

	client := gosseract.NewClient()
	client.SetLanguage("ron", "rus", "eng")

	return s
}

func (s *Server) Run() error {
	slog.Info("Warming up...")
	defer s.tess.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", s.uploadHandler)
	mux.HandleFunc("/file/{id}", s.serveHandler)
	mux.HandleFunc("POST /send-email", s.mailHandler)
	mux.HandleFunc("/login", handleGoogleLogin)       // Google login route
    mux.HandleFunc("/callback", handleGoogleCallback)  // Google OAuth callback
	mux.HandleFunc("/add-event", handleAddEvent)

	slog.Info("Registered handlers and serving")
	return http.ListenAndServe(s.listenAddr, mux)
}

func (s *Server) imgOcr(path string) (string, error) {
	s.tess.SetImage(path)
	if text, err := s.tess.Text(); err != nil {
		return "", err
	} else {

		return text, nil
	}

}

func (s *Server) serveHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	path := fmt.Sprintf("%s%s.pdf", "./public/", id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := "./uploads/" + handler.Filename
	// Create a new file in the uploads directory
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if text, err := s.imgOcr(filePath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {

		w.Write([]byte(text))
	}
}


func (s *Server) mailHandler(w http.ResponseWriter, r *http.Request) {
	
    to := r.FormValue("email")
	slog.Info("mail: " + to)
    // Validate the email address
    if to == "" {
        http.Error(w, "Email is required", http.StatusBadRequest)
        return
    }

    // Call the mailService function
    err := mailService(to)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        return
    }

    // Respond with a success message
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Email sent successfully!"))
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, "session-id")
    if err == nil && session.Values["accessToken"] != nil {
        // If a valid session is found, redirect to your desired page
        log.Println("User already logged in, redirecting...")
        http.Redirect(w, r, "http://localhost:1337", http.StatusSeeOther)
        return
    }

    // If the session is invalid or doesn't exist, initiate Google login
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

    // Store the access token in the session
    session, _ := store.Get(r, "session-id") // Ignore the error for simplicity
    session.Values["accessToken"] = token.AccessToken // Store the new access token
    session.Options.MaxAge = 3600 // Set the session expiration time, e.g., 1 hour
    session.Save(r, w) // Save the session

    // Send a plain text response instead of redirecting
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Login successful!")) // Send a success message
}




func getTokenFromSession(r *http.Request) (*oauth2.Token, error) {
    // Retrieve the session from the request
    session, err := store.Get(r, "session-id") // Replace "session-name" with your actual session name
    if err != nil {
        return nil, err
    }

    // Retrieve the access token from the session
    accessToken, ok := session.Values["accessToken"].(string)
    if !ok {
        return nil, fmt.Errorf("access token not found in session")
    }

    // Create and return an oauth2.Token
    token := &oauth2.Token{AccessToken: accessToken}
    return token, nil
}


func handleAddEvent(w http.ResponseWriter, r *http.Request) {
    // Retrieve the token from session
    token, err := getTokenFromSession(r)
    if err != nil {
        http.Error(w, "Failed to retrieve access token: "+err.Error(), http.StatusUnauthorized)
        return
    }

    // Try to add the event to the calendar
    err = addEventToCalendar(token)
    if err != nil {
        if strings.Contains(err.Error(), "session expired") {
            http.Error(w, "Session expired. Please log in again.", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Failed to add event: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with success if the event was added successfully
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Event added successfully."))
}
