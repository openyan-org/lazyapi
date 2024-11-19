package lazyapi

type Endpoint struct {
	Method         HTTPMethod
	Path           string
	QueryParams    map[string]string // e.g., {"userId": "int"}
	BodySchema     interface{}       // Reference to a model for body schema
	ResponseSchema interface{}       // Reference to a model for response schema
	Action         ActionType        // e.g., "create_record", "get_records"
}

func NewEndpoint(method HTTPMethod, path string) *Endpoint {
	return &Endpoint{
		Method:      method,
		Path:        path,
		QueryParams: map[string]string{},
		Action:      "none",
	}
}

func (e *Endpoint) SetBodySchema(schema interface{}) {
	e.BodySchema = schema
}

func (e *Endpoint) SetResponseSchema(schema interface{}) {
	e.ResponseSchema = schema
}

func (e *Endpoint) SetAction(action ActionType) {
	e.Action = action
}
