package types

type RspTrainResult struct {
	CTime  int64 `json:"catch_time"`
	Result int   `json:"catch_result"`
}

type RspTrain struct {
	Txhash string `json:"txhash"`
}
