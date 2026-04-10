package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

type UpdateVirtualDirectoryRequest struct {
	Path         string `json:"path,omitempty"`
	PhysicalPath string `json:"physical_path,omitempty"`
}

func (client Client) UpdateVirtualDirectory(ctx context.Context, id string, req UpdateVirtualDirectoryRequest) (*VirtualDirectory, error) {
	url := fmt.Sprintf("/api/webserver/virtual-directories/%s", id)
	res, err := httpPatch(ctx, client, url, req)
	if err != nil {
		return nil, err
	}
	var vdir VirtualDirectory
	if err = json.Unmarshal(res, &vdir); err != nil {
		return nil, err
	}
	return &vdir, nil
}
