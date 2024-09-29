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