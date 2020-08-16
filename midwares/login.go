package midwares

import (
	"article_share/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandlerManger struct {
	UserNameKey string        //UserName
	PWdNameKey  string        // PassWord
	Admins      []utils.Admin // admins
	Domain      string        // domain
	Secure      bool          //cookie secure
}

func (a *AuthHandlerManger) Init(cfg *utils.Configure, userNameKey, pwdNameKey string) {
	a.UserNameKey = userNameKey
	a.PWdNameKey = pwdNameKey
	a.Admins = cfg.Admins
	a.Domain = cfg.Http.Domain
	a.Secure = cfg.Http.Secure
}

// 保存登录信息
func (a AuthHandlerManger) CreateAuthAndSaveHandle(auth_failed_callback func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		userName := c.PostForm("UserName")
		passWord := c.PostForm("PassWord")
		pwd := ""
		for _, admin := range a.Admins {
			if admin.UserName == userName {
				pwd = admin.PassWord // 明文对比
				break
			}
		}
		if pwd != "" && passWord != pwd {
			auth_failed_callback(c)
			c.Abort()
			return
		}
		c.SetCookie(a.UserNameKey, userName, 3600, "/", a.Domain, a.Secure, false)
		c.SetCookie(a.PWdNameKey, passWord, 3600, "/", a.Domain, a.Secure, false)
		c.Next()
	}
}

// 认证登录信息handle
func (a AuthHandlerManger) CreateAuthByCookieHandle(auth_failed_callback func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取保存的登录信息
		un, err0 := c.Cookie(a.UserNameKey)
		pw, err1 := c.Cookie(a.PWdNameKey)

		// 验证登录信息是否获取成功
		if err0 != nil || err1 != nil || un == "" || pw == "" {
			auth_failed_callback(c)
			c.Abort()
			return
		}

		// 账户名对应的真实密码获取
		pwd := ""
		for _, admin := range a.Admins {
			if admin.UserName == un {
				pwd = admin.PassWord // 明文对比
				break
			}
		}

		// 密码比较
		if pwd != "" && pw != pwd {
			auth_failed_callback(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
