package server

import (
	"testing"

	"github.com/rack-ac/circadian/internal/whoop"
)

func TestSummarizeProfile(t *testing.T) {
	t.Parallel()

	summary := summarizePayload(&whoop.Profile{
		UserID:    42,
		Email:     "test@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
	})

	if summary["type"] != "profile" {
		t.Fatalf("type = %v, want profile", summary["type"])
	}
	if summary["full_name"] != "Jane Doe" {
		t.Fatalf("full_name = %v, want Jane Doe", summary["full_name"])
	}
}

func TestSummarizeWorkoutCollection(t *testing.T) {
	t.Parallel()

	summary := summarizePayload(&whoop.Collection[whoop.Workout]{
		Records: []whoop.Workout{
			{Start: "2026-04-10T10:00:00Z"},
			{Start: "2026-04-08T10:00:00Z"},
		},
		NextToken: "next",
	})

	if summary["type"] != "workouts" {
		t.Fatalf("type = %v, want workouts", summary["type"])
	}
	if summary["record_count"] != 2 {
		t.Fatalf("record_count = %v, want 2", summary["record_count"])
	}
	if summary["has_more"] != true {
		t.Fatalf("has_more = %v, want true", summary["has_more"])
	}
}
