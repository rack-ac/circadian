package server

import (
	"fmt"
	"time"
)

func validateCollectionArgs(args collectionArgs) error {
	if args.Limit < 0 {
		return fmt.Errorf("limit must be >= 0")
	}
	if args.Limit > 25 {
		return fmt.Errorf("limit must be <= 25")
	}
	if err := validateRFC3339("start", args.Start); err != nil {
		return err
	}
	if err := validateRFC3339("end", args.End); err != nil {
		return err
	}
	return nil
}

func validateCycleArgs(args cycleArgs) error {
	if err := validateCollectionArgs(args.collectionArgs); err != nil {
		return err
	}
	if args.CycleID > 0 && hasCollectionFilters(args.collectionArgs) {
		return fmt.Errorf("cycle_id cannot be combined with limit, start, end, or next_token")
	}
	return nil
}

func validateRecoveryArgs(args recoveryArgs) error {
	if err := validateCollectionArgs(args.collectionArgs); err != nil {
		return err
	}
	if args.CycleID > 0 && hasCollectionFilters(args.collectionArgs) {
		return fmt.Errorf("cycle_id cannot be combined with limit, start, end, or next_token")
	}
	return nil
}

func validateSleepArgs(args sleepArgs) error {
	if err := validateCollectionArgs(args.collectionArgs); err != nil {
		return err
	}
	if args.SleepID != "" && hasCollectionFilters(args.collectionArgs) {
		return fmt.Errorf("sleep_id cannot be combined with limit, start, end, or next_token")
	}
	return nil
}

func validateWorkoutArgs(args workoutArgs) error {
	if err := validateCollectionArgs(args.collectionArgs); err != nil {
		return err
	}
	if args.WorkoutID != "" && hasCollectionFilters(args.collectionArgs) {
		return fmt.Errorf("workout_id cannot be combined with limit, start, end, or next_token")
	}
	return nil
}

func hasCollectionFilters(args collectionArgs) bool {
	return args.Limit != 0 || args.Start != "" || args.End != "" || args.NextToken != ""
}

func validateRFC3339(name, value string) error {
	if value == "" {
		return nil
	}
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return fmt.Errorf("%s must be a valid RFC3339 timestamp", name)
	}
	return nil
}
