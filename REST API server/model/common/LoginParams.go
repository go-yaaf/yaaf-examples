package model

// LoginParams model used for authorize user by email or by SMS
// @Data
type LoginParams struct {
	Email  string `json:"email"`  // User email for login authentication or mail verification
	Mobile string `json:"mobile"` // User mobile phone for SMS verification
	Token  string `json:"token"`  // User token if he wsa already authenticated
}
