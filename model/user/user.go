package user

type User struct {
	ID       int    `gorm:"primaryKey;colum:id;autoIncrement"`
	username string `gorm:"colum:username;type:varchar(255);not null"`
	password string `gorm:"colum:password;type:varchar(255);not null"`
}
