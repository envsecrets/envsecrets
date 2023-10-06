package secrets

import (
	"github.com/envsecrets/envsecrets/cli/internal/dotenv"
	"github.com/envsecrets/envsecrets/dto"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/payload"
)

type Service interface {
	Init(context.ServiceContext, *clients.GQLClient, *RemoteConfig) (*dto.Secret, error)
	Set(context.ServiceContext, *clients.GQLClient, *dto.Secret) error
	Get(context.ServiceContext, *clients.GQLClient, *GetOptions) (*dto.Secret, error)
	List(context.ServiceContext, *clients.GQLClient, *ListOptions) (*dto.KPMap, error)
	Delete(context.ServiceContext, *clients.GQLClient, *DeleteOptions) (*dto.Secret, error)
}

type DefaultService struct{}

func (d *DefaultService) Init(ctx context.ServiceContext, client *clients.GQLClient, options *RemoteConfig) (*dto.Secret, error) {

	if options != nil {

		//	Get the environments service.
		service := environments.GetService()

		//	Fetch the ID of the environment first.
		environment, err := service.GetByNameAndProjectID(ctx, client, options.EnvironmentName, options.ProjectID)
		if err != nil {
			return nil, err
		}

		return &dto.Secret{
			EnvID: environment.ID,
		}, nil
	}

	return &dto.Secret{}, nil
}

// --- Flow ---
//
// If the secret has a remote environment ID saved, then set the values in the remote environment.
// Else, update the values in the local environment.
func (*DefaultService) Set(ctx context.ServiceContext, client *clients.GQLClient, secret *dto.Secret) error {

	if secret.EnvID != "" {

		data := make(map[string]*payload.Payload)
		mapping := secret.Data.ToKVMap().GetMapping()
		for key, value := range mapping {
			data[key] = &payload.Payload{
				Value: value,
			}
		}

		result, err := secrets.Set(ctx, client, &secretCommons.SetOptions{
			EnvID: secret.EnvID,
			Data:  data,
		})
		if err != nil {
			return err
		}

		secret.Version = result.Version
		return nil
	}

	//	Write to local environment file.
	return dotenv.Save(secret.Data)
}

// --- Flow ---
//
// If the secret has a remote environment ID saved, then get the values from the remote environment.
// Else, get the values from the local environment.
func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, options *GetOptions) (*dto.Secret, error) {

	if options.EnvID != "" {

		secret, err := secrets.Get(ctx, client, &secretCommons.GetOptions{
			EnvID:   options.EnvID,
			Key:     options.Key,
			Version: options.Version,
		})
		if err != nil {
			return nil, err
		}

		var mapping dto.KPMap
		for key, value := range secret.Data {
			mapping.Set(key, &dto.Payload{
				Value:     value.Value,
				Exposable: value.Exposable,
			})
		}

		//	temporarily mark all values in the mapping as encoded
		mapping.MarkAllEncoded()

		return &dto.Secret{
			EnvID:   options.EnvID,
			Version: secret.Version,
			Data:    &mapping,
		}, nil
	}

	//	Read from local environment file.
	data, err := dotenv.Load()
	if err != nil {
		return nil, err
	}

	return &dto.Secret{
		Data: data,
	}, nil
}

// --- Flow ---
//
// If the secret has a remote environment ID saved, then get the array of keys from the remote environment.
// Else, get the array of keys from the local environment.
func (*DefaultService) List(ctx context.ServiceContext, client *clients.GQLClient, options *ListOptions) (*dto.KPMap, error) {

	if options.EnvID != "" {

		secret, err := secrets.List(ctx, client, &secretCommons.ListRequestOptions{
			EnvID:   options.EnvID,
			Version: options.Version,
		})
		if err != nil {
			return nil, err
		}

		var mapping dto.KPMap
		for key, value := range secret.Data {
			mapping.Set(key, &dto.Payload{
				Value:     value.Value,
				Exposable: value.Exposable,
			})
		}

		return &mapping, nil
	}

	//	Read from local environment file.
	data, err := dotenv.Load()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// --- Flow ---
//
// If the secret has a remote environment ID saved, then delete the values from the remote environment.
// Else, delete the values from the local environment.
func (*DefaultService) Delete(ctx context.ServiceContext, client *clients.GQLClient, options *DeleteOptions) (*dto.Secret, error) {

	if options.EnvID != "" {

		secret, err := secrets.Delete(ctx, client, &secretCommons.DeleteSecretOptions{
			EnvID:   options.EnvID,
			Key:     options.Key,
			Version: options.Version,
		})
		if err != nil {
			return nil, err
		}

		return &dto.Secret{
			EnvID:   options.EnvID,
			Version: secret.Version,
		}, nil
	}

	//	Read from local environment file.
	data, err := dotenv.Load()
	if err != nil {
		return nil, err
	}

	//	Delete the key from the local environment file.
	data.Delete(options.Key)

	//	Save the updated local environment file.
	if err := dotenv.Write(data); err != nil {
		return nil, err
	}

	return &dto.Secret{}, nil
}
