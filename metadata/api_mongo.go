package metadata

import (
	"fourth.com/apibutler/config"
	"fourth.com/apibutler/mongo"
	"labix.org/v2/mgo/bson"
)

type MongoApiStore struct {
	store *mongo.MongoStore
}

func NewMongoApiStoreFromConfig() ApiStore {
	return &MongoApiStore{
		store: &mongo.MongoStore{
			MongoUrl:    config.Options.MongoUrl,
			MongoDbName: config.Options.MongoDbName,
			NewItemCtor: func() interface{} {
				return &Api{}
			},
			CollectionName: "apis",
		},
	}
}

func (m *MongoApiStore) AddApi(a *Api) error {
	a.ID = bson.NewObjectId()
	err := m.store.Add(a)
	return err
}

func (m *MongoApiStore) Apis() ([]*Api, error) {
	items, err := m.store.ItemIter()

	if err != nil {
		return nil, err
	}

	apis := make([]*Api, 0, 100)

	for i := range items {
		newapi := i.(*Api)
		apis = append(apis, newapi)
	}

	return apis, nil
}

func (m *MongoApiStore) Forget(path string) {
}
