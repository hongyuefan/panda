package arithmetic

import (
	"time"
)

const (
	Result_Multiple = 1.5
	Result_Normal   = 1
)

func Rule(TrimN int, hash string) (result float32, err error) {

	var b byte

	if b, err = SplitTx_Trim_N(hash, TrimN); err != nil {
		return 0, err
	}

	return income(b), nil
}

func income(b byte) float32 {
	if int(b) == int(time.Now().Weekday()) {
		return Result_Multiple
	}
	return Result_Normal
}
