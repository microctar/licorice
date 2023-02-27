package parser

import "github.com/microctar/licorice/app/utils"

type Proxy interface {
	Parse(string, utils.REQueryer) error
	GetName() string
}

func NewParser(queryer utils.REQueryer) *Parser {
	return &Parser{reQueryer: queryer}
}
