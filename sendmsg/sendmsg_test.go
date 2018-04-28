package sendmsg

import (
	"encoding/json"
	"testing"
)

type RspCommon struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type Data struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

func TestSigMsg(t *testing.T) {

	t.Log(SigMsg("15600199768", "8bd9c73577e97decf4772defa48d1066", "2356", "1678945"))

}

func TestSendMsg(t *testing.T) {

	if err := SendMsg("1400081071", "8bd9c73577e97decf4772defa48d1066", "86", "15600199768", []string{"2345"}, 102918); err != nil {
		t.Log(err)
	}

}

func TestJson(t *testing.T) {

	d := Data{
		Account: "1234",
		Balance: "333",
	}

	b, _ := json.Marshal(d)

	r := RspCommon{
		Success: true,
		Message: "success",
		Data:    b,
	}

	m, _ := json.Marshal(r)

	t.Log(string(m))
}
