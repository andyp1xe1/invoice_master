package main

import (
    "context"
	"net/http"
	"fmt"
	"time"
	

	"google.golang.org/api/googleapi"
    "google.golang.org/api/calendar/v3"
	"golang.org/x/oauth2"

)

func addEventToCalendar(token *oauth2.Token, dueDate time.Time) error {
    ctx := context.Background()
    client := googleOAuthConfig.Client(ctx, token)

    calendarService, err := calendar.New(client)
    if err != nil {
        return err
    }

    // Create a fixed zone for GMT+3
    gmtPlus3 := time.FixedZone("GMT+3", 3*60*60) // GMT+3 timezone

    // Create start and end times in GMT+3
    startTime := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, gmtPlus3) // Midnight start in GMT+3
    endTime := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 23, 59, 0, 0, gmtPlus3) // 23:59 end in GMT+3



	// Set up reminders
	reminders:=&calendar.EventReminders{
    Overrides: []*calendar.EventReminder{
        {
			Method: "email",
			Minutes: 1440, // 1 day before
		},
		{
			Method: "popup",
			Minutes: 10,   // 10 minutes before
		},
    },
    UseDefault:      false,
    ForceSendFields: []string{"UseDefault"},
}

	event := &calendar.Event{
	Summary:  "Invoice Payment Due",
	Start:    &calendar.EventDateTime{DateTime: startTime.Format(time.RFC3339)}, // Format for DateTime
	End:      &calendar.EventDateTime{DateTime: endTime.Format(time.RFC3339)},     // Format for DateTime
	Reminders: reminders,
	}

    _, err = calendarService.Events.Insert("primary", event).Do()
    if err != nil {
        if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == http.StatusUnauthorized {
            return fmt.Errorf("session expired: %v", gErr.Message)
        }
        return err
    }
    return nil
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