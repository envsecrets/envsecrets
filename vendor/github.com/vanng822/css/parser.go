package css

import (
	"log"

	"github.com/gorilla/css/scanner"
)

/*
	stylesheet  : [ CDO | CDC | S | statement ]*;
	statement   : ruleset | at-rule;
	at-rule     : ATKEYWORD S* any* [ block | ';' S* ];
	block       : '{' S* [ any | block | ATKEYWORD S* | ';' S* ]* '}' S*;
	ruleset     : selector? '{' S* declaration? [ ';' S* declaration? ]* '}' S*;
	selector    : any+;
	declaration : property S* ':' S* value;
	property    : IDENT;
	value       : [ any | block | ATKEYWORD S* ]+;
	any         : [ IDENT | NUMBER | PERCENTAGE | DIMENSION | STRING
	              | DELIM | URI | HASH | UNICODE-RANGE | INCLUDES
	              | DASHMATCH | ':' | FUNCTION S* [any|unused]* ')'
	              | '(' S* [any|unused]* ')' | '[' S* [any|unused]* ']'
	              ] S*;
	unused      : block | ATKEYWORD S* | ';' S* | CDO S* | CDC S*;
*/

type State int

const (
	STATE_NONE State = iota
	STATE_SELECTOR
	STATE_PROPERTY
	STATE_VALUE
)

type parserContext struct {
	State             State
	NowSelector       []*scanner.Token
	NowRuleType       RuleType
	CurrentNestedRule *CSSRule
}

func (context *parserContext) resetContextStyleRule() {
	context.NowSelector = make([]*scanner.Token, 0)
	context.NowRuleType = STYLE_RULE
	context.State = STATE_NONE
}

func parseRule(context *parserContext, s *scanner.Scanner, css *CSSStyleSheet) {
	rule := NewRule(context.NowRuleType)
	selector := append(context.NowSelector, parseSelector(s)...)
	rule.Style.Selector = &CSSValue{Tokens: selector}
	rule.Style.Styles = parseBlock(s)
	if context.CurrentNestedRule != nil {
		context.CurrentNestedRule.Rules = append(context.CurrentNestedRule.Rules, rule)
	} else {
		css.CssRuleList = append(css.CssRuleList, rule)
	}
	context.resetContextStyleRule()
}

// Parse takes a string of valid css rules, stylesheet,
// and parses it. Be aware this function has poor error handling
// so you should have valid syntax in your css
func Parse(csstext string) *CSSStyleSheet {
	context := &parserContext{
		State:             STATE_NONE,
		NowSelector:       make([]*scanner.Token, 0),
		NowRuleType:       STYLE_RULE,
		CurrentNestedRule: nil,
	}

	css := &CSSStyleSheet{}
	css.CssRuleList = make([]*CSSRule, 0)
	s := scanner.New(csstext)

	for {
		token := s.Next()

		if token.Type == scanner.TokenEOF || token.Type == scanner.TokenError {
			break
		}

		switch token.Type {
		case scanner.TokenCDO:
			break
		case scanner.TokenCDC:
			break
		case scanner.TokenComment:
			break
		case scanner.TokenS:
			break
		case scanner.TokenAtKeyword:
			switch token.Value {
			case "@media":
				context.NowRuleType = MEDIA_RULE
			case "@font-face":
				// Parse as normal rule, would be nice to parse according to syntax
				// https://developer.mozilla.org/en-US/docs/Web/CSS/@font-face
				context.NowRuleType = FONT_FACE_RULE
				parseRule(context, s, css)
			case "@import":
				// No validation
				// https://developer.mozilla.org/en-US/docs/Web/CSS/@import
				rule := parseAtNoBody(s, IMPORT_RULE)
				if rule != nil {
					css.CssRuleList = append(css.CssRuleList, rule)
				}
				context.resetContextStyleRule()
			case "@charset":
				// No validation
				// https://developer.mozilla.org/en-US/docs/Web/CSS/@charset
				rule := parseAtNoBody(s, CHARSET_RULE)
				if rule != nil {
					css.CssRuleList = append(css.CssRuleList, rule)
				}
				context.resetContextStyleRule()

			case "@page":
				context.NowRuleType = PAGE_RULE
				parseRule(context, s, css)
			case "@keyframes":
				context.NowRuleType = KEYFRAMES_RULE
			case "@-webkit-keyframes":
				context.NowRuleType = WEBKIT_KEYFRAMES_RULE
			case "@counter-style":
				context.NowRuleType = COUNTER_STYLE_RULE
				parseRule(context, s, css)
			default:
				log.Printf("Skip unsupported atrule: %s", token.Value)
				skipRules(s)
				context.resetContextStyleRule()
			}
		default:
			if context.State == STATE_NONE {
				if token.Value == "}" && context.CurrentNestedRule != nil {
					// close media/keyframe/… rule
					css.CssRuleList = append(css.CssRuleList, context.CurrentNestedRule)
					context.CurrentNestedRule = nil
					break
				}
			}

			if context.NowRuleType == MEDIA_RULE || context.NowRuleType == KEYFRAMES_RULE || context.NowRuleType == WEBKIT_KEYFRAMES_RULE {
				context.CurrentNestedRule = NewRule(context.NowRuleType)
				sel := append([]*scanner.Token{token}, parseSelector(s)...)
				context.CurrentNestedRule.Style.Selector = &CSSValue{Tokens: sel}
				context.resetContextStyleRule()
				break
			} else {
				context.NowSelector = append(context.NowSelector, token)
				parseRule(context, s, css)
				break
			}
		}
	}
	return css
}
