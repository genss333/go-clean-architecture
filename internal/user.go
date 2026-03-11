package internal

type User struct {
	ID        int
	UserName  string
	Firstname string
	Lastname  string
}

type UserRepository interface {
	GetUserByID(id int) (*User, error)
}
