package controllers

import (
	"fmt"
	"panda/arithmetic"
	"panda/models"
	"panda/sendmsg"
	"panda/types"
	"strings"

	"github.com/astaxie/beego"
)

type MsgController struct {
	beego.Controller
}

func (c *MsgController) ValidateMsgCode() {
	var (
		rspEmail types.RspEmail
	)
	mobile := c.GetString("mobile")
	mobile_code := c.GetString("code")

	if validMobileCode(mobile_code, getSessionString(c.GetSession(mobile))) {
		goto errDeal
	}
	rspEmail = types.RspEmail{
		Success: true,
		Message: "validate code success",
	}

	c.Ctx.Output.JSON(rspEmail, false, false)
	return

errDeal:
	rspEmail = types.RspEmail{
		Success: false,
		Message: "validate code failed",
	}
	c.Ctx.Output.JSON(rspEmail, false, false)
	return
}

func (c *MsgController) SendMsgCode() {

	var (
		err      error
		code     string
		rspEmail types.RspEmail
		conf     models.Config
		params   []string
	)
	codeId := c.GetString("codeId")
	verifyValue := c.GetString("verifyValue")
	mobile := c.GetString("mobile")
	nation := c.GetString("nation")
	if nation == "" {
		nation = "86"
	}

	if err = ValidMobile(mobile); err != nil {
		goto errDeal
	}

	if !VCodeValidate(codeId, verifyValue) {
		err = fmt.Errorf("picture validate failed")
		goto errDeal
	}

	code = arithmetic.GetRandLimit(Email_Code_Len)

	params = []string{code}

	conf = GetConfigData()

	if err = sendmsg.SendMsg(conf.AppId, conf.AppKey, nation, mobile, params, conf.TplId); err != nil {
		goto errDeal
	}

	c.SetSession(mobile, code)

	rspEmail = types.RspEmail{
		Success: true,
		Message: "send success",
	}
	c.Ctx.Output.JSON(rspEmail, false, false)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return

}

func validMobileCode(email_code, validcode string) bool {

	if len(email_code) != Email_Code_Len {
		return false
	}

	if !strings.EqualFold(validcode, email_code) {
		return false
	}
	return true
}
