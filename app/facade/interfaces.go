package facade

import "github.com/patrickmn/go-cache"

type Generator interface {
	// function to gather intel
	Collect(encSubscription string, basedir string, ruleFilename string) error

	// export data to bytes
	Export() (data []byte, err error)

	// set up generator with cache
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
