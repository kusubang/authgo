package usr

type User struct {
	Id    string
	Email string
}

type UserService interface {
	Find(id, pw string) (User, error)
	IsValid(id, pw string) bool
}
