package arithmetic

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
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
func SplitTx_Trim_N(tx string, n int) (b int, err error) {

	hTx := strings.TrimPrefix(strings.ToLower(tx), "0x")

	return HexToI(hTx[len(hTx)-n : len(hTx)-n+1])

}

func HexToI(a string) (b int, err error) {

	var ok bool

	m := make(map[string]int, 16)

	m["a"] = 10
	m["b"] = 11
	m["c"] = 12
	m["d"] = 13
	m["e"] = 14
	m["f"] = 15

	b, err = strconv.Atoi(a)
	if err != nil {
		err = nil
		if b, ok = m[a]; !ok {
			err = fmt.Errorf("hex not right")
		}
	}
	return
}

//获取第n个数字
func SplitTx_N(tx string, n int) (b int, err error) {
	hTx := strings.TrimPrefix(strings.ToLower(tx), "0x")
	if n > len(hTx) {
		err = fmt.Errorf("n longer than tx")
		return
	}
	return HexToI(hTx[n-1 : n])
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
