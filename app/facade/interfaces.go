package facade

import "github.com/patrickmn/go-cache"

type Generator interface {
	GetDefaultConfig() any
	Collect(enc_subscribtion string, basedir string, rule_filename string) error
	Merge(name string, data any)
	Export() (data []byte, err error)
	Setup(client string, cache *cache.Cache)
}

func NewCachedGenerator(client string, cache *cache.Cache) Generator {
	var customGenerator Generator

	switch client {
	case "clash":
		customGenerator = &ClashConfig{}
	}

	customGenerator.Setup(client, cache)

	return customGenerator
}

func NewGenerator(client string) Generator {
	return NewCachedGenerator(client, (*cache.Cache)(nil))
}
