package main

import (
    "fmt"
    "log"
    "gopkg.in/gomail.v2"

	"github.com/joho/godotenv"
)

func mailService() {
    // Replace these with your own details
    from := "your-email@gmail.com"
    password := "your-app-password" // App Password from Gmail
    to := "recipient-email@example.com"
    subject := "Test Email with PDF Attachment"
    body := "Please find the attached PDF file."

    // Set up the message
    msg := gomail.NewMessage()
    msg.SetHeader("From", from)
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/plain", body)

    // Attach a PDF file
    pdfPath := "path/to/your/file.pdf"
    msg.Attach(pdfPath)

    // Set up the SMTP dialer
    dialer := gomail.NewDialer("smtp.gmail.com", 587, from, password)

    // Send the email
    err := dialer.DialAndSend(msg)
    if err != nil {
        log.Fatalf("Failed to send the email: %v", err)
    }

    fmt.Println("Email sent successfully with PDF attachment!")
}