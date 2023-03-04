package organisation

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/machinebox/graphql"
)

//	Insert new permissions.
func Insert(ctx context.ServiceContext, client *clients.GQLClient, options *commons.OrganisationPermissionsInsertOptions) *errors.Error {

	req := graphql.NewRequest(`
	mutation MyMutation($org_id: uuid!, $user_id: uuid!, $permissions: jsonb!) {
		insert_org_level_permissions(objects: {user_id: $user_id, org_id: $org_id, permissions: $permissions}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("org_id", options.OrgID)
	req.Var("user_id", options.UserID)
	permissions, _ := options.Permissions.Marshal()
	req.Var("permissions", string(permissions))

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_org_level_permissions"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, "failed to insert permissions", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
