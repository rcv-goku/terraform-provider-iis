package iis

import (
	"fmt"
	"net/http"
)

type Client struct {
	HttpClient http.Client
	Host       string
	AccessKey  string
	AgentMode  bool // true = iis-agent (simple Bearer), false = MS IIS Admin API (NTLM)

	// NTLM Authentication fields (legacy MS IIS Admin API only)
	NTLMUsername string
	NTLMPassword string
	NTLMDomain   string

	// API path prefix — "/api/webserver" for MS API, "/api" for iis-agent
	APIPrefix string
}

func (c Client) appPoolsPath() string {
	if c.AgentMode {
		return "/api/app-pools"
	}
	return "/api/webserver/application-pools"
}

func (c Client) appPoolPath(id string) string {
	if c.AgentMode {
		return fmt.Sprintf("/api/app-pools?name=%s", id)
	}
	return fmt.Sprintf("/api/webserver/application-pools/%s", id)
}

func (c Client) websitesPath() string {
	if c.AgentMode {
		return "/api/websites"
	}
	return "/api/webserver/websites"
}

func (c Client) websitePath(id string) string {
	if c.AgentMode {
		return fmt.Sprintf("/api/websites?name=%s", id)
	}
	return fmt.Sprintf("/api/webserver/websites/%s", id)
}
