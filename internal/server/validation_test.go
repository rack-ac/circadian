package server

import "testing"

func TestValidateCollectionArgs(t *testing.T) {
	t.Parallel()

	if err := validateCollectionArgs(collectionArgs{Limit: 26}); err == nil {
		t.Fatal("expected limit validation error")
	}
	if err := validateCollectionArgs(collectionArgs{Start: "not-a-time"}); err == nil {
		t.Fatal("expected RFC3339 validation error")
	}
	if err := validateCollectionArgs(collectionArgs{Limit: 25, Start: "2026-04-12T10:00:00Z"}); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateSingleRecordArgsRejectCollectionFilters(t *testing.T) {
	t.Parallel()

	if err := validateSleepArgs(sleepArgs{
		SleepID: "abc",
		collectionArgs: collectionArgs{
			Limit: 10,
		},
	}); err == nil {
		t.Fatal("expected sleep_id filter conflict")
	}
}
