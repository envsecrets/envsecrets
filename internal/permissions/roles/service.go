package roles

import (
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/machinebox/graphql"
)

//	Insert new permissions.
func Insert(ctx context.ServiceContext, client *clients.GQLClient, options *commons.RoleInsertOptions) *errors.Error {

	req := graphql.NewRequest(`
	mutation MyMutation($name: String!, $org_id: uuid!, $permissions: jsonb!) {
		insert_roles(objects: {name: $name, org_id: $org_id, permissions: $permissions}) {
		  affected_rows
		}
	  }	   
	`)

	req.Var("org_id", options.OrgID)
	req.Var("name", options.Name)
	req.Var("permissions", options.Permissions)

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_roles"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, "failed to insert role", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
