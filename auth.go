package main

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/gocontrib/auth"
)

func authAPI(mux chi.Router) {
	config := makeAuthConfig()
	mux.Post("/api/login", auth.LoginHandlerFunc(config))
}

func makeAuthConfig() *auth.Config {
	return &auth.Config{
		UserStore: makeUserStore(),
	}
}

func makeUserStore() auth.UserStore {
	return &userStoreImpl{}
}

var fakeUsers = []*userInfo{
	&userInfo{
		ID:    "1",
		Name:  "admin",
		Email: "stodyshev@gmail.com",
		Admin: true,
	},
	&userInfo{
		ID:    "2",
		Name:  "sergeyt",
		Email: "stodyshev@gmail.com",
	},
}

type userStoreImpl struct {
}

func (s *userStoreImpl) ValidateCredentials(username, password string) (auth.User, error) {
	for _, u := range fakeUsers {
		if u.Name == username && password == username {
			return u, nil
		}
	}
	return nil, fmt.Errorf("%s user not found", username)
}

func (s *userStoreImpl) FindUserByID(userID string) (auth.User, error) {
	for _, u := range fakeUsers {
		if u.ID == userID {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *userStoreImpl) Close() {
}

type userInfo struct {
	ID    string
	Name  string
	Email string
	Admin bool
}

func (u *userInfo) GetID() string    { return u.ID }
func (u *userInfo) GetName() string  { return u.Name }
func (u *userInfo) GetEmail() string { return u.Email }
func (u *userInfo) IsAdmin() bool    { return u.Admin }
