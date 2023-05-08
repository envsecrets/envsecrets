package circle

import secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"

func toKVMap(data map[string]secretCommons.Payload) map[string]interface{} {
	response := make(map[string]interface{})
	for key, payload := range data {
		response[key] = payload.Value
	}
	return response
}
