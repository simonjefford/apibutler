package metadata

import "labix.org/v2/mgo/bson"

type Api struct {
	Path      string        `json:"path"`
	App       string        `json:"app"`
	NeedsAuth bool          `json:"-"`
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Stack     string        `json:"stack"`
}

type ApiStore interface {
	AddApi(a *Api) error
	Apis() ([]*Api, error)
	Forget(path string)
}
