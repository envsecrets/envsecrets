package commons

type InvitationOptions struct {
	ID            string
	Key           string
	ReceiverEmail string
	SenderID      string
	OrgID         string
}

type SendKeyOptions struct {
	UserID  string
	Key     string
	OrgName string
}
