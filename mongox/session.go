// Copyright 2018 Kuei-chun Chen. All rights reserved.

package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Session -
type Session struct {
	collection *mongo.Collection
	filter     any
	limit      int64
	project    any
	skip       int64
	sort       any
}

// Limit sets sorting
func (s *Session) Limit(limit int64) *Session {
	s.limit = limit
	return s
}

// Project sets sorting
func (s *Session) Project(project any) *Session {
	s.project = project
	return s
}

// Skip sets sorting
func (s *Session) Skip(skip int64) *Session {
	s.skip = skip
	return s
}

// Sort sets sorting
func (s *Session) Sort(sort any) *Session {
	s.sort = sort
	return s
}

func (s *Session) All(result any) (err error) {
	opts := options.Find()
	if s.limit > 0 {
		opts.SetLimit(s.limit)
	}
	if s.project != nil {
		opts.SetProjection(s.project)
	}
	if s.sort != nil {
		opts.SetSort(s.sort)
	}
	cur, err := s.collection.Find(context.TODO(), s.filter, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())
	if err := cur.All(context.TODO(), result); err != nil {
		return err
	}
	return nil
}

// Decode returns all docs
func (s *Session) Page(result any) (count int64, err error) {
	count, _ = s.collection.CountDocuments(context.TODO(), s.filter)

	opts := options.Find()
	if s.limit > 0 {
		opts.SetLimit(s.limit)
	}
	if s.project != nil {
		opts.SetProjection(s.project)
	}
	if s.skip > 0 {
		opts.SetSkip(s.skip)
	}
	if s.sort != nil {
		opts.SetSort(s.sort)
	}
	cur, err := s.collection.Find(context.TODO(), s.filter, opts)
	if err != nil {
		return count, err
	}
	defer cur.Close(context.TODO())
	if err := cur.All(context.TODO(), result); err != nil {
		return count, err
	}
	return count, nil
}

func (s *Session) One(result interface{}) (err error) {
	opts := options.FindOne()
	if s.project != nil {
		opts.SetProjection(s.project)
	}
	if err := s.collection.FindOne(context.TODO(), s.filter, opts).Decode(result); err != nil {
		return err
	}
	return nil
}

func (s *Session) Count() (int64, error) {
	count, err := s.collection.CountDocuments(context.TODO(), s.filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
