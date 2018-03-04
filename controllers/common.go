package controllers

import (
	"errors"
	"panda/types"
	"strings"
)

func ParseAndValidToken(auth string) (userId int64, err error) {

	var token string

	if token, err = ParseToken(auth); err != nil {
		return
	}
	return TokenValidate(token)
}

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
	case types.Trans_Status_Waiting:
		return "交易验证中"
	case types.Trans_Status_Failed:
		return "交易失败"
	case types.Trans_Status_Success:
		return "交易成功"
	}
	return "unknown status "
}
