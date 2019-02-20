package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gocontrib/auth"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"
)

func makeUserStore() *userStore {
	return &userStore{}
}

type userStore struct {
}

func (s *userStore) ValidateCredentials(ctx context.Context, username, password string) (auth.User, error) {
	// TODO detect phone and normalize it
	query := fmt.Sprintf(`query users($username: string, $password: string) {
        users(func: has(%s)) @filter(eq(email, $username) OR eq(login, $username) OR eq(phone, $username)) {
			uid
			name
			email
			role
            checkpwd(password, $password)
        }
	}`, userLabel())

	vars := map[string]string{
		"$username": username,
		"$password": password,
	}

	return s.FindUser(ctx, query, vars, username, true)
}

func (s *userStore) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	query := fmt.Sprintf(`query users($id: string) {
        users(func: has(%s)) @filter(eq(email, $id)) {
			uid
			name
			email
			role
        }
	}`, userLabel())

	vars := map[string]string{
		"$id": email,
	}

	return s.FindUser(ctx, query, vars, email, false)
}

func (s *userStore) FindUserByID(ctx context.Context, userID string) (auth.User, error) {
	query := fmt.Sprintf(`query users($id: string) {
        users(func: uid($id)) @filter(has(%s)) {
			uid
			name
			email
			role
        }
	}`, userLabel())

	vars := map[string]string{
		"$id": userID,
	}

	return s.FindUser(ctx, query, vars, userID, false)
}

func userLabel() string {
	return dgraph.NodeLabel("user")
}

func (s *userStore) Close() {
}

func (s *userStore) FindUser(ctx context.Context, query string, vars map[string]string, userID string, checkPwd bool) (auth.User, error) {
	client, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, query, vars)
	if err != nil {
		log.Errorf("dgraph.Txn.QueryWithVars fail: %v", err)
		return nil, err
	}

	var result struct {
		Users []struct {
			ID       string `json:"uid"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Role     string `json:"role"`
			CheckPwd *bool  `json:"checkpwd(password)"`
		} `json:"users"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return nil, err
	}

	if len(result.Users) == 0 {
		return nil, fmt.Errorf("user not found by %s", userID)
	}

	user := result.Users[0]
	if checkPwd && user.CheckPwd != nil && !*user.CheckPwd {
		return nil, fmt.Errorf("wrong password: %s", userID)
	}

	return &auth.UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Admin: user.Role == "admin",
		Claims: map[string]interface{}{
			"email": user.Email,
			"role":  user.Role,
		},
	}, nil
}

func (s *userStore) CreateUser(ctx context.Context, account auth.UserData) (auth.User, error) {
	client, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	tx := client.NewTxn()
	defer tx.Discard(ctx)

	// TODO fill JSON using reflection
	in := make(utils.OrderedJSON)
	in["name"] = account.Name
	in["first_name"] = account.FirstName
	in["last_name"] = account.LastName
	in["email"] = account.Email
	in["avatar"] = account.AvatarURL
	in["location"] = account.Location

	_, err = dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: userLabel(),
		By:        "system",
	})
	if err != nil {
		return nil, err
	}

	// TODO optimize decode map to auth.UserInfo
	// result := results[0]

	return s.FindUserByEmail(ctx, account.Email)
}
