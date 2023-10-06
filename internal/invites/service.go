package invites

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/machinebox/graphql"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*Invite, error)
	Send(context.ServiceContext, *clients.GQLClient, *SendOptions) error
	Accept(context.ServiceContext, *clients.GQLClient, string) error
	Update(context.ServiceContext, *clients.GQLClient, string, *UpdateOptions) error
}

type DefaultService struct{}

func (*DefaultService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*Invite, error) {

	req := graphql.NewRequest(`
	query MyQuery($id: uuid!) {
		invites_by_pk(id: $id) {
			id
			key
			org_id
			role_id
			email
			accepted
		}
	  }	  
	`)

	req.Var("id", id)

	var response struct {
		Invite Invite `json:"invites_by_pk"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response.Invite, nil
}

func (*DefaultService) Send(ctx context.ServiceContext, client *clients.GQLClient, options *SendOptions) error {

	//	Add the admin header to graphql client to be able to fetch the invitee's public key.
	client.Headers = append(client.Headers, clients.XHasuraAdminSecretHeader)

	//	Fetch the invitee user's public key.
	public_key, err := keys.GetPublicKeyByUserEmail(ctx, client, options.InviteeEmail)
	if err != nil {
		return err
	}

	//	Encrypt the passed decrypted key with the invitee's public key.
	var key [32]byte
	copy(key[:], public_key)
	inviteeKeyCopy, err := keys.SealAsymmetricallyAnonymous(options.Key, key)
	if err != nil {
		return err
	}

	//	Insert the invite.
	req := graphql.NewRequest(`
	mutation MyMutation($objects: [invites_insert_input!]!) {
		insert_invites(objects: $objects) {
		  affected_rows
		}
	  }			
	`)

	req.Var("objects", []InsertOptions{
		{
			Key:    base64.StdEncoding.EncodeToString(inviteeKeyCopy),
			OrgID:  options.OrgID,
			RoleID: options.RoleID,
			Email:  options.InviteeEmail,
			UserID: options.InviterID,
		},
	})

	var response struct {
		Query struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"insert_invites"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	//	Validate the mutation as been written to the database
	if response.Query.AffectedRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

// ---	Flow ---
// 1. Fetch the invite row from database using it's ID.
// 2. Copy the encrypted key copy from the invite row.
// 3. Insert new membership in the organisation for the invitee, their key copy and the assigned role from invite row.
// 4. Mark the invite accepted.
func (d *DefaultService) Accept(ctx context.ServiceContext, client *clients.GQLClient, id string) error {

	//	Get the invite
	invite, err := d.Get(ctx, client, id)
	if err != nil {
		return err
	}

	//	If the invite is already marked accepted, return an error.
	if invite.Accepted {
		return fmt.Errorf("invite already accepted")
	}

	//	Get the invitee user.
	user, err := users.GetByEmail(ctx, client, invite.Email)
	if err != nil {
		return err
	}

	//	Add the admin header to the client to be able to add membership to the organisation.
	client.Headers = append(client.Headers, clients.XHasuraAdminSecretHeader)

	//	Insert the membership.
	if err := memberships.CreateWithUserID(ctx, client, &memberships.CreateOptions{
		UserID: user.ID,
		OrgID:  invite.OrgID,
		RoleID: invite.RoleID,
		Key:    invite.Key,
	}); err != nil {
		return err
	}

	//	Mark the invite accepted.
	if err := d.Update(ctx, client, id, &UpdateOptions{
		Set: SetUpdateOptions{
			Accepted: true,
		},
	}); err != nil {
		return err
	}

	//	Update the invite limit of the organisation.
	if err := organisations.GetService().UpdateInviteLimit(ctx, client, &organisations.UpdateInviteLimitOptions{
		ID:               invite.OrgID,
		IncrementLimitBy: -1,
	}); err != nil {
		return err
	}

	return nil
}

func (*DefaultService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *UpdateOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($id: uuid!, $set: invites_set_input) {
		update_invites(where: {id: {_eq: $id}}, _set: $set) {
		  affected_rows
		}
	  }			 
	`)

	req.Var("id", id)
	req.Var("set", options.Set)

	var response struct {
		Query struct {
			AffectedRows int `json:"affected_rows"`
		} `json:"update_invites"`
	}

	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	//	Validate the mutation as been written to the database
	if response.Query.AffectedRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
