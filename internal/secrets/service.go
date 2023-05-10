package secrets

import (
	"encoding/base64"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
)

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) (*commons.Secret, *errors.Error) {
	return graphql.Set(ctx, client, options)
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) *errors.Error {

	//	Directly delete the key-value in Hasura.
	if err := graphql.Delete(ctx, client, options); err != nil {
		return err
	}

	return nil
}

// Cleanup entries from `secrets` table.
func Cleanup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CleanupSecretOptions) *errors.Error {
	return graphql.Cleanup(ctx, client, options)
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.GetResponse, *errors.Error) {

	errMessage := "Failed to fetch value"

	//	Inittialize our payload
	var payload commons.Payload

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		resp, err := graphql.GetByKeyByVersion(ctx, client, options)
		if err != nil {
			return nil, err
		}

		payload = resp.Data[options.Key]

	} else {

		resp, err := graphql.GetByKey(ctx, client, options)
		if err != nil {
			return nil, err
		}

		payload = resp.Data[options.Key]
		options.Version = &resp.Version
	}

	if payload.Value != nil {
		return &commons.GetResponse{
			Data: map[string]commons.Payload{
				options.Key: payload,
			},
			Version: options.Version,
		}, nil
	}

	return nil, errors.New(internalErrors.New(errMessage), errMessage, errors.ErrorTypeBadResponse, errors.ErrorSourceGraphQL)
}

func GetAll(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.GetResponse, *errors.Error) {

	var data commons.GetResponse

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		resp, err := graphql.GetByVersion(ctx, client, options)
		if err != nil {
			return nil, err
		}

		data = *resp

	} else {

		resp, err := graphql.Get(ctx, client, options)
		if err != nil {
			return nil, err
		}

		data = *resp
	}

	return &data, nil
}

// Fetches only the keys of a secret row.
func List(ctx context.ServiceContext, client *clients.GQLClient, options *commons.ListRequestOptions) (*commons.GetResponse, *errors.Error) {

	var data commons.GetResponse

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		resp, err := graphql.GetByVersion(ctx, client, &commons.GetSecretOptions{
			EnvID:   options.EnvID,
			Version: options.Version,
		})
		if err != nil {
			return nil, err
		}

		data = *resp

	} else {

		resp, err := graphql.Get(ctx, client, &commons.GetSecretOptions{
			EnvID: options.EnvID,
		})
		if err != nil {
			return nil, err
		}

		data = *resp
	}

	//	Remove the values from payload.
	//	Only keep the type.
	for key, item := range data.Data {
		data.Data[key] = commons.Payload{
			Type: item.Type,
		}
	}
	return &data, nil
}

// Pulls all secret key-value pairs from the source environment,
// and overwrites them in the target environment.
// It creates a new secret version.
func Merge(ctx context.ServiceContext, client *clients.GQLClient, options *commons.MergeSecretOptions) (*commons.Secret, *errors.Error) {

	//	Fetch all key-value pairs of the source environment.
	sourceVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		EnvID:   options.SourceEnvID,
		Version: options.SourceVersion,
	})
	if err != nil {
		return nil, err
	}

	//	Fetch all key-value pairs of the target environment.
	targetVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		EnvID: options.TargetEnvID,
	})
	if err != nil {
		return nil, err
	}

	//	If the target variables is nil,
	//	then no pairs were fetched.
	if targetVariables.Data == nil {
		targetVariables.Data = make(map[string]commons.Payload)
	}

	//	Iterate through the target pairs,
	//	and overwrite the matching ones from the source pairs.
	for key, payload := range sourceVariables.Data {
		targetVariables.Data[key] = payload
	}

	//	Set the updated pairs in Hasura.
	return graphql.Set(ctx, client, &commons.SetSecretOptions{
		EnvID: options.TargetEnvID,
		Data:  targetVariables.Data,
	})
}

func Decrypt(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DecryptSecretOptions) (map[string]commons.Payload, *errors.Error) {

	//	Get the server's copy of org-key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, options.OrgID)
	if err != nil {
		return nil, err
	}
	copy(orgKey[:], orgKeyBytes)

	//	Decrypt the value of every secret.
	for key, payload := range options.Data {

		//	Base64 decode the secret value
		decoded, er := base64.StdEncoding.DecodeString(payload.Value.(string))
		if er != nil {
			return nil, errors.New(er, "Failed to base64 decode value for secret "+key, errors.ErrorTypeBase64Decode, errors.ErrorSourceGo)
		}

		//	If the secret is of type `ciphertext`,
		//	we will need to decode it first.
		if payload.Type == commons.Ciphertext {

			//	Decrypt the value using org-key.
			decrypted, err := keys.OpenSymmetrically(decoded, orgKey)
			if err != nil {
				return nil, err
			}

			payload.Value = string(decrypted)
		} else {
			payload.Value = string(decoded)
		}

		options.Data[key] = payload
	}
	return options.Data, nil
}
