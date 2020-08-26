package main

import (
	"article_share/midwares"
	"article_share/models"
	"article_share/utils"
	"encoding/base64"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/gin-gonic/gin"
)

//dir
type UserGinHandlers struct {
	Model       *models.ArticleModel
	FileSaveDir string
	Engine      *regexp.Regexp
	BeiAn       string
}

// init
func (u *UserGinHandlers) Init(cfg *utils.Configure) {
	if runtime.GOOS == "windows" {
		u.FileSaveDir = cfg.Dir.Win
	} else {
		u.FileSaveDir = cfg.Dir.Linux
	}
	u.Model = &models.ArticleModel{}
	u.Model.Init(cfg)
	u.Engine = regexp.MustCompile("\\.([^.]+$)")
	u.BeiAn = cfg.BeiAn
}

func main() {
	// 载入配置文件
	cfg := utils.LoadConfig("conf.json")

	// 初始化 session初始化handle， session鉴定及重置handle
	sm := midwares.SessionHandlerManager{}
	sm.Init(cfg, "_session_uid", "_session_utk", "/")
	SessionExistCheck := sm.SessionExistCheck
	SessionValid := sm.SessionValid

	//其它接口
	u := UserGinHandlers{}
	u.Init(cfg)

	//初始化登录handle及鉴权hanlde
	authm := midwares.AuthHandlerManger{}
	authm.Init(cfg, "UserName", "PassWord")

	// gin setting
	gin.SetMode(gin.ReleaseMode)
	fd, _ := os.Create("/var/gin.log")
	gin.DefaultWriter = io.MultiWriter(fd)
	r := gin.Default()
	r.Use(gin.Recovery(), SessionExistCheck)

	r.SetFuncMap(template.FuncMap{
		"user_define_safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	r.StaticFS("/static", http.Dir("./static"))
	r.LoadHTMLGlob("./templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"beian": u.BeiAn})
	})

	AuthAndSave := authm.CreateAuthAndSaveHandle(func(c *gin.Context) {
		c.Redirect(302, "/admin/login")
	})
	Auth := authm.CreateAuthByCookieHandle(func(c *gin.Context) {
		next := c.Query("next")
		if next == "" {
			next = c.Request.RequestURI
		}
		c.Redirect(302, "/admin/login?next="+next)
	})

	r.GET("/admin/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"beian": u.BeiAn})
	})

	r.POST("/admin/login", SessionValid, AuthAndSave, func(c *gin.Context) {
		url := c.Query("next")
		if url == "" {
			url = "/"
		}
		c.Redirect(302, url)
	})

	r.GET("/admin/edit", SessionValid, Auth, func(c *gin.Context) {
		c.HTML(http.StatusOK, "editor.html", gin.H{"beian": u.BeiAn})
	})

	r.POST("/admin/edit", SessionValid, Auth, func(c *gin.Context) {
		msg := u.Model.UpSertArticle(c)
		c.String(http.StatusOK, msg)
	})

	r.GET("/admin/re_edit", SessionValid, Auth, func(c *gin.Context) {
		category := c.Query("category")
		title := c.Query("title")
		article := u.Model.LoadArticleDetail(category, title)
		c.HTML(http.StatusOK, "re_edit.html", gin.H{"category": article.Category, "title": article.Title,
			"content": article.Content, "brief": article.Brief, "beian": u.BeiAn})
	})

	r.POST("/admin/re_edit", SessionValid, Auth, func(c *gin.Context) {
		msg := u.Model.UpdateArticle(c)
		c.String(http.StatusOK, msg)
	})

	r.POST("/files", SessionValid, Auth, func(c *gin.Context) {
		file, _ := c.FormFile("file")
		engines := u.Engine.FindAllString(file.Filename, 1)
		var engine string
		if engines != nil {
			engine = engines[0]
		}
		name := base64.StdEncoding.EncodeToString([]byte(file.Filename)) + engine
		c.SaveUploadedFile(file, u.FileSaveDir+"/"+name)
		c.JSON(http.StatusOK, gin.H{"location": "/files?filename=" + name})
	})
	r.GET("/files", func(c *gin.Context) {
		file := c.Query("filename")
		if file == "" {
			c.String(http.StatusForbidden, "文件不存在")
			return
		}
		file = u.FileSaveDir + "/" + file
		c.File(file)
	})

	r.GET("/article/detail", func(c *gin.Context) {
		title := c.Query("title")
		category := c.Query("category")
		ret := u.Model.LoadArticleDetail(category, title)
		c.HTML(http.StatusOK, "detail.html", gin.H{"category": category, "title": title, "content": ret.Content, "beian": u.BeiAn})
	})

	r.GET("/article/list", func(c *gin.Context) {
		ret := u.Model.LoadArticlesBrief()
		c.JSON(http.StatusOK, ret)
	})

	r.GET("/article/delete", SessionValid, Auth, func(c *gin.Context) {
		ok, s := u.Model.DeleteArticle(c)
		if ok {
			c.Redirect(302, "/")
		} else {
			c.String(403, s)
		}
	})

	r.Run(cfg.Http.Listen)
}
