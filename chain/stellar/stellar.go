package stellar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	n "github.com/stellar/go/network"
	"github.com/stellar/go/xdr"
)

type OpChain struct {
	netWork string
	provide *horizon.Client
}

type RspChain struct {
	ResultXdr string `json:"result_xdr"`
	Ledger    int64  `json:"ledger"`
}

type Balance struct {
	Amount    string `json:"balance"`
	AssetType string `json:"asset_type"`
}
type Balances struct {
	ArryBalance []*Balance `json:"balances"`
}

func NewOpChain(netWork string) *OpChain {
	switch netWork {
	case "dev":
		return &OpChain{
			netWork: n.TestNetworkPassphrase,
			provide: horizon.DefaultTestNetClient,
		}
	case "prod":
		return &OpChain{
			netWork: n.PublicNetworkPassphrase,
			provide: horizon.DefaultPublicNetClient,
		}
	}
	return &OpChain{}
}

func (o *OpChain) GenKeyPair() (pub, priv string, err error) {

	pair, err := keypair.Random()
	if err != nil {
		return
	}
	pub = pair.Address()
	priv = pair.Seed()
	return
}

func (o *OpChain) GetBalance(address string) (balance string, err error) {

	var balances Balances

	url := o.provide.URL + "/accounts/" + address

	rsp, err := http.Get(url)
	if err != nil {
		return
	}
	if rsp.StatusCode != http.StatusOK {
		err = fmt.Errorf("query account information %v error code %v", address, rsp.StatusCode)
		return
	}
	buffer := new(bytes.Buffer)

	if _, err = io.Copy(buffer, rsp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(buffer.Bytes(), &balances); err != nil {
		return
	}

	for _, b := range balances.ArryBalance {
		if b.AssetType == "native" {
			balance = b.Amount
			break
		}
	}
	return
}

func (o *OpChain) DoTransaction(sPrivkey, dPublic, amount string) (txhash string, err error) {
	tx, err := b.Transaction(
		b.Network{o.netWork},
		b.SourceAccount{sPrivkey},
		b.AutoSequence{o.provide},
		b.Payment(
			b.Destination{dPublic},
			b.NativeAmount{amount},
		),
	)

	if err != nil {
		return
	}

	txe, err := tx.Sign(sPrivkey)
	if err != nil {
		return
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return
	}
	rsp, err := o.provide.SubmitTransaction(txeB64)
	if err != nil {
		return
	}

	txhash = rsp.Hash
	return
}

func (o *OpChain) parseResultXDR(data string) (ok bool, err error) {

	rawr := strings.NewReader(data)

	b64r := base64.NewDecoder(base64.StdEncoding, rawr)

	var tx xdr.TransactionResult

	_, err = xdr.Unmarshal(b64r, &tx)
	if err != nil {
		return
	}

	if tx.Result.Code == xdr.TransactionResultCodeTxSuccess {
		ok = true
	}
	err = fmt.Errorf("transaction is failed")
	return
}

/*result: 0 网络错误或hash尚未完成，下次重试
  result：1  交易成功
  result: -1 交易失败
*/
func (o *OpChain) QueryTransaction(txhash string) (result int, err error) {

	var (
		rspChain RspChain
		buffer   *bytes.Buffer
		ok       bool
	)
	url := o.provide.URL + "/transactions/" + txhash

	rsp, err := http.Get(url)
	if err != nil {
		return
	}
	if rsp.StatusCode != http.StatusOK {
		err = fmt.Errorf("query trans txhash %v error code %v", txhash, rsp.StatusCode)
		return
	}

	buffer = new(bytes.Buffer)

	if _, err = io.Copy(buffer, rsp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(buffer.Bytes(), &rspChain); err != nil {
		return
	}

	if ok, err = o.parseResultXDR(rspChain.ResultXdr); !ok {
		return -1, err
	}
	return 1, nil
}
