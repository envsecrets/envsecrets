package actions

import (
	"encoding/json"
)

func MapToStruct(source any, target interface{}) error {
	entity, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(entity, &target)
}
