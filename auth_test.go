package main

import (
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	c := setup(t)
	defer c.Close()

	resp := c.expect.POST("/api/login").WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusOK).
		JSON()

	printJSON(resp.Raw())
}
