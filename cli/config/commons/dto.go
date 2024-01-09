package commons

import (
	"encoding/base64"
	"encoding/json"

	"github.com/envsecrets/envsecrets/internal/users"
)

type Account struct {
	AccessToken  string     `json:"access_token" yaml:"accessToken"`
	RefreshToken string     `json:"refresh_token" yaml:"refreshToken"`
	User         users.User `json:"user" yaml:"user"`
}

type Project struct {
	//OrgID     string `json:"org_id" yaml:"org_id"`
	ProjectID string `json:"project_id" yaml:"project_id"`
	Key       []byte `json:"key" yaml:"key"`
	//AutoCapitalize bool   `json:"auto_capitalize" yaml:"auto_capitalize"`
}

// Custom json marshalling.
// We need to do this because we want to store the org key as a base64 encoded string in the config file.
// We can't store it as a byte array because yaml marshalling will convert it to a string.
func (p *Project) MarshalJSON() ([]byte, error) {

	var structure struct {
		//OrgID     string `json:"org_id" yaml:"org_id"`
		ProjectID string `json:"project_id" yaml:"project_id"`
		Key       string `json:"key" yaml:"key"`
		//AutoCapitalize bool   `json:"auto_capitalize" yaml:"auto_capitalize"`
	}

	//structure.OrgID = p.OrgID
	structure.ProjectID = p.ProjectID
	structure.Key = base64.StdEncoding.EncodeToString(p.Key)
	//structure.AutoCapitalize = p.AutoCapitalize
	return json.Marshal(structure)
}

// Custom json unmarshalling.
// We need to do this because we want to store the org key as a base64 encoded string in the config file.
func (p *Project) UnmarshalJSON(data []byte) error {

	var structure struct {
		OrgID     string `json:"org_id" yaml:"org_id"`
		ProjectID string `json:"project_id" yaml:"project_id"`
		Key       string `json:"key" yaml:"key"`
		//AutoCapitalize bool   `json:"auto_capitalize" yaml:"auto_capitalize"`
	}

	if err := json.Unmarshal(data, &structure); err != nil {
		return err
	}

	key, err := base64.StdEncoding.DecodeString(structure.Key)
	if err != nil {
		return err
	}

	*p = Project{
		//OrgID:     structure.OrgID,
		ProjectID: structure.ProjectID,
		Key:       key,
		//AutoCapitalize: structure.AutoCapitalize,
	}

	return nil
}

/*
	 func (p *ProjectStringified) Unstringify() (*Project, error) {

		key, err := base64.StdEncoding.DecodeString(p.Key)
		if err != nil {
			return nil, err
		}

		return &Project{
			Organisation:   p.Organisation,
			Project:        p.Project,
			Key:         key,
			AutoCapitalize: p.AutoCapitalize,
		}, nil
	}

	func (p *Project) Stringify() *ProjectStringified {
		return &ProjectStringified{
			Organisation:   p.Organisation,
			Project:        p.Project,
			Key:         base64.StdEncoding.EncodeToString(p.Key),
			AutoCapitalize: p.AutoCapitalize,
		}
	}
*/
type Keys struct {
	Public  []byte `json:"public_key" yaml:"public_key"`
	Private []byte `json:"private_key" yaml:"private_key"`
	Sync    []byte `json:"sync_key" yaml:"sync_key"`
}

type KeysStringified struct {
	Public  string `json:"public_key" yaml:"public_key"`
	Private string `json:"private_key" yaml:"private_key"`
	Sync    string `json:"sync_key" yaml:"sync_key"`
}

func (k *Keys) Stringify() *KeysStringified {
	return &KeysStringified{
		Public:  base64.StdEncoding.EncodeToString(k.Public),
		Private: base64.StdEncoding.EncodeToString(k.Private),
		Sync:    base64.StdEncoding.EncodeToString(k.Sync),
	}
}

func (k *KeysStringified) Unstringify() (*Keys, error) {

	public, err := base64.StdEncoding.DecodeString(k.Public)
	if err != nil {
		return nil, err
	}

	private, err := base64.StdEncoding.DecodeString(k.Private)
	if err != nil {
		return nil, err
	}

	sync, err := base64.StdEncoding.DecodeString(k.Sync)
	if err != nil {
		return nil, err
	}

	return &Keys{
		Public:  public,
		Private: private,
		Sync:    sync,
	}, nil
}
