package types

import (
	"encoding/json"
	"fmt"
)

var (
	Error_Player_Balance = fmt.Errorf("账户余额不足")
)

var (
	Svg_File_Path = "/svg"
)

type RspCommon struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}
