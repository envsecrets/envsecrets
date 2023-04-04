package css

import (
	"github.com/gorilla/css/scanner"
)

func parseSelector(s *scanner.Scanner) []*scanner.Token {
	/*
		selector    : any+;
		any         : [ IDENT | NUMBER | PERCENTAGE | DIMENSION | STRING
		              | DELIM | URI | HASH | UNICODE-RANGE | INCLUDES
		              | DASHMATCH | ':' | FUNCTION S* [any|unused]* ')'
		              | '(' S* [any|unused]* ')' | '[' S* [any|unused]* ']'
		              ] S*;
	*/

	result := make([]*scanner.Token, 0)

Loop:
	for {
		token := s.Next()

		switch token.Type {
		case scanner.TokenError, scanner.TokenEOF:
			break Loop
		case scanner.TokenChar:
			if token.Value == "{" {
				break Loop
			}
			fallthrough
		case scanner.TokenIdent:
			fallthrough
		case scanner.TokenS:
			fallthrough
		case scanner.TokenNumber:
			fallthrough
		case scanner.TokenPercentage:
			fallthrough
		case scanner.TokenDimension:
			fallthrough
		case scanner.TokenString:
			fallthrough
		case scanner.TokenURI:
			fallthrough
		case scanner.TokenHash:
			fallthrough
		case scanner.TokenUnicodeRange:
			fallthrough
		case scanner.TokenIncludes:
			fallthrough
		case scanner.TokenDashMatch:
			fallthrough
		case scanner.TokenFunction:
			fallthrough
		case scanner.TokenSuffixMatch:
			fallthrough
		case scanner.TokenPrefixMatch:
			fallthrough
		case scanner.TokenSubstringMatch:
			result = append(result, token)
		}
	}

	return result
}
