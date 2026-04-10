package iis

import (
	"context"
	"encoding/json"
)

type CreateVirtualDirectoryRequest struct {
	Path         string    `json:"path"`
	PhysicalPath string    `json:"physical_path"`
	WebApp       Reference `json:"webapp,omitempty"`
	Website      Reference `json:"website"`
}

func (client Client) CreateVirtualDirectory(ctx context.Context, req CreateVirtualDirectoryRequest) (*VirtualDirectory, error) {
	res, err := httpPost(ctx, client, "/api/webserver/virtual-directories", req)
	if err != nil {
		return nil, err
	}
	var vdir VirtualDirectory
	if err = json.Unmarshal(res, &vdir); err != nil {
		return nil, err
	}
	return &vdir, nil
}
