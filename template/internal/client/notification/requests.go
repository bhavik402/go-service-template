package notification

type SendWelcomeEmailRequest struct {
	ToEmail string `json:"to_email"`
	ToName  string `json:"to_name"`
}
