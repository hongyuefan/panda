package controllers

import (
	"testing"
)

func TestVCodeGenerate(t *testing.T) {

	id, png := VCodeGenerate(60, 240, Cap_Metic_Mod)

	t.Log(png)

	t.Log(VCodeValidate(id, "8"))

	return
}

func TestParseToken(t *testing.T) {
	token, _ := ParseToken("Bearer ddfsdfadf")
	t.Log(token)
}

//func TestSendEmail(t *testing.T) {

//	err := SendEmail("1027350999@qq.com", "太上熊猫平台", "1027350999@qq.com", "", "太上熊猫平台", "大是大非似懂非懂")
//	if err != nil {
//		t.Log(err)
//	}
//	return
//}
