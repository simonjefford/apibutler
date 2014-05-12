package mongo

import "labix.org/v2/mgo"

type MongoStore struct {
	MongoUrl       string
	MongoDbName    string
	CollectionName string
	NewItemCtor    func() interface{}
}

func (m *MongoStore) openSession() (*mgo.Session, error) {
	sess, err := mgo.Dial(m.MongoUrl)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (m *MongoStore) collection(s *mgo.Session) *mgo.Collection {
	db := s.DB(m.MongoDbName)
	return db.C(m.CollectionName)
}

func (m *MongoStore) Add(i interface{}) error {
	sess, err := m.openSession()
	if err != nil {
		return err
	}

	defer sess.Close()

	c := m.collection(sess)

	err = c.Insert(i)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoStore) ItemIter() (chan interface{}, error) {
	sess, err := m.openSession()
	if err != nil {
		return nil, err
	}
	c := m.collection(sess)

	iter := c.Find(nil).Iter()

	items := make(chan interface{})

	go func() {
		defer sess.Close()
		item := m.NewItemCtor()
		for iter.Next(item) {
			items <- item
			item = m.NewItemCtor()
		}
		close(items)
	}()

	return items, nil
}
