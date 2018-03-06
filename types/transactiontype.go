package types

import (
	"fmt"
)

type RspTransQ struct {
	Total int          `json:"total"`
	Data  []TransQData `json:"data"`
}

var (
	Attr_Type_Minjie  int64 = 1
	Attr_Type_Liliang int64 = 2
	Attr_Type_Zhili   int64 = 3
)

var (
	Train_Multi_Ratio = 10
)

var (
	Trans_Type_WithDrawal int64 = 1
	Trans_Type_Catch      int64 = 2
	Trans_Type_Train      int64 = 3
	Trans_Type_Bonus      int64 = 4
)

var (
	Trans_Status_Waiting = 0
	Trans_Status_Failed  = -1
	Trans_Status_Success = 1
)

var (
	Error_Trans_AmountOver    = fmt.Errorf("转账数值大于账户总额")
	Error_Trans_MisType       = fmt.Errorf("转账类型不匹配")
	Error_Trans_CatchIntervel = fmt.Errorf("捕捉冷却时间不足")

	Error_Train_AmountOver = fmt.Errorf("训练额度大于当日最大限度")
)

type TransQData struct {
	TxId   string `json:"txid"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
	TxHash string `json:"txhash"`
	Time   string `json:"time"`
	Status string `json:"status"`
}
