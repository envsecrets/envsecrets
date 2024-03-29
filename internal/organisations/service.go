package organisations

import (
	"encoding/base64"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	keyCommons "github.com/envsecrets/envsecrets/internal/keys/commons"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/roles"
	"github.com/envsecrets/envsecrets/utils"
	"github.com/machinebox/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Organisation, error)
	GetByEnvironment(context.ServiceContext, *clients.GQLClient, string) (*Organisation, error)
	GetByProject(context.ServiceContext, *clients.GQLClient, string) (*Organisation, error)
	GetInviteLimit(context.ServiceContext, *clients.GQLClient, string) (*int, error)
	Create(context.ServiceContext, *clients.GQLClient, *CreateOptions) (*Organisation, error)
	List(context.ServiceContext, *clients.GQLClient) (*[]Organisation, error)
	UpdateInviteLimit(context.ServiceContext, *clients.GQLClient, *UpdateInviteLimitOptions) error
}

type DefaultService struct{}

func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
			id
			name
		}
	  }	  
	`)

	req.Var("id", id)

	var response struct {
		Organisation *Organisation `json:"organisations_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return response.Organisation, nil
}

func (*DefaultService) GetByProject(ctx context.ServiceContext, client *clients.GQLClient, project_id string) (*Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations(where: {projects: {id: {_eq: $id}}}, limit: 1) {
		  id
		}
	  }		 
	`)

	req.Var("id", project_id)

	var response struct {
		Organisations []Organisation `json:"organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Organisations[0], nil
}

func (*DefaultService) GetByEnvironment(ctx context.ServiceContext, client *clients.GQLClient, env_id string) (*Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery($env_id: uuid!) {
		organisations(where: {projects: {environments: {id: {_eq: $env_id}}}}, limit: 1) {
		  id
		}
	  }	   
	`)

	req.Var("env_id", env_id)

	var response struct {
		Organisations []Organisation `json:"organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Organisations[0], nil
}

func (*DefaultService) GetInviteLimit(ctx context.ServiceContext, client *clients.GQLClient, id string) (*int, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		organisations_by_pk(id: $id) {
			invite_limit
		}
	  }	  
	`)

	req.Var("id", id)

	var response struct {
		Organisation *Organisation `json:"organisations_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return response.Organisation.InviteLimit, nil
}

func (*DefaultService) List(ctx context.ServiceContext, client *clients.GQLClient) (*[]Organisation, error) {

	req := graphql.NewRequest(`
	query MyQuery {
		organisations {
		  id
		  name
		}
	  }	   
	`)

	var response struct {
		Organisations []Organisation `json:"organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Organisations, nil
}

func (*DefaultService) Create(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (organisation *Organisation, err error) {

	if options.UserID != "" {
		organisation, err = createWithUserID(ctx, client, options)
		if err != nil {
			return nil, err
		}
	} else {
		organisation, err = create(ctx, client, options.Name)
		if err != nil {
			return nil, err
		}
	}

	//	Generate default roles for the organisation.
	if _, err := roles.Insert(ctx, client, &roles.RoleInsertOptions{
		OrgID: organisation.ID,
		Name:  "viewer",
		Permissions: roles.Permissions{
			Projects: roles.CRUD{
				Read: true,
			},
		},
	}); err != nil {
		return nil, err
	}

	if _, err := roles.Insert(ctx, client, &roles.RoleInsertOptions{
		OrgID: organisation.ID,
		Name:  "editor",
		Permissions: roles.Permissions{
			Projects: roles.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Environments: roles.CRUD{
				Create: true,
				Update: true,
				Delete: true,
			},
		},
	}); err != nil {
		return nil, err
	}

	adminRole, err := roles.Insert(ctx, client, &roles.RoleInsertOptions{
		OrgID: organisation.ID,
		Name:  "admin",
		Permissions: roles.Permissions{
			Integrations: roles.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Permissions: roles.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Projects: roles.CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
			Environments: roles.CRUD{
				Create: true,
				Update: true,
				Delete: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	//	Generate a symmetric key for cryptographic operations in this organisation.
	keyBytes, err := utils.GenerateRandomBytes(keyCommons.KEY_BYTES)
	if err != nil {
		return nil, err
	}

	//	Encrypt the key using owner's public key
	publicKeyBytes, err := keys.GetPublicKeyByUserID(ctx, client, organisation.UserID)
	if err != nil {
		return nil, err
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyBytes)
	result, err := keys.SealAsymmetricallyAnonymous(keyBytes, publicKey)
	if err != nil {
		return nil, err
	}

	if err := memberships.CreateWithUserID(ctx, client, &memberships.CreateOptions{
		UserID: organisation.UserID,
		OrgID:  organisation.ID,
		RoleID: adminRole.ID,
		Key:    base64.StdEncoding.EncodeToString(result),
	}); err != nil {
		return nil, err
	}

	return organisation, nil
}

func (*DefaultService) UpdateInviteLimit(ctx context.ServiceContext, client *clients.GQLClient, options *UpdateInviteLimitOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $limit: Int!) {
		update_organisations(where: {id: {_eq: $id}}, _inc: {invite_limit: $limit}) {
			affected_rows
		  }
	  }			
	`)

	req.Var("id", options.ID)
	req.Var("limit", options.IncrementLimitBy)

	var response struct {
		Result struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	if response.Result.AffectedRows == 0 {
		return fmt.Errorf("failed to update the invite limit")
	}

	return nil
}

//
//	--- GraphQL ---
//

// Create a new organisation
func create(ctx context.ServiceContext, client *clients.GQLClient, name string) (*Organisation, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!) {
		insert_organisations(objects: {name: $name}) {
		  returning {
			id
			name
			user_id
		  }
		}
	  }
	`)

	req.Var("name", name)

	var response struct {
		Result struct {
			Organisations []Organisation `json:"returning"`
		} `json:"insert_organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Result.Organisations[0], nil
}

// Create a new organisation
func createWithUserID(ctx context.ServiceContext, client *clients.GQLClient, options *CreateOptions) (*Organisation, error) {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $user_id: uuid!) {
		insert_organisations(objects: {name: $name, user_id: $user_id}) {
		  returning {
			id
			name
			user_id
		  }
		}
	  }
	`)

	req.Var("name", options.Name)
	req.Var("user_id", options.UserID)

	var response struct {
		Result struct {
			Organisations []Organisation `json:"returning"`
		} `json:"insert_organisations"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Result.Organisations[0], nil
}
