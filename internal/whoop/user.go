package whoop

import "context"

func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	var profile Profile
	if err := c.getJSON(ctx, "/developer/v2/user/profile/basic", nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (c *Client) GetBodyMeasurement(ctx context.Context) (*BodyMeasurement, error) {
	var measurement BodyMeasurement
	if err := c.getJSON(ctx, "/developer/v2/user/measurement/body", nil, &measurement); err != nil {
		return nil, err
	}
	return &measurement, nil
}
