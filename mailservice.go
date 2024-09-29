package main

import (
    "fmt"
    "gopkg.in/gomail.v2"
	"os"

	"github.com/joho/godotenv"
)

func mailService(to string) error {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        return fmt.Errorf("error loading .env file: %v", err)
    }

    // Retrieve email and password from environment variables
    from := os.Getenv("EMAIL")
    password := os.Getenv("PASSWORD")

    // Subject and body of the email
    subject := "Test Email with PDF Attachment"
    body := "Please find the attached PDF file."
    pdfPath := "/Users/mcittkmims/Downloads/Weekly_Report-3.pdf" // Example PDF file

    // Set up the email message
    msg := gomail.NewMessage()
    msg.SetHeader("From", from)
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/plain", body)
    msg.Attach(pdfPath)

    // Set up the SMTP dialer
    dialer := gomail.NewDialer("smtp.gmail.com", 587, from, password)

    // Send the email
    return dialer.DialAndSend(msg)
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