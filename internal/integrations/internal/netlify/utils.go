package netlify

import secretCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"

func transform(data map[string]secretCommons.Payload) (result []map[string]interface{}) {
	for key, payload := range data {
		result = append(result, map[string]interface{}{
			"key":    key,
			"scopes": []string{"builds", "functions", "runtime", "post-processing"},
			"values": []map[string]interface{}{
				{
					"value":   payload.Value,
					"context": "all",
				},
			},
		})
	}
	return
}
