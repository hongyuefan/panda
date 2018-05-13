package controllers

import (
	"errors"
	"strings"
)

const (
	Len_Mobile = 11
)

func ValidMobile(mobile string) (err error) {
	if len(mobile) != Len_Mobile {
		return errors.New("手机长度不足")
	}
	if !strings.HasPrefix(mobile, "1") {
		return errors.New("手机格式不正确")
	}
	return
}

func ValidateEmail(userName string) (err error) {
	if len(userName) < 2 {
		err = errors.New("username length not enough")
		return err
	}
	if !strings.Contains(userName, "@") {
		err = errors.New("username format not right")
		return err
	}
	return nil
}

func ValidatePassWord(passWord string) (err error) {
	if len(passWord) < 6 {
		err = errors.New("密码长度至少6位")
		return err
	}
	return nil
}

func ValidatePublicKey(pub string) (err error) {

	if len(pub) < 32 {
		err = errors.New("publickey is not right")
		return
	}
	return nil
}

func ValidatePrivKey(priv string) (err error) {

	if len(priv) < 32 {
		err = errors.New("privkey is not right")
		return
	}
	return nil
}
