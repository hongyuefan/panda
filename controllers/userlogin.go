package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"panda/models"
	t "panda/transaction"
	"panda/types"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	UserType_Normal = "2"
	UserType_Frozen = "4"

	User_Exist     = 1
	User_Not_Exist = 0
)

type UserLoginController struct {
	beego.Controller
}

func (c *UserLoginController) VerifyUser() {

	var (
		reqVerify types.ReqVerifyCode
		rspVerify types.RspVerifyCode
		mUser     models.Player
		orm       *models.Common
		err       error
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	reqVerify.UserName = c.Ctx.Request.FormValue("userName")
	reqVerify.TimeStamp = c.Ctx.Request.FormValue("timeStamp")

	mUser.Mobile = reqVerify.UserName

	orm = models.NewCommon()

	if err = orm.CommonGetOne(&mUser, "Mobile"); err != nil {
		goto errDeal
	}

	if mUser.Id != 0 {
		rspVerify = types.RspVerifyCode{
			types.RspBase{
				MemberIsExist: User_Exist,
			},
		}

	} else {
		rspVerify = types.RspVerifyCode{
			types.RspBase{
				MemberIsExist: User_Not_Exist,
			},
		}
	}

	SuccessHandler(c.Ctx, rspVerify)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return

}

func (c *UserLoginController) RegistUserByMobile() {
	var (
		reqRgt          types.ReqRegist
		rspRgt          types.RspRegist
		invitationCode  string
		mUser           *models.Player
		orm             *models.Common
		conf            models.Config
		public, privkey string
		token           string
		uid             int64
		err             error
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	reqRgt.NickName = c.Ctx.Request.FormValue("nickName")
	reqRgt.Password = c.Ctx.Request.FormValue("passWord")
	reqRgt.TimeStamp = c.Ctx.Request.FormValue("timeStamp")
	reqRgt.UserName = c.Ctx.Request.FormValue("userName")
	reqRgt.VerifyCode = c.Ctx.Request.FormValue("verifyCode")
	invitationCode = c.Ctx.Request.FormValue("code")

	if err = ValidMobile(reqRgt.UserName); err != nil {
		goto errDeal
	}
	if err = ValidatePassWord(reqRgt.Password); err != nil {
		goto errDeal
	}

	if len(reqRgt.VerifyCode) != Email_Code_Len || !validMobileCode(reqRgt.VerifyCode, getSessionString(c.GetSession(reqRgt.UserName))) {
		err = errors.New("手机验证码错误")
		goto errDeal
	}

	mUser = &models.Player{
		Mobile: reqRgt.UserName,
	}

	orm = models.NewCommon()

	if err = orm.CommonGetOne(mUser, "Mobile"); err != nil {
		goto errDeal
	}
	if mUser.Id != 0 {
		err = errors.New("手机号已经存在，请登录！")
		goto errDeal
	}

	if public, privkey, err = t.Genkey(); err != nil {
		goto errDeal
	}

	conf = GetConfigData()

	mUser = &models.Player{
		Nickname:    reqRgt.NickName,
		Mobile:      reqRgt.UserName,
		Password:    reqRgt.Password,
		Createtime:  time.Now().Unix(),
		PubPublic:   public,
		PubPrivkey:  privkey,
		GamblingNum: types.Gambling_Num_Default,
		Avatar:      conf.HostUrl + types.Pic_File_Path + "/Default.jpg",
	}

	if uid, err = orm.CommonInsert(mUser); err != nil {
		err = fmt.Errorf("用户注册失败，昵称 %v 已存在！", mUser.Nickname)
		goto errDeal
	}

	if token, err = TokenGenerate(uid); err != nil {
		goto errDeal
	}

	rspRgt = types.RspRegist{
		RspBase: types.RspBase{
			MemberIsExist: 0,
		},
		Data: types.User{
			MemberId:      fmt.Sprintf("%v", uid),
			NickName:      reqRgt.NickName,
			UserName:      reqRgt.UserName,
			UserType:      UserType_Normal,
			Token:         token,
			Avatar:        mUser.Avatar,
			Balance:       "0",
			Freeze:        "0",
			WalletAddress: "",
			Mypets:        "",
		},
	}

	if invitationCode != "" {
		models.UpdateInvitationCount(invitationCode)
	}
	SuccessHandler(c.Ctx, rspRgt)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) RegistUser() {
	var (
		reqRgt          types.ReqRegist
		rspRgt          types.RspRegist
		invitationCode  string
		mUser           *models.Player
		orm             *models.Common
		public, privkey string
		token           string
		uid             int64
		err             error
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	reqRgt.NickName = c.Ctx.Request.FormValue("nickName")
	reqRgt.Password = c.Ctx.Request.FormValue("passWord")
	reqRgt.TimeStamp = c.Ctx.Request.FormValue("timeStamp")
	reqRgt.UserName = c.Ctx.Request.FormValue("userName")
	reqRgt.VerifyCode = c.Ctx.Request.FormValue("verifyCode")
	invitationCode = c.Ctx.Request.FormValue("code")

	if err = ValidateEmail(reqRgt.UserName); err != nil {
		goto errDeal
	}
	if err = ValidatePassWord(reqRgt.Password); err != nil {
		goto errDeal
	}

	if len(reqRgt.VerifyCode) != Email_Code_Len || !validEmailCode(reqRgt.VerifyCode, getSessionString(c.GetSession(reqRgt.UserName))) {
		err = errors.New("email verify code not right")
		goto errDeal
	}

	if public, privkey, err = t.Genkey(); err != nil {
		goto errDeal
	}

	mUser = &models.Player{
		Nickname:    reqRgt.NickName,
		Email:       reqRgt.UserName,
		Password:    reqRgt.Password,
		Createtime:  time.Now().Unix(),
		PubPublic:   public,
		PubPrivkey:  privkey,
		GamblingNum: types.Gambling_Num_Default,
	}
	orm = models.NewCommon()

	//	if err = orm.BeginTx(); err != nil {
	//		goto errDeal
	//	}

	if uid, err = orm.CommonInsert(mUser); err != nil {
		goto errDeal
	}

	if token, err = TokenGenerate(uid); err != nil {
		//orm.Rollback()
		goto errDeal
	}

	//	if err = orm.Commit(); err != nil {
	//		goto errDeal
	//	}
	rspRgt = types.RspRegist{
		RspBase: types.RspBase{
			MemberIsExist: 0,
		},
		Data: types.User{
			MemberId:      fmt.Sprintf("%v", uid),
			NickName:      reqRgt.NickName,
			UserName:      reqRgt.UserName,
			UserType:      UserType_Normal,
			Token:         token,
			Avatar:        "",
			Balance:       "0",
			Freeze:        "0",
			WalletAddress: "",
			Mypets:        "",
		},
	}

	if invitationCode != "" {
		models.UpdateInvitationCount(invitationCode)
	}
	SuccessHandler(c.Ctx, rspRgt)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) GetNewToken() {

	var (
		uId   int64
		err   error
		token string
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	if uId, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if token, err = TokenGenerate(uId); err != nil {
		goto errDeal
	}
	SuccessHandler(c.Ctx, token)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) ModifyNickName() {

	var (
		nickName string
		mUser    models.Player
		err      error
		orm      *models.Common
	)

	type Result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	nickName = c.Ctx.Request.FormValue("nickname")

	if mUser.Id, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	orm = models.NewCommon()

	mUser.Nickname = nickName

	if _, err = orm.CommonUpdateById(&mUser, "nickname"); err != nil {
		goto errDeal
	}

	c.Ctx.Output.JSON(Result{
		Success: true,
		Message: "update nickname success",
	}, false, false)

	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) UserLoginByMobile() {
	var (
		reqLogin       types.ReqLogin
		rspLogin       types.RspRegist
		mUser          *models.Player
		orm            *models.Common
		token, balance string
		err            error
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	reqLogin.UserName = c.Ctx.Request.FormValue("userName")
	reqLogin.PassWord = c.Ctx.Request.FormValue("passWord")

	mUser = &models.Player{
		Mobile:   reqLogin.UserName,
		Password: reqLogin.PassWord,
	}

	orm = models.NewCommon()

	if err = orm.CommonGetOne(mUser, "Mobile", "PassWord"); err != nil || mUser.Id <= 0 {
		err = fmt.Errorf("用户名或密码错误！")
		goto errDeal
	}

	if mUser.Id > 0 {
		if token, err = TokenGenerate(mUser.Id); err != nil {
			goto errDeal
		}

		if balance, err = t.GetBalance(mUser.PubPublic); err != nil {
			balance = "0"
		}

		rspLogin = types.RspRegist{
			RspBase: types.RspBase{
				MemberIsExist: User_Exist,
			},
			Data: types.User{
				MemberId:      fmt.Sprintf("%v", mUser.Id),
				NickName:      mUser.Nickname,
				UserName:      mUser.Mobile,
				UserType:      UserType_Normal,
				Token:         token,
				Avatar:        mUser.Avatar,
				Balance:       balance,
				Freeze:        fmt.Sprintf("%v", mUser.Isdel),
				WalletAddress: mUser.Pubkey,
				Mypets:        "",
			},
		}
	}
	SuccessHandler(c.Ctx, rspLogin)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) UserLogin() {
	var (
		reqLogin       types.ReqLogin
		rspLogin       types.RspRegist
		mUser          *models.Player
		orm            *models.Common
		token, balance string
		err            error
	)

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	reqLogin.UserName = c.Ctx.Request.FormValue("userName")
	reqLogin.PassWord = c.Ctx.Request.FormValue("passWord")

	mUser = &models.Player{
		Email:    reqLogin.UserName,
		Password: reqLogin.PassWord,
	}

	orm = models.NewCommon()

	if err = orm.CommonGetOne(mUser, "Email", "PassWord"); err != nil {
		goto errDeal
	}

	if mUser.Id > 0 {
		if token, err = TokenGenerate(mUser.Id); err != nil {
			goto errDeal
		}

		if balance, err = t.GetBalance(mUser.PubPublic); err != nil {
			balance = "0"
		}

		rspLogin = types.RspRegist{
			RspBase: types.RspBase{
				MemberIsExist: User_Exist,
			},
			Data: types.User{
				MemberId:      fmt.Sprintf("%v", mUser.Id),
				NickName:      mUser.Nickname,
				UserName:      mUser.Email,
				UserType:      UserType_Normal,
				Token:         token,
				Avatar:        mUser.Avatar,
				Balance:       balance,
				Freeze:        fmt.Sprintf("%v", mUser.Isdel),
				WalletAddress: mUser.Pubkey,
				Mypets:        "",
			},
		}
	} else {
		rspLogin = types.RspRegist{
			RspBase: types.RspBase{
				MemberIsExist: User_Not_Exist,
			},
			Data: types.User{},
		}
	}
	SuccessHandler(c.Ctx, rspLogin)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) UploadPic() {
	type RspUpload struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	var (
		err         error
		base64Pic   string
		mUser       models.Player
		orm         *models.Common
		conf        models.Config
		strFileName string
		bytFile     []byte
	)

	subStr := "base64,"

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	if mUser.Id, err = ParseAndValidToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}

	conf = GetConfigData()

	base64Pic = c.Ctx.Request.FormValue("avatar")

	base64Pic = base64Pic[strings.Index(base64Pic, subStr)+len(subStr):]

	if len(base64Pic) > 210000 {
		err = fmt.Errorf("file is too large")
		goto errDeal
	}

	if bytFile, err = base64.StdEncoding.DecodeString(base64Pic); err != nil {
		goto errDeal
	}

	if err = WriteToFile(beego.AppConfig.String("pic_path")+fmt.Sprintf("%v.jpg", mUser.Id), bytFile); err != nil {
		goto errDeal
	}

	strFileName = conf.HostUrl + types.Pic_File_Path + "/" + fmt.Sprintf("%v.jpg", mUser.Id)

	orm = models.NewCommon()

	mUser.Avatar = strFileName

	if _, err = orm.CommonUpdateById(&mUser, "avatar"); err != nil {
		goto errDeal
	}

	c.Ctx.Output.JSON(RspUpload{
		Success: true,
		Message: "upload pic success",
	}, false, false)

	return
errDeal:
	c.Ctx.Output.JSON(RspUpload{
		Success: false,
		Message: err.Error(),
	}, false, false)
	return
}

func (c *UserLoginController) UpdatePassWord() {
	type RspPass struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	var (
		err         error
		newPassWord string
		email       string
		verifyCody  string
		user        models.Player
		com         *models.Common
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	email = c.Ctx.Request.FormValue("userName")
	newPassWord = c.Ctx.Request.FormValue("newPassword")
	verifyCody = c.Ctx.Request.FormValue("verifyCode")

	if err = ValidatePassWord(newPassWord); err != nil {
		goto errDeal
	}

	if !validEmailCode(verifyCody, getSessionString(c.GetSession(email))) {
		err = errors.New("email validate code not right")
		goto errDeal
	}

	user.Email = email

	com = models.NewCommon()

	if err = com.CommonGetOne(&user, "Email"); err != nil {
		goto errDeal
	}

	user.Password = newPassWord

	if _, err = com.CommonUpdateById(&user, "password"); err != nil {
		goto errDeal
	}

	c.Ctx.Output.JSON(RspPass{
		Success: true,
		Message: "update password success",
	}, false, false)

	return
errDeal:
	c.Ctx.Output.JSON(RspPass{
		Success: false,
		Message: err.Error(),
	}, false, false)
	return
}

func (c *UserLoginController) UpdatePassWordByMobile() {
	type RspPass struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	var (
		err         error
		newPassWord string
		mobile      string
		verifyCody  string
		user        models.Player
		com         *models.Common
	)
	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	mobile = c.Ctx.Request.FormValue("userName")
	newPassWord = c.Ctx.Request.FormValue("newPassword")
	verifyCody = c.Ctx.Request.FormValue("verifyCode")

	if err = ValidatePassWord(newPassWord); err != nil {
		goto errDeal
	}

	if !validMobileCode(verifyCody, getSessionString(c.GetSession(mobile))) {
		err = errors.New("短信验证码失败")
		goto errDeal
	}

	user.Mobile = mobile

	com = models.NewCommon()

	if err = com.CommonGetOne(&user, "Mobile"); err != nil {
		err = fmt.Errorf("用户不存在")
		goto errDeal
	}

	user.Password = newPassWord

	if _, err = com.CommonUpdateById(&user, "password"); err != nil {
		err = fmt.Errorf("密码更新失败，请重试！")
		goto errDeal
	}

	c.Ctx.Output.JSON(RspPass{
		Success: true,
		Message: "密码修改成功",
	}, false, false)

	return
errDeal:
	c.Ctx.Output.JSON(RspPass{
		Success: false,
		Message: err.Error(),
	}, false, false)
	return
}
