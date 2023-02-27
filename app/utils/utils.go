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
	if cacheStore == (*cache.Cache)(nil) {
		return &regExp{}
	}

	return &cachedRE{
		cacheStore: cacheStore,
	}
}

func NewRegexpQueryer() REQueryer {
	return NewCachedRegexpQueryer((*cache.Cache)(nil))
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

// Regexp utilities

func ReGetOne(re *regexp.Regexp, text string) string {

	if result := re.FindStringSubmatch(text); result != nil {

		// e.g. result => ["matched string", "substring"]

		return result[1]
	}

	return ""
}

// convert string to boolean

func Str2Bool(valBool string) bool {
	return valBool == "true"
}
