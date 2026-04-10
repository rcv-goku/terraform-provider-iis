package iis

import (
	"context"
	"encoding/json"
)

// createAppPoolRequest is the minimal payload for POST
type createAppPoolRequest struct {
	Name                  string `json:"name"`
	ManagedRuntimeVersion string `json:"managed_runtime_version,omitempty"`
}

func (client Client) CreateAppPool(ctx context.Context, req ApplicationPool) (*ApplicationPool, error) {
	endpoint := "/api/webserver/application-pools"
	if client.AgentMode {
		endpoint = "/api/app-pools"
	}

	createReq := createAppPoolRequest{
		Name:                  req.Name,
		ManagedRuntimeVersion: req.ManagedRuntimeVersion,
	}

	res, err := httpPost(ctx, client, endpoint, createReq)
	if err != nil {
		if IsConflictError(err) {
			pool, getErr := client.GetAppPoolByName(ctx, req.Name)
			if getErr == nil {
				return pool, nil
			}
		}
		return nil, err
	}
	var pool ApplicationPool
	err = json.Unmarshal(res, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}
