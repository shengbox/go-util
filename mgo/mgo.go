package mgo

import (
	"log"
	"os"
	"strings"

	"github.com/globalsign/mgo"
	_ "github.com/joho/godotenv/autoload"
)

var (
	DbName  string
	globalS *mgo.Session
)

func init() {
	host := os.Getenv("mongodb.host")
	authDb := os.Getenv("mongodb.authdb")
	if authDb == "" {
		authDb = "admin"
	}
	if host != "" {
		DbName = os.Getenv("mongodb.database")
		diaInfo := &mgo.DialInfo{
			Addrs:          strings.Split(host, ","),
			Database:       authDb,
			Username:       os.Getenv("mongodb.username"),
			Password:       os.Getenv("mongodb.password"),
			PoolLimit:      3,
			ReplicaSetName: os.Getenv("mongodb.replicaSet"),
		}
		//mongolog := log.New(os.Stderr, "MONGO ", log.LstdFlags)
		//mgo.SetLogger(mongolog)
		//mgo.SetDebug(true)
		s, err := mgo.DialWithInfo(diaInfo)
		if err != nil {
			log.Fatalln("create session error", err)
		}
		globalS = s
	}
}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB(db).C(collection)
	return s, c
}

func Insert(db, collection string, docs ...interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindPage(db, collection string, query, selector interface{}, pageNum, pageSize int, sort []string, result interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, err := c.Find(query).Count()
	if err != nil {
		return 0, err
	}
	err = c.Find(query).Select(selector).Sort(sort...).Skip((pageNum - 1) * pageSize).Limit(pageSize).All(result)
	return count, err
}

func FindLimit(db, collection string, query, selector interface{}, limit int, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Limit(limit).All(result)
}

func FindAll(db, collection string, query, selector interface{}, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Limit(2000).All(result)
}

func FindAllSort(db, collection string, query, selector interface{}, pageNum int, result interface{}, sortFields ...string) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Sort(sortFields...).Skip((pageNum - 1) * 100).Limit(100).All(result)
}

func CountAll(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

func Distinct(db, collection string, query, selector interface{}, distinctKey string, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Distinct(distinctKey, result)
}

func Update(db, collection string, query, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Update(query, update)
}

func UpdateAll(db, collection string, query, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	_, err := c.UpdateAll(query, update)
	return err
}

func Upsert(db, collection string, query, update interface{}) (info *mgo.ChangeInfo, err error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Upsert(query, update)
}

func Remove(db, collection string, query interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Remove(query)
}

func RemoveAll(db, collection string, query interface{}) (*mgo.ChangeInfo, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.RemoveAll(query)
}

func Pipe(db, collection string, query interface{}, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Pipe(query).All(result)
}
