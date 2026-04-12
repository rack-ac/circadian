package whoop

import (
	"encoding/json"
	"fmt"
	"math"
)

type WholeNumber int

func (n *WholeNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*n = WholeNumber(intValue)
		return nil
	}

	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err == nil {
		rounded := math.Round(floatValue)
		if math.Abs(floatValue-rounded) > 1e-9 {
			return fmt.Errorf("expected whole number, got %s", data)
		}
		*n = WholeNumber(int(rounded))
		return nil
	}

	return fmt.Errorf("expected whole number, got %s", data)
}

type Collection[T any] struct {
	Records   []T    `json:"records"`
	NextToken string `json:"next_token,omitempty"`
}

type Profile struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type BodyMeasurement struct {
	HeightMeter    float64 `json:"height_meter"`
	WeightKilogram float64 `json:"weight_kilogram"`
	MaxHeartRate   int     `json:"max_heart_rate"`
}

type Cycle struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	Start          string      `json:"start"`
	End            string      `json:"end"`
	TimezoneOffset string      `json:"timezone_offset"`
	ScoreState     string      `json:"score_state"`
	Score          *CycleScore `json:"score,omitempty"`
}

type CycleScore struct {
	Strain           float64 `json:"strain"`
	Kilojoule        float64 `json:"kilojoule"`
	AverageHeartRate int     `json:"average_heart_rate"`
	MaxHeartRate     int     `json:"max_heart_rate"`
}

type Recovery struct {
	CycleID    int64          `json:"cycle_id"`
	SleepID    string         `json:"sleep_id"`
	UserID     int64          `json:"user_id"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	ScoreState string         `json:"score_state"`
	Score      *RecoveryScore `json:"score,omitempty"`
}

type RecoveryScore struct {
	UserCalibrating  bool        `json:"user_calibrating"`
	RecoveryScore    WholeNumber `json:"recovery_score"`
	RestingHeartRate WholeNumber `json:"resting_heart_rate"`
	HRVRMSSDMilli    float64     `json:"hrv_rmssd_milli"`
	SPO2Percentage   float64     `json:"spo2_percentage"`
	SkinTempCelsius  float64     `json:"skin_temp_celsius"`
}

type Sleep struct {
	ID             string      `json:"id"`
	CycleID        int64       `json:"cycle_id"`
	V1ID           int64       `json:"v1_id"`
	UserID         int64       `json:"user_id"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	Start          string      `json:"start"`
	End            string      `json:"end"`
	TimezoneOffset string      `json:"timezone_offset"`
	Nap            bool        `json:"nap"`
	ScoreState     string      `json:"score_state"`
	Score          *SleepScore `json:"score,omitempty"`
}

type SleepScore struct {
	StageSummary               *SleepStageSummary `json:"stage_summary,omitempty"`
	SleepNeeded                *SleepNeeded       `json:"sleep_needed,omitempty"`
	RespiratoryRate            float64            `json:"respiratory_rate"`
	SleepPerformancePercentage float64            `json:"sleep_performance_percentage"`
	SleepConsistencyPercentage float64            `json:"sleep_consistency_percentage"`
	SleepEfficiencyPercentage  float64            `json:"sleep_efficiency_percentage"`
}

type SleepStageSummary struct {
	TotalInBedTimeMilli      int64 `json:"total_in_bed_time_milli"`
	TotalAwakeTimeMilli      int64 `json:"total_awake_time_milli"`
	TotalNoDataTimeMilli     int64 `json:"total_no_data_time_milli"`
	TotalLightSleepTimeMilli int64 `json:"total_light_sleep_time_milli"`
	TotalSlowWaveSleepMilli  int64 `json:"total_slow_wave_sleep_time_milli"`
	TotalREMSleepTimeMilli   int64 `json:"total_rem_sleep_time_milli"`
	SleepCycleCount          int   `json:"sleep_cycle_count"`
	DisturbanceCount         int   `json:"disturbance_count"`
}

type SleepNeeded struct {
	BaselineMilli             int64 `json:"baseline_milli"`
	NeedFromSleepDebtMilli    int64 `json:"need_from_sleep_debt_milli"`
	NeedFromRecentStrainMilli int64 `json:"need_from_recent_strain_milli"`
	NeedFromRecentNapMilli    int64 `json:"need_from_recent_nap_milli"`
}

type Workout struct {
	ID             string        `json:"id"`
	V1ID           int64         `json:"v1_id"`
	UserID         int64         `json:"user_id"`
	CreatedAt      string        `json:"created_at"`
	UpdatedAt      string        `json:"updated_at"`
	Start          string        `json:"start"`
	End            string        `json:"end"`
	TimezoneOffset string        `json:"timezone_offset"`
	SportName      string        `json:"sport_name"`
	SportID        int           `json:"sport_id"`
	ScoreState     string        `json:"score_state"`
	Score          *WorkoutScore `json:"score,omitempty"`
}

type WorkoutScore struct {
	Strain              float64        `json:"strain"`
	AverageHeartRate    int            `json:"average_heart_rate"`
	MaxHeartRate        int            `json:"max_heart_rate"`
	Kilojoule           float64        `json:"kilojoule"`
	PercentRecorded     float64        `json:"percent_recorded"`
	DistanceMeter       float64        `json:"distance_meter"`
	AltitudeGainMeter   float64        `json:"altitude_gain_meter"`
	AltitudeChangeMeter float64        `json:"altitude_change_meter"`
	ZoneDurations       *ZoneDurations `json:"zone_durations,omitempty"`
}

type ZoneDurations struct {
	ZoneZeroMilli  int64 `json:"zone_zero_milli"`
	ZoneOneMilli   int64 `json:"zone_one_milli"`
	ZoneTwoMilli   int64 `json:"zone_two_milli"`
	ZoneThreeMilli int64 `json:"zone_three_milli"`
	ZoneFourMilli  int64 `json:"zone_four_milli"`
	ZoneFiveMilli  int64 `json:"zone_five_milli"`
}
