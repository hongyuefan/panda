package controllers

import (
	"errors"
	"panda/types"

	"github.com/astaxie/beego"
)

type TransactionContoller struct {
	beego.Controller
}

var (
	PublicKey_Setting string = "0xfab9898d9898s988c898989876756787879809-0-90898"
	PrivKey_Setting   string
)

func (t *TransactionContoller) Transactions(ntype int, sourPub, sourPriv string, amount uint64) (txhash string, err error) {

	if err = ValidatePublicKey(sourPub); err != nil {
		return
	}

	switch ntype {
	case types.Type_Catch:
		if err = ValidatePrivKey(sourPriv); err != nil {
			return
		}
		return DoTransaction(sourPub, sourPriv, PublicKey_Setting, amount)
	case types.Type_Train:
		if err = ValidatePrivKey(sourPriv); err != nil {
			return
		}
		return DoTransaction(sourPub, sourPriv, PublicKey_Setting, amount)
	case types.Type_Bonus:
		return DoTransaction(PublicKey_Setting, PrivKey_Setting, sourPub, amount)
	}
	return "", errors.New("transaction type not right")
}

func (t *TransaccountController) IsTransactionSuccess(txhash string) (err error) {
	return nil
}

func DoTransaction(sourPub, sourPriv, desPub string, amount uint64) (txhash string, err error) {

	return "0x12345678", nil
}
