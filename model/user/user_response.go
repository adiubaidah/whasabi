package user

type UserLoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type UserLogoutResponse struct {
	Message string `json:"message"`
}
