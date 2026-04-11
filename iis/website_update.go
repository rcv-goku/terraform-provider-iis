package iis

import (
	"context"
	"encoding/json"
)

func (client Client) UpdateWebsite(ctx context.Context, update Website) (*Website, error) {
	body := BuildPatchBody(update)
	delete(body, "id")
	res, err := httpPatch(ctx, client, client.websitePath(update.ID), body)
	if err != nil {
		return nil, err
	}
	var site Website
	err = json.Unmarshal(res, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}
