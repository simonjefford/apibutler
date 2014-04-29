package metadata

type Api struct {
	Path      string `json:"path"`
	ID        int64  `json:"id"`
	App       string `json:"app"`
	NeedsAuth bool   `json:"needsAuth"`
}

type ApiStorage interface {
	AddApi(a *Api)
	Apis() ([]*Api, error)
	Forget(path string)
}
