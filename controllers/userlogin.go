package controllers

import (
	"errors"
	"fmt"
	"panda/models"
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

	//	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqVerify); err != nil {
	//		goto errDeal
	//	}

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
		reqRgt types.ReqRegist
		rspRgt types.RspRegist
		mUser  *models.Player
		orm    *models.Common
		token  string
		uid    int64
		err    error
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

	if err = ValidateUserName(reqRgt.UserName); err != nil {
		goto errDeal
	}
	if err = ValidatePassWord(reqRgt.Password); err != nil {
		goto errDeal
	}

	if len(reqRgt.VerifyCode) != Email_Code_Len || !validEmailCode(reqRgt.VerifyCode, getSessionString(c.GetSession(reqRgt.UserName))) {
		err = errors.New("email verify code not right")
		goto errDeal
	}

	mUser = &models.Player{
		Nickname:   reqRgt.NickName,
		Email:      reqRgt.UserName,
		Password:   reqRgt.Password,
		Createtime: time.Now().Unix(),
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
				Success:       true,
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

func ValidateUserName(userName string) (err error) {
	if len(userName) < 2 {
		err = errors.New("username length not enough")
		return err
	}
	if !strings.Contains(userName, "@") {
		err = errors.New("username format not right")
		return err
	}
	return nil
}

func ValidatePassWord(passWord string) (err error) {
	if len(passWord) < 6 {
		err = errors.New("password length not enough")
		return err
	}
	return nil
}
