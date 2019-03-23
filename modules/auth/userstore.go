package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/dgo"
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
		Admin: isAdmin(user.Role),
		Claims: map[string]interface{}{
			"email": user.Email,
			"role":  user.Role,
		},
	}, nil
}

func isAdmin(role string) bool {
	roles := splitRoles(role)
	isAdmin := false
	for _, r := range roles {
		if r == "admin" {
			isAdmin = true
			break
		}
	}
	return isAdmin
}

func splitRoles(s string) []string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0)
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if len(t) > 0 {
			result = append(result, strings.ToLower(t))
		}
	}
	return result
}

func (s *userStore) CreateUser(ctx context.Context, account auth.UserData) (auth.User, error) {
	client, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	tx := client.NewTxn()
	defer tx.Discard(ctx)

	// TODO generate unique login
	in := make(utils.OrderedJSON)
	in["name"] = account.Name
	in["first_name"] = account.FirstName
	in["last_name"] = account.LastName
	in["email"] = account.Email
	in["avatar"] = account.AvatarURL
	in["location"] = account.Location
	in["registered_at"] = time.Now()

	results, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: userLabel(),
		By:        "system",
	})
	if err != nil {
		return nil, err
	}

	user := mapUser(results[0])

	// create account for given oauth provider
	in = makeAccount(account)
	in["registered_at"] = time.Now()

	results, err = dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: dgraph.NodeLabel("account"),
		By:        "system",
	})
	if err != nil {
		return nil, err
	}

	// link account with user record
	acc := results[0]
	accountID := getString(acc, "uid")
	err = linkAccount(ctx, tx, user.ID, accountID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func makeAccount(account auth.UserData) utils.OrderedJSON {
	in := make(utils.OrderedJSON)
	in["provider"] = account.Provider
	in["email"] = account.Email
	in["name"] = account.Name
	in["first_name"] = account.FirstName
	in["last_name"] = account.LastName
	in["nick_name"] = account.NickName
	in["user_id"] = account.UserID
	in["description"] = account.Description
	in["avatar"] = account.AvatarURL
	in["location"] = account.Location
	in["role"] = account.Role
	return in
}

func linkAccount(ctx context.Context, tx *dgo.Txn, userID, accountID string) error {
	in := make(utils.OrderedJSON)
	in["account"] = map[string]interface{}{
		"uid": accountID,
	}

	_, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: userLabel(),
		ID:        userID,
		By:        "system",
	})

	return err
}

func (s *userStore) UpdateAccount(ctx context.Context, user auth.User, data auth.UserData) error {
	client, err := dgraph.NewClient()
	if err != nil {
		return err
	}

	tx := client.NewTxn()
	defer tx.Discard(ctx)

	query := `query($provider: string, $email: string) {
		accounts(has(Account)) @filter(eq(provider, $provider) AND eq(email, $email)) {
			uid
		}
	}`

	vars := map[string]string{
		"$provider": data.Provider,
		"$email":    data.Email,
	}

	resp, err := tx.QueryWithVars(ctx, query, vars)
	if err != nil {
		log.Errorf("dgraph.Txn.QueryWithVars fail: %v", err)
		return err
	}

	var result struct {
		Accounts []struct {
			ID string `json:"uid"`
		} `json:"accounts"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return err
	}

	accountID := ""
	if len(result.Accounts) > 0 {
		accountID = result.Accounts[0].ID
	}

	in := makeAccount(data)

	results, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: dgraph.NodeLabel("account"),
		ID:        accountID,
		By:        "system",
	})
	if err != nil {
		return err
	}

	// link account with user record
	acc := results[0]
	accountID = getString(acc, "uid")
	return linkAccount(ctx, tx, user.GetID(), accountID)
}

func mapUser(raw map[string]interface{}) *auth.UserInfo {
	email := getString(raw, "email")
	role := getString(raw, "role")
	return &auth.UserInfo{
		ID:    getString(raw, "uid"),
		Name:  getString(raw, "name"),
		Email: email,
		Admin: isAdmin(role),
		Claims: map[string]interface{}{
			"email": email,
			"role":  role,
		},
	}
}

func getString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if ok {
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}
