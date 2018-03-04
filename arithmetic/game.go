package arithmetic

import (
	"strings"
	"time"
)

const (
	Result_Multiple = 1.5
	Result_Normal   = 1
)

func CatchResult(hash string, comp string) (result int, err error) {

	comp = strings.Replace(comp, ".", "", -1)

	h1, err := SplitTx_Trim_N(hash, 1)
	if err != nil {
		return
	}
	h2, err := SplitTx_Trim_N(hash, 2)
	if err != nil {
		return
	}

	c1, err := SplitTx_N(comp, 1)
	if err != nil {
		return
	}
	c2, err := SplitTx_N(comp, 2)
	if err != nil {
		return
	}

	if c1 > h2 || (c1 == h2 && c2 > h1) {
		result = 1
	}
	if c1 == h2 && c2 == h1 {
		result = 2
	}
	return
}

func income(b byte) float32 {
	if int(b) == int(time.Now().Weekday()) {
		return Result_Multiple
	}
	return Result_Normal
}
