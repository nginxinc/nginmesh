package nginx

// Config represents complete NGINX configuration.
type Config struct {
	Main        Main
	HTTPConfigs []HTTPConfig
	TCPConfigs  []TCPConfig
}

// Main represents the main configuration (nginx.conf file).
type Main struct {
	Mixer           *MainMixer
	DestinationMaps []DestinationMap
	PodIP           string
	ServiceNode     string
	ServiceCluster  string
	Tracing         bool
	CollectorServer string
	LOGLEVEL		string
}

// DestinationMap is a map between remote destination and the corresponding local destination.
type DestinationMap struct {
	Remote string
	Local  string
}

// HTTPConfig is HTTP NGINX configuration -- a combination of virtual servers, upstreams and maps.
type HTTPConfig struct {
	Name           string
	Upstreams      []Upstream
	VirtualServers []VirtualServer
	Mixer          *HTTPMixer
	Maps           []Map
}

// VirtualServer is configuration for one virtual server.
type VirtualServer struct {
	Address      string
	Port         string
	Names        []string
	Locations    []Location
	IsTarget     bool
	SplitClients []SplitClient
	SSL          *VirtualServerSSL
}

// VirtualServerSSL holds SSL-related configuration for a virtual server.
type VirtualServerSSL struct {
	Certificate        string
	Key                string
	TrustedCertificate string
}

// HTTPMixer contains Mixer configuration for a virtual server.
type HTTPMixer struct {
	SourceIP           string
	SourceUID          string
	SourceLabels	   map[string]string
	DestinationIP      string
	DestinationService string
	DestinationUID     string
	DestinationLabels  map[string]string
}

// MainMixer is the global Mixer configuration.
type MainMixer struct {
	MixerServer string
	MixerPort   string
}

// Location represents a location.
type Location struct {
	Internal          bool
	Path              string
	Expressions       []*Expression
	Upstream          string
	SSL               *LocationSSL
	Sets              []Set
	ConnectTimeout    int64
	ProxyNextUpstream *ProxyNextUpstream
	Redirect          *Redirect
	Rewrite           *Rewrite
	Host              string
	MixerCheck        bool
	MixerReport       bool
	Tracing           bool
	CollectorTopic	  string
}

// Rewrite is configuration for rewriting a URL.
type Rewrite struct {
	Prefix      string
	Replacement string
}

// Redirect is configuration for redirecting an HTTP request.
type Redirect struct {
	Code int
	URL  string
}

// ProxyNextUpstream holds configuration for the proxy_next_upstream set of directives.
type ProxyNextUpstream struct {
	Condition string
	Timeout   int64
	Tries     int
}

// Set represents the set directive.
type Set struct {
	Variable string
	Value    string
}

// LocationSSL is SSL configuration for connections between NGINX and upstreams.
type LocationSSL struct {
	Certificate        string
	Key                string
	TrustedCertificate string
	Name               string
}

// Expression represents the if directive.
type Expression struct {
	Result    string
	Condition string
}

// Upstream is configuration for an upstream.
type Upstream struct {
	Name     string
	Servers  []Server
	LBMethod string
}

// These constants represent load balancing algorithms.
const (
	LBMethodRR        = ""
	LBMethodRandom    = "hash $pid$request_id"
	LBMethodLeastConn = "least_conn"
)

// Server represents an upstream server.
type Server struct {
	Address     string
	Port        string
	MaxConns    int
	MaxFails    int
	FailTimeout int64
}

// SplitClient is configuration for the split_client directive.
type SplitClient struct {
	Variable      string
	Distributions []Distribution
}

// Distribution represents a distribution in the split_client directive.
type Distribution struct {
	Weight int
	Value  string
}

// TCPConfig is TCP configuration - a combination of servers and upstreams.
type TCPConfig struct {
	Name      string
	Upstreams []Upstream
	Servers   []TCPServer
}

// TCPServer holds configuration for a TCP server.
type TCPServer struct {
	Address        string
	Port           string
	Upstream       string
	ConnectTimeout int64
}

// Map is configuration for the map directive.
type Map struct {
	Source   string
	Variable string
	Params   map[string]string
}
