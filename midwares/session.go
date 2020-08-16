package midwares

import (
	"article_share/utils"
	"encoding/base64"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

// uuid str
func GetUUID() string {
	b, _ := uuid.New()
	return base64.StdEncoding.EncodeToString(b[:])
}

// session data
type SessionAbstract interface {
	GetID() string
	Get() (value interface{}, err error)
	Set(value interface{})
	Del()
	Exist() bool
}

// session manager
type SessionManagerAbstract interface {
	GetOrCreateSession(uid string) SessionAbstract
	IsExists(uid string) bool
	Init(v utils.ValidateInfo)
}

// session
type Session struct {
	pool utils.RedisPool
	Uid  string
}

//set key
func (s *Session) GetID() string {
	return s.Uid
}

// session token
func (s *Session) Get() (value interface{}, err error) {
	return s.pool.GetCon().Do("GET", s.Uid)
}

// set token
func (s *Session) Set(value interface{}) {
	s.pool.GetCon().Do("SETEX", s.Uid, 3600, value)
}

// del token
func (s *Session) Del() {
	s.pool.GetCon().Do("DEL", s.Uid)
}

// exist
func (s *Session) Exist() bool {
	v, _ := s.pool.GetCon().Do("EXISTS", s.Uid)
	r := v.(bool)
	return r
}

// session manager
type RSessionManager struct {
	Pool utils.RedisPool
}

// 获取session
func (s *RSessionManager) GetOrCreateSession(uid string) SessionAbstract {
	return &Session{pool: s.Pool, Uid: uid}
}

// 判断session是否存在
func (s *RSessionManager) IsExists(uid string) bool {
	v, _ := s.Pool.GetCon().Do("EXISTS", uid)
	r := v.(bool)
	return r
}

// 初始化
func (s *RSessionManager) Init(v utils.ValidateInfo) {
	s.Pool = &utils.RedisPoolInstance{}
	s.Pool.InitRedis(v)
}

type SessionHandlerManager struct {
	CookieSessionKey string                 // cookiename
	CookieTokenKey   string                 // cookie key Token
	Redirect         string                 //重定向位置
	Domain           string                 // domain
	Secure           bool                   //cookie secure
	manager          SessionManagerAbstract // 全局session 管理器
}

func (s *SessionHandlerManager) Init(cfg *utils.Configure, session_id_key, token_key, redirect string) {
	s.CookieSessionKey = session_id_key
	s.CookieTokenKey = token_key
	s.Redirect = redirect
	s.Domain = cfg.Http.Domain
	s.Secure = cfg.Http.Secure
	s.manager = &RSessionManager{}
	s.manager.Init(cfg.Redis)
}

// 不采用指针，SessionHandlerManager不允许被修改
func (s SessionHandlerManager) SessionExistCheck(c *gin.Context) {
	sessionId, err0 := c.Cookie(s.CookieSessionKey)
	_, err1 := c.Cookie(s.CookieTokenKey)
	if err0 != nil || err1 != nil { // session not exists
		u := GetUUID()
		u1 := GetUUID()
		sessionId = u
		session := s.manager.GetOrCreateSession(sessionId)
		token := sessionId + u1
		session.Set(token)
		c.SetCookie(s.CookieTokenKey, token, 3600, "/", s.Domain, s.Secure, false)
		c.SetCookie(s.CookieSessionKey, sessionId, 3600, "/", s.Domain, s.Secure, false)
	}
	c.Next()
}

//通过cookies获取sessionid键值，和token值
//通过sessionid从服务维持的session中获取token值
//两种途径的token进行校验
func (s SessionHandlerManager) SessionValid(c *gin.Context) {
	sessionId, err0 := c.Cookie(s.CookieSessionKey)
	token, err1 := c.Cookie(s.CookieTokenKey)
	if err0 != nil {
		c.Redirect(302, s.Redirect)
		c.Abort()
		return
	}

	if err1 != nil {
		c.SetCookie(s.CookieSessionKey, sessionId, -1, "/", s.Domain, s.Secure, false)
		c.Redirect(302, s.Redirect)
		c.Abort()
		return
	}
	session := s.manager.GetOrCreateSession(sessionId)
	v, _ := session.Get()
	token_0, ok := v.([]byte)

	if !ok {
		c.SetCookie(s.CookieTokenKey, token, -1, "/", s.Domain, s.Secure, false)
		c.SetCookie(s.CookieSessionKey, sessionId, -1, "/", s.Domain, s.Secure, false)
		c.Redirect(302, s.Redirect)
		c.Abort()
		return
	}

	if token != string(token_0) {
		c.SetCookie(s.CookieTokenKey, token, 3600, "/", s.Domain, s.Secure, false)
		c.SetCookie(s.CookieSessionKey, sessionId, 3600, "/", s.Domain, s.Secure, false)
		c.Redirect(302, s.Redirect)
		c.Abort()
		return
	}
	c.Next()
}
