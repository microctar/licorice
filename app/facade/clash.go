package facade

import (
	"errors"
	"os"

	"github.com/microctar/licorice/app/parser"
	"github.com/microctar/licorice/app/utils"
	"github.com/microctar/licorice/app/utils/acl"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

var _ Generator = (*ClashConfig)(nil)

// Generate Default Clash Configuration
type ClashConfig struct {
	RawConfig RawConfig
	aclr      acl.ACLReader
	reQueryer utils.REQueryer
}

func (cc *ClashConfig) getDefaultConfig() {

	cc.RawConfig = RawConfig{
		Port:               7890,
		SocksPort:          7891,
		AllowLan:           true,
		Mode:               Rule,
		LogLevel:           INFO,
		ExternalController: ":9090",
	}

}

func (cc *ClashConfig) Collect(encSubscription string, basedir string, ruleFilename string) error {

	cc.getDefaultConfig()

	data := parser.NewParser(cc.reQueryer)
	if err := data.Parse(encSubscription); err != nil {
		return err
	}

	cc.RawConfig.Proxy = data.Proxies

	if _, status := os.Stat(basedir); status == nil {

		if readErr := cc.aclr.ReadFile(basedir, ruleFilename); readErr != nil {
			return readErr
		}

		diverter := cc.aclr.Expose().(*acl.ClashDiverter)

		// append ruleset to RawConfig
		for _, ruleset := range diverter.Ruleset {
			cc.RawConfig.Rule = append(cc.RawConfig.Rule, ruleset...)
		}

		// replace ".*" with real group name
		for _, proxyGrp := range diverter.CustomProxyGroup {
			tmpdata := proxyGrp["proxies"].([]string)
			if tail := len(tmpdata) - 1; tmpdata[tail] == ".*" {
				tmpdata = tmpdata[:tail]
				tmpdata = append(tmpdata, data.Groups...)
				proxyGrp["proxies"] = tmpdata
			}
		}

		// add proxy-groups to RawConfig
		cc.RawConfig.ProxyGroup = diverter.CustomProxyGroup
	} else {
		return errors.New("cannot find acl4ssr config directory")
	}

	return nil
}

// Export Clash Config

func (cc *ClashConfig) Export() ([]byte, error) {

	out, err := yaml.Marshal(cc.RawConfig)

	if err != nil {
		return nil, err
	}

	return out, nil

}

func (cc *ClashConfig) Setup(client string, cachestore *cache.Cache) {
	if cachestore == (*cache.Cache)(nil) {
		cc.aclr = acl.NewACLR(client)
		cc.reQueryer = utils.NewRegexpQueryer()
		return
	}

	cc.aclr = acl.NewCachedACLR(client, cachestore)
	cc.reQueryer = utils.NewCachedRegexpQueryer(cachestore)
}
