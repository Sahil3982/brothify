package services

import (
	"testing"
	"time"
)

func TestFeedbackReminderDue(t *testing.T) {
	reservationAt := time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{
			name: "before six hours",
			now:  reservationAt.Add(6*time.Hour - time.Second),
			want: false,
		},
		{
			name: "at six hours",
			now:  reservationAt.Add(6 * time.Hour),
			want: true,
		},
		{
			name: "after six hours",
			now:  reservationAt.Add(7 * time.Hour),
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := feedbackReminderDue(reservationAt, test.now); got != test.want {
				t.Fatalf("feedbackReminderDue() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestParseReservationDateTimeUsesConfiguredLocation(t *testing.T) {
	location := time.FixedZone("reservation", 5*60*60+30*60)

	got, err := parseReservationDateTime("2026-10-11", "7:45 PM", location)
	if err != nil {
		t.Fatalf("parseReservationDateTime() error = %v", err)
	}
	if got.Location() != location {
		t.Fatalf("parseReservationDateTime() location = %v, want %v", got.Location(), location)
	}
	if got.Hour() != 19 || got.Minute() != 45 {
		t.Fatalf("parseReservationDateTime() = %v, want 19:45", got)
	}
}
