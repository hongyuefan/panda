package controllers

import (
	"fmt"
	"panda/models"
	"panda/types"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	Expire_Time = 30
	Secret_Key  = "this is panda block chain !"
)

func TokenGenerate(userId int64) (token string, err error) {

	ojwt := jwt.New(jwt.SigningMethodHS256)

	ojwt.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * Expire_Time).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    fmt.Sprintf("%v", userId),
	}

	return ojwt.SignedString([]byte(Secret_Key))
}

func TokenValidate(token string) (ok bool, userId int64, err error) {

	ojwt, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret_Key), nil
	})

	if err != nil {
		return false, 0, err
	}

	if claims, ok := ojwt.Claims.(*jwt.StandardClaims); ok && ojwt.Valid {

		userId, _ = strconv.ParseInt(claims.Issuer, 10, 64)

		return ok, userId, nil
	}

	return false, 0, nil

}

func TokenSelect(userId int64) (token string, err error) {

	oToken := &types.Token{
		Uid: userId,
	}

	orm := models.NewCommon()

	if err = orm.CommonGetOne(oToken, "uid"); err != nil {
		return "", err
	}
	return oToken.Token, nil
}

func TokenInsert(token string, userId int64) (err error) {

	return
}
