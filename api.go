package lazyapi

type API struct {
	PackageName       string
	Language          Language
	WebFramework      WebFramework
	DatabaseEngine    DatabaseEngine
	PathPrefix        string
	GlobalMiddlewares []Middleware
	Endpoints         []Endpoint
	Models            []Model
}

type Middleware struct {
	Name string
	// configs for necessary values such as JWT secret
	Config map[string]interface{}
}

func NewAPI(packageName string, language Language, webFramework WebFramework) *API {
	return &API{
		PackageName:       packageName,
		Language:          Language(Go),
		WebFramework:      webFramework,
		DatabaseEngine:    DatabaseEngine(None),
		PathPrefix:        "",
		GlobalMiddlewares: []Middleware{},
		Endpoints:         []Endpoint{},
		Models:            []Model{},
	}
}

func (api *API) SetPathPrefix(pathPrefix string) {
	api.PathPrefix = pathPrefix
}

func (api *API) SetDatabase(dbEngine DatabaseEngine) {
	api.DatabaseEngine = dbEngine
}

func (api *API) AddMiddleware(m *Middleware) {
	api.GlobalMiddlewares = append(api.GlobalMiddlewares, *m)
}

func (api *API) AddModel(m *Model) {
	api.Models = append(api.Models, *m)
}

func (api *API) AddEndpoint(e *Endpoint) {
	api.Endpoints = append(api.Endpoints, *e)
}
