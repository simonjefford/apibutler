package middleware

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoStackStore struct {
	MongoUrl    string
	MongoDbName string
}

func (m *MongoStackStore) openSession() (*mgo.Session, error) {
	sess, err := mgo.Dial(m.MongoUrl)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (m *MongoStackStore) stackMongoCollection(s *mgo.Session) *mgo.Collection {
	db := s.DB(m.MongoDbName)
	return db.C("stacks")
}

func (m *MongoStackStore) AddStack(s *Stack) error {
	sess, err := m.openSession()
	if err != nil {
		return err
	}

	defer sess.Close()

	c := m.stackMongoCollection(sess)
	s.ID = bson.NewObjectId()
	err = c.Insert(s)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStackStore) Stacks() ([]*Stack, error) {
	sess, err := m.openSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	c := m.stackMongoCollection(sess)
	stacks := make([]*Stack, 0, 100)
	iter := c.Find(nil).Iter()
	stack := &Stack{}
	for iter.Next(stack) {
		stacks = append(stacks, stack)
		stack = &Stack{}
	}

	return stacks, nil
}
