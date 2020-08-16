package models

import (
	"article_share/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	Category string `bson:"category" json:"category"`
	Title    string `bson:"title" json:"title"`
	Content  string `bson:"content" json:"content"`
	Brief    string `bson:"brief" json:"brief"`
}

type ArticleModel struct {
	Client utils.MongoClient // mongo
}

func (a *ArticleModel) Init(cfg *utils.Configure) {
	a.Client.Connect(cfg.Mongo)
}

// 获取文件简图
func (a *ArticleModel) LoadArticlesBrief() []Article {
	collection := a.Client.Collection("blog", "article")
	ret := []Article{}
	collection.Find(bson.M{}).Select(bson.M{"content": 0}).All(&ret)
	return ret
}

//获取文件详情
func (a *ArticleModel) LoadArticleDetail(category, title string) Article {
	collection := a.Client.Collection("blog", "article")
	ret := Article{}
	collection.Find(bson.M{"title": title, "category": category}).One(&ret)
	return ret
}

//修改或者新建文件
func (a *ArticleModel) UpSertArticle(c *gin.Context) string {
	article := Article{
		Title:    c.PostForm("title"),
		Category: c.PostForm("category"),
		Brief:    c.PostForm("brief"),
		Content:  c.PostForm("content")}
	msg := "failed"
	if article.Title != "" && article.Category != "" {
		collection := a.Client.Collection("blog", "article")
		_, err := collection.Upsert(bson.M{"title": article.Title, "category": article.Category}, article)
		if err == nil {
			msg = "success"
		} else {
			msg = "failed to save"
		}
	}
	return msg
}

//修改文件
func (a *ArticleModel) UpdateArticle(c *gin.Context) string {
	article := Article{
		Title:    c.PostForm("title"),
		Category: c.PostForm("category"),
		Brief:    c.PostForm("brief"),
		Content:  c.PostForm("content")}
	msg := "failed"
	if article.Title != "" && article.Category != "" {
		collection := a.Client.Collection("blog", "article")
		err := collection.Update(bson.M{"title": article.Title, "category": article.Category}, article)
		if err == nil {
			msg = "success"
		} else {
			msg = "failed to save"
		}
	}
	return msg
}

//删除文件
func (a *ArticleModel) DeleteArticle(c *gin.Context) (bool, string) {
	title := c.Query("title")
	category := c.Query("category")
	collection := a.Client.Collection("blog", "article")
	err := collection.Remove(bson.M{"title": title, "category": category})
	if err == nil {
		return true, "success"
	}
	return false, err.Error()
}
