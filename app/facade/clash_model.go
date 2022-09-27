package facade

import "github.com/microctar/licorice/app/config"

type DNSMode int

const (
	DNSNormal DNSMode = iota
	DNSFakeIP
	DNSMapping
)

type RawDNS struct {
	Enable            bool              `yaml:"enable"`
	IPv6              bool              `yaml:"ipv6"`
	UseHosts          bool              `yaml:"use-hosts"`
	NameServer        []string          `yaml:"nameserver"`
	Fallback          []string          `yaml:"fallback"`
	FallbackFilter    RawFallbackFilter `yaml:"fallback-filter"`
	Listen            string            `yaml:"listen"`
	EnhancedMode      DNSMode           `yaml:"enhanced-mode"`
	FakeIPRange       string            `yaml:"fake-ip-range"`
	FakeIPFilter      []string          `yaml:"fake-ip-filter"`
	DefaultNameserver []string          `yaml:"default-nameserver"`
	NameServerPolicy  map[string]string `yaml:"nameserver-policy"`
}

type RawFallbackFilter struct {
	GeoIP     bool     `yaml:"geoip"`
	GeoIPCode string   `yaml:"geoip-code"`
	IPCIDR    []string `yaml:"ipcidr"`
	Domain    []string `yaml:"domain"`
}

type Experimental struct{}

type Profile struct {
	StoreSelected bool `yaml:"store-selected"`
	StoreFakeIP   bool `yaml:"store-fake-ip"`
}

type TunnelMode int

const (
	Global TunnelMode = iota
	Rule
	Direct
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	SILENT
)

type RawConfig struct {
	Port               int        `yaml:"port,omitempty"`
	SocksPort          int        `yaml:"socks-port,omitempty"`
	RedirPort          int        `yaml:"redir-port,omitempty"`
	TProxyPort         int        `yaml:"tproxy-port,omitempty"`
	MixedPort          int        `yaml:"mixed-port,omitempty"`
	Authentication     []string   `yaml:"authentication,omitempty"`
	AllowLan           bool       `yaml:"allow-lan,omitempty"`
	BindAddress        string     `yaml:"bind-address,omitempty"`
	Mode               TunnelMode `yaml:"mode,omitempty"`
	LogLevel           LogLevel   `yaml:"log-level,omitempty"`
	IPv6               bool       `yaml:"ipv6,omitempty"`
	ExternalController string     `yaml:"external-controller,omitempty"`
	ExternalUI         string     `yaml:"external-ui,omitempty"`
	Secret             string     `yaml:"secret,omitempty"`
	Interface          string     `yaml:"interface-name,omitempty"`
	RoutingMark        int        `yaml:"routing-mark,omitempty"`

	ProxyProvider map[string]map[string]any `yaml:"proxy-providers,omitempty"`
	Hosts         map[string]string         `yaml:"hosts,omitempty"`
	DNS           RawDNS                    `yaml:"dns,omitempty"`
	Experimental  Experimental              `yaml:"experimental,omitempty"`
	Profile       Profile                   `yaml:"profile,omitempty"`
	Proxy         []config.Proxy            `yaml:"proxies,omitempty"`
	ProxyGroup    []map[string]any          `yaml:"proxy-groups,omitempty"`
	Rule          []string                  `yaml:"rules,omitempty"`
}

// Custom MarshalYAML Function

func (dns_mode DNSMode) String() string {
	switch dns_mode {
	case DNSNormal:
		return "DNSNormal"
	case DNSFakeIP:
		return "DNSFakeIP"
	case DNSMapping:
		return "DNSMapping"
	default:
		return "Unknown"
	}
}

func (t TunnelMode) String() string {
	switch t {
	case Global:
		return "global"
	case Rule:
		return "rule"
	case Direct:
		return "direct"
	default:
		return "unknown"
	}
}

func (loglevel LogLevel) String() string {
	switch loglevel {
	case INFO:
		return "info"
	case WARNING:
		return "warning"
	case ERROR:
		return "error"
	case DEBUG:
		return "debug"
	case SILENT:
		return "silent"
	default:
		return "unknown"
	}
}

func (dns_mode DNSMode) MarshalYAML() (any, error) {
	return dns_mode.String(), nil
}

func (t TunnelMode) MarshalYAML() (any, error) {
	return t.String(), nil
}

func (loglevel LogLevel) MarshalYAML() (any, error) {
	return loglevel.String(), nil
}
