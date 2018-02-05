package arithmetic

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
	"time"
)

//获取（min，max）之间大小的随机数
func GetRand(min float64, max float64) (result float64) {
	source := rand.NewSource(time.Now().Unix())
	nRand := rand.New(source)
	return nRand.Float64()*(max-min) + min
}

//截取tx倒数n个byte
func SplitTx(tx string, n int) (arry []byte, err error) {
	hTx := strings.TrimPrefix(strings.ToLower(tx), "0x")
	b, err := hex.DecodeString(hTx)
	if err != nil {
		return
	}
	arry = b[len(b)-n:]
	return
}

//b对m取模
func GetMod(b byte, m byte) (byte, error) {
	if m <= 0 {
		err := errors.New("cant mod zero")
		return 0, err
	}
	return b % m, nil
}

//翻转byte数组
func OverturnArray(arry []byte) (arry_result []byte) {

	nLen := len(arry)

	arry_result = make([]byte, nLen)

	for i, b := range arry {
		arry_result[nLen-i-1] = b
	}
	return
}
