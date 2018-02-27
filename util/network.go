package util

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Get(url string) (b []byte, err error) {
	body, err := http.Get(url)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, body.Body)
	b = buf.Bytes()
	return
}
