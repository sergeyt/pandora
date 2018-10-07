package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gocontrib/auth"
	"github.com/sergeyt/pandora/modules/dgraph"
)

func makeUserStore() auth.UserStore {
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
			Password []struct {
				CheckPwd bool `json:"checkpwd"`
			} `json:"password"`
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
	if checkPwd && !user.Password[0].CheckPwd {
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
