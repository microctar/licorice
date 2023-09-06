package utils

import (
	"regexp"

	"github.com/patrickmn/go-cache"
)

type REQueryer interface {
	Query(pattern string) *regexp.Regexp
}

type (
	cachedRE struct {
		cacheStore *cache.Cache
	}

	regExp struct{}
)

func NewCachedRegexpQueryer(cacheStore *cache.Cache) REQueryer {
	return &cachedRE{
		cacheStore: cacheStore,
	}
}

func NewRegexpQueryer() REQueryer {
	return &regExp{}
}

func (cre *cachedRE) set(pattern string) *regexp.Regexp {
	compiledRE := regexp.MustCompile(pattern)
	cre.cacheStore.Set(pattern, compiledRE, cache.NoExpiration)

	return compiledRE
}

func (cre *cachedRE) Query(pattern string) *regexp.Regexp {
	if compiledRE, found := cre.cacheStore.Get(pattern); found {
		return compiledRE.(*regexp.Regexp)
	}

	return cre.set(pattern)
}

func (re *regExp) Query(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
