package iis

import (
	"context"
	"encoding/json"
)

func (client Client) UpdateAppPool(ctx context.Context, id string, pool ApplicationPool) (*ApplicationPool, error) {
	body := BuildPatchBody(pool)
	delete(body, "id")
	delete(body, "name")
	res, err := httpPatch(ctx, client, client.appPoolPath(id), body)
	if err != nil {
		return nil, err
	}
	var updated ApplicationPool
	err = json.Unmarshal(res, &updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
