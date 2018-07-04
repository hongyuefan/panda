package ont

import (
	"fmt"
	"panda/arithmetic"
	"strconv"

	"github.com/ontio/ontology-crypto/keypair"
	s "github.com/ontio/ontology-crypto/signature"
	osdk "github.com/ontio/ontology-go-sdk"
	rpc "github.com/ontio/ontology-go-sdk/rpc"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
)

type OpChain struct {
	ontSDK  *osdk.OntologySdk
	aClient account.Client
}

func NewOpChain(walletFile, rpcAddress string) (opC *OpChain, err error) {
	opC = &OpChain{
		ontSDK: osdk.NewOntologySdk(),
	}
	if opC.aClient, err = opC.ontSDK.OpenOrCreateWallet(walletFile); err != nil {
		return nil, err
	}
	opC.ontSDK.Rpc = rpc.NewRpcClient()

	opC.ontSDK.Rpc.SetAddress(rpcAddress)

	return opC, nil
}

func (o *OpChain) GenKeyPair(lable string) (pub, priv string, err error) {

	priv = arithmetic.GenCode(6)

	accMeta := o.aClient.GetAccountMetadataByLabel(lable)
	if accMeta == nil {
		_, err = o.aClient.NewAccount(lable, keypair.PK_ECDSA, keypair.P256, s.SHA256withECDSA, []byte(priv))
		if err != nil {
			err = fmt.Errorf("NewAccount error:%s", err)
			return
		}
	}
	account, err := o.aClient.GetAccountByLabel(lable, []byte(priv))
	if err != nil {
		return
	}
	pub = account.Address.ToBase58()
	return
}

func (o *OpChain) GetBalance(address string) (balance string, err error) {

	public, err := common.AddressFromBase58(address)
	if err != nil {
		return
	}
	balan, err := o.ontSDK.Rpc.GetBalance(public)
	if err != nil {
		return
	}
	balance = fmt.Sprintf("%v", balan.Ont)
	return
}

func (o *OpChain) DoTransaction(lable, sPrivkey, dPublic, amount string) (txhash string, err error) {

	nAmount, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return
	}

	account, err := o.aClient.GetAccountByLabel(lable, []byte(sPrivkey))
	if err != nil {
		return
	}
	destAddr, err := common.AddressFromBase58(dPublic)
	if err != nil {
		return
	}

	hash, err := o.ontSDK.Rpc.TransferFrom(0, 20000, "ont", account, account.Address, destAddr, nAmount)
	if err != nil {
		return
	}
	txhash = hash.ToHexString()
	return
}

/*result: 0 网络错误或hash尚未完成，下次重试
  result：1  交易成功
  result: -1 交易失败
*/
func (o *OpChain) QueryTransaction(txhash string) (result int, err error) {

	hash, err := common.Uint256FromHexString(txhash)
	if err != nil {
		return -1, err
	}
	trans, err := o.ontSDK.Rpc.GetRawTransaction(hash)
	if err != nil {
		return 0, err
	}
	if trans != nil {
		return 1, nil
	} else {
		return -1, nil
	}
}

func (o *OpChain) ValidatePublicKey(publickey string) (err error) {
	return
}
