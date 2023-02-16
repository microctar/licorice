package facade

import (
	"errors"
	"os"

	"github.com/microctar/licorice/app/parser"
	"github.com/microctar/licorice/app/utils/acl"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

var _ Generator = (*ClashConfig)(nil)

// Generate Default Clash Configuration
type ClashConfig struct {
	RawConfig RawConfig
	aclr      acl.ACLReader
}

func (cc *ClashConfig) getDefaultConfig() any {

	cc.RawConfig = RawConfig{
		Port:               7890,
		SocksPort:          7891,
		AllowLan:           true,
		Mode:               Rule,
		LogLevel:           INFO,
		ExternalController: ":9090",
	}

	return cc.RawConfig
}

func (cc *ClashConfig) Collect(enc_subcribtion string, basedir string, rule_filename string) error {

	cc.getDefaultConfig()

	data := parser.NewParser()
	if err := data.Parse(enc_subcribtion); err != nil {
		return err
	}

	cc.RawConfig.Proxy = data.Proxies

	if _, status := os.Stat(basedir); status == nil {

		if read_err := cc.aclr.ReadFile(basedir, rule_filename); read_err != nil {
			return read_err
		}

		diverter := cc.aclr.Expose().(*acl.ClashDiverter)

		// append ruleset to RawConfig
		for _, ruleset := range diverter.Ruleset {
			cc.RawConfig.Rule = append(cc.RawConfig.Rule, ruleset...)
		}

		// replace ".*" with real group name
		for _, proxies_group := range diverter.CustomProxyGroup {
			tmpdata := proxies_group["proxies"].([]string)
			if tail := len(tmpdata) - 1; tmpdata[tail] == ".*" {
				tmpdata = tmpdata[:tail]
				tmpdata = append(tmpdata, data.Groups...)
				proxies_group["proxies"] = tmpdata
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
		return
	}

	cc.aclr = acl.NewCachedACLR(client, cachestore)
}
