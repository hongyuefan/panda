package controllers

import (
	"errors"
	"fmt"
	"panda/models"
	t "panda/transaction"
	"panda/types"
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

	mUser.Email = reqVerify.UserName

	orm = models.NewCommon()

	if err = orm.CommonGetOne(&mUser, "Email"); err != nil {
		goto errDeal
	}

	if mUser.Id != 0 {
		rspVerify = types.RspVerifyCode{
			types.RspBase{
				MemberIsExist: User_Exist,
				Success:       true,
				Message:       types.USER_EXIST,
			},
		}

	} else {
		rspVerify = types.RspVerifyCode{
			types.RspBase{
				MemberIsExist: User_Not_Exist,
				Success:       true,
				Message:       "",
			},
		}
	}

	c.Ctx.Output.JSON(rspVerify, false, false)

	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return

}

func (c *UserLoginController) RegistUser() {
	var (
		reqRgt          types.ReqRegist
		rspRgt          types.RspRegist
		mUser           *models.Player
		orm             *models.Common
		public, privkey string
		token           string
		uid             int64
		err             error
	)
	//	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqRgt); err != nil {
	//		goto errDeal
	//	}

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	reqRgt.NickName = c.Ctx.Request.FormValue("nickName")
	reqRgt.Password = c.Ctx.Request.FormValue("passWord")
	reqRgt.TimeStamp = c.Ctx.Request.FormValue("timeStamp")
	reqRgt.UserName = c.Ctx.Request.FormValue("userName")
	reqRgt.VerifyCode = c.Ctx.Request.FormValue("verifyCode")

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
		Nickname:   reqRgt.NickName,
		Email:      reqRgt.UserName,
		Password:   reqRgt.Password,
		Createtime: time.Now().Unix(),
		PubPublic:  public,
		PubPrivkey: privkey,
	}
	orm = models.NewCommon()

	if err = orm.BeginTx(); err != nil {
		goto errDeal
	}

	if uid, err = orm.CommonInsert(mUser); err != nil {
		goto errDeal
	}

	if token, err = TokenGenerate(uid); err != nil {
		orm.Rollback()
		goto errDeal
	}

	if err = orm.Commit(); err != nil {
		goto errDeal
	}
	rspRgt = types.RspRegist{
		RspBase: types.RspBase{
			MemberIsExist: 0,
			Success:       true,
			Message:       types.USER_REGIST_OK,
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
	c.Ctx.Output.JSON(rspRgt, false, false)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
	return
}

func (c *UserLoginController) ModifyNickName() {

	var (
		nickName string
		userId   int64
		mUser    models.Player
		err      error
		orm      *models.Common
		token    string
	)

	type Result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err = c.Ctx.Request.ParseForm(); err != nil {
		goto errDeal
	}

	nickName = c.Ctx.Request.FormValue("nickname")

	if token, err = ParseToken(c.Ctx.Input.Header("Authorization")); err != nil {
		goto errDeal
	}
	if userId, err = TokenValidate(token); err != nil {
		goto errDeal
	}

	orm = models.NewCommon()

	mUser.Id = userId

	//	if err = orm.CommonGetOne(&mUser, "id"); err != nil {
	//		goto errDeal
	//	}

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

func (c *UserLoginController) UserLogin() {
	var (
		reqLogin types.ReqLogin
		rspLogin types.RspRegist
		mUser    *models.Player
		orm      *models.Common
		token    string
		err      error
	)

	//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqLogin); err != nil {
	//		goto errDeal
	//	}

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

		rspLogin = types.RspRegist{
			RspBase: types.RspBase{
				MemberIsExist: User_Exist,
				Success:       true,
				Message:       types.USER_LOGIN_SUCCESS,
			},
			Data: types.User{
				MemberId:      fmt.Sprintf("%v", mUser.Id),
				NickName:      mUser.Nickname,
				UserName:      mUser.Email,
				UserType:      UserType_Normal,
				Token:         token,
				Avatar:        mUser.Avatar,
				Balance:       mUser.Balance,
				Freeze:        fmt.Sprintf("%v", mUser.Isdel),
				WalletAddress: mUser.Pubkey,
				Mypets:        "",
			},
		}
	} else {
		rspLogin = types.RspRegist{
			RspBase: types.RspBase{
				MemberIsExist: User_Not_Exist,
				Success:       false,
				Message:       types.USER_LOGIN_FAILED,
			},
			Data: types.User{},
		}
	}
	c.Ctx.Output.JSON(rspLogin, false, false)
	return
errDeal:
	ErrorHandler(c.Ctx, err)
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
	verifyCody = c.Ctx.Request.FormValue("verifyCody")

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
