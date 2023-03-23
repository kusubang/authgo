package usr

import "errors"

type User struct {
	Id    string
	Email string
}

type UserServiceMem struct {
}

func (u *UserServiceMem) Find(id, pw string) (User, error) {

	if id != "user1" || pw != "1234" {
		return User{}, errors.New("user not found")
	}

	return User{
		"user1",
		"user1@naver.com",
	}, nil
}

type UserService interface {
	Find(id, pw string) (User, error)
}
