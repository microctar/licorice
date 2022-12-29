package facade

import (
	"errors"
	"os"

	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/parser"
	"github.com/microctar/licorice/app/utils"
	"gopkg.in/yaml.v3"
)

// Generate Default Clash Configuration
type ClashConfig struct {
	RawConfig RawConfig
}

func (clash_config *ClashConfig) GetDefaultConfig() any {

	clash_config.RawConfig = RawConfig{
		Port:               7890,
		SocksPort:          7891,
		AllowLan:           false,
		Mode:               Rule,
		LogLevel:           INFO,
		ExternalController: ":9090",
	}

	return clash_config.RawConfig
}

func (clash_config *ClashConfig) Merge(name string, data any) {

	switch name {
	case "proxies":
		clash_config.RawConfig.Proxy = data.([]config.Proxy)
	case "proxy-groups":
		clash_config.RawConfig.ProxyGroup = data.([]map[string]any)
	case "rules":
		clash_config.RawConfig.Rule = append(clash_config.RawConfig.Rule, data.([]string)...)
	}
}

func (clash_config *ClashConfig) Collect(enc_subcribtion string, basedir string, rule_filename string) error {

	clash_config.GetDefaultConfig()

	data := parser.Parser{}
	if err := data.Parse(enc_subcribtion); err != nil {
		return err
	}

	clash_config.Merge("proxies", data.Proxies)

	if _, status := os.Stat(basedir); status == nil {

		diverter := utils.Diverter{}
		read_err := diverter.ReadFile(basedir, rule_filename)

		if read_err != nil {
			return read_err
		}

		// append ruleset to RawConfig
		for _, ruleset := range diverter.Ruleset {
			clash_config.Merge("rules", ruleset)
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
		clash_config.Merge("proxy-groups", diverter.CustomProxyGroup)
	} else {
		return errors.New("cannot find acl4ssr config directory")
	}

	return nil
}

// Export Clash Config

func (clash_config ClashConfig) Export() ([]byte, error) {

	out, err := yaml.Marshal(clash_config.RawConfig)

	if err != nil {
		return nil, err
	}

	return out, nil

}
