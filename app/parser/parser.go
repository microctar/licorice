package parser

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/microctar/licorice/app/config"
	"github.com/microctar/licorice/app/parser/protocol"
	"github.com/microctar/licorice/app/utils"
)

type Parser struct {
	Proxies []config.Proxy
	Groups  []string
}

func (target *Parser) Parse(enc_subscription string) error {

	metadata, b64err := base64.StdEncoding.DecodeString(enc_subscription)

	if b64err != nil {
		return b64err
	}

	//  metadata => []byte with character '\n' as line ending

	subscription := strings.Split(string(metadata), "\n")

	for _, fragment := range subscription {
		proto := utils.Get("(.*?):\\/\\/", fragment)

		var proxy config.Proxy

		switch proto {

		// if the protocol is shadowsocks
		case "ss":
			proxy = &protocol.ProxyShadowsocks{}
		default:
			// skip
			continue
		}

		if err := proxy.Parse(fragment); err != nil {
			log.Fatal(err)
		} else {
			target.Proxies = append(target.Proxies, proxy)
			target.Groups = append(target.Groups, proxy.GetName())
		}

	}

	return nil
}

func (target *Parser) GetName() string {
	return ""
}
