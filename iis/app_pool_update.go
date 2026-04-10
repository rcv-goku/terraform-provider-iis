package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (client Client) UpdateAppPool(ctx context.Context, id string, pool ApplicationPool) (*ApplicationPool, error) {
	var url string
	if client.AgentMode {
		url = fmt.Sprintf("/api/app-pools?name=%s", id)
	} else {
		url = fmt.Sprintf("/api/webserver/application-pools/%s", id)
	}
	body := BuildPatchBody(pool)
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
