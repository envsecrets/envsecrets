package commons

import (
	"encoding/base64"

	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
)

type Account struct {
	AccessToken  string           `json:"access_token" yaml:"accessToken"`
	RefreshToken string           `json:"refresh_token" yaml:"refreshToken"`
	User         userCommons.User `json:"user" yaml:"user"`
}

type ProjectStringified struct {
	Version        int    `json:"version" yaml:"version"`
	Organisation   string `json:"organisation" yaml:"organisation"`
	Project        string `json:"project" yaml:"project"`
	Environment    string `json:"environment" yaml:"environment"`
	OrgKey         string `json:"org_key" yaml:"org_key"`
	AutoCapitalize bool   `json:"auto_capitalize" yaml:"auto_capitalize"`
}

type Project struct {
	Version        int    `json:"version" yaml:"version"`
	Organisation   string `json:"organisation" yaml:"organisation"`
	Project        string `json:"project" yaml:"project"`
	Environment    string `json:"environment" yaml:"environment"`
	OrgKey         []byte `json:"org_key" yaml:"org_key"`
	AutoCapitalize bool   `json:"auto_capitalize" yaml:"auto_capitalize"`
}

func (p *ProjectStringified) Unstringify() (*Project, error) {

	key, err := base64.StdEncoding.DecodeString(p.OrgKey)
	if err != nil {
		return nil, err
	}

	return &Project{
		Version:        p.Version,
		Organisation:   p.Organisation,
		Project:        p.Project,
		Environment:    p.Environment,
		OrgKey:         key,
		AutoCapitalize: p.AutoCapitalize,
	}, nil
}

func (p *Project) Stringify() *ProjectStringified {
	return &ProjectStringified{
		Version:        p.Version,
		Organisation:   p.Organisation,
		Project:        p.Project,
		Environment:    p.Environment,
		OrgKey:         base64.StdEncoding.EncodeToString(p.OrgKey),
		AutoCapitalize: p.AutoCapitalize,
	}
}

type Keys struct {
	Public  []byte `json:"public_key" yaml:"public_key"`
	Private []byte `json:"private_key" yaml:"private_key"`
}

type KeysStringified struct {
	Public  string `json:"public_key" yaml:"public_key"`
	Private string `json:"private_key" yaml:"private_key"`
}

func (k *Keys) Stringify() *KeysStringified {
	return &KeysStringified{
		Public:  base64.StdEncoding.EncodeToString(k.Public),
		Private: base64.StdEncoding.EncodeToString(k.Private),
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

	return &Keys{
		Public:  public,
		Private: private,
	}, nil
}
