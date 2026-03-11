package infrastructure

import "github.com/genss333/go-clean-architecture/internal"


type UserModel struct {
	ID        int
	UserName  string
	Firstname string
	Lastname  string
}

func (model *UserModel) ToEntity() internal.User {
	return internal.User{
		ID:        model.ID,
		UserName:  model.UserName,
		Firstname: model.Firstname,
		Lastname:  model.Lastname,
	}
}
