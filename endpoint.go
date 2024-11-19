package lazyapi

type Endpoint struct {
	Method         string
	Path           string
	QueryParams    map[string]string // e.g., {"userId": "int"}
	BodySchema     *Model            // Reference to a model for body schema
	ResponseSchema *Model            // Reference to a model for response schema
	Action         string            // e.g., "create_record", "get_records"
}
