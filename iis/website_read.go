package iis

import "context"

func (client Client) ReadWebsite(ctx context.Context, id string) (*Website, error) {
	var site Website
	if err := getJson(ctx, client, client.websitePath(id), &site); err != nil {
		return nil, err
	}
	return &site, nil
}
