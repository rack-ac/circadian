package whoop

import (
	"context"
	"net/url"
)

func (c *Client) GetSleep(ctx context.Context, sleepID string) (*Sleep, error) {
	var sleep Sleep
	if err := c.getJSON(ctx, "/developer/v2/activity/sleep/"+url.PathEscape(sleepID), nil, &sleep); err != nil {
		return nil, err
	}
	return &sleep, nil
}

func (c *Client) ListSleeps(ctx context.Context, params CollectionParams) (*Collection[Sleep], error) {
	var collection Collection[Sleep]
	if err := c.getJSON(ctx, "/developer/v2/activity/sleep", collectionQuery(params), &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (c *Client) GetWorkout(ctx context.Context, workoutID string) (*Workout, error) {
	var workout Workout
	if err := c.getJSON(ctx, "/developer/v2/activity/workout/"+url.PathEscape(workoutID), nil, &workout); err != nil {
		return nil, err
	}
	return &workout, nil
}

func (c *Client) ListWorkouts(ctx context.Context, params CollectionParams) (*Collection[Workout], error) {
	var collection Collection[Workout]
	if err := c.getJSON(ctx, "/developer/v2/activity/workout", collectionQuery(params), &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}
