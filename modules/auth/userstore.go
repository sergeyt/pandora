package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gocontrib/auth"
	"github.com/markbates/goth"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
)

func makeUserStore() *UserStore {
	return &UserStore{}
}

type UserStore struct {
}

func (s *UserStore) ValidateCredentials(ctx context.Context, username, password string) (auth.User, error) {
	query := fmt.Sprintf(`{
        users(func: has(%s)) @filter(eq(email, %q) OR eq(login, %q)) {
			uid
			name
			email
			role
            checkpwd(password, %q)
        }
	}`, userLabel(), username, username, password)

	return s.FindUser(ctx, query, username, true)
}

func (s *UserStore) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	query := fmt.Sprintf(`{
        users(func: has(%s)) @filter(eq(email, %q)) {
			uid
			name
			email
			role
        }
	}`, userLabel(), email)

	return s.FindUser(ctx, query, email, false)
}

func (s *UserStore) FindUserByID(ctx context.Context, userID string) (auth.User, error) {
	query := fmt.Sprintf(`{
        users(func: uid(%s)) @filter(has(%s)) {
			uid
			name
			email
			role
        }
	}`, userID, userLabel())

	return s.FindUser(ctx, query, userID, false)
}

func userLabel() string {
	return dgraph.NodeLabel("user")
}

func (s *UserStore) Close() {
}

func (s *UserStore) FindUser(ctx context.Context, query, userID string, checkPwd bool) (auth.User, error) {
	client, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	txn := client.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.Query(ctx, query)
	if err != nil {
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
		return nil, err
	}

	if len(result.Users) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
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

type CreateUserData struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Location  string `json:"location"`
}

func (s *UserStore) CreateUser(ctx context.Context, account goth.User) (auth.User, error) {
	data := CreateUserData{
		Name:      account.Name,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		Email:     account.Email,
		// TODO login should be unique
		// Login: account.NickName,
	}

	client, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	tx := client.NewTxn()
	defer tx.Discard(ctx)

	// TODO fill JSON using reflection
	in := make(utils.OrderedJSON)
	in["name"] = data.Name
	in["first_name"] = data.FirstName
	in["last_name"] = data.LastName
	in["email"] = data.Email
	in["avatar"] = data.Avatar
	in["location"] = data.Location

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

	return s.FindUserByEmail(ctx, data.Email)
}
