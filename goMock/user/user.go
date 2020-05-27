package user

import "genosha/goMock/person"

type User struct {
  Person person.Job
}

func NewUser(p person.Job) *User {
  return &User{Person: p}
}

func (u *User) GetUserInfo(id int64) error {
  return u.Person.Get(id)
}
