package cloudstore

import "os"

func noop() {}

func env(name, defval string) string {
	val := os.Getenv(name)
	if val == "" {
		val = defval
	}
	return val
}
