package nhost

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/internal/context"
)

func Signup(ctx context.ServiceContext, options *SignupOptions) *Error {

	body, err := options.Marshal()
	if err != nil {
		return &Error{Message: err.Error(), Code: http.StatusText(http.StatusBadRequest)}
	}

	//	Initialize a new request
	req, err := http.NewRequest(http.MethodPost, os.Getenv("NHOST_AUTH_URL")+"/signup/email-password", bytes.NewBuffer(body))
	if err != nil {
		return &Error{Message: err.Error(), Code: http.StatusText(http.StatusBadRequest)}
	}

	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Error{Message: err.Error(), Code: http.StatusText(http.StatusBadRequest)}
	}

	if resp.StatusCode != 200 {

		defer resp.Body.Close()

		var response Error
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &Error{Message: err.Error(), Code: http.StatusText(http.StatusBadRequest)}
		}

		if err := json.Unmarshal(result, &response); err != nil {
			return &Error{Message: err.Error(), Code: http.StatusText(http.StatusBadRequest)}
		}

		return &response
	}

	return nil
}
