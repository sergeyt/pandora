package apiutil

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Pagination struct {
	Offset int
	Limit  int
}

func ParsePagination(r *http.Request) (result Pagination, err error) {
	offset, err := parseIntParam(r, "offset", 0, true, false)
	if err != nil {
		return result, err
	}

	limit, err := parseIntParam(r, "limit", 100, true, false)
	if err != nil {
		return result, err
	}

	result.Offset = int(offset)
	result.Limit = int(limit)

	return result, nil
}

func parseIntParam(r *http.Request, name string, defval int, nonNegative, positive bool) (int, error) {
	s := strings.TrimSpace(r.URL.Query().Get(name))
	if len(s) == 0 {
		return defval, nil
	}
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s param is not valid: expect integer number. %s", name, err)
	}
	if nonNegative && val < 0 {
		return 0, fmt.Errorf("%s param is not valid: expect non negative integer number", name)
	}
	if positive && val <= 0 {
		return 0, fmt.Errorf("%s param is not valid: expect positive integer number", name)
	}
	return int(val), nil
}
