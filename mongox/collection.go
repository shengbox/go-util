// Copyright 2018 Kuei-chun Chen. All rights reserved.

package mongox

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection contains mongo.Collection
type Collection struct {
	collection *mongo.Collection
}

// Find finds docs by given filter
func (c *Collection) Find(filter interface{}) *Session {
	return &Session{filter: filter, collection: c.collection}
}

func (s *Collection) Update(filter, update any) error {
	_, err := s.collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func (s *Collection) UpdateAll(filter, update any) error {
	_, err := s.collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func (s *Collection) Insert(docs ...any) (*mongo.InsertManyResult, error) {
	return s.collection.InsertMany(context.TODO(), docs)
}

func (s *Collection) Upsert(query, update any) (result *mongo.UpdateResult, err error) {
	replaceOptions := options.Replace().SetUpsert(true)
	return s.collection.ReplaceOne(context.TODO(), query, update, replaceOptions)
}

func (c *Collection) Remove(filter any) error {
	result, err := c.collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("delete failed, expected 1 but got 0")
	}
	return nil
}

func (c *Collection) Aggregate(pipeline, result any) error {
	cur, err := c.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())
	if err = cur.All(context.TODO(), result); err != nil {
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
