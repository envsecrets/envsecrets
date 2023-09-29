package invites

import (
	"encoding/base64"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites/commons"
	"github.com/envsecrets/envsecrets/internal/invites/graphql"
	"github.com/envsecrets/envsecrets/internal/keys"
)

type Service interface {
	Get(context.ServiceContext, *clients.GQLClient, string) (*commons.Invite, error)
	Send(context.ServiceContext, *clients.GQLClient, *commons.SendOptions) error
	Update(context.ServiceContext, *clients.GQLClient, string, *commons.UpdateOptions) (*commons.Invite, error)
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
	copy(key[:], options.Key)
	inviteeKeyCopy, err := keys.SealAsymmetricallyAnonymous(public_key, key)
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

func (*DefaultInviteService) Update(ctx context.ServiceContext, client *clients.GQLClient, id string, options *commons.UpdateOptions) (*commons.Invite, error) {
	return graphql.Update(ctx, client, id, options)
}
