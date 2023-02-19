package integrations

import (
	"net/url"

	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
	"github.com/envsecrets/envsecrets/internal/integrations/github"
)

type Service interface {
	//Save(interface{}, commons.IntegrationType) error
	Delete(commons.IntegrationType) *errors.Error
	Callback(commons.IntegrationType, url.Values) *errors.Error
}

type DefaultIntegrationService struct{}

/* func (*DefaultIntegrationService) Save(payload interface{}, integrationType commons.IntegrationType) error {
	switch integrationType {
	case commons.ProjectIntegration:

		integration, ok := payload.(commons.Project)
		if !ok {
			return errors.New("failed type assertion to project integration")
		}
		return project.Save(&integration)

	case commons.AccountIntegration:

		integration, ok := payload.(commons.Account)
		if !ok {
			return errors.New("failed type assertion to account integration")
		}
		return account.Save(&integration)
	}
	return nil
}
*/
func (*DefaultIntegrationService) Delete(integrationType commons.IntegrationType) *errors.Error {
	switch integrationType {
	case commons.Github:
		return nil
	case commons.Vercel:
		return nil
	}

	return nil
}

func (*DefaultIntegrationService) Callback(integrationType commons.IntegrationType, params url.Values) *errors.Error {
	switch integrationType {
	case commons.Github:
		return github.Callback(&commons.GithubCallbackRequest{
			Code:           params.Get("code"),
			InstallationID: params.Get("installation_id"),
			SetupAction:    params.Get("setup_action"),
			State:          params.Get("state"),
		})
	case commons.Vercel:
		return nil
	}

	return nil
}
