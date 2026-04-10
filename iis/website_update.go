package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

// websiteUpdateRequest uses omitempty on all fields so only user-specified values are sent.
// The IIS Admin API returns 500 when receiving zero-value fields it doesn't expect.
type websiteLimitsUpdate struct {
	ConnectionTimeout *int64 `json:"connection_timeout,omitempty"`
	MaxBandwidth      *int64 `json:"max_bandwidth,omitempty"`
	MaxConnections    *int64 `json:"max_connections,omitempty"`
	MaxUrlSegments    *int64 `json:"max_url_segments,omitempty"`
}

type websiteUpdateRequest struct {
	Name             string               `json:"name,omitempty"`
	Status           string               `json:"status,omitempty"`
	PhysicalPath     string               `json:"physical_path,omitempty"`
	Bindings         []WebsiteBinding     `json:"bindings,omitempty"`
	ApplicationPool  *ApplicationReference `json:"application_pool,omitempty"`
	EnabledProtocols string               `json:"enabled_protocols,omitempty"`
	ServerAutoStart  *bool                `json:"server_auto_start,omitempty"`
	Limits           *websiteLimitsUpdate `json:"limits,omitempty"`
}

// toWebsiteUpdateRequest converts a Website to an update request, only including non-zero fields.
func toWebsiteUpdateRequest(site Website) websiteUpdateRequest {
	req := websiteUpdateRequest{
		Name:         site.Name,
		Status:       site.Status,
		PhysicalPath: site.PhysicalPath,
	}

	if len(site.Bindings) > 0 {
		req.Bindings = site.Bindings
	}

	if site.ApplicationPool.ID != "" {
		appPool := site.ApplicationPool
		req.ApplicationPool = &appPool
	}

	if site.EnabledProtocols != "" {
		req.EnabledProtocols = site.EnabledProtocols
	}

	if site.ServerAutoStart {
		req.ServerAutoStart = &site.ServerAutoStart
	}

	emptyLimits := WebsiteLimits{}
	if site.Limits != emptyLimits {
		limits := &websiteLimitsUpdate{}
		if site.Limits.ConnectionTimeout > 0 {
			limits.ConnectionTimeout = &site.Limits.ConnectionTimeout
		}
		if site.Limits.MaxBandwidth > 0 {
			limits.MaxBandwidth = &site.Limits.MaxBandwidth
		}
		if site.Limits.MaxConnections > 0 {
			limits.MaxConnections = &site.Limits.MaxConnections
		}
		if site.Limits.MaxUrlSegments > 0 {
			limits.MaxUrlSegments = &site.Limits.MaxUrlSegments
		}
		req.Limits = limits
	}

	return req
}

func (client Client) UpdateWebsite(ctx context.Context, update Website) (*Website, error) {
	url := fmt.Sprintf("/api/webserver/websites/%s", update.ID)
	updateReq := toWebsiteUpdateRequest(update)
	res, err := httpPatch(ctx, client, url, updateReq)
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
