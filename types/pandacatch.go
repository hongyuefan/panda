package types

type RspCatchResult struct {
	CTime  int64 `json:"catch_time"`
	Result int   `json:"catch_result"`
}

type RspCatch struct {
	Txhash string `json:"txhash"`
}
