package mongo

import (
	"errors"
	"fmt"
	"testing"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	testDBName         = "mongotest"
	testCollectionName = "records"
)

func clear(t *testing.T) {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	_, err = sess.DB(testDBName).C(testCollectionName).RemoveAll(bson.M{})
}

type record struct {
	Name string        `json:"name"`
	Age  int           `json:"age"`
	ID   bson.ObjectId `bson:"_id" json:"id"`
}

func storeUnderTest() *MongoStore {
	return &MongoStore{
		MongoUrl:    "localhost:27017",
		MongoDbName: testDBName,
		NewItemCtor: func() interface{} {
			return &record{}
		},
		CollectionName: testCollectionName,
	}
}

func verifyRecord(t *testing.T, r *record) error {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	var result []*record
	err = sess.DB(testDBName).C(testCollectionName).Find(nil).All(&result)
	if err != nil {
		return err
	}

	ok := len(result) == 1 && result[0].Age == r.Age || result[0].Name == r.Name

	if !ok {
		return errors.New(fmt.Sprintf("Unexpected result %v", result))
	}

	return nil
}

func Test_Update(t *testing.T) {
	clear(t)
	s := storeUnderTest()
	r := &record{
		Name: "alice",
		Age:  30,
		ID:   bson.NewObjectId(),
	}

	err := s.Add(r)

	if err != nil {
		t.Fatal(err)
	}

	err = verifyRecord(t, r)

	if err != nil {
		t.Fatal(err)
	}

	r.Age = 31
	err = s.Update(r.ID, r)

	if err != nil {
		t.Fatal(err)
	}

	err = verifyRecord(t, r)

	if err != nil {
		t.Fatal(err)
	}
}
