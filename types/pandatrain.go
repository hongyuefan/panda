package types

type RspTrainResult struct {
	Success bool   `json:"success"`
	CTime   int64  `json:"catch_time"`
	Result  int    `json:"catch_result"`
	Message string `json:"message"`
}

type RspTrain struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Txhash  string `json:"txhash"`
}
