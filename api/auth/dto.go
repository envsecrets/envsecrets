package auth

import "github.com/golang-jwt/jwt/v4"

type SigninOptions struct {
	Email    string `json:"email,omitempty"`
	OTP      string `json:"otp,omitempty"`
	Ticket   string `json:"ticket,omitempty"`
	Password string `json:"password,omitempty"`
}

type ToggleMFAOptions struct {
	Code string `json:"code"`
}

type UpdatePasswordOptions struct {
	NewPassword string `json:"newPassword,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
}

type SignupOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

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
