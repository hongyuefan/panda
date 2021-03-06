package controllers

import (
	"errors"
	"strings"
)

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
		err = errors.New("password length not enough")
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
