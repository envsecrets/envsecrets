// Package css is for parsing css stylesheet.
//
//	import (
//		"github.com/vanng822/css"
//		"fmt"
//	)
//	func main() {
//		csstext = "td {width: 100px; height: 100px;}"
//		ss := css.Parse(csstext)
//		rules := ss.GetCSSRuleList()
//		for _, rule := range rules {
//			fmt.Println(rule.Style.Selector.Text())
//			fmt.Println(rule.Style.Styles)
//		}
//	}
package css
