package types

type RspGenCode struct {
	CodeId    string `json:"codeId"`
	PngBase64 string `json:"pngBase64"`
}

type ReqVerify struct {
	CodeId      string
	VerifyValue string
}
