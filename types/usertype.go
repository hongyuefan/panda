package types

var (
	USER_EXIST     = "user already exist"
	USER_REGIST_OK = "user regist success"

	USER_LOGIN_SUCCESS = "user login success"
	USER_LOGIN_FAILED  = "user login failed"
)

type ReqVerifyCode struct {
	UserName  string `json:"userName"`
	TimeStamp string `json:"timeStamp"`
}

type RspBase struct {
	MemberIsExist uint8  `json:"memberIsExist"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
}

type RspVerifyCode struct {
	RspBase
}

type ReqRegist struct {
	UserName   string `json:"userName"`
	Password   string `json:"passWord"`
	VerifyCode string `json:"verifyCode"`
	NickName   string `json:"nickName"`
	TimeStamp  string `json:"timeStamp"`
}

type RspRegist struct {
	RspBase
	Data User `json:"data"`
}

type ReqLogin struct {
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
}

type User struct {
	MemberId      string `json:"memberId"`
	NickName      string `json:"nickName"`
	UserType      string `json:"userType"`
	Token         string `json:"token"`
	UserName      string `json:"userName"`
	Avatar        string `json:"avatar"`
	Balance       string `json:"balance"`
	Freeze        string `json:"freeze"`
	WalletAddress string `json:"walletaddress"`
	Mypets        string `json:"mypets"`
}

type Token struct {
	Uid      int64  `orm:"column(uid)"`
	Token    string `orm:"column(token);size(256)"`
	Creatime int64  `orm:"column(creatime)"`
}
