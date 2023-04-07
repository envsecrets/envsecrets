package css

import (
	"github.com/gorilla/css/scanner"
)

func parseAtNoBody(s *scanner.Scanner, ruleType RuleType) *CSSRule {
	/*

		Syntax:
		@charset charset;

		Example:
		@charset "UTF-8";


		Syntax:
		@import url;                      or
		@import url list-of-media-queries;

		Example:
		@import url("fineprint.css") print;
		@import url("bluish.css") projection, tv;
		@import 'custom.css';
		@import url("chrome://communicator/skin/");
		@import "common.css" screen, projection;
		@import url('landscape.css') screen and (orientation:landscape);

	*/

	parsed := make([]*scanner.Token, 0)
Loop:
	for {
		token := s.Next()

		if token.Type == scanner.TokenEOF || token.Type == scanner.TokenError {
			return nil
		}
		// take everything for now
		switch token.Type {
		case scanner.TokenEOF, scanner.TokenError:
			break Loop
		case scanner.TokenChar:
			if token.Value == ";" {
				break Loop
			}
			fallthrough
		default:
			parsed = append(parsed, token)
		}
	}

	rule := NewRule(ruleType)
	rule.Style.Selector = &CSSValue{Tokens: parsed}
	return rule
}
