package supabase

import secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"

func transform(data map[string]secretCommons.Payload) (result []map[string]interface{}) {
	for key, payload := range data {
		result = append(result, map[string]interface{}{
			"name":  key,
			"value": payload.Value,
		})
	}
	return
}
