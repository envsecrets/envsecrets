package css

import (
	"github.com/gorilla/css/scanner"
)

type blockParserContext struct {
	State        State
	NowProperty  string
	NowValue     []*scanner.Token
	NowImportant bool
}

func (context *blockParserContext) extractDeclaration() *CSSStyleDeclaration {
	decl := CSSStyleDeclaration{
		Property:  context.NowProperty,
		Value:     &CSSValue{Tokens: context.NowValue},
		Important: context.NowImportant,
	}
	context.NowProperty = ""
	context.NowValue = make([]*scanner.Token, 0)
	context.NowImportant = false
	return &decl
}

// ParseBlock take a string of a css block,
// parses it and returns a map of css style declarations.
func ParseBlock(csstext string) []*CSSStyleDeclaration {
	s := scanner.New(csstext)
	return parseBlock(s)
}

func parseBlock(s *scanner.Scanner) []*CSSStyleDeclaration {
	/* block       : '{' S* [ any | block | ATKEYWORD S* | ';' S* ]* '}' S*;
	property    : IDENT;
	value       : [ any | block | ATKEYWORD S* ]+;
	any         : [ IDENT | NUMBER | PERCENTAGE | DIMENSION | STRING
	              | DELIM | URI | HASH | UNICODE-RANGE | INCLUDES
	              | DASHMATCH | ':' | FUNCTION S* [any|unused]* ')'
	              | '(' S* [any|unused]* ')' | '[' S* [any|unused]* ']'
	              ] S*;
	*/
	decls := make([]*CSSStyleDeclaration, 0)

	context := &blockParserContext{
		State:        STATE_NONE,
		NowProperty:  "",
		NowValue:     make([]*scanner.Token, 0),
		NowImportant: false,
	}

	for {
		token := s.Next()

		//fmt.Printf("BLOCK(%d): %s:'%s'\n", context.State, token.Type.String(), token.Value)

		if token.Type == scanner.TokenError {
			break
		}

		if token.Type == scanner.TokenEOF {
			if context.State == STATE_VALUE {
				// we are ending without ; or }
				// this can happen when we parse only css declaration
				decl := context.extractDeclaration()
				decls = append(decls, decl)
			}
			break
		}

		switch token.Type {

		case scanner.TokenS:
			if context.State == STATE_VALUE {
				context.NowValue = append(context.NowValue, token)
			}
		case scanner.TokenIdent:
			if context.State == STATE_NONE {
				context.State = STATE_PROPERTY
				context.NowProperty += token.Value
				break
			}
			if token.Value == "important" {
				context.NowImportant = true
			} else {
				context.NowValue = append(context.NowValue, token)
			}
		case scanner.TokenChar:
			if context.State == STATE_NONE {
				if token.Value == "{" {
					break
				}
			}
			if context.State == STATE_PROPERTY {
				if token.Value == ":" {
					context.State = STATE_VALUE
				}
				// CHAR and STATE_PROPERTY but not : then weird
				// break to ignore it
				break
			}
			// should be no state or value
			if token.Value == ";" {
				decl := context.extractDeclaration()
				decls = append(decls, decl)
				context.State = STATE_NONE
			} else if token.Value == "}" { // last property in a block can have optional ;
				if context.State == STATE_VALUE {
					// only valid if state is still VALUE, could be ;}
					decl := context.extractDeclaration()
					decls = append(decls, decl)
				}
				// we are done
				return decls
			} else if token.Value != "!" {
				context.NowValue = append(context.NowValue, token)
			}
			break

		// any
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
		case scanner.TokenSubstringMatch:
			context.NowValue = append(context.NowValue, token)
		}
	}

	return decls
}
