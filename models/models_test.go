package models

import (
	"testing"
)

func TestGetAllAgree(t *testing.T) {
	init()
	query := make(map[string]string, 0)
	result, _ := GetAllAgree(query, []string{""}, []string{"id"}, []string{"asc"}, 0, 100)
	t.Log(result)
}
