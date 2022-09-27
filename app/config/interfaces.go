package config

type Proxy interface {
	Parse(string) error
	GetName() string
}

type Generator interface {
	GetDefaultConfig() any
	Collect(enc_subscribtion string, basedir string, rule_filename string) error
	Merge(name string, data any)
	Export() (data []byte, err error)
}

type ACLReader interface {
	// e.g. basedir => /usr/local/etc rule_filename => rules/ACL4SSR/Clash/config/example.ini
	ReadFile(basedir string, rule_filename string) error
}
