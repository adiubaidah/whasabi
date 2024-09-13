package user

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogoutRequest struct {
	Token string `json:"token"`
}
