package protocol

type ProxyType uint16

const (
	Unknown ProxyType = iota
	Shadowsocks
)

func (proxy_type ProxyType) String() string {
	switch proxy_type {
	case Shadowsocks:
		return "ss"
	default:
		return "unknown"
	}
}

func (target ProxyType) MarshalYAML() (any, error) {
	return target.String(), nil
}
