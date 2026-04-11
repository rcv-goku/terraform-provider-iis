package iis

type Website struct {
	Name             string               `json:"name"`
	ID               string               `json:"id"`
	Status           string               `json:"status"`
	PhysicalPath     string               `json:"physical_path"`
	Bindings         []WebsiteBinding     `json:"bindings"`
	ApplicationPool  ApplicationReference `json:"application_pool"`
	EnabledProtocols string               `json:"enabled_protocols"`
	ServerAutoStart  bool                 `json:"server_auto_start"`
	Limits           WebsiteLimits        `json:"limits"`
}

type WebsiteBinding struct {
	Protocol    string             `json:"protocol"`
	Port        int                `json:"port"`
	IPAddress   string             `json:"ip_address"`
	Hostname    string             `json:"hostname"`
	Certificate BindingCertificate `json:"certificate"`
}

type BindingCertificate struct {
	ID                string `json:"id"`
	CertificateHash   string `json:"certificate_hash,omitempty"`
	CertificateStore  string `json:"certificate_store_name,omitempty"`
}

type WebsiteLimits struct {
	ConnectionTimeout int64 `json:"connection_timeout"`
	MaxBandwidth      int64 `json:"max_bandwidth"`
	MaxConnections    int64 `json:"max_connections"`
	MaxUrlSegments    int64 `json:"max_url_segments"`
}
