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

