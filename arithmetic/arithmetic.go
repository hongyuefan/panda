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

func ParseLimit(strLimit string) (min float64, max float64, err error) {

	arrysLimit := strings.Split(strLimit, "_")

	if len(arrysLimit) != 2 {
		err = fmt.Errorf("conf attrivalue limit set error")
		return
	}
	if min, err = strconv.ParseFloat(arrysLimit[0], 64); err != nil {
		return
	}
	if max, err = strconv.ParseFloat(arrysLimit[1], 64); err != nil {
		return
	}
	return

}

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

func SplitTx_Trim_N_S(tx string) string {
	return tx[len(tx)-6:]
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

func GetChar_Num() (c string) {
	return string(byte(GetRand(48, 58)))
}
func GetChar_Cap() (c string) {
	return string(byte(GetRand(65, 91)))
}
func GetChar_Low() (c string) {
	return string(byte(GetRand(97, 123)))
}

func GenCode(n int) (code string) {
	var (
		rand int
		str  string
	)
	for i := 0; i < n; i++ {
		rand = int(GetRand(0, 3))
		time.Sleep(time.Nanosecond)
		switch rand {
		case 0:
			str += GetChar_Num()
			continue
		case 1:
			str += GetChar_Low()
			continue
		case 2:
			str += GetChar_Cap()
			continue
		}
	}
	return str
}

//解析 1_2 格式
func SplitAndParseToFloat(s string) (small, big float64, err error) {

	arry := strings.Split(s, "_")

	if len(arry) != 2 {
		err = fmt.Errorf("SplitAndParseToFloat params not right")
		return
	}

	if small, err = strconv.ParseFloat(arry[0], 64); err != nil {
		return
	}
	if big, err = strconv.ParseFloat(arry[1], 64); err != nil {
		return
	}
	return
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

//求x的n次方
func Powerf(x float64, n int) float64 {

	ans := 1.0

	for n != 0 {
		ans *= x
		n--
	}
	return ans
}

//秒装换为时分秒
func ParseSecond(second int64) (hour, min, sec int64) {
	sec = second % 60
	min = second / 60
	if min >= 60 {
		hour = min / 60
		min = min % 60
	}
	return
}

//取小数点后两位
func ParseFloat(p float64) (result string) {

	i_p := int64(p * 100)

	f_p := float64(i_p) / float64(100)

	return fmt.Sprintf("%v", f_p)
}
