package routes

type Route struct {
	Path            string
	ApplicationName string
	IsPrefix        bool
}

func Get() []Route {
	return []Route{
		Route{
			Path:            "/recipes",
			ApplicationName: "Test node backend",
			IsPrefix:        true,
		},
		Route{
			Path:            "/recipes/other",
			ApplicationName: "Another test node backend",
			IsPrefix:        true,
		},
		Route{
			Path:            "/stock",
			ApplicationName: "Another test node backend",
			IsPrefix:        true,
		},
	}
}
