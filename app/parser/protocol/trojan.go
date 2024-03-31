package protocol

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/microctar/licorice/app/utils"
)

type GrpcOptions struct {
	GrpcServiceName string `yaml:"grpc-service-name,omitempty"`
}

type WSOptions struct {
	Path                string            `yaml:"path,omitempty"`
	Headers             map[string]string `yaml:"headers,omitempty"`
	MaxEarlyData        int               `yaml:"max-early-data,omitempty"`
	EarlyDataHeaderName string            `yaml:"early-data-header-name,omitempty"`
}

type ProxyTrojan struct {
	Name           string      `yaml:"name"`
	Server         string      `yaml:"server"`
	Port           uint16      `yaml:"port"`
	Type           ProxyType   `yaml:"type"`
	Password       string      `yaml:"password"`
	ALPN           []string    `yaml:"alpn,omitempty"`
	SNI            string      `yaml:"sni,omitempty"`
	SkipCertVerify bool        `yaml:"skip-cert-verify,omitempty"`
	UDP            bool        `yaml:"udp,omitempty"`
	Network        string      `yaml:"network,omitempty"`
	GrpcOpts       GrpcOptions `yaml:"grpc-opts,omitempty"`
	WSOpts         WSOptions   `yaml:"ws-opts,omitempty"`
}

func (proxy *ProxyTrojan) Parse(uriScheme string, reQueryer utils.REQueryer) error {

	trojanURL, urlerr := url.Parse(uriScheme)

	if urlerr != nil {
		return urlerr
	}

	proxy.Name = strings.TrimSpace(trojanURL.Fragment)
	proxy.Server = trojanURL.Hostname()
	proxy.Type = Trojan
	proxy.UDP = true

	// extract and verify port

	strPort := trojanURL.Port()

	{
		port, convErr := strconv.ParseUint(strPort, 10, 16)
		if convErr != nil {
			return errors.New("parser => cannot parse port of trojan config")
		}

		if port == 0 {
			return errors.New("parser(trojan) => invalid port")
		}

		proxy.Port = uint16(port)
	}

	// extract password
	{
		passwd := utils.ReGetFirst(reQueryer.Query("trojan:\\/\\/(.*)@"), uriScheme)

		proxy.Password = passwd
	}

	proxy.SkipCertVerify = (trojanURL.Query().Get("allowInsecure") == "1")
	proxy.SNI = trojanURL.Query().Get("sni")

	return nil
}

func (proxy *ProxyTrojan) GetName() string {
	return proxy.Name
}
