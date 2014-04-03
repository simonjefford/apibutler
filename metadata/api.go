package metadata

type Api struct {
	Fragment  string `json:"fragment"`
	Limit     int    `json:"limit"`
	Seconds   int    `json:"seconds"`
	ID        int64  `json:"id"`
	App       string `json:"app"`
	IsPrefix  bool   `json:"isPrefix"`
	NeedsAuth bool   `json:"needsAuth"`
}

type ApiStorage interface {
	AddApi(a *Api)
	Apis() ([]*Api, error)
	Forget(path string)
}
