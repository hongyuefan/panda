package types

var (
	Type_Train = 0x01
	Type_Catch = 0x02
	Type_Bonus = 0x03
)

type RspTransQ struct {
	Total int          `json:"total"`
	Data  []TransQData `json:"data"`
}

type TransQData struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Amount string `json:"amount"`
	TxHash string `json:"txhash"`
	Time   string `json:"time"`
	Status string `json:"status"`
}
