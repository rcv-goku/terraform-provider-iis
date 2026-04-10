package iis

import (
	"context"
	"encoding/json"
)

// createAppPoolRequest is the minimal payload for POST - IIS Admin API only accepts name on create
type createAppPoolRequest struct {
	Name                  string `json:"name"`
	ManagedRuntimeVersion string `json:"managed_runtime_version,omitempty"`
}

func (client Client) CreateAppPool(ctx context.Context, req ApplicationPool) (*ApplicationPool, error) {
	// IIS Admin API only accepts name (and optionally managed_runtime_version) on POST
	// All other fields must be set via PATCH after creation
	createReq := createAppPoolRequest{
		Name:                  req.Name,
		ManagedRuntimeVersion: req.ManagedRuntimeVersion,
	}

	res, err := httpPost(ctx, client, "/api/webserver/application-pools", createReq)
	if err != nil {
		// If we get a 409 Conflict, the app pool already exists
		// Try to retrieve it by name instead
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
