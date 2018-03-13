// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"panda/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/v1/tsxm/verifycode", &controllers.UserLoginController{}, "post:VerifyUser")
	beego.Router("/v1/tsxm/regist", &controllers.UserLoginController{}, "post:RegistUser")
	beego.Router("/v1/tsxm/login", &controllers.UserLoginController{}, "post:UserLogin")
	beego.Router("/v1/tsxm/genverifycode", &controllers.VerifyController{}, "get:GenerateCode")
	beego.Router("/v1/tsxm/validateCode", &controllers.VerifyController{}, "get:ValidateCode")
	beego.Router("/v1/tsxm/sendemail", &controllers.EmailController{}, "get:SendEmailCode")
	beego.Router("/v1/tsxm/validemail", &controllers.EmailController{}, "get:ValidateEmailCode")
	beego.Router("/v1/tsxm/agreement", &controllers.AgreeContoller{}, "get:GetAgreement")
	beego.Router("/v1/tsxm/balance", &controllers.BalanceConroller{}, "get:GetBalance")
	beego.Router("/v1/tsxm/recharge", &controllers.QRCodeController{}, "get:GenCode")
	beego.Router("/v1/tsxm/modifyname", &controllers.UserLoginController{}, "post:ModifyNickName")
	beego.Router("/v1/tsxm/transaction", &controllers.TransQContoller{}, "get:GetTransQ")
	beego.Router("/v1/tsxm/forgetpassword", &controllers.UserLoginController{}, "post:UpdatePassWord")
	beego.Router("/v1/tsxm/loadconfig", &controllers.ConfigDataController{}, "get:LoadConfig")
	beego.Router("/v1/tsxm/capture", &controllers.PandaCatchController{}, "get:HandlerPandaCatch")
	beego.Router("/v1/tsxm/capture/result", &controllers.PandaCatchController{}, "get:HandlerGetPandaCatch")
	beego.Router("/v1/tsxm/setwallet", &controllers.WalletController{}, "get:SetWalletAddress")
	beego.Router("/v1/tsxm/trainpet", &controllers.TrainController{}, "get:HandlerTrainPet")
	beego.Router("/v1/tsxm/getPets", &controllers.PetController{}, "get:HandlerGetPets")
	beego.Router("/v1/tsxm/getPetAttribute", &controllers.PetController{}, "get:HandlerGetPetAttribute")
	beego.Router("/v1/tsxm/bonus", &controllers.BonusController{}, "get:HandlerBonus")
	beego.Router("/v1/tsxm/withdrawal", &controllers.WithDrawalController{}, "get:HandlerWithDrawal")
	beego.Router("/v1/tsxm/offer", &controllers.OfferController{}, "get:HandlerDoOffer")
	beego.Router("/v1/tsxm/offer/query", &controllers.OfferController{}, "get:HandlerGetOffer")
	beego.Router("/v1/tsxm/offer/update", &controllers.OfferController{}, "get:HandlerUpdatePrice")
	beego.Router("/v1/tsxm/offer/delete", &controllers.OfferController{}, "get:HandlerDeleteOffer")
	beego.Router("/v1/tsxm/offer/buy", &controllers.OfferController{}, "get:HandlerBuyPet")
	beego.Router("/v1/tsxm/luckdraw", &controllers.GamblingController{}, "get:HandlerGambling")
	beego.Router("/v1/tsxm/notice", &controllers.NoticeController{}, "get:HandlerNotice")
	beego.Router("/v1/tsxm/notice", &controllers.NoticeController{}, "get:HandlerNews")
	beego.Router("/v1/tsxm/invitationcode", &controllers.InvitationController{}, "get:HandlerGenerateInvitationCode")

	beego.Router("/v1/tsxm/test", &controllers.GeneratesvgfileController{}, "get:HandlerGenerate")

	//	ns := beego.NewNamespace("/v1",

	//		beego.NSNamespace("/income",
	//			beego.NSInclude(
	//				&controllers.IncomeController{},
	//			),
	//		),

	//		beego.NSNamespace("/attribute",
	//			beego.NSInclude(
	//				&controllers.AttributeController{},
	//			),
	//		),

	//		beego.NSNamespace("/order",
	//			beego.NSInclude(
	//				&controllers.OrderController{},
	//			),
	//		),

	//		beego.NSNamespace("/pet",
	//			beego.NSInclude(
	//				&controllers.PetController{},
	//			),
	//		),

	//		beego.NSNamespace("/trade",
	//			beego.NSInclude(
	//				&controllers.TradeController{},
	//			),
	//		),

	//		beego.NSNamespace("/incometype",
	//			beego.NSInclude(
	//				&controllers.IncometypeController{},
	//			),
	//		),

	//		beego.NSNamespace("/catch",
	//			beego.NSInclude(
	//				&controllers.CatchController{},
	//			),
	//		),

	//		beego.NSNamespace("/feedeffect",
	//			beego.NSInclude(
	//				&controllers.FeedeffectController{},
	//			),
	//		),

	//		beego.NSNamespace("/ordereffect",
	//			beego.NSInclude(
	//				&controllers.OrdereffectController{},
	//			),
	//		),

	//		beego.NSNamespace("/player",
	//			beego.NSInclude(
	//				&controllers.PlayerController{},
	//			),
	//		),

	//		beego.NSNamespace("/transaccount",
	//			beego.NSInclude(
	//				&controllers.TransaccountController{},
	//			),
	//		),

	//		beego.NSNamespace("/attrvalue",
	//			beego.NSInclude(
	//				&controllers.AttrvalueController{},
	//			),
	//		),
	//	)
	//	beego.AddNamespace(ns)
}
