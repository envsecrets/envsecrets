package tokens

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/tokens/commons"
	"github.com/envsecrets/envsecrets/internal/tokens/graphql"
	"github.com/envsecrets/envsecrets/internal/tokens/internal"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

func Create(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CreateServiceOptions) (string, *errors.Error) {

	//	Export the key.
	key, err := keys.ExportKey(ctx, &keyCommons.KeyExportOptions{
		Type: string(keyCommons.EncryptionKey),
		Name: options.OrgID,
	})
	if err != nil {
		return "", err
	}

	version := len(key.Data.Keys)
	base64EncodedKey := key.Data.Keys[fmt.Sprint(version)].(string)

	//	Base64 decode the key
	value, er := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(base64EncodedKey)
	if err != nil {
		return "", errors.New(er, "Failed to create token", errors.ErrorTypeInvalidKey, errors.ErrorSourceVault)
	}

	//	Generate unique UUID for this token
	id := uuid.NewString()

	now := time.Now()
	exp := now.Add(options.Expiry)

	//	Create the token
	token, err := internal.Create(&commons.CreateOptions{
		ID:            id,
		Key:           value,
		EnvID:         options.EnvID,
		Expiry:        exp,
		IssuedAt:      now,
		NotBeforeTime: now,
	})
	if err != nil {
		return "", err
	}

	//	Insert record in the database
	_, err = graphql.Create(ctx, client, &commons.CreateGraphQLOptions{
		ID:     id,
		EnvID:  options.EnvID,
		Expiry: exp,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func Decrypt(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DecryptServiceOptions) (*paseto.JSONToken, *errors.Error) {

	//	Export the key.
	key, err := keys.ExportKey(ctx, &keyCommons.KeyExportOptions{
		Type: string(keyCommons.EncryptionKey),
		Name: options.OrgID,
	})
	if err != nil {
		return nil, err
	}

	version := len(key.Data.Keys)
	base64EncodedKey := key.Data.Keys[fmt.Sprint(version)].(string)

	//	Base64 decode the key
	value, er := base64.StdEncoding.DecodeString(base64EncodedKey)
	if err != nil {
		return nil, errors.New(er, "Failed to decrypt token", errors.ErrorTypeInvalidKey, errors.ErrorSourceVault)
	}

	return internal.Decrypt(&commons.DecryptOptions{
		Key:   value,
		Token: options.Token,
	})
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Token, *errors.Error) {
	return graphql.Get(ctx, client, id)
}
