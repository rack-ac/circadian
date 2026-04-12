package whoop

import (
	"encoding/json"
	"testing"
)

func TestRecoveryScoreAcceptsWholeNumberFloats(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
		"cycle_id": 123,
		"sleep_id": "sleep-1",
		"user_id": 456,
		"created_at": "2026-04-12T06:00:00Z",
		"updated_at": "2026-04-12T06:05:00Z",
		"score_state": "SCORED",
		"score": {
			"user_calibrating": false,
			"recovery_score": 51.0,
			"resting_heart_rate": 52.0,
			"hrv_rmssd_milli": 42.7,
			"spo2_percentage": 97.0,
			"skin_temp_celsius": 36.6
		}
	}`)

	var recovery Recovery
	if err := json.Unmarshal(payload, &recovery); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got := int(recovery.Score.RecoveryScore); got != 51 {
		t.Fatalf("recovery.Score.RecoveryScore = %d, want 51", got)
	}
	if got := int(recovery.Score.RestingHeartRate); got != 52 {
		t.Fatalf("recovery.Score.RestingHeartRate = %d, want 52", got)
	}
}

func TestWholeNumberRejectsFractionalFloats(t *testing.T) {
	t.Parallel()

	var value WholeNumber
	if err := json.Unmarshal([]byte(`51.5`), &value); err == nil {
		t.Fatal("json.Unmarshal() error = nil, want non-nil")
	}
}
