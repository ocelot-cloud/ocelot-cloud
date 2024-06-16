package main

type UserManager interface {
	DoesUserExist(username string) bool
	CreateUser(username, password string) error
	LoginUser(username, password string) (bool, error)
	DeleteUser(username, password string) error
}
