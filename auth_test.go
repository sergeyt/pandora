package main

import (
	"net/http"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	c := setup(t)
	defer c.Close()

	resp := c.expect.POST("/api/login").
		WithBasicAuth("admin", os.Getenv("ADMIN_PWD")).
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())
}
