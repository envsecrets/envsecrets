package subscriptions

type Status string

const (
	StatusActive            Status = "active"
	StatusIncomplete        Status = "incomplete"
	StatusIncompleteExpired Status = "incomplete_expired"
	StatusCanceled          Status = "canceled"
	StatusUnpaid            Status = "unpaid"
	StatusTrialing          Status = "trialing"
	StatusPastDue           Status = "past_due"
	StatusPaused            Status = "paused"
	StatusResumed           Status = "resumed"
)
