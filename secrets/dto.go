package secrets

import "fmt"

type Secrets []Secret

func (s *Secrets) Map() map[string]interface{} {

	response := make(map[string]interface{})

	for _, item := range *s {
		response[item.Key] = item.Value
	}

	return response
}

type Secret struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (s *Secret) String() string {
	return fmt.Sprintf("%s=%v", s.Key, s.Value)
}

func (s *Secret) Map() map[string]interface{} {
	return map[string]interface{}{
		s.Key: s.Value,
	}
}

type Path struct {
	Organisation string `json:"org"`
	Project      string `json:"project"`
	Environment  string `json:"env"`
}

func (p *Path) Location() string {
	return fmt.Sprintf("%s/%s/%s", p.Organisation, p.Project, p.Environment)
}
