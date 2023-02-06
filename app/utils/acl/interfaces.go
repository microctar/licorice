package acl

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

var _ ACLReader = (*CachedACLR)(nil)

type ACLReader interface {
	// e.g. basedir => /usr/local/etc, rule_filename => rules/ACL4SSR/Clash/config/example.ini
	ReadFile(basedir string, rule_filename string) error
	Expose() any
}

type CachedACLR struct {
	aclr  ACLReader
	cache *cache.Cache
}

func NewACLR(client string) ACLReader {
	var aclr ACLReader

	switch client {
	case "clash":
		aclr = &clashDiverter{
			Ruleset: make(map[string][]string),
		}
	}

	return aclr
}

func NewCachedACLR(client string, cache *cache.Cache) ACLReader {
	var aclr ACLReader

	switch client {
	case "clash":
		aclr = &clashDiverter{
			Ruleset: make(map[string][]string),
		}
	}

	return &CachedACLR{
		aclr:  aclr,
		cache: cache,
	}
}

func (c *CachedACLR) ReadFile(basedir string, rule_filename string) error {
	rulefile := fmt.Sprintf("%s/%s", basedir, rule_filename)

	if data, found := c.cache.Get(rulefile); found {
		c.aclr = data.(ACLReader)
		return nil
	}

	err := c.aclr.ReadFile(basedir, rule_filename)

	if err == nil {
		c.cache.Set(rulefile, c.aclr, cache.DefaultExpiration)
	}

	return err
}

func (c *CachedACLR) Expose() any {
	return c.aclr.Expose()
}
