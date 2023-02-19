package github

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/errors"
	"github.com/envsecrets/envsecrets/internal/integrations/commons"
)

func Callback(request *commons.GithubCallbackRequest) *errors.Error {

	//	Prepare the request body.
	reqBody, err := json.Marshal(map[string]interface{}{
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"code":          request.Code,
		"redirect_uri":  "https://webhook.site/664ecabd-9a72-4876-a491-8373ef5ca3a7",
	})
	if err != nil {
		return errors.New(err, "failed to marshal oauth request body", errors.ErrorTypeJSONMarshal, errors.ErrorSourceGo)
	}

	//	Get user's access token from Github API.
	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.New(err, "failed prepare oauth access token request", errors.ErrorTypeBadRequest, errors.ErrorSourceGo)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New(err, "failed send oauth access token request", errors.ErrorTypeBadRequest, errors.ErrorSourceGithub)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.New(err, "failed to read oauth response body", errors.ErrorTypeBadResponse, errors.ErrorSourceGo)
	}

	var response commons.OauthAuthResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return errors.New(err, "failed to unmarshal oauth response body", errors.ErrorTypeJSONUnmarshal, errors.ErrorSourceGo)
	}

	//	TODO: Create a new record in Hasura.

	return nil
}
