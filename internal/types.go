package internal

type Configuration struct {
	HTTP interface{} `json:"http,omitempty"`
	TCP  interface{} `json:"tcp,omitempty"`
	UDP  interface{} `json:"udp,omitempty"`
	TLS  interface{} `json:"tls,omitempty"`
}

type Configurations map[string]*Configuration
