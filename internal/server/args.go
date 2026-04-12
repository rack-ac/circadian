package server

import "github.com/rack-ac/circadian/internal/whoop"

type collectionArgs struct {
	Limit     int    `json:"limit,omitempty" jsonschema:"Max records to return, up to 25."`
	Start     string `json:"start,omitempty" jsonschema:"RFC3339 timestamp inclusive lower bound."`
	End       string `json:"end,omitempty" jsonschema:"RFC3339 timestamp exclusive upper bound."`
	NextToken string `json:"next_token,omitempty" jsonschema:"Pagination token from the prior response."`
}

type cycleArgs struct {
	CycleID int64 `json:"cycle_id,omitempty" jsonschema:"Optional WHOOP cycle ID. If omitted, returns a collection."`
	collectionArgs
}

type recoveryArgs struct {
	CycleID int64 `json:"cycle_id,omitempty" jsonschema:"Optional cycle ID to fetch a single recovery. If omitted, returns a collection."`
	collectionArgs
}

type sleepArgs struct {
	SleepID string `json:"sleep_id,omitempty" jsonschema:"Optional WHOOP sleep UUID. If omitted, returns a collection."`
	collectionArgs
}

type workoutArgs struct {
	WorkoutID string `json:"workout_id,omitempty" jsonschema:"Optional WHOOP workout UUID. If omitted, returns a collection."`
	collectionArgs
}

func toCollectionParams(args collectionArgs) whoop.CollectionParams {
	return whoop.CollectionParams{
		Limit:     args.Limit,
		Start:     args.Start,
		End:       args.End,
		NextToken: args.NextToken,
	}
}
