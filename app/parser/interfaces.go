package parser

type Proxy interface {
	Parse(string) error
	GetName() string
}

func NewParser() *Parser {
	return &Parser{}
}
