package auth

import (
	"encoding/json"
	"errors"
)

type SigninWithPasswordOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninWithMFAOptions struct {
	OTP    string `json:"otp"`
	Ticket string `json:"ticket"`
}

// Custom marshaller for the SigninWithMFAOptions.
func (o *SigninWithMFAOptions) JSONMarshal() ([]byte, error) {

	//	OTP should be 6 digits.
	if len(o.OTP) != 6 {
		return nil, errors.New("otp should be 6 digits")
	}

	return json.Marshal(o)
}

type ToggleMFAOptions struct {
	Code string `json:"code"`

	//	The value of "null" will deactivate MFA.
	ActiveMFAType MFAType `json:"activeMfaType"`
}

// Custom marshaller for the ToggleMFAOptions.
func (o *ToggleMFAOptions) JSONMarshal() ([]byte, error) {

	//	OTP should be 6 digits.
	if len(o.Code) != 6 {
		return nil, errors.New("code should be 6 digits")
	}

	return json.Marshal(o)
}

type GenerateQRResponse struct {
	Secret string `json:"secret"`
	Image  string `json:"image"`
}

type SigninResponse struct {
	MFA     map[string]interface{} `json:"mfa"`
	Session map[string]interface{} `json:"session"`
	Keys    map[string]string      `json:"keys"`
}

type NhostSigninResponse struct {
	MFA     map[string]interface{} `json:"mfa"`
	Session map[string]interface{} `json:"session"`
}

type NhostSession struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UpdatePasswordOptions struct {
	NewPassword string `json:"newPassword,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
}

type DecryptKeysFromSessionOptions struct {
	Password string `json:"password"`
	Session  map[string]interface{}
}
