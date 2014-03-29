package metadata

type Path struct {
	Fragment string `json:"fragment"`
	Limit    int    `json:"limit"`
	Seconds  int    `json:"seconds"`
	ID       int64  `json:"id"`
}
