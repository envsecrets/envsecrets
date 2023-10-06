package secrets

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/events"
	eventCommons "github.com/envsecrets/envsecrets/internal/events/commons"
	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/graphql"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
)

// Returns a new initialized 'Secret' object.
func New() *commons.Secret {
	version := 1
	return &commons.Secret{
		Version: &version,
		Data:    keypayload.KPMap{},
	}
}

func ParseAndInitialize(data []byte) (*commons.Secret, error) {

	var result commons.Secret
	if data == nil {
		return nil, fmt.Errorf("invalid inputs")
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	result.MarkEncoded()
	return &result, nil
}

// Cleanup entries from `secrets` table.
func Cleanup(ctx context.ServiceContext, client *clients.GQLClient, options *commons.CleanupSecretOptions) error {
	return graphql.Delete(ctx, client, &graphql.DeleteOptions{
		EnvID:   options.EnvID,
		Version: options.Version,
	})
}

func Get(ctx context.ServiceContext, client *clients.GQLClient, options *commons.GetOptions) (*commons.Secret, error) {
	return graphql.Get(ctx, client, &graphql.GetOptions{
		EnvID:   options.EnvID,
		Key:     options.Key,
		Version: options.Version,
	})
}

// Fetches only the keys of a secret row.
func List(ctx context.ServiceContext, client *clients.GQLClient, options *commons.ListRequestOptions) (*commons.Secret, error) {

	result, err := graphql.Get(ctx, client, &graphql.GetOptions{
		EnvID:   options.EnvID,
		Version: options.Version,
	})
	if err != nil {
		return nil, err
	}

	//	Remove the values from payload.
	result.DeleteValues()
	return result, nil
}

func Set(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SetOptions) (*commons.Secret, error) {

	//	Fetch the secret of latest version.
	secret, err := Get(ctx, client, &commons.GetOptions{
		EnvID: options.EnvID,
	})
	if err != nil {
		if strings.Compare(err.Error(), string(clients.ErrorTypeRecordNotFound)) == 0 {
			secret = New()
		} else {
			return nil, err
		}
	} else {

		//	We need to create an incremented version.
		secret.IncrementVersion()
	}

	//	Set or overwrite values in the secret.
	secret.Overwrite(options.Data)

	return graphql.Set(ctx, client, &graphql.SetOptions{
		EnvID:   options.EnvID,
		Data:    secret.Data,
		Version: secret.Version,
	})
}

func Sync(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SyncOptions) error {

	var allEvents *eventCommons.Events
	var err error

	// Fetch the integration event for the environment
	if options.IntegrationType != "" {

		allEvents, err = events.GetByEnvironmentAndIntegrationType(ctx, client, options.EnvID, integrations.Type(options.IntegrationType))
		if err != nil {
			return err
		}

	} else {

		allEvents, err = events.GetByEnvironment(ctx, client, options.EnvID)
		if err != nil {
			return err
		}

	}

	//	Get the integration service
	integrationService := integrations.GetService()

	for _, event := range *allEvents {
		if err := integrationService.Sync(ctx, client, &integrations.SyncOptions{
			IntegrationID: event.Integration.ID,
			EventID:       event.ID,
			EntityDetails: event.EntityDetails,
			Data:          options.Data,
		}); err != nil {
			return err
		}
	}

	return nil
}

func Delete(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DeleteSecretOptions) (*commons.Secret, error) {

	//	Fetch the secret with ALL the latest values.
	existing, err := Get(ctx, client, &commons.GetOptions{
		EnvID:   options.EnvID,
		Version: options.Version,
	})
	if err != nil {
		return nil, err
	}

	//	Create a shallow copy of the existing secret.
	new := existing

	//	Delete our key=value pair.
	new.Delete(options.Key)

	//	We need to create an incremented version.
	new.IncrementVersion()

	return graphql.Set(ctx, client, &graphql.SetOptions{
		EnvID:   options.EnvID,
		Data:    new.Data,
		Version: new.Version,
	})
}

// Pulls all secret key=value pairs from the source environment,
// and overwrites them in the target environment.
// It creates a new secret version.
func Merge(ctx context.ServiceContext, client *clients.GQLClient, options *commons.MergeOptions) (*commons.Secret, error) {

	//	Fetch all key=value pairs of the target environment.
	target, err := Get(ctx, client, &commons.GetOptions{
		EnvID: options.TargetEnvID,
	})
	if err != nil {
		if strings.Compare(err.Error(), string(clients.ErrorTypeRecordNotFound)) == 0 {
			target = New()
		} else {
			return nil, err
		}
	}

	//	Fetch all key=value pairs of the source environment.
	source, err := Get(ctx, client, &commons.GetOptions{
		EnvID:   options.SourceEnvID,
		Version: options.SourceVersion,
	})
	if err != nil {
		return nil, err
	}

	//	Iterate through the target pairs,
	//	and overwrite the matching ones from the source pairs.
	target.Overwrite(source.Data)

	//	We need to create an incremented version.
	target.IncrementVersion()

	//	Set the updated pairs in Hasura.
	return graphql.Set(ctx, client, &graphql.SetOptions{
		EnvID:   options.TargetEnvID,
		Data:    target.Data,
		Version: target.Version,
	})
}

func Decrypt(ctx context.ServiceContext, client *clients.GQLClient, options *commons.DecryptOptions) (*commons.Secret, error) {

	//	Get the server's copy of org-key.
	var k [32]byte
	copy(k[:], options.Key)

	//	Decrypt the value of every secret.
	if err := options.Secret.Decrypt(k); err != nil {
		return nil, err
	}

	return options.Secret, nil
}
