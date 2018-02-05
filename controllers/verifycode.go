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
	rspCode.CodeId, rspCode.PngBase64 = VCodeGenerate(60, 240, Cap_NumChar_Mod)

	c.Ctx.Output.JSON(rspCode, false, false)

	return
}

func VCodeGenerate(heigt, width, mode int) (capId, pngBase64 string) {

	configChar := cap.ConfigCharacter{
		Height:             heigt,
		Width:              width,
		Mode:               mode,
		ComplexOfNoiseText: cap.CaptchaComplexLower,
		ComplexOfNoiseDot:  cap.CaptchaComplexLower,
		IsUseSimpleFont:    true,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
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
		IsUseSimpleFont:    true,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    true,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}

	_, charCap := cap.GenerateCaptcha(capId, configChar)

	pngBase64 = cap.CaptchaWriteToBase64Encoding(charCap)

	return
}

func VCodeValidate(identifier, verifyValue string) bool {
	return cap.VerifyCaptcha(identifier, verifyValue)
}
