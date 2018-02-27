package arithmetic

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//获取（min，max）之间大小的随机数
func GetRand(min float64, max float64) (result float64) {
	source := rand.NewSource(time.Now().UnixNano())
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

//获取tx倒数第n个数字
func SplitTx_Trim_N(tx string, n int) (b byte, err error) {
	hTx := strings.TrimPrefix(strings.ToLower(tx), "0x")
	arry, err := hex.DecodeString(hTx)
	if err != nil {
		return
	}
	b = arry[len(arry)-n] % 0x0f
	return
}

//获取第n个数字
func SplitTx_N(value string, n int) (b byte, err error) {
	hTx := strings.TrimPrefix(strings.ToLower(tx), "0x")
	arry, err := hex.DecodeString(hTx)
	if err != nil {
		return
	}
	b = arry[len(arry)-n] % 0x0f
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

//获取位数为count的随机数据字符串
func GetRandLimit(count int) (result string) {

	for i := 0; i < count; i++ {

		srand := rand.Intn(9)

		result += fmt.Sprintf("%v", int(srand))
	}
	return
}
