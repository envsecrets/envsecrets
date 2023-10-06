package invites

type SendOptions struct {
	OrgID           string `json:"org_id,omitempty"`
	RoleID          string `json:"role_id,omitempty"`
	InviteeEmail    string `json:"invitee_email,omitempty"`
	InviterPassword string `json:"inviter_password,omitempty"`
}

type AcceptOptions struct {
	ID string `json:"id,omitempty"`
}
