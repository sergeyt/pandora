package auth

import (
	"context"
	"os"

	goauth "github.com/gocontrib/auth"
	log "github.com/sirupsen/logrus"
)

// InitUsers ensures system users
func InitUsers() {
	ctx := context.Background()
	system := goauth.UserData{
		Name:     "system",
		NickName: "system",
		Email:    os.Getenv("SYSTEM_EMAIL"),
		Password: os.Getenv("SYSTEM_PWD"),
	}
	admin := goauth.UserData{
		Name:     "admin",
		NickName: "admin",
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PWD"),
		Role:     "admin",
	}
	users := []goauth.UserData{
		system,
		admin,
	}
	store := makeUserStore()
	for _, u := range users {
		log.Infof("init user %s", u.Name)
		_, err := store.CreateUser(ctx, u)
		if err != nil {
			log.Errorf("CreateUser fail: %v", err)
		}
	}
}
