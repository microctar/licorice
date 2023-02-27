package acl

import (
	"fmt"

	"github.com/microctar/licorice/app/utils"
	"github.com/patrickmn/go-cache"
)

var _ ACLReader = (*cachedACLR)(nil)

type ACLReader interface {
	// e.g. basedir => /usr/local/etc, ruleFilename => rules/ACL4SSR/Clash/config/example.ini
	ReadFile(basedir string, ruleFilename string) error

	// Expose exposes underlying type of ACLReader
	Expose() any

	// SetQueryer setup regexp queryer with/without cache layer
	SetQueryer(queryer utils.REQueryer)
}

type cachedACLR struct {
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

	aclr.SetQueryer(utils.NewRegexpQueryer())

	return aclr
}

func NewCachedACLR(client string, cache *cache.Cache) ACLReader {
	aclr := &cachedACLR{
		aclr:  NewACLR(client),
		cache: cache,
	}

	aclr.SetQueryer(utils.NewCachedRegexpQueryer(cache))

	return aclr
}

func (c *cachedACLR) ReadFile(basedir string, ruleFilename string) error {
	rulefile := fmt.Sprintf("%s/%s", basedir, ruleFilename)

	if data, found := c.cache.Get(rulefile); found {
		c.aclr = data.(ACLReader)
		return nil
	}

	if err := c.aclr.ReadFile(basedir, ruleFilename); err != nil {
		return err
	}

	c.cache.Set(rulefile, c.aclr, cache.DefaultExpiration)

	return nil
}

func (c *cachedACLR) SetQueryer(queryer utils.REQueryer) {
	c.aclr.SetQueryer(queryer)
}

func (c *cachedACLR) Expose() any {
	return c.aclr.Expose()
}
