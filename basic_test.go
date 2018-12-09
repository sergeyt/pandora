package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergeyt/pandora/modules/dgraph"

	"github.com/sergeyt/pandora/modules/config"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

type TC struct {
	t      *testing.T
	server *httptest.Server
	expect *httpexpect.Expect
}

func (c *TC) Close() {
	c.server.Close()
}

func setup(t *testing.T) *TC {
	config.Parse()

	// TODO separate config for testing
	// HACK to run test from host we have to use host dgraph address
	config.DB.Addr = "localhost:9080"

	dgraph.InitSchema()

	server := httptest.NewServer(makeAPIHandler())

	return &TC{
		t:      t,
		server: server,
		expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL:  server.URL,
			Reporter: httpexpect.NewRequireReporter(t),
			Printers: []httpexpect.Printer{
				httpexpect.NewDebugPrinter(t, true),
			},
		}),
	}
}

type TestUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestCRUD(t *testing.T) {
	c := setup(t)
	defer c.Close()

	in := &TestUser{
		Name: "bob",
		Age:  39,
	}

	fmt.Println("CREATE")

	authorization := "local_admin"

	c.expect.POST("/api/data/user").
		WithJSON(in).
		Expect().
		Status(http.StatusUnauthorized)

	resp := c.expect.POST("/api/data/user").
		WithHeader("Authorization", authorization).
		WithJSON(in).
		Expect().
		Status(http.StatusOK).
		JSON()

	id := resp.Path("$.uid").String().Raw()

	printJSON(resp.Raw())

	in = &TestUser{
		Name: "joe",
		Age:  40,
	}

	resp = c.expect.POST("/api/data/user").
		WithHeader("Authorization", authorization).
		WithJSON(in).
		Expect().
		Status(http.StatusOK).
		JSON()

	id2 := resp.Path("$.uid").String().Raw()

	printJSON(resp.Raw())

	fmt.Println("GET LIST")

	resp = c.expect.GET("/api/data/user/list").
		WithHeader("Authorization", authorization).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("GET BY ID")

	resp = c.expect.GET("/api/data/user/"+id).
		WithHeader("Authorization", authorization).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("QUERY")

	query := `{
		data(func: eq(name, "bob")) @filter(has(_user)) {
			uid
			name
			age
		}
	}`
	resp = c.expect.POST("/api/query").
		WithHeader("Authorization", authorization).
		WithBytes([]byte(query)).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("UPDATE")

	in = &TestUser{
		Name: "rob",
		Age:  42,
	}

	resp = c.expect.PUT("/api/data/user/"+id).
		WithHeader("Authorization", authorization).
		WithJSON(in).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("GET BY ID")

	resp = c.expect.GET("/api/data/user/"+id).
		WithHeader("Authorization", authorization).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("DELETE")

	c.expect.DELETE("/api/data/user/"+id).
		WithHeader("Authorization", authorization).
		Expect().
		Status(http.StatusOK)

	c.expect.DELETE("/api/data/user/"+id2).
		WithHeader("Authorization", authorization).
		Expect().
		Status(http.StatusOK)
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
