package invites

import (
	"encoding/base64"
	"fmt"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/invites/graphql"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/users"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Invite, error)
	Send(context.ServiceContext, *clients.GQLClient, *commons.SendOptions) error
	Accept(context.ServiceContext, *clients.GQLClient, string) error
	Update(context.ServiceContext, *clients.GQLClient, string, *commons.UpdateOptions) error
}

type DefaultInviteService struct{}

func (*DefaultInviteService) Get(ctx context.ServiceContext, client *clients.GQLClient, id string) (*commons.Invite, error) {
	return graphql.Get(ctx, client, id)
}

func (*DefaultInviteService) Send(ctx context.ServiceContext, client *clients.GQLClient, options *commons.SendOptions) error {

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

	return graphql.Insert(ctx, client, []commons.InsertOptions{
		{
			Key:    base64.StdEncoding.EncodeToString(inviteeKeyCopy),
			OrgID:  options.OrgID,
			RoleID: options.RoleID,
			Email:  options.InviteeEmail,
			UserID: options.InviterID,
		},
	})
}

// ---	Flow ---
// 1. Fetch the invite row from database using it's ID.
// 2. Copy the encrypted key copy from the invite row.
// 3. Insert new membership in the organisation for the invitee, their key copy and the assigned role from invite row.
// 4. Mark the invite accepted.
func (d *DefaultInviteService) Accept(ctx context.ServiceContext, client *clients.GQLClient, id string) error {

	//	Get the invite
	invite, err := graphql.Get(ctx, client, id)
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
	return d.Update(ctx, client, id, &commons.UpdateOptions{
		Set: commons.SetUpdateOptions{
			Accepted: true,
		},
	})
}

func (*DefaultInviteService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *commons.UpdateOptions) error {
	return graphql.Update(ctx, client, id, options)
}
