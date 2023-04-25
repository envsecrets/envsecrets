package auth

import "github.com/golang-jwt/jwt/v4"

type HasuraClaims struct {
	AllowedRoles []string `json:"x-hasura-allowed-roles,omitempty"`
	DefaultRole  string   `json:"x-hasura-default-role,omitempty"`
	UserID       string   `json:"x-hasura-user-id,omitempty"`
	UserEmail    string   `json:"x-hasura-user-email,omitempty"`
	IsAnonymous  string   `json:"x-hasura-user-is-anonymous,omitempty"`
}

type Claims struct {
	Hasura HasuraClaims `json:"https://hasura.io/jwt/claims"`
	jwt.RegisteredClaims
}
