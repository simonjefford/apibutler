package metadata

type Route struct {
	Path            string
	ApplicationName string
	IsPrefix        bool
	NeedsAuth       bool
}

func NewRoute(path, applicationName string) Route {
	return Route{
		Path:            path,
		ApplicationName: applicationName,
		IsPrefix:        true,
		NeedsAuth:       true,
	}
}

func NewPublicRoute(path, applicationName string) Route {
	r := NewRoute(path, applicationName)
	r.NeedsAuth = false
	return r
}

func GetRoutes() []Route {
	return []Route{
		NewRoute("/recipes", "Test node backend"),
		NewRoute("/recipes/other", "Another test node backend"),
		NewRoute("/stock", "Another test node backend"),
		NewPublicRoute("/stock/public", "Another test node backend"),
	}
}
