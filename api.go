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
	"time"
	"encoding/json"

	"golang.org/x/oauth2"
	"github.com/otiai10/gosseract/v2"
	"github.com/gorilla/sessions"

)

var (
    store = sessions.NewCookieStore([]byte("zel!@N7-U$NUjw9BQj+S%8DMS1XA?z%1cgJp-sE0IVY2G6P9Fq?TDImfbqnX")) 
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
    // Retrieve the session from the request
    session, err := store.Get(r, "session-id") 
    if err != nil {
        http.Error(w, "Failed to retrieve session: "+err.Error(), http.StatusUnauthorized)
        return
    }

    // Extract the email from the session
    email, ok := session.Values["email"].(string)
    if !ok || email == "" {
        http.Error(w, "Email is not found in session", http.StatusUnauthorized)
        return
    }

    slog.Info("mail: " + email)

    // mailService function
    err = mailService(email)
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Email sent successfully!"))
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




func getTokenFromSession(r *http.Request) (*oauth2.Token, error) {
    // Retrieve the session from the request
    session, err := store.Get(r, "session-id") 
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

    // Decode the JSON request body to retrieve the due date
    var requestBody map[string]interface{}
    err = json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Get the due date from the request body
    dueDateStr, ok := requestBody["dueDate"].(string)
    if !ok || dueDateStr == "" {
        http.Error(w, "Due date not found or invalid", http.StatusBadRequest)
        return
    }

    // Parse the due date
    dueDate, err := time.Parse("2006-01-02", dueDateStr) // YYYY-MM-DD format
    if err != nil {
        http.Error(w, "Invalid date format: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Try to add the event to the calendar
    err = addEventToCalendar(token, dueDate)
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