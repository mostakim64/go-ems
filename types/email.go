package types

type EmailPayload struct {
	MailTo  string      `json:"mail_to"`
	Subject string      `json:"subject"`
	Body    interface{} `json:"body"`
}
