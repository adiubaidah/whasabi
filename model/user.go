package model

type User struct {
	ID       int     `gorm:"primaryKey;colum:id;autoIncrement"`
	Username string  `gorm:"colum:username;type:varchar(255);not null"`
	Password string  `gorm:"colum:password;type:varchar(255);not null"`
	IsActive bool    `gorm:"colum:is_active;default:false"`
	Role     string  `gorm:"colum:role;default:user;not null"`
	Process  Process `gorm:"foreignKey:user_id;references:id"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	Role     string `json:"role"`
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
	IsActive bool   `json:"is_active" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type UserRegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,eqfield=Password"`
}

type UserUpdateRequest struct {
	ID       int    `json:"id" validate:"required"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
	Role     string `json:"role"`
}

type UserLoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type UserLogoutResponse struct {
	Message string `json:"message"`
}
