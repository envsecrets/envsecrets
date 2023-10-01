package commons

type InvitationOptions struct {
	ID            string
	ReceiverEmail string
	SenderID      string
	OrgID         string
}

type SendKeyOptions struct {
	UserID  string
	Key     string
	OrgName string
}
