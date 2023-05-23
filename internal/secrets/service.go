package secrets

import (
	internalErrors "errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
)

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) (*commons.Row, *errors.Error) {
	return graphql.Set(ctx, client, options)
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) *errors.Error {

	//	Directly delete the key=value in Hasura.
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

	if payload.Value != "" {
		secrets := make(commons.Secrets)
		secrets.Set(options.Key, payload)
		return &commons.GetResponse{
			Secrets: secrets,
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
	for key, item := range data.Secrets {
		data.Secrets[key] = commons.Payload{
			Type: item.Type,
		}
	}
	return &data, nil
}

// Pulls all secret key=value pairs from the source environment,
// and overwrites them in the target environment.
// It creates a new secret version.
func Merge(ctx context.ServiceContext, client *clients.GQLClient, options *commons.MergeSecretOptions) (*commons.Row, *errors.Error) {

	//	Fetch all key=value pairs of the source environment.
	sourceVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		EnvID:   options.SourceEnvID,
		Version: options.SourceVersion,
	})
	if err != nil {
		return nil, err
	}

	//	Fetch all key=value pairs of the target environment.
	targetVariables, err := GetAll(ctx, client, &commons.GetSecretOptions{
		EnvID: options.TargetEnvID,
	})
	if err != nil {
		return nil, err
	}

	//	If the target variables is nil,
	//	then no pairs were fetched.
	if targetVariables.Secrets == nil {
		targetVariables.Secrets = make(commons.Secrets)
	}

	//	Iterate through the target pairs,
	//	and overwrite the matching ones from the source pairs.
	targetVariables.Secrets.Overwrite(sourceVariables.Secrets)

	//	Set the updated pairs in Hasura.
	return graphql.Set(ctx, client, &commons.SetSecretOptions{
		EnvID:   options.TargetEnvID,
		Secrets: targetVariables.Secrets,
	})
}

func Decrypt(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DecryptSecretOptions) (*commons.Secrets, *errors.Error) {

	//	Get the server's copy of org-key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, options.OrgID)
	if err != nil {
		return nil, err
	}
	copy(orgKey[:], orgKeyBytes)

	//	Decrypt the value of every secret.
	if err := options.Secrets.Decrypt(orgKey); err != nil {
		return nil, errors.New(err, "Failed to decrypt secrets", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	return &options.Secrets, nil
}
