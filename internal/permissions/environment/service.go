package environment

import (
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/permissions/commons"
	"github.com/machinebox/graphql"
)

//	Insert new permissions.
func Insert(ctx context.ServiceContext, client *client.GQLClient, options *commons.EnvironmentPermissionsInsertOptions) *errors.Error {

	req := graphql.NewRequest(`
	mutation MyMutation($env_id: uuid!, $user_id: uuid!, $permissions: jsonb!) {
		insert_env_level_permissions(objects: {user_id: $user_id, env_id: $env_id, permissions: $permissions}) {
		  affected_rows
		}
	  }	  
	`)

	req.Var("env_id", options.EnvID)
	req.Var("user_id", options.UserID)
	permissions, _ := options.Permissions.Marshal()
	req.Var("permissions", string(permissions))

	var response map[string]interface{}
	if err := client.Do(ctx, req, &response); err != nil {
		return err
	}

	returned := response["insert_env_level_permissions"].(map[string]interface{})

	affectedRows := returned["affected_rows"].(float64)
	if affectedRows == 0 {
		return errors.New(nil, "failed to insert permissions", errors.ErrorTypeInvalidResponse, errors.ErrorSourceGraphQL)
	}

	return nil
}
