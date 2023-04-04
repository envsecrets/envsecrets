package css

import (
	"fmt"
	"strings"
)

type CSSStyleRule struct {
	Selector *CSSValue
	Styles   []*CSSStyleDeclaration
}

func (sr *CSSStyleRule) Text() string {
	decls := make([]string, 0, len(sr.Styles))

	for _, s := range sr.Styles {
		decls = append(decls, s.Text())
	}

	return fmt.Sprintf("%s {\n%s\n}", sr.Selector.Text(), strings.Join(decls, ";\n"))
}
