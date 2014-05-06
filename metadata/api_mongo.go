package metadata

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoApiStore struct {
	MongoUrl    string
	MongoDbName string
}

func (m *MongoApiStore) openSession() (*mgo.Session, error) {
	sess, err := mgo.Dial(m.MongoUrl)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (m *MongoApiStore) apiMongoCollection(s *mgo.Session) *mgo.Collection {
	db := s.DB(m.MongoDbName)
	return db.C("apis")
}

func (m *MongoApiStore) AddApi(a *Api) error {
	sess, err := m.openSession()
	if err != nil {
		return err
	}

	defer sess.Close()

	c := m.apiMongoCollection(sess)
	a.ID = bson.NewObjectId()
	err = c.Insert(a)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoApiStore) Apis() ([]*Api, error) {
	sess, err := m.openSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	c := m.apiMongoCollection(sess)
	apis := make([]*Api, 0, 100)
	iter := c.Find(nil).Iter()
	api := &Api{}
	for iter.Next(&api) {
		apis = append(apis, api)
		api = &Api{}
	}

	return apis, nil
}

func (m *MongoApiStore) Forget(path string) {
}
