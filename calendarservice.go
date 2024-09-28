package main

import (
    "context"
	"net/http"
	"fmt"

	"google.golang.org/api/googleapi"
    "google.golang.org/api/calendar/v3"
	"golang.org/x/oauth2"

)

func addEventToCalendar(token *oauth2.Token) error {
    ctx := context.Background()
    client := googleOAuthConfig.Client(ctx, token)

    calendarService, err := calendar.New(client)
    if err != nil {
        return err
    }

    event := &calendar.Event{
        Summary: "Sample Event",
        Start: &calendar.EventDateTime{
            DateTime: "2024-09-30T10:00:00-07:00", // Example start time
            TimeZone: "America/Los_Angeles",
        },
        End: &calendar.EventDateTime{
            DateTime: "2024-09-30T11:00:00-07:00", // Example end time
            TimeZone: "America/Los_Angeles",
        },
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
