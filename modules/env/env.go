package env

import "os"

// Get environment variable or return given default value if variable is not defined
func Get(name, defval string) string {
	val := os.Getenv(name)
	if val == "" {
		val = defval
	}
	return val
}
