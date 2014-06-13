package middleware

import (
	"fourth.com/apibutler/config"
	"fourth.com/apibutler/mongo"
	"labix.org/v2/mgo/bson"
)

type MongoStackStore struct {
	store *mongo.MongoStore
}

func NewMongoStackStoreFromConfig() StackStore {
	return NewMongoStackStore(config.Options.MongoUrl, config.Options.MongoDbName)
}

func NewMongoStackStore(mongoUrl, mongoDbName string) StackStore {
	return &MongoStackStore{
		store: &mongo.MongoStore{
			MongoUrl:    mongoUrl,
			MongoDbName: mongoDbName,
			NewItemCtor: func() interface{} {
				return &Stack{}
			},
			CollectionName: "stacks",
		},
	}
}

func (m *MongoStackStore) AddStack(s *Stack) error {
	s.ID = bson.NewObjectId()
	return m.store.Add(s)
}

func (m *MongoStackStore) Stacks() ([]*Stack, error) {
	items, err := m.store.ItemIter()

	if err != nil {
		return nil, err
	}

	stacks := make([]*Stack, 0, 100)

	for i := range items {
		stack := i.(*Stack)
		stacks = append(stacks, stack)
	}

	return stacks, nil
}
