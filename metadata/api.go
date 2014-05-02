package metadata

import "labix.org/v2/mgo/bson"

type Api struct {
	Path      string        `json:"path"`
	App       string        `json:"app"`
	NeedsAuth bool          `json:"needsAuth"`
	ID        bson.ObjectId `bson:"_id" json:"id"`
}

type ApiStorage interface {
	AddApi(a *Api)
	Apis() ([]*Api, error)
	Forget(path string)
}
