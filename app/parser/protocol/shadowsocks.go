package protocol

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/microctar/licorice/app/utils"
)

// SIP002 URI SCHEME
// Following [RFC3986](https://www.ietf.org/rfc/rfc3986.txt)
// SS-URI = "ss://" userinfo "@" hostname ":" port [ "/" ] [ "?" plugin ] [ "#" tag ]
// userinfo = websafe-base64-encode-utf8(method  ":" password)
//            method ":" password

type ProxyShadowsocks struct {
	Name       string                 `yaml:"name"`
	Server     string                 `yaml:"server"`
	Port       uint16                 `yaml:"port"`
	Type       ProxyType              `yaml:"type"`
	Method     string                 `yaml:"cipher"`
	Password   string                 `yaml:"password"`
	Plugin     string                 `yaml:"plugin,omitempty"`
	PluginOpts map[string]interface{} `yaml:"plugin-opts,omitempty"`
	UDP        bool                   `yaml:"udp"`
}

func (proxy *ProxyShadowsocks) Parse(uriScheme string) error {
	ssURL, urlerr := url.Parse(uriScheme)

	if urlerr != nil {
		return urlerr
	}

	proxy.Name = strings.TrimSpace(ssURL.Fragment)
	proxy.Server = ssURL.Hostname()
	proxy.Type = Shadowsocks
	proxy.UDP = true

	// extract and verify port

	strPort := ssURL.Port()

	{
		port, convErr := strconv.ParseUint(strPort, 10, 16)
		if convErr != nil {
			return errors.New("parser => cannot parser port of shadowsocks config")
		}
		proxy.Port = uint16(port)
	}

	if proxy.Port == 0 {
		return errors.New("parser => port of shadowsocks config is incorrect")
	}

	{
		// extract userinfo

		userinfo := utils.Get("ss:\\/\\/(.*?)@", uriScheme)

		mp, b64err := base64.RawURLEncoding.DecodeString(userinfo)

		if b64err != nil {
			return b64err
		}

		// info => [method (aka cipher), password]

		info := strings.Split(string(mp), ":")

		proxy.Method, proxy.Password = info[0], info[1]
	}

	// "TOR_PT_SERVER_TRANSPORT_OPTIONS" -- A semicolon-separated list
	//  of <key>:<value> pairs, where <key> is a transport name and
	//  <value> is a k=v string value with options that are to be passed
	//  to the transport. Colons, semicolons, equal signs and backslashes
	//  MUST be escaped with a backslash. TOR_PT_SERVER_TRANSPORT_OPTIONS
	//  is optional and might not be present in the environment of the
	//  proxy if no options are need to be passed to transports.

	// try to extract plugin options

	plugin := ssURL.Query().Get("plugin")

	if plugin != "" {
		proxy.PluginOpts = make(map[string]interface{})
		pluginAndOpts := strings.Split(plugin, ";")
		proxy.Plugin = pluginAndOpts[0]

		for _, opts := range pluginAndOpts[1:] {
			kvpair := strings.Split(opts, "=")
			proxy.PluginOpts[kvpair[0]] = kvpair[1]
		}
	}

	return nil
}

// GetName, the implementation of shadowsocks for Proxy interface
func (proxy *ProxyShadowsocks) GetName() string {
	return proxy.Name
}
