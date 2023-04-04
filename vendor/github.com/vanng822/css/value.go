package css

import (
	"strings"

	"github.com/gorilla/css/scanner"
)

type CSSValue struct {
	Tokens []*scanner.Token
}

func NewCSSValue(csstext string) *CSSValue {
	sc := scanner.New(csstext)
	val := CSSValue{Tokens: make([]*scanner.Token, 0)}
Loop:
	for {
		token := sc.Next()
		switch token.Type {
		case scanner.TokenError, scanner.TokenEOF:
			break Loop
		default:
			val.Tokens = append(val.Tokens, token)
		}
	}
	return &val
}

func NewCSSValueString(data string) *CSSValue {
	data = strings.ReplaceAll(data, `\`, `\\`)
	data = strings.ReplaceAll(data, `"`, `\"`)
	data = `"` + data + `"`
	token := scanner.Token{scanner.TokenString, data, 0, 0}
	return &CSSValue{Tokens: []*scanner.Token{&token}}
}

func (v *CSSValue) SplitOnToken(split *scanner.Token) []*CSSValue {
	res := make([]*CSSValue, 0)
	current := make([]*scanner.Token, 0)
	for _, tok := range v.Tokens {
		if tok.Type == split.Type && tok.Value == split.Value {
			res = append(res, &CSSValue{Tokens: current})
			current = make([]*scanner.Token, 0)
		} else {
			current = append(current, tok)
		}
	}
	res = append(res, &CSSValue{Tokens: current})
	return res
}

func (v *CSSValue) Text() string {
	var b strings.Builder
	for _, t := range v.Tokens {
		b.WriteString(t.Value)
	}
	return strings.TrimSpace(b.String())
}

func (v *CSSValue) ParsedText() string {
	var b strings.Builder
	for _, t := range v.Tokens {
		switch t.Type {
		case scanner.TokenString:
			val := t.Value[1 : len(t.Value)-1] // remove trailing / leading quotes
			val = strings.ReplaceAll(val, `\"`, `"`)
			val = strings.ReplaceAll(val, `\'`, `'`)
			val = strings.ReplaceAll(val, `\\`, `\`)
			// \A9 should be replaced by the corresponding rune
			b.WriteString(val)
		default:
			b.WriteString(t.Value)
		}
	}
	return strings.TrimSpace(b.String())
}
