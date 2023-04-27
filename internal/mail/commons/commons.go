package commons

import "github.com/matcornic/hermes/v2"

type MailType string

var (

	// Configure hermes by setting a theme and your product info
	Hermes = hermes.Hermes{
		// Optional Theme
		//	Theme: new(Default),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name:        "envsecrets.com",
			Link:        "https://envsecrets.com",
			TroubleText: "If the {ACTION}-button is not working for you, just copy and paste the URL below into your web browser.",
			Copyright:   "GSTIN: 09AFVPW6899R2ZD",
			// Optional product logo
			//	Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
)
