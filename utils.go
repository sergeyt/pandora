package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OrderedJSON map[string]interface{}

func (m OrderedJSON) ToJSON(order ...string) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{', '\n'})

	written := make(map[string]struct{})
	i := 0

	for _, k := range order {
		v, ok := m[k]
		if !ok {
			continue
		}
		if i > 0 {
			buf.Write([]byte{',', '\n'})
		}
		m.writeEntry(buf, k, v)
		written[k] = struct{}{}
		i = i + 1
	}

	for k, v := range m {
		if _, ok := written[k]; ok {
			continue
		}
		if i > 0 {
			buf.Write([]byte{',', '\n'})
		}
		m.writeEntry(buf, k, v)
		i = i + 1
	}

	buf.Write([]byte{'\n', '}', '\n'})

	return buf.Bytes(), nil
}

func (m OrderedJSON) writeEntry(buf *bytes.Buffer, k string, v interface{}) error {
	fmt.Fprintf(buf, "  \"%s\": ", k)
	b, err := json.MarshalIndent(v, "  ", "  ")
	if err != nil {
		return err
	}
	_, err = buf.Write(b)
	if err != nil {
		return err
	}
	return nil
}
