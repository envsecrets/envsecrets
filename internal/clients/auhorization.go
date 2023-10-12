package clients

type Authorization struct {
	Token     string
	TokenType TokenType
}

type TokenType string

const (
	PAT    TokenType = "pat"
	Bearer TokenType = "Bearer"
)
