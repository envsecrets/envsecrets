package auth

type SigninOptions struct {
	Email    string `json:"email,omitempty"`
	OTP      string `json:"otp,omitempty"`
	Ticket   string `json:"ticket,omitempty"`
	Password string `json:"password,omitempty"`
}

type ToggleMFAOptions struct {
	Code string `json:"code"`
}

type UpdatePasswordOptions struct {
	NewPassword string `json:"newPassword,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
}

type SignupOptions struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
