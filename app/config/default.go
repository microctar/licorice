package config

import (
	"os"
	"runtime"
)

var (
	DefaultConfigDirectory []string
	SystemwideDirectory    string = "/usr/local/etc/licorice"
	DefaultClashConfigPath string = "rules/ACL4SSR/Clash/config"
	DefaultClashRule       string = "ACL4SSR.ini"
)

func init() {
	if runtime.GOOS == "freebsd" || runtime.GOOS == "linux" {

		if userConfdir, cdErr := os.UserConfigDir(); cdErr == nil {
			DefaultConfigDirectory = append(DefaultConfigDirectory, userConfdir+"/licorice")
		}

		if userHomedir, hdErr := os.UserHomeDir(); hdErr == nil {
			DefaultConfigDirectory = append(DefaultConfigDirectory, userHomedir+"/.licorice")
		}

	}
}

func GetDefaultConfigDirectory() string {

	for _, directory := range DefaultConfigDirectory {
		_, status := os.Stat(directory)
		if status == nil {
			return directory
		}
	}

	return SystemwideDirectory
}
