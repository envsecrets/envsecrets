package secrets

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
)

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetSecretOptions) (*commons.Secret, error) {
	return graphql.Set(ctx, client, options)
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) error {

	//	Directly delete the key=value in Hasura.
	if err := graphql.Delete(ctx, client, options); err != nil {
		return err
	}

	return nil
}

// Cleanup entries from `secrets` table.
func Cleanup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CleanupSecretOptions) error {
	return graphql.Cleanup(ctx, client, options)
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Secret, error) {

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {
		return graphql.GetByKeyByVersion(ctx, client, options)
	}

	return graphql.GetByKey(ctx, client, options)
}

func GetAll(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetSecretOptions) (*commons.Secret, error) {

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {
		return graphql.GetByVersion(ctx, client, options)
	}

	return graphql.Get(ctx, client, options)
}

// Fetches only the keys of a secret row.
func List(ctx context.ServiceContext, client *clients.GQLClient, options *commons.ListRequestOptions) (result *commons.Secret, err error) {

	//	If the request has a specific version specified,
	//	make the call for only that version
	if options.Version != nil {

		result, err = graphql.GetByVersion(ctx, client, &commons.GetSecretOptions{
			EnvID:   options.EnvID,
			Version: options.Version,
		})
		if err != nil {
			return nil, err
		}

	} else {

		result, err = graphql.Get(ctx, client, &commons.GetSecretOptions{
			EnvID: options.EnvID,
		})
		if err != nil {
			return nil, err
		}

	}

	//	Remove the values from payload.
	//	Only keep the type.
	result.Empty()
	return
}

// Pulls all secret key=value pairs from the source environment,
// and overwrites them in the target environment.
// It creates a new secret version.
func Merge(ctx context.ServiceContext, client *clients.GQLClient, options *commons.MergeSecretOptions) (*commons.Secret, error) {

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

	//	Iterate through the target pairs,
	//	and overwrite the matching ones from the source pairs.
	targetVariables.Overwrite(sourceVariables.Data)

	//	Set the updated pairs in Hasura.
	return targetVariables, nil
}

func Decrypt(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DecryptSecretOptions) (*commons.Secret, error) {

	//	Get the server's copy of org-key.
	var orgKey [32]byte
	orgKeyBytes, err := keys.GetOrgKeyServerCopy(ctx, options.OrgID)
	if err != nil {
		return nil, err
	}
	copy(orgKey[:], orgKeyBytes)

	//	Decrypt the value of every secret.
	if err := options.Secret.Decrypt(orgKey); err != nil {
		return nil, err
	}

	return &options.Secret, nil
}
