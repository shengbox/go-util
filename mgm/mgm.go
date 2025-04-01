package mgm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/joho/godotenv/autoload"
)

var (
	client *mongo.Client
	DbName string
)

// 初始化
func init() {
	client = SetConnect()
	DbName = os.Getenv("mongodb_database")
}

// 连接设置
func SetConnect() *mongo.Client {
	username := os.Getenv("mongodb_username")
	password := os.Getenv("mongodb_password")
	host := os.Getenv("mongodb_host")
	authSource := os.Getenv("mongodb_authSource")

	uri := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s&direct=true&authMechanism=SCRAM-SHA-1", username, password, host, authSource)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(10)) // 连接池
	if err != nil {
		log.Fatal("MongoDB Connect failed:", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	return client
}

func FindOne(DB, collectionName string, filter, selector, result any) error {
	collection := client.Database(DbName).Collection(collectionName)
	opts := options.FindOne().SetProjection(selector)
	if err := collection.FindOne(context.TODO(), filter, opts).Decode(result); err != nil {
		return err
	}
	return nil
}

func FindAll(DB, collectionName string, filter, selector, result any) error {
	collection, _ := client.Database(DB).Collection(collectionName).Clone()
	opts := options.Find()
	opts.SetProjection(selector)
	if filter == nil {
		filter = bson.M{}
	}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), result); err != nil {
		return err
	}
	return nil
}

func FindPage(DB, collectionName string, filter, selector any, page, size int, sort, result any) (count int64, err error) {
	collection, _ := client.Database(DB).Collection(collectionName).Clone()
	count, _ = collection.CountDocuments(context.TODO(), filter)
	skip := (int64(page) - 1) * int64(size)
	opts := options.Find().SetLimit(int64(size))
	opts.SetSkip(skip)

	if sort != nil {
		switch sort := sort.(type) {
		case []string:
			sortBson := bson.D{}
			for _, key := range sort {
				if strings.HasPrefix(key, "-") {
					sortBson = append(sortBson, bson.E{Key: key[1:], Value: -1})
				} else {
					sortBson = append(sortBson, bson.E{Key: key, Value: 1})
				}
			}
			opts.SetSort(sortBson)
		default:
			opts.SetSort(sort)
		}
	}
	if selector != nil {
		opts.SetProjection(selector)
	}

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return count, err
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), result); err != nil {
		return 0, err
	}
	return count, nil
}

func Count(DB, collectionName string, filter any) (int64, error) {
	db := client.Database(DB)
	collection, _ := db.Collection(collectionName).Clone()
	if filter == nil {
		filter = bson.M{}
	}
	return collection.CountDocuments(context.TODO(), filter)
}

func CountAll(DB, collectionName string, filter any) (int, error) {
	db := client.Database(DB)
	collection, _ := db.Collection(collectionName).Clone()
	if filter == nil {
		filter = bson.M{}
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func InsertOne(DB, collectionName string, doc any) (*mongo.InsertOneResult, error) {
	collection := client.Database(DB).Collection(collectionName)
	return collection.InsertOne(context.TODO(), doc)
}

func Insert(DB, collectionName string, docs ...any) (*mongo.InsertManyResult, error) {
	collection := client.Database(DB).Collection(collectionName)
	return collection.InsertMany(context.TODO(), docs)
}

func Upsert(DB, collectionName string, query, update any) (result *mongo.UpdateResult, err error) {
	collection := client.Database(DB).Collection(collectionName)
	replaceOptions := options.Replace().SetUpsert(true)
	return collection.ReplaceOne(context.TODO(), query, update, replaceOptions)
}

func Remove(DB, collectionName string, filter any) error {
	collection := client.Database(DB).Collection(collectionName)
	result, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("delete failed, expected 1 but got 0")
	}
	return nil
}

func RemoveAll(DB, collectionName string, filter any) (int64, error) {
	collection := client.Database(DB).Collection(collectionName)
	count, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count.DeletedCount, nil
}

func Update(DB, collectionName string, filter, update any) error {
	collection, _ := client.Database(DB).Collection(collectionName).Clone()
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func UpdateAll(DB, collectionName string, filter, update any) error {
	collection, _ := client.Database(DB).Collection(collectionName).Clone()
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func Aggregate(DB, collectionName string, pipeline, result any) error {
	collection := client.Database(DB).Collection(collectionName)
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())
	if err = cursor.All(context.TODO(), result); err != nil {
		return err
	}
	return nil
}

func MongoPipeline(str string) mongo.Pipeline {
	var pipeline = []bson.D{}
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.D
		bson.UnmarshalExtJSON([]byte(str), false, &doc)
		pipeline = append(pipeline, doc)
	} else {
		bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	}
	return pipeline
}
