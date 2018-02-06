package types

type ReqEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type RspEmail struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
