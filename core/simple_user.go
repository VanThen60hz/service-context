package core

// SimpleUser only contains public infos
type SimpleUser struct {
	SQLModel
	Email     string `json:"email" gorm:"column:email" db:"email"`
	LastName  string `json:"last_name" gorm:"column:last_name;" db:"last_name"`
	FirstName string `json:"first_name" gorm:"column:first_name;" db:"first_name"`
	AvatarId  int    `json:"-" gorm:"column:avatar_id" db:"avatar_id"`
	Avatar    *Image `json:"avatar" gorm:"-" db:"-"`
}

func (SimpleUser) TableName() string {
	return "users"
}

func NewSimpleUser(id int, email, firstName, lastName string, avatarId int) SimpleUser {
	return SimpleUser{
		SQLModel:  SQLModel{Id: id},
		Email:     email,
		LastName:  lastName,
		FirstName: firstName,
		AvatarId:  avatarId,
	}
}
