package middleware

import (
	"testing"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"fourth.com/apibutler/jsonconfig"
)

func clear(t *testing.T) {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	_, err = sess.DB("stack_test").C("stacks").RemoveAll(bson.M{})
}

func Test_InsertAndRetrieve(t *testing.T) {
	clear(t)
	store := NewMongoStackStore("localhost:27017", "stack_test")

	s := NewStack()
	s.Name = "default"
	s.AddMiddleware("mongo_teststack", jsonconfig.Obj{
		"header": "foo",
		"life":   42,
	})

	store.AddStack(s)
	stacks, err := store.Stacks()

	if err != nil {
		t.Fatal(err)
	}

	count := len(stacks)

	if count != 1 {
		t.Fatalf("Unexpected number of stacks: %d (stacks = %v)", count, stacks)
	}

	h := stacks[0].Middlewares[0].Config["header"].(string)

	if h != "foo" {
		t.Errorf("Unexpected config value. %v", stacks[0].Middlewares[0].Config)
	}

	m := stacks[0].Middlewares[0]

	if m.Name != "mongo_teststack" {
		t.Errorf("Unexpected middleware name %s", m)
	}
}
