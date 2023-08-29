package payloads

type Order struct {
	ID             int
	RecipientEmail string
	Status         string
}

type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// EmailContent represents the content of an email to be sent.
type EmailContent struct {
	Sender      Sender      `json:"sender"`
	To          []Recipient `json:"to"`
	Subject     string      `json:"subject"`
	HtmlContent string      `json:"htmlContent"`
}
