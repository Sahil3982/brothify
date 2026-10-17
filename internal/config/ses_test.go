package config

import (
	"strings"
	"testing"

	"github.com/brothify/internal/models"
	"github.com/google/uuid"
)

func TestBuildFeedbackURLAddsReservationID(t *testing.T) {
	reservationID := uuid.New()

	got, err := buildFeedbackURL("https://example.com/feedback?source=email", reservationID)
	if err != nil {
		t.Fatalf("buildFeedbackURL() error = %v", err)
	}
	want := "https://example.com/feedback?reservation_id=" + reservationID.String() + "&source=email"
	if got != want {
		t.Fatalf("buildFeedbackURL() = %q, want %q", got, want)
	}
}

func TestBuildFeedbackURLRejectsNonHTTPURL(t *testing.T) {
	if _, err := buildFeedbackURL("javascript:alert(1)", uuid.New()); err == nil {
		t.Fatal("buildFeedbackURL() error = nil, want an absolute HTTP(S) URL error")
	}
}

func TestBuildFeedbackReminderHTMLRendersTemplateAndFrontendLink(t *testing.T) {
	reservationID := uuid.New()
	reservation := &models.Reservation{
		ID:                    reservationID,
		RESERVATIONPERSONNAME: "Avery",
		RESERVATIONDATE:       "2026-10-11",
		RESERVATIONTIME:       "19:30",
	}

	html, err := buildFeedbackReminderHTML(reservation, "https://example.com/feedback")
	if err != nil {
		t.Fatalf("buildFeedbackReminderHTML() error = %v", err)
	}
	for _, want := range []string{
		"Hi Avery,",
		"Review Your Visit",
		"https://example.com/feedback?reservation_id=" + reservationID.String(),
		"2026-10-11",
		"19:30",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered feedback email does not contain %q", want)
		}
	}
}
