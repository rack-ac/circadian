package server

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rack-ac/circadian/internal/whoop"
)

func shapedResult(payload any, err error) (*mcp.CallToolResult, any, error) {
	if err != nil {
		return nil, nil, err
	}

	envelope := map[string]any{
		"summary": summarizePayload(payload),
		"data":    payload,
	}

	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal tool output: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func summarizePayload(payload any) map[string]any {
	switch v := payload.(type) {
	case *whoop.Profile:
		return map[string]any{
			"type":      "profile",
			"full_name": stringsTrimSpace(v.FirstName + " " + v.LastName),
			"email":     v.Email,
			"user_id":   v.UserID,
		}
	case *whoop.BodyMeasurement:
		return map[string]any{
			"type":              "body_measurement",
			"height_meter":      v.HeightMeter,
			"weight_kilogram":   v.WeightKilogram,
			"max_heart_rate":    v.MaxHeartRate,
			"weight_pounds":     round(v.WeightKilogram * 2.20462),
			"height_centimeter": round(v.HeightMeter * 100),
		}
	case *whoop.Cycle:
		return map[string]any{
			"type":     "cycle",
			"cycle_id": v.ID,
			"start":    v.Start,
			"end":      v.End,
			"strain":   valueOrZero(v.Score, func(s *whoop.CycleScore) float64 { return s.Strain }),
			"score":    v.ScoreState,
			"duration": durationSummary(v.Start, v.End),
		}
	case *whoop.Recovery:
		return map[string]any{
			"type":               "recovery",
			"cycle_id":           v.CycleID,
			"sleep_id":           v.SleepID,
			"recovery_score":     valueOrZero(v.Score, func(s *whoop.RecoveryScore) int { return int(s.RecoveryScore) }),
			"resting_heart_rate": valueOrZero(v.Score, func(s *whoop.RecoveryScore) int { return int(s.RestingHeartRate) }),
			"hrv_rmssd_milli":    valueOrZero(v.Score, func(s *whoop.RecoveryScore) float64 { return s.HRVRMSSDMilli }),
			"score":              v.ScoreState,
		}
	case *whoop.Sleep:
		return map[string]any{
			"type":                         "sleep",
			"sleep_id":                     v.ID,
			"cycle_id":                     v.CycleID,
			"start":                        v.Start,
			"end":                          v.End,
			"duration":                     durationSummary(v.Start, v.End),
			"sleep_performance_percentage": valueOrZero(v.Score, func(s *whoop.SleepScore) float64 { return s.SleepPerformancePercentage }),
			"sleep_efficiency_percentage":  valueOrZero(v.Score, func(s *whoop.SleepScore) float64 { return s.SleepEfficiencyPercentage }),
			"nap":                          v.Nap,
			"score":                        v.ScoreState,
		}
	case *whoop.Workout:
		return map[string]any{
			"type":               "workout",
			"workout_id":         v.ID,
			"sport_name":         v.SportName,
			"start":              v.Start,
			"end":                v.End,
			"duration":           durationSummary(v.Start, v.End),
			"strain":             valueOrZero(v.Score, func(s *whoop.WorkoutScore) float64 { return s.Strain }),
			"average_heart_rate": valueOrZero(v.Score, func(s *whoop.WorkoutScore) int { return s.AverageHeartRate }),
			"distance_meter":     valueOrZero(v.Score, func(s *whoop.WorkoutScore) float64 { return s.DistanceMeter }),
			"score":              v.ScoreState,
		}
	case *whoop.Collection[whoop.Cycle]:
		return summarizeCollection("cycles", len(v.Records), v.NextToken, firstLastCycleWindow(v.Records))
	case *whoop.Collection[whoop.Recovery]:
		return summarizeCollection("recoveries", len(v.Records), v.NextToken, nil)
	case *whoop.Collection[whoop.Sleep]:
		return summarizeCollection("sleeps", len(v.Records), v.NextToken, firstLastSleepWindow(v.Records))
	case *whoop.Collection[whoop.Workout]:
		return summarizeCollection("workouts", len(v.Records), v.NextToken, firstLastWorkoutWindow(v.Records))
	case map[string]any:
		return v
	default:
		return map[string]any{"type": "result"}
	}
}

func summarizeCollection(kind string, count int, nextToken string, extra map[string]any) map[string]any {
	summary := map[string]any{
		"type":         kind,
		"record_count": count,
		"has_more":     nextToken != "",
	}
	if nextToken != "" {
		summary["next_token"] = nextToken
	}
	for key, value := range extra {
		summary[key] = value
	}
	return summary
}

func firstLastCycleWindow(records []whoop.Cycle) map[string]any {
	if len(records) == 0 {
		return nil
	}
	return map[string]any{
		"latest_start": records[0].Start,
		"oldest_start": records[len(records)-1].Start,
	}
}

func firstLastSleepWindow(records []whoop.Sleep) map[string]any {
	if len(records) == 0 {
		return nil
	}
	return map[string]any{
		"latest_start": records[0].Start,
		"oldest_start": records[len(records)-1].Start,
	}
}

func firstLastWorkoutWindow(records []whoop.Workout) map[string]any {
	if len(records) == 0 {
		return nil
	}
	return map[string]any{
		"latest_start": records[0].Start,
		"oldest_start": records[len(records)-1].Start,
	}
}

func durationSummary(start, end string) string {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return ""
	}
	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return ""
	}
	if endTime.Before(startTime) {
		return ""
	}
	return endTime.Sub(startTime).String()
}

func stringsTrimSpace(value string) string {
	return strings.TrimSpace(value)
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}

func valueOrZero[T any, R int | float64](value *T, getter func(*T) R) R {
	var zero R
	if value == nil {
		return zero
	}
	return getter(value)
}
