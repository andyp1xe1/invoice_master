package main

import (
	"fmt"
	//"github.com/joho/godotenv"
	"github.com/otiai10/gosseract/v2"
	"io"
	"net/http"
	"os"
)

type Response map[string]interface{}

type Server struct {
	ListenAddr string
	Tess       gosseract.Client
}

func (s *Server) New() {
	//err := godotenv.Load()
	//if err != nil {
	//	slog.Error(err.Error())
	//}

	s.ListenAddr = "localhost:1337"

	client := gosseract.NewClient()
	client.SetLanguage("ron", "rus", "eng")
}

func (s *Server) Run() {
	defer s.Tess.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", s.uploadHandler)
	mux.HandleFunc("/file/{id}", s.serveHandler)

	http.ListenAndServe(s.ListenAddr, mux)
}

func (s *Server) ImgOcr(path string) (string, error) {
	s.Tess.SetImage(path)
	if text, err := s.Tess.Text(); err != nil {
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

	if text, err := s.ImgOcr(filePath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {

		w.Write([]byte(text))
	}
}
