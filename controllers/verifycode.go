package controllers

import (
	"panda/types"

	"github.com/astaxie/beego"
	cap "github.com/mojocn/base64Captcha"
)

const (
	Cap_Num_Mod     = cap.CaptchaModeNumber
	Cap_Char_Mod    = cap.CaptchaModeAlphabet
	Cap_Metic_Mod   = cap.CaptchaModeArithmetic
	Cap_NumChar_Mod = cap.CaptchaModeNumberAlphabet
)

type VerifyController struct {
	beego.Controller
}

func (c *VerifyController) GenerateCode() {
	var (
		rspCode types.RspGenCode
	)
	rspCode.CodeId, rspCode.PngBase64 = VCodeGenerate(60, 240, Cap_Metic_Mod)

	c.Ctx.Output.JSON(rspCode, false, false)

	return
}

func (c *VerifyController) ValidateCode() {

	var (
		rspVerify types.RspVerify
		err       error
	)

	codeId := c.GetString("codeId")
	value := c.GetString("verifyValue")

	if len(codeId) == 0 || len(value) == 0 {
		err = types.Error_Params_Empty
		goto errDeal
	}

	if !VCodeValidate(codeId, value) {
		rspVerify = types.RspVerify{
			Success: false,
			Message: types.VERIFY_VALID_FAILED,
		}
	} else {
		rspVerify = types.RspVerify{
			Success: true,
			Message: types.VERIFY_VALID_SUCCESS,
		}
	}
	c.Ctx.Output.JSON(rspVerify, false, false)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func VCodeGenerate(heigt, width, mode int) (capId, pngBase64 string) {

	configChar := cap.ConfigCharacter{
		Height:             heigt,
		Width:              width,
		Mode:               mode,
		ComplexOfNoiseText: cap.CaptchaComplexLower,
		ComplexOfNoiseDot:  cap.CaptchaComplexLower,
		IsUseSimpleFont:    false,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}

	capId, charCap := cap.GenerateCaptcha("", configChar)

	pngBase64 = cap.CaptchaWriteToBase64Encoding(charCap)

	return
}

func VCodeGenerateByCapId(heigt, width, mode int, capId string) (pngBase64 string) {

	configChar := cap.ConfigCharacter{
		Height:             heigt,
		Width:              width,
		Mode:               mode,
		ComplexOfNoiseText: cap.CaptchaComplexLower,
		ComplexOfNoiseDot:  cap.CaptchaComplexLower,
		IsUseSimpleFont:    false,
		IsShowHollowLine:   true,
		IsShowNoiseDot:     true,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    true,
		IsShowSineLine:     true,
		CaptchaLen:         4,
	}

	_, charCap := cap.GenerateCaptcha(capId, configChar)

	pngBase64 = cap.CaptchaWriteToBase64Encoding(charCap)

	return
}

func VCodeValidate(identifier, verifyValue string) bool {
	return cap.VerifyCaptcha(identifier, verifyValue)
}
