package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Response map[string]interface{}

type Server struct {
	listenAddr string
	scanner    *Scanner
	db         *gorm.DB
}

func NewServer(addr string) (*Server, error) {
	err := godotenv.Load()
	if err != nil {
		slog.Error("loading env failed:", err)
	}
	apiKey = os.Getenv("GROQ_API_KEY")

	s := &Server{}
	s.listenAddr = addr
	s.scanner = NewScanner()

	db, err := initDb()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	s.db = db

	systemPrompt, err = readSys("sys.md")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return s, nil
}

func (s *Server) Run() error {
	slog.Info("Warming up...")
	defer s.scanner.tess.Close()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.Handle("/", http.FileServer(http.Dir("./views")))

	mux.HandleFunc("POST /upload", s.uploadHandler)
	mux.HandleFunc("/file/{id}", s.serveHandler)

	slog.Info("Registered handlers and serving")
	return http.ListenAndServe(s.listenAddr, mux)
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

	if text, err := s.scanner.extSwitch(filePath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if result, err := llama(text); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if json, err := json.Marshal(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(json)
	}
}
