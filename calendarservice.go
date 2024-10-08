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


