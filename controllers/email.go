package controllers

import (
	"panda/arithmetic"
	"panda/types"

	"github.com/astaxie/beego"
	"github.com/go-gomail/gomail"
)

var (
	Email_Sender  = "1027350999@qq.com"
	Email_Name    = "太上熊猫"
	Email_Pass    = "bznzexsyprsibegf"
	Email_Stmp    = "smtp.qq.com"
	Email_Port    = 465
	Email_Sub     = "太上熊猫验证码"
	Email_Content = "您的本次邮箱验证码为："
)

type EmailController struct {
	beego.Controller
}

func (c *EmailController) SendEmailCode() {

	var (
		err      error
		code     string
		rspEmail types.RspEmail
	)
	email := c.GetString("email")

	if err = ValidateUserName(email); err != nil {
		goto errDeal
	}

	code = arithmetic.GetRandLimit(4)

	if err = SendEmail(Email_Sender, Email_Name, email, "", Email_Sub, Email_Content+code); err != nil {
		goto errDeal
	}

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

func SendEmail(from, from_name, to, to_name, subject, content string) (err error) {

	m := gomail.NewMessage()

	m.SetAddressHeader("From", from, from_name) // 发件人
	m.SetHeader("To",                           // 收件人
		m.FormatAddress(to, to_name),
	)
	m.SetHeader("Subject", subject) // 主题
	m.SetBody("text/html", content) // 正文

	d := gomail.NewPlainDialer(Email_Stmp, Email_Port, Email_Sender, Email_Pass) // 发送邮件服务器、端口、发件人账号、发件人密码

	return d.DialAndSend(m)
}
