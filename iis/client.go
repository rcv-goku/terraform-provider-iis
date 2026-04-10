package iis

import "net/http"

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
