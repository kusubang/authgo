package usr

type User struct {
	Id    string
	Email string
}

type UserServiceMem struct {
}

func (u *UserServiceMem) Find(id, pw string) (User, error) {
	return User{
		"user1",
		"user1@naver.com",
	}, nil
}

type UserService interface {
	Find(id, pw string) (User, error)
}
