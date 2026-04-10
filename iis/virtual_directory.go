package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

type VirtualDirectory struct {
	ID           string               `json:"id"`
	Path         string               `json:"path"`
	PhysicalPath string               `json:"physical_path"`
	Location     string               `json:"location"`
	WebApp       ApplicationReference `json:"webapp"`
	Website      ApplicationReference `json:"website"`
}

func (v *VirtualDirectory) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

func (client Client) ReadVirtualDirectory(ctx context.Context, id string) (*VirtualDirectory, error) {
	url := fmt.Sprintf("/api/webserver/virtual-directories/%s", id)
	var vdir VirtualDirectory
	if err := getJson(ctx, client, url, &vdir); err != nil {
		return nil, err
	}
	return &vdir, nil
}

func (client Client) DeleteVirtualDirectory(ctx context.Context, id string) error {
	url := fmt.Sprintf("/api/webserver/virtual-directories/%s", id)
	return httpDelete(ctx, client, url)
}
