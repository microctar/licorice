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

		user_configdir, os_err_cd := os.UserConfigDir()
		user_homedir, os_err_hd := os.UserHomeDir()

		if os_err_cd == nil {
			DefaultConfigDirectory = append(DefaultConfigDirectory, user_configdir+"/licorice")
		}

		if os_err_hd == nil {
			DefaultConfigDirectory = append(DefaultConfigDirectory, user_homedir+"/.licorice")
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
