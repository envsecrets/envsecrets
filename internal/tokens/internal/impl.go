package internal

import (
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"

	"github.com/o1egl/paseto"
)

func Create(options *commons.CreateOptions) (string, *errors.Error) {

	errorMessage := "Failed to create token"

	jsonToken := paseto.JSONToken{
		Audience:  options.EnvID,
		Issuer:    "envsecrets",
		Jti:       options.ID,
		Subject:   "Environment Token",
		IssuedAt:  options.IssuedAt,
		NotBefore: options.NotBeforeTime,
	}

	if !options.Expiry.IsZero() {
		jsonToken.Expiration = options.Expiry
	}

	// Add claims to the token
	jsonToken.Set("env_id", options.EnvID)

	//	Prepare token footer
	footer := "envsecrets environment token"

	// Encrypt data
	token, err := paseto.NewV2().Encrypt(options.Key, jsonToken, footer)
	if err != nil {
		return "", errors.New(err, errorMessage, errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return token, nil
}

func Decrypt(options *commons.DecryptOptions) (*paseto.JSONToken, *errors.Error) {

	var token paseto.JSONToken
	if err := paseto.NewV2().Decrypt(options.Token, options.Key, &token, nil); err != nil {
		return nil, errors.New(err, "Failed to descrypt token", errors.ErrorTypeInvalidToken, errors.ErrorSourceGo)
	}

	return &token, nil
}
