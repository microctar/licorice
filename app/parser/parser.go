package parser

import (
	"encoding/base64"
	"strings"

	"github.com/microctar/licorice/app/parser/protocol"
	"github.com/microctar/licorice/app/utils"
)

type Parser struct {
	Proxies   []Proxy
	Groups    []string
	reQueryer utils.REQueryer
}

func (target *Parser) Parse(encSubscription string) error {

	metadata, b64err := base64.StdEncoding.DecodeString(encSubscription)

	if b64err != nil {
		return b64err
	}

	//  metadata => []byte with character '\n' as line ending

	subscription := strings.Split(string(metadata), "\n")

	for _, fragment := range subscription {
		proto := utils.ReGetFirst(target.reQueryer.Query("(.*?):\\/\\/"), fragment)

		var proxy Proxy

		switch proto {

		// if the protocol is shadowsocks
		case "ss":
			proxy = &protocol.ProxyShadowsocks{}
		case "trojan":
			proxy = &protocol.ProxyTrojan{}
		default:
			// skip
			continue
		}

		if err := proxy.Parse(fragment, target.reQueryer); err != nil {
			return err
		}

		target.Proxies = append(target.Proxies, proxy)
		target.Groups = append(target.Groups, proxy.GetName())
	}

	return nil
}
