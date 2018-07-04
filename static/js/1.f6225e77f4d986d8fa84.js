webpackJsonp([1],{188:function(e,t,n){"use strict";t.__esModule=!0;var a=n(85),i=function(e){return e&&e.__esModule?e:{default:e}}(a);t.default=i.default||function(e){for(var t=1;t<arguments.length;t++){var n=arguments[t];for(var a in n)Object.prototype.hasOwnProperty.call(n,a)&&(e[a]=n[a])}return e}},190:function(e,t,n){"use strict";t.__esModule=!0,t.default=function(e){return e.name="van-"+e.name,e.install=e.install||o,e.mixins=e.mixins||[],e.mixins.push(i.default),e},n(44);var a=n(204),i=function(e){return e&&e.__esModule?e:{default:e}}(a),o=function(e){e.component(this.name,this)}},195:function(e,t,n){"use strict";t.__esModule=!0;var a=n(190),i=function(e){return e&&e.__esModule?e:{default:e}}(a);t.default=(0,i.default)({render:function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"van-loading",class:["van-loading--"+e.type,"van-loading--"+e.color],style:e.style},[n("span",{staticClass:"van-loading__spinner",class:"van-loading__spinner--"+e.type},[e._l("spinner"===e.type?12:0,function(e){return n("i")}),"circular"===e.type?n("svg",{staticClass:"van-loading__circular",attrs:{viewBox:"25 25 50 50"}},[n("circle",{attrs:{cx:"50",cy:"50",r:"20",fill:"none"}})]):e._e()],2)])},name:"loading",props:{size:String,type:{type:String,default:"circular"},color:{type:String,default:"black"}},computed:{style:function(){return this.size?{width:this.size,height:this.size}:{}}}})},204:function(e,t,n){"use strict";t.__esModule=!0;var a=n(86);t.default={computed:{$t:function(){var e=this.$options.name,t=e?(0,a.camelize)(e)+".":"",n=this.$vantMessages[this.$vantLang];return function(e){for(var i=arguments.length,o=Array(i>1?i-1:0),s=1;s<i;s++)o[s-1]=arguments[s];var l=(0,a.get)(n,t+e)||(0,a.get)(n,e);return"function"==typeof l?l.apply(null,o):l}}}}},212:function(e,t,n){n(84)},225:function(e,t,n){e.exports=n.p+"static/img/avatar.7925ec6.png"},261:function(e,t,n){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var a=n(212),i=(n.n(a),n(195)),o=n.n(i),s=n(188),l=n.n(s),c=n(83),r=n(43);n.n(r);t.default={data:function(){return{phoneState:"",passwordState:"",nickNameState:"",phone:"",nickName:"",username:"",password:"",passwordAgain:"",captcha:"",isLogin:!0,isForget:!1,agreeValue:!1,popupVisible:!1,phoneCode:"",loading:!1,base64Data:"",codeId:"",time60s:0,timer:null}},computed:l()({},n.i(c.a)({english:function(e){return e.english},userName:function(e){return e.user.userName},agreement:function(e){return e.agreement}}),{code:function(){var e=window.sessionStorage.getItem("code");if(e)return e}}),mounted:function(){this.code&&(this.isLogin=!1,this.isForget=!1),this.userName&&(this.phone=this.uesrName),this.getAgreement()},watch:{isLogin:function(e,t){this.time60s=0,clearInterval(this.timer),this.timer=null,e||this.getCodeSrc(),this.phone="",this.password="",this.phoneCode="",this.captcha=""},isForget:function(e,t){this.time60s=0,clearInterval(this.timer),this.timer=null,this.getCodeSrc(),this.phone="",this.password="",this.phoneCode="",this.captcha=""}},methods:{getCodeSrc:function(){var e=this;this.loading=!0,this.$store.dispatch("getCodebase64",{cb:function(t,n){n&&(e.base64Data=n.pngBase64,e.codeId=n.codeId,e.loading=!1)}})},loginFn:function(){this.phone?this.password?this.setLoginFn():n.i(r.MessageBox)("","请输入账号密码"):n.i(r.MessageBox)("","请输入账号手机")},toRegist:function(){this.isLogin=!1,this.isForget=!1},tips:function(e){if(!/^[A-Za-z0-9]{6,20}$/.test(e))return void n.i(r.MessageBox)("提示","密码需由6-20位数字和字母组合而成！")},toast:function(e,t){n.i(r.Toast)({message:e,position:"middle",duration:t})},registFn:function(){var e=this;if(this.tips(this.password),""===this.nickName)return void this.toast("昵称不能为空",2e3);this.agreeValue?(this.loading=!0,e.$store.dispatch("setRegist",{params:{userName:e.phone,passWord:e.password,verifyCode:e.phoneCode,nickName:e.nickName,code:e.code?e.code:""},cb:function(t,n){e.loading=!1,n?(e.setLoginFn(),window.sessionStorage.setItem("code","")):t&&e.toast(t.message,2e3)}})):this.toast("请先同意协议",2e3)},forgetFn:function(){this.isForget=!0},forgetBackLogin:function(){this.isForget=!1},registBackLogin:function(){this.isLogin=!0,this.isForget=!1},changeLanguage:function(){this.$store.commit("ENGLISH",!this.english)},forgetOKFn:function(){return""===this.phoneCode?void this.toast("短信码不能为空",2e3):(this.tips(this.password),this.password!==this.passwordAgain?void this.toast("密码不一致",2e3):void this.setForgetPassword())},setForgetPassword:function(){var e=this;this.$store.dispatch("setForgetpassword",{params:{userName:e.phone,newPassword:e.password,verifyCode:e.phoneCode},cb:function(t,n){e.loading=!1,n?e.setLoginFn():t&&this.toast(t.message,2e3)}})},sendPhoneCode:function(){var e=this;if(""===this.phone)return void this.toast("请输入手机",1e3);if(this.checkPhone(this.phone)){if(""===this.captcha)return void this.toast("请输入图形验证码",1e3);if(this.time60s>0)return void this.toast("请"+e.time60s+"秒后再试！",1e3);this.time60s=60,this.timer=setInterval(function(){e.time60s>0?e.time60s=e.time60s-1:(clearInterval(e.timer),e.timer=null)},1e3),""!==this.phone&&""!==this.captcha&&(this.loading=!0,this.$store.dispatch("getSendMsg",{params:{mobile:e.phone,codeId:e.codeId,verifyValue:e.captcha},cb:function(t,n){e.loading=!1,n?e.toast(n.message,2e3):t&&(e.toast(t.message,2e3),e.time60s=0,clearInterval(e.timer),e.tiemr=null)}}))}},setLoginFn:function(){var e=this;this.loading=!0,this.$store.dispatch("setLogin",{userName:e.phone,passWord:e.password,cb:function(t,n){e.loading=!1,n?n.memberIsExist?(window.sessionStorage.setItem("token",n.info.token),e.$router.push({meta:{auth:!0},path:"/home"})):(console.log(n,"data里没有success和message"),e.toast(n.message,2e3)):t&&e.toast(t.message,2e3)}})},checkPhone:function(e){return!!/^((1[3-8][0-9])+\d{8})$/.test(e)||(this.toast("您输入的手机格式不正确",2e3),!1)},getAgreement:function(){this.$store.dispatch("getAgreement",{params:{type:"regist"}})}},components:{vanLoading:o.a}}},301:function(e,t,n){var a=n(178);t=e.exports=n(176)(!0),t.push([e.i,".login[data-v-72ab13c0]{background:#00182c url("+a(n(353))+") no-repeat;background-size:cover;display:-webkit-box;display:-ms-flexbox;display:flex;justify-content:center;-webkit-box-align:center;-ms-flex-align:center;align-items:center;-webkit-box-pack:center;-webkit-box-orient:vertical;-webkit-box-direction:normal;-ms-flex-direction:column;flex-direction:column;padding-bottom:0}.login .lang[data-v-72ab13c0]{position:absolute;right:10px;top:10px}.login .login-con[data-v-72ab13c0]{width:80%}.login .login-con-header-title[data-v-72ab13c0]{color:#fff;font-size:28px}.login .login-logo[data-v-72ab13c0]{display:block;margin:0 auto;width:100px;height:100px;line-height:100px;border-radius:50px;border:1px solid #ccc;background:#fff url("+a(n(225))+") 50% no-repeat;background-size:100% 100%}.login .agree-popup[data-v-72ab13c0]{width:100%;height:100%;background-color:#fff}.login .login-agreement-ck[data-v-72ab13c0]{display:inline-block;vertical-align:middle}.login .login-agreement-con[data-v-72ab13c0]{margin:20px auto}.login .login-agreement-con input[data-v-72ab13c0]{vertical-align:middle}.login .login-agreement-con .login-viewprotocol[data-v-72ab13c0]{cursor:pointer}.login .login-field[data-v-72ab13c0]{margin:10px 0;border-radius:5px;opacity:.8}.login .login-mt-button[data-v-72ab13c0]{opacity:.8}.login .login-btn[data-v-72ab13c0]{margin:20px auto}.login .login-phoneBtn[data-v-72ab13c0],.login .login-phoneCode[data-v-72ab13c0]{display:inline-block;vertical-align:middle}.login .login-phoneCode[data-v-72ab13c0]{width:68%}.login .login-phoneBtn[data-v-72ab13c0]{width:30%}.login .login-protocol-con[data-v-72ab13c0]{padding:20px;height:calc(100% - 41px);overflow:auto}","",{version:3,sources:["/Users/huigeek/lxh/myPro/github/Hcat/src/pages/login.vue"],names:[],mappings:"AACA,wBACE,2DAAyD,AACzD,sBAAuB,AACvB,oBAAqB,AACrB,oBAAqB,AACrB,aAAc,AACd,uBAAwB,AACxB,yBAA0B,AAC1B,sBAAuB,AACvB,mBAAoB,AACpB,wBAAyB,AACzB,4BAA6B,AAC7B,6BAA8B,AAC9B,0BAA2B,AAC3B,sBAAuB,AACvB,gBAAkB,CACnB,AACD,8BACI,kBAAmB,AACnB,WAAY,AACZ,QAAU,CACb,AACD,mCACI,SAAW,CACd,AACD,gDACI,WAAY,AACZ,cAAgB,CACnB,AACD,oCACI,cAAe,AACf,cAAe,AACf,YAAa,AACb,aAAc,AACd,kBAAmB,AACnB,mBAAoB,AACpB,sBAAuB,AACvB,4DAAgE,AAChE,yBAA2B,CAC9B,AACD,qCACI,WAAY,AACZ,YAAa,AACb,qBAAuB,CAC1B,AACD,4CACI,qBAAsB,AACtB,qBAAuB,CAC1B,AACD,6CACI,gBAAkB,CACrB,AACD,mDACM,qBAAuB,CAC5B,AACD,iEACM,cAAgB,CACrB,AACD,qCACI,cAAe,AACf,kBAAmB,AACnB,UAAa,CAChB,AACD,yCACI,UAAa,CAChB,AACD,mCACI,gBAAkB,CACrB,AACD,iFAEI,qBAAsB,AACtB,qBAAuB,CAC1B,AACD,yCACI,SAAW,CACd,AACD,wCACI,SAAW,CACd,AACD,4CACI,aAAc,AACd,yBAA0B,AAC1B,aAAe,CAClB",file:"login.vue",sourcesContent:["\n.login[data-v-72ab13c0] {\n  background: #00182c url(../assets/img/bg2.png) no-repeat;\n  background-size: cover;\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  justify-content: center;\n  -webkit-box-align: center;\n  -ms-flex-align: center;\n  align-items: center;\n  -webkit-box-pack: center;\n  -webkit-box-orient: vertical;\n  -webkit-box-direction: normal;\n  -ms-flex-direction: column;\n  flex-direction: column;\n  padding-bottom: 0;\n}\n.login .lang[data-v-72ab13c0] {\n    position: absolute;\n    right: 10px;\n    top: 10px;\n}\n.login .login-con[data-v-72ab13c0] {\n    width: 80%;\n}\n.login .login-con-header-title[data-v-72ab13c0] {\n    color: #fff;\n    font-size: 28px;\n}\n.login .login-logo[data-v-72ab13c0] {\n    display: block;\n    margin: 0 auto;\n    width: 100px;\n    height: 100px;\n    line-height: 100px;\n    border-radius: 50px;\n    border: 1px solid #ccc;\n    background: #fff url(../assets/img/avatar.png) center no-repeat;\n    background-size: 100% 100%;\n}\n.login .agree-popup[data-v-72ab13c0] {\n    width: 100%;\n    height: 100%;\n    background-color: #fff;\n}\n.login .login-agreement-ck[data-v-72ab13c0] {\n    display: inline-block;\n    vertical-align: middle;\n}\n.login .login-agreement-con[data-v-72ab13c0] {\n    margin: 20px auto;\n}\n.login .login-agreement-con input[data-v-72ab13c0] {\n      vertical-align: middle;\n}\n.login .login-agreement-con .login-viewprotocol[data-v-72ab13c0] {\n      cursor: pointer;\n}\n.login .login-field[data-v-72ab13c0] {\n    margin: 10px 0;\n    border-radius: 5px;\n    opacity: 0.8;\n}\n.login .login-mt-button[data-v-72ab13c0] {\n    opacity: 0.8;\n}\n.login .login-btn[data-v-72ab13c0] {\n    margin: 20px auto;\n}\n.login .login-phoneCode[data-v-72ab13c0],\n  .login .login-phoneBtn[data-v-72ab13c0] {\n    display: inline-block;\n    vertical-align: middle;\n}\n.login .login-phoneCode[data-v-72ab13c0] {\n    width: 68%;\n}\n.login .login-phoneBtn[data-v-72ab13c0] {\n    width: 30%;\n}\n.login .login-protocol-con[data-v-72ab13c0] {\n    padding: 20px;\n    height: calc(100% - 41px);\n    overflow: auto;\n}\n"],sourceRoot:""}])},339:function(e,t,n){var a=n(301);"string"==typeof a&&(a=[[e.i,a,""]]),a.locals&&(e.exports=a.locals);n(177)("6a855236",a,!0,{})},353:function(e,t,n){e.exports=n.p+"static/img/bg2.17def5d.png"},407:function(e,t){e.exports={render:function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",{staticClass:"login page"},[e.isLogin?n("div",{staticClass:"login-con"},[n("div",{staticClass:"lang"},[n("mt-button",{staticClass:"iconfont login-mt-button",class:{"icon-yingwenyuyan":e.english,"icon-zhongwenyuyan":!e.english},nativeOn:{click:function(t){return e.changeLanguage(t)}}})],1),e._v(" "),e.isForget?n("div",[n("div",[n("div",{staticClass:"tac login-btn"},[n("div",{staticClass:"login-logo"}),e._v(" "),n("div",{staticClass:"login-con-header-title"},[e.english?n("span",[e._v("Forget")]):n("span",[e._v("忘记密码")]),e._v(" "),n("br"),e._v(" "),e.english?n("span",[e._v("tsxm")]):n("span",[e._v("太上熊猫")])])]),e._v(" "),n("div",{staticClass:"login-con-body"},[e.english?n("div",[n("mt-field",{staticClass:"login-field",attrs:{label:"phone",placeholder:"Please enter the phone"},model:{value:e.phone,callback:function(t){e.phone=t},expression:"phone"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{placeholder:"verification code",label:"verification code"},model:{value:e.captcha,callback:function(t){e.captcha=t},expression:"captcha"}},[n("img",{attrs:{src:e.base64Data,height:"45px",width:"100px"},on:{click:function(t){e.getCodeSrc()}}})]),e._v(" "),n("mt-field",{staticClass:"login-field login-phoneCode",attrs:{label:"phone code",placeholder:"Please enter the phone code"},model:{value:e.phoneCode,callback:function(t){e.phoneCode=t},expression:"phoneCode"}}),e._v(" "),n("mt-button",{staticClass:"login-phoneBtn login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.sendPhoneCode(t)}}},[e.english?n("span",[e._v("send")]):n("span",[e._v("发送")])]),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"password",placeholder:"Please enter a new password",type:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"password",placeholder:"Please enter your new password again",type:"password"},model:{value:e.passwordAgain,callback:function(t){e.passwordAgain=t},expression:"passwordAgain"}})],1):n("div",[n("mt-field",{staticClass:"login-field",attrs:{label:"手机",placeholder:"请输入手机"},model:{value:e.phone,callback:function(t){e.phone=t},expression:"phone"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{placeholder:"验证码",label:"验证码"},model:{value:e.captcha,callback:function(t){e.captcha=t},expression:"captcha"}},[n("img",{attrs:{src:e.base64Data,height:"45px",width:"100px"},on:{click:function(t){e.getCodeSrc()}}})]),e._v(" "),n("mt-field",{staticClass:"login-field login-phoneCode",attrs:{label:"手机验证码",placeholder:"验证码"},model:{value:e.phoneCode,callback:function(t){e.phoneCode=t},expression:"phoneCode"}}),e._v(" "),n("mt-button",{staticClass:"login-phoneBtn login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.sendPhoneCode(t)}}},[e.english?n("span",[e._v("send")]):n("span",[e._v("发送")])]),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"密码",placeholder:"请输入新密码",type:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"密码",placeholder:"请再次输入新密码",type:"password"},model:{value:e.passwordAgain,callback:function(t){e.passwordAgain=t},expression:"passwordAgain"}})],1),e._v(" "),n("div",{staticClass:"tac login-btn"},[n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.forgetOKFn(t)}}},[e.english?n("span",[e._v("OK")]):n("span",[e._v("确定")])]),e._v(" "),n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.forgetBackLogin(t)}}},[e.english?n("span",[e._v("Back")]):n("span",[e._v("返回登录")])])],1)])])]):n("div",[n("div",{staticClass:"tac login-btn"},[n("div",{staticClass:"login-logo tac"}),e._v(" "),n("div",{staticClass:"login-con-header-title"},[e.english?n("span",[e._v("Login")]):n("span",[e._v("登录")]),e._v(" "),n("br"),e._v(" "),e.english?n("span",[e._v("tsxm")]):n("span",[e._v("太上熊猫")])])]),e._v(" "),n("div",[e.english?n("div",[n("mt-field",{staticClass:"login-field",attrs:{label:"Phone",placeholder:"Please enter the phone",state:e.phoneState},model:{value:e.phone,callback:function(t){e.phone=t},expression:"phone"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"Password",placeholder:"Please enter the password",type:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}})],1):n("div",[n("mt-field",{staticClass:"login-field",attrs:{label:"手机",placeholder:"请输入手机",state:e.phoneState},model:{value:e.phone,callback:function(t){e.phone=t},expression:"phone"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"密码",placeholder:"请输入密码",type:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}})],1),e._v(" "),n("div",{staticClass:"tac login-btn"},[n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.loginFn(t)}}},[e.english?n("span",[e._v("Login")]):n("span",[e._v("登录")])]),e._v(" "),n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.toRegist(t)}}},[e.english?n("span",[e._v("Regist")]):n("span",[e._v("注册")])]),e._v(" "),n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.forgetFn(t)}}},[e.english?n("span",[e._v("Forget")]):n("span",[e._v("忘记密码")])])],1)])])]):n("div",{staticClass:"login-con"},[n("div",[n("div",{staticClass:"tac login-btn"},[n("div",{staticClass:"login-logo"}),e._v(" "),n("div",{staticClass:"login-con-header-title"},[e.english?n("span",[e._v("Regist")]):n("span",[e._v("注册")]),e._v(" "),n("br"),e._v(" "),e.english?n("span",[e._v("tsxm")]):n("span",[e._v("太上熊猫")])])]),e._v(" "),n("div",{staticClass:"login-con-body"},[e.english?n("div",[n("mt-field",{staticClass:"login-field",attrs:{label:"phone",placeholder:"Please enter the phone"},model:{value:e.phone,callback:function(t){e.phone=t},expression:"phone"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"verification code",placeholder:"verification code"},model:{value:e.captcha,callback:function(t){e.captcha=t},expression:"captcha"}},[n("img",{attrs:{src:e.base64Data,height:"45px",width:"100px"},on:{click:function(t){e.getCodeSrc()}}})]),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"phone code",placeholder:"Please enter the phone code"},model:{value:e.phoneCode,callback:function(t){e.phoneCode=t},expression:"phoneCode"}}),e._v(" "),n("mt-button",{staticClass:"login-phoneBtn login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.sendPhoneCode(t)}}},[e.english?n("span",[e._v("send")]):n("span",[e._v("发送")])]),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"nickname",placeholder:"Please enter the nickname",attr:{maxlength:10}},model:{value:e.username,callback:function(t){e.username=t},expression:"username"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"password",placeholder:"Please enter the password",type:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}})],1):n("div",[n("mt-field",{staticClass:"login-field",attrs:{label:"手机",placeholder:"请输入手机"},model:{value:e.phone,callback:function(t){e.phone=t},expression:"phone"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"验证码",placeholder:"验证码"},model:{value:e.captcha,callback:function(t){e.captcha=t},expression:"captcha"}},[n("img",{attrs:{src:e.base64Data,height:"45px",width:"100px"},on:{click:function(t){e.getCodeSrc()}}})]),e._v(" "),n("mt-field",{staticClass:"login-field login-phoneCode",attrs:{label:"手机验证码",placeholder:"验证码"},model:{value:e.phoneCode,callback:function(t){e.phoneCode=t},expression:"phoneCode"}}),e._v(" "),n("mt-button",{staticClass:"login-phoneBtn login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.sendPhoneCode(t)}}},[e.english?n("span",[e._v("send")]):n("span",[e._v("发送")])]),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"昵称",placeholder:"请输入昵称",attr:{maxlength:10}},model:{value:e.nickName,callback:function(t){e.nickName=t},expression:"nickName"}}),e._v(" "),n("mt-field",{staticClass:"login-field",attrs:{label:"密码",placeholder:"请输入密码",type:"password"},model:{value:e.password,callback:function(t){e.password=t},expression:"password"}})],1),e._v(" "),n("div",{staticClass:"login-agreement-con"},[n("input",{directives:[{name:"model",rawName:"v-model",value:e.agreeValue,expression:"agreeValue"}],attrs:{type:"checkbox"},domProps:{checked:Array.isArray(e.agreeValue)?e._i(e.agreeValue,null)>-1:e.agreeValue},on:{change:function(t){var n=e.agreeValue,a=t.target,i=!!a.checked;if(Array.isArray(n)){var o=e._i(n,null);a.checked?o<0&&(e.agreeValue=n.concat([null])):o>-1&&(e.agreeValue=n.slice(0,o).concat(n.slice(o+1)))}else e.agreeValue=i}}}),e._v(" "),n("span",[e._v("同意协议")]),e._v(" "),n("span",{staticClass:"login-field login-viewprotocol",on:{click:function(t){e.popupVisible=!0}}},[e.english?n("span",[e._v("right popup")]):n("span",[e._v("查看协议内容")])])]),e._v(" "),n("div",{staticClass:"tac login-btn"},[n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.registFn(t)}}},[e.english?n("span",[e._v("Regist")]):n("span",[e._v("注册")])]),e._v(" "),n("mt-button",{staticClass:"login-mt-button",attrs:{type:"default"},nativeOn:{click:function(t){return e.registBackLogin(t)}}},[e.english?n("span",[e._v("Return")]):n("span",[e._v("返回登录")])])],1)])])]),e._v(" "),e.loading?n("div",{staticClass:"loading"},[n("van-loading",{attrs:{type:"spinner",color:"black"}})],1):e._e(),e._v(" "),n("mt-popup",{staticClass:"agree-popup",attrs:{position:"right",modal:!1},model:{value:e.popupVisible,callback:function(t){e.popupVisible=t},expression:"popupVisible"}},[n("mt-button",{staticClass:"login-mt-button",attrs:{size:"large",type:"default"},nativeOn:{click:function(t){e.popupVisible=!1}}},[e.english?n("span",[e._v("Close Protocol")]):n("span",[e._v("关闭协议")])]),e._v(" "),e.agreement?n("div",{staticClass:"login-protocol-con"},[e.english?n("span",[e._v("Protocol content")]):n("span",[n("div",{domProps:{innerHTML:e._s(e.agreement)}})])]):e._e()],1)],1)},staticRenderFns:[]}},92:function(e,t,n){function a(e){n(339)}var i=n(82)(n(261),n(407),a,"data-v-72ab13c0",null);e.exports=i.exports}});
//# sourceMappingURL=1.f6225e77f4d986d8fa84.js.map