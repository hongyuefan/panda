package transaction

import (
	"panda/chain"
	"panda/chain/stellar"

	"github.com/astaxie/beego"
)

var trans chain.ChainOp

//初始化实例，根据需求定制实例对象
func init() {
	trans = stellar.NewOpChain(beego.BConfig.RunMode)
}

func Genkey() (pub, priv string, err error) {
	return trans.GenKeyPair()
}

func DoTransaction(sPrivkey, dPublic, amount string) (txhash string, err error) {
	return trans.DoTransaction(sPrivkey, dPublic, amount)
}

func QueryTransaction(txhash string) (int64, error) {
	return trans.QueryTransaction(txhash)
}
