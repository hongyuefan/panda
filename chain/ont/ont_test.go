package ont

import (
	"fmt"
	"testing"
)

var (
	opC *OpChain
)

func TestMain(t *testing.M) {
	var err error
	opC, err = NewOpChain("./wallet.dat", "http://localhost:20336")
	if err != nil {
		fmt.Errorf("newopchain error:%s", err)
		return
	}
	t.Run()
}

//func TestGenKey(t *testing.T) {

//	pub, priv, err := opC.GenKeyPair("fanhongyue")

//	t.Log(pub, priv, err)

//	return
//}

func TestTransfer(t *testing.T) {

	txhash, err := opC.DoTransaction("ont", "123456", "ARFwSSfVgsabq11EbhvDdtZrivtGzdp5gi", 0, 20000, 1)

	t.Log(txhash, err)

	return
}

func TestGetBalance(t *testing.T) {

	balance, err := opC.GetBalance("ont", "123456")

	t.Log(balance, err)

	balance, err = opC.GetBalance("fanhongyue", "X7F5yJ")

	t.Log(balance, err)

	return
}

func TestQueryTransaction(t *testing.T) {

	ret, err := opC.QueryTransaction("10351bcc31093d732a767b40e190d198e2bcbe6f20da67f4d9b81ff3a9934640")

	t.Log(ret, err)

	return
}
