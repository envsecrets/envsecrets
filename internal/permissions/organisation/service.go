package organisation

import (
	"errors"

	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/machinebox/graphql"
)

// Insert new permissions.
func Insert(ctx context.ServiceContext, client *clients.GQLClient, options *commons.OrganisationPermissionsInsertOptions) error {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!, $user_id: uuid!, $role_id: uuid!, $key: String!) {
		insert_org_has_user(objects: {user_id: $user_id, org_id: $org_id, role_id: $role_id, key: $key}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)
	req.Var("user_id", options.UserID)
	req.Var("role_id", options.RoleID)
	req.Var("key", options.Key)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_org_has_user"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New("failed to insert permission")
	}

	return nil
}
