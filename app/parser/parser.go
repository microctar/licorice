package parser

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/microctar/licorice/app/parser/protocol"
	"github.com/microctar/licorice/app/utils"
)

var _ Proxy = (*Parser)(nil)

type Parser struct {
	Proxies []Proxy
	Groups  []string
}

func (target *Parser) Parse(encSubscription string) error {

	metadata, b64err := base64.StdEncoding.DecodeString(encSubscription)

	if b64err != nil {
		return b64err
	}

	//  metadata => []byte with character '\n' as line ending

	subscription := strings.Split(string(metadata), "\n")

	for _, fragment := range subscription {
		proto := utils.Get("(.*?):\\/\\/", fragment)

		var proxy Proxy

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
		}

		target.Proxies = append(target.Proxies, proxy)
		target.Groups = append(target.Groups, proxy.GetName())
	}

	return nil
}

func (target *Parser) GetName() string {
	return ""
}
