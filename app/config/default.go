package config

import (
	"os"
	"runtime"
)

var defaultConfigDirectory []string

const (
	systemwideDirectory = "/usr/local/etc/licorice"

	DefaultClashRulePath = "rules/ACL4SSR/Clash/config"
	DefaultClashRuleFile = "ACL4SSR.ini"
)

func init() {
	if runtime.GOOS == "freebsd" || runtime.GOOS == "linux" {

		if userConfdir, cdErr := os.UserConfigDir(); cdErr == nil {
			defaultConfigDirectory = append(defaultConfigDirectory, userConfdir+"/licorice")
		}

		if userHomedir, hdErr := os.UserHomeDir(); hdErr == nil {
			defaultConfigDirectory = append(defaultConfigDirectory, userHomedir+"/.licorice")
		}

	}
}

func GetDefaultConfigDirectory() string {

	for _, directory := range defaultConfigDirectory {
		_, status := os.Stat(directory)
		if status == nil {
			return directory
		}
	}

	return systemwideDirectory
}
