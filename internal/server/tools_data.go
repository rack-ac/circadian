package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDataTools(server *mcp.Server, svc services) {
	mcp.AddTool(server, &mcp.Tool{Name: "whoop_profile", Description: "Get the authenticated WHOOP user's basic profile."}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		resp, err := svc.client.GetProfile(ctx)
		return shapedResult(resp, err)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_body_measurements", Description: "Get the authenticated WHOOP user's height, weight, and max heart rate."}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		resp, err := svc.client.GetBodyMeasurement(ctx)
		return shapedResult(resp, err)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_cycles", Description: "Get one cycle by cycle_id or list cycles."}, func(ctx context.Context, req *mcp.CallToolRequest, args cycleArgs) (*mcp.CallToolResult, any, error) {
		if err := validateCycleArgs(args); err != nil {
			return nil, nil, err
		}
		if args.CycleID > 0 {
			resp, err := svc.client.GetCycle(ctx, args.CycleID)
			return shapedResult(resp, err)
		}
		resp, err := svc.client.ListCycles(ctx, toCollectionParams(args.collectionArgs))
		return shapedResult(resp, err)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_recoveries", Description: "Get one recovery by cycle_id or list recoveries."}, func(ctx context.Context, req *mcp.CallToolRequest, args recoveryArgs) (*mcp.CallToolResult, any, error) {
		if err := validateRecoveryArgs(args); err != nil {
			return nil, nil, err
		}
		if args.CycleID > 0 {
			resp, err := svc.client.GetRecoveryForCycle(ctx, args.CycleID)
			return shapedResult(resp, err)
		}
		resp, err := svc.client.ListRecoveries(ctx, toCollectionParams(args.collectionArgs))
		return shapedResult(resp, err)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_sleeps", Description: "Get one sleep by sleep_id or list sleeps."}, func(ctx context.Context, req *mcp.CallToolRequest, args sleepArgs) (*mcp.CallToolResult, any, error) {
		if err := validateSleepArgs(args); err != nil {
			return nil, nil, err
		}
		if args.SleepID != "" {
			resp, err := svc.client.GetSleep(ctx, args.SleepID)
			return shapedResult(resp, err)
		}
		resp, err := svc.client.ListSleeps(ctx, toCollectionParams(args.collectionArgs))
		return shapedResult(resp, err)
	})

	mcp.AddTool(server, &mcp.Tool{Name: "whoop_workouts", Description: "Get one workout by workout_id or list workouts."}, func(ctx context.Context, req *mcp.CallToolRequest, args workoutArgs) (*mcp.CallToolResult, any, error) {
		if err := validateWorkoutArgs(args); err != nil {
			return nil, nil, err
		}
		if args.WorkoutID != "" {
			resp, err := svc.client.GetWorkout(ctx, args.WorkoutID)
			return shapedResult(resp, err)
		}
		resp, err := svc.client.ListWorkouts(ctx, toCollectionParams(args.collectionArgs))
		return shapedResult(resp, err)
	})
}
