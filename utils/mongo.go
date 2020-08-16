package utils

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

// mongo
type MongoClient struct {
	Client *mgo.Session
}

// connect
func (m *MongoClient) Connect(v ValidateInfo) {
	uri := ""
	if v.UserName == "" {
		uri = fmt.Sprintf("mongodb://%s:%d", v.Host, v.Port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", v.UserName, v.PassWord, v.Host, v.Port)
	}
	client, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal("连接mongo失败", err)
	}
	m.Client = client
}

// collection
func (m *MongoClient) Collection(dbname, collectname string) *mgo.Collection {
	return m.Client.DB(dbname).C(collectname)
}
