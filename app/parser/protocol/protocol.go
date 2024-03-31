package protocol

// types of protocol used in Virtual Private Network
type ProxyType uint16

const (
	// unknown protocol
	Unknown ProxyType = iota

	// a fast tunnel proxy [Shadowsocks](https://shadowsocks.org)
	Shadowsocks

	// An unidentifiable mechanism
	Trojan
)

func (target ProxyType) String() string {
	switch target {
	case Shadowsocks:
		return "ss"
	case Trojan:
		return "trojan"
	default:
		return "unknown"
	}
}

func (target ProxyType) MarshalYAML() (any, error) {
	return target.String(), nil
}
