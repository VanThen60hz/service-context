package core

// SimpleUser only contains public infos
type SimpleUser struct {
	SQLModel
	Email     string `json:"email" gorm:"column:email" db:"email"`
	LastName  string `json:"last_name" gorm:"column:last_name;" db:"last_name"`
	FirstName string `json:"first_name" gorm:"column:first_name;" db:"first_name"`
	Avatar    *Image `json:"avatar" gorm:"column:avatar;" db:"avatar"`
}

func (SimpleUser) TableName() string {
	return "users"
}

func NewSimpleUser(id int, email, firstName, lastName string, avatar *Image) SimpleUser {
	return SimpleUser{
		SQLModel:  SQLModel{Id: id},
		Email:     email,
		LastName:  lastName,
		FirstName: firstName,
		Avatar:    avatar,
	}
}
