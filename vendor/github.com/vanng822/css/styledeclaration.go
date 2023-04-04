package css

import (
	"fmt"
)

type CSSStyleDeclaration struct {
	Property  string
	Value     *CSSValue
	Important bool
}

func NewCSSStyleDeclaration(property, value string, important bool) *CSSStyleDeclaration {
	return &CSSStyleDeclaration{
		Property:  property,
		Value:     NewCSSValue(value),
		Important: important,
	}
}

func (decl *CSSStyleDeclaration) Text() string {
	res := fmt.Sprintf("%s: %s", decl.Property, decl.Value.Text())
	if decl.Important {
		res += " !important"
	}
	return res
}
