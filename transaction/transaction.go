package transaction

import (
	"panda/chain"
	//"panda/chain/stellar"
	"panda/chain/ont"

	"github.com/astaxie/beego"
)

var (
	Trans_Success = 1
	Trans_Failed  = -1
	Trans_Retry   = 0
)
var trans chain.ChainOp

//初始化实例，根据需求定制实例对象
func init() {
	//trans = stellar.NewOpChain(beego.BConfig.RunMode)
	trans, _ = ont.NewOpChain("./wallet.dat", beego.AppConfig.String("chain_rpc"))
}

func Genkey(lable string) (pub, priv string, err error) {
	return trans.GenKeyPair(lable)
}

func DoTransaction(lable, sPrivkey, dPublic, amount string) (txhash string, err error) {
	return trans.DoTransaction(lable, sPrivkey, dPublic, amount)
}

func QueryTransaction(txhash string) (int, error) {
	return trans.QueryTransaction(txhash)
}

func GetBalance(address string) (balance string, err error) {
	return trans.GetBalance(address)
}

func ValidatePublicKey(address string) (err error) {
	return trans.ValidatePublicKey(address)
}
