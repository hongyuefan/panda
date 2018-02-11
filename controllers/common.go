package controllers

import (
	"errors"
	"strings"
)

func ParseToken(auth string) (token string, err error) {
	if !strings.HasPrefix(auth, "Bearer ") {
		err = errors.New("token format not right")
		return
	}
	tokens := strings.Fields(auth)
	return tokens[1], nil
}

func TranstateString(status int) string {
	switch status {
	case 0:
		return "transaction failed"
	case 1:
		return "transaction waiting"
	case 2:
		return "transaction success"
	}
	return "unknown status "
}

func TransTypeString(ntype int) string {
	switch ntype {
	case 0:
		return "catch "
	case 1:
		return "train "
	case 2:
		return "transfer "
	}
	return "unknown type "
}
