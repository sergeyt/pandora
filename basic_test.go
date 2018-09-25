package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	parseConfig()
	initSchema()

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

func TestCRUD(t *testing.T) {
	c := setup(t)
	defer c.Close()

	in := &struct {
		UID  string `json:"uid,omitempty"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "bob",
		Age:  39,
	}

	fmt.Println("CREATE")

	resp := c.expect.POST("/api/data/user").WithJSON(in).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	id := resp.Path("$.uid").String().Raw()

	fmt.Println("GET BY ID")

	resp = c.expect.GET("/api/data/user/" + id).
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
	resp = c.expect.POST("/api/query").WithBytes([]byte(query)).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("UPDATE")

	in = &struct {
		UID  string `json:"uid,omitempty"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "rob",
		Age:  42,
	}

	resp = c.expect.PUT("/api/data/user/" + id).WithJSON(in).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("GET BY ID")

	resp = c.expect.GET("/api/data/user/" + id).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("DELETE")

	resp = c.expect.DELETE("/api/data/user/" + id).
		Expect().
		Status(http.StatusOK).
		JSON()
	printJSON(resp.Raw())
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
