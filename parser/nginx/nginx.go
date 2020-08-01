package nginx

import "github.com/eggsampler/certgot/parser"

type SimpleDirective struct {
	Name       string
	Parameters []string
}

type BlockDirective struct {
	Name             string
	Parameters       []string
	SimpleDirectives map[string]SimpleDirective
	BlockDirectives  map[string]BlockDirective
}

var (
	directiveName      = parser.OneOrMore(parser.Or(parser.Letter(), parser.Rune('_')))
	directiveParameter = parser.OneOrMore(
		parser.And(
			parser.NotRune(' '),
			parser.NotRune(';'),
			parser.NotRune('{')))

	simpleDirectiveParser = parser.And(
		parser.ZeroOrMore(parser.Space()),
		directiveName,
		parser.ZeroOrMore(
			parser.And(
				parser.OneOrMore(parser.Space()),
				directiveParameter)),
		parser.ZeroOrMore(parser.Space()),
		parser.Rune(';'))
)
