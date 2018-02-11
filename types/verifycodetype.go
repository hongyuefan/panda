package types

type RspGenCode struct {
	CodeId    string `json:"codeId"`
	PngBase64 string `json:"pngBase64"`
}

type ReqVerify struct {
	CodeId      string `json:"codeId"`
	VerifyValue string `json:"verifyValue"`
}

type RspVerify struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
