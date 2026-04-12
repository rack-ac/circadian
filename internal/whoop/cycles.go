package whoop

import (
	"context"
	"fmt"
)

func (c *Client) GetCycle(ctx context.Context, cycleID int64) (*Cycle, error) {
	var cycle Cycle
	if err := c.getJSON(ctx, fmt.Sprintf("/developer/v2/cycle/%d", cycleID), nil, &cycle); err != nil {
		return nil, err
	}
	return &cycle, nil
}

func (c *Client) ListCycles(ctx context.Context, params CollectionParams) (*Collection[Cycle], error) {
	var collection Collection[Cycle]
	if err := c.getJSON(ctx, "/developer/v2/cycle", collectionQuery(params), &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}

func (c *Client) GetRecoveryForCycle(ctx context.Context, cycleID int64) (*Recovery, error) {
	var recovery Recovery
	if err := c.getJSON(ctx, fmt.Sprintf("/developer/v2/cycle/%d/recovery", cycleID), nil, &recovery); err != nil {
		return nil, err
	}
	return &recovery, nil
}

func (c *Client) ListRecoveries(ctx context.Context, params CollectionParams) (*Collection[Recovery], error) {
	var collection Collection[Recovery]
	if err := c.getJSON(ctx, "/developer/v2/recovery", collectionQuery(params), &collection); err != nil {
		return nil, err
	}
	return &collection, nil
}
