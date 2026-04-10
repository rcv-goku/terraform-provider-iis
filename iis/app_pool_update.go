package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (client Client) UpdateAppPool(ctx context.Context, id string, pool ApplicationPool) (*ApplicationPool, error) {
	url := fmt.Sprintf("/api/webserver/application-pools/%s", id)
	// BuildPatchBody strips zero/empty values to prevent IIS API 500 errors
	body := BuildPatchBody(pool)
	// Remove fields the API doesn't accept on PATCH
	delete(body, "id")
	delete(body, "name")
	res, err := httpPatch(ctx, client, url, body)
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
