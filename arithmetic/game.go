package arithmetic

import (
	"strconv"
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

func TrainResult(hash string, N int) (statue int, err error) {

	p, err := SplitTx_Trim_N(hash, N)
	if err != nil {
		return
	}
	weekDay := int(time.Now().Weekday())

	if p == weekDay {
		statue = 2
		return
	}
	if p >= 0 && p <= 6 {
		statue = 1
		return
	}
	return
}

func CatchCold(ratio float64, catchTime float64, ctimes int, years float64, zhili string) (coldTm float64, err error) {

	fzl, err := strconv.ParseFloat(zhili, 64)
	if err != nil {
		return 0, err
	}

	coldTm = catchTime * Powerf(ratio, ctimes) * ((1 + years) * 7 / 10) / (1 + (fzl*fzl)/10)

	return
}
