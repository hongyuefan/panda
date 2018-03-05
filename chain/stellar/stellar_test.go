package stellar

import (
	"testing"
	"time"
)

func TestQueryTransaction(t *testing.T) {

	opc := NewOpChain("dev")

	from := "SBJ7OP6ET27OKU4PZK5D2ABCVJJWBHT3RUXMEKUY5W6PLV22V6KKI2SI"
	to, _, _ := opc.GenKeyPair()
	to = "GCMBTRU2BBGCHOXBERPY2U42VRG2FLBMV3423HGAL4NS3XI2CM2DMIGS"

	txhash, err := opc.DoTransaction(from, to, "1")
	if err != nil {
		t.Log(err)
		return
	}

	t.Log("txhash:", txhash)

	time.Sleep(3 * time.Second)

	leger, err := opc.QueryTransaction(txhash)

	if err != nil {
		t.Log(err)
		return
	}
	t.Log(leger)

	balance, err := opc.GetBalance("GBD3XKDFVKSRFSL7YXYVE6EVUX7ZD2X3FPJQUDEVJM7D7BXUNH3QJOQX")

	if err != nil {
		t.Log(err)
		return
	}
	t.Log(balance)

	return
}
