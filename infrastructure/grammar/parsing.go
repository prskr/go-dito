package grammar

import "github.com/alecthomas/participle/v2"

// Parse takes a raw rule and parses it into the given target instance
// currently only ResponsePipeline and Check are supported for parsing
func Parse[T any](rule string) (*T, error) {
	parser, err := participle.Build[T](
		participle.Unquote("String"),
		participle.Unquote("RawString"),
	)
	if err != nil {
		return nil, err
	}

	return parser.ParseString("", rule)
}
