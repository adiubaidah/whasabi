package model

type User struct {
	ID       int    `gorm:"primaryKey;colum:id;autoIncrement"`
	Username string `gorm:"colum:username;type:varchar(255);not null"`
	Password string `gorm:"colum:password;type:varchar(255);not null"`
	Role     string `gorm:"colum:role;not null"`
	Service  Ai     `gorm:"foreignKey:user_id;references:id"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLogoutRequest struct {
	Token string `json:"token"`
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type UserLoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type UserLogoutResponse struct {
	Message string `json:"message"`
}
