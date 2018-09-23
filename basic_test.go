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
			Reporter: httpexpect.NewAssertReporter(t),
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
		UID  string `json:"uid"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "Michael",
		Age:  39,
	}

	fmt.Println("CREATE")

	resp := c.expect.POST("/api/data/user").WithJSON(in).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	uid := resp.Path("$.uids[\"blank-0\"]").String().Raw()

	fmt.Println("QUERY")

	query := `{
		michael(func: eq(name, "Michael")) {
			uid
			name
			age
		}
	}`
	resp = c.expect.POST("/api/data/query").WithBytes([]byte(query)).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())

	fmt.Println("DELETE")

	resp = c.expect.DELETE("/api/data/user/" + uid).
		Expect().
		Status(http.StatusOK).
		JSON()
	printJSON(resp.Raw())
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
