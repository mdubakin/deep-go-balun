package main

import "fmt"

type Option func(*User)

func WithEmail(email string) Option {
	return func(user *User) {
		if nil == user.Email {
			user.Email = new(string)
		}
		*user.Email = email
	}
}

func WithPhone(phone string) Option {
	return func(user *User) {
		if nil == user.Phone {
			user.Phone = new(string)
		}
		*user.Phone = phone
	}
}

func WithAddress(address string) Option {
	return func(user *User) {
		if nil == user.Address {
			user.Address = new(string)
		}
		*user.Address = address
	}
}

type User struct {
	Name    string
	Surname string
	Email   *string
	Phone   *string
	Address *string
}

func NewUser(name string, surname string, options ...Option) User {
	user := User{
		Name:    name,
		Surname: surname,
	}

	for _, option := range options {
		option(&user)
	}

	return user
}

func main() {
	user1 := NewUser("Ivan", "Ivanov", WithEmail("ivanov@yandex.ru"))
	user2 := NewUser("Petr", "Petrov", WithEmail("petrov@yandex.ru"), WithPhone("+67453"))

	fmt.Printf("%#v\n", user1)
	fmt.Printf("%#v\n", user2)
}
