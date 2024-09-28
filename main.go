package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /upload", uploadHandler)
	mux.HandleFunc("/file/{id}", serveHandler)

	http.ListenAndServe("localhost:1337", mux)
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	path := fmt.Sprintf("%s%s.pdf", "./public/", id)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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

	// Create a new file in the uploads directory
	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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

	w.Write([]byte("File uploaded successfully"))
}
