package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// 配置文件载入
type Configure struct {
	Mongo  ValidateInfo `json:"mongo"`
	Redis  ValidateInfo `json:"redis"`
	Log    PathOrFile   `json:"log"`
	Dir    PathOrFile   `json:"dir"`
	Admins []Admin      `json:"admin"`
	Http   HttpCf       `json:"http"`
	BeiAn  string       `json:"beian"`
}

// Mongo, redis
type ValidateInfo struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// File or DirConf
type PathOrFile struct {
	Win   string `json:"win32"`
	Linux string `json:"linux"`
}

// users
type Admin struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

// http
type HttpCf struct {
	Secure bool   `json:"secure"`
	Listen string `json:"listen"`
	Domain string `json:"domain"`
}

// 载入配置文件
func LoadConfig(filename string) *Configure {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("读取配置文件失败", err)
	}
	ret := &Configure{}
	err = json.Unmarshal(b, ret)
	if err != nil {
		log.Fatal("解析配置文件失败", err)
	}
	return ret
}
