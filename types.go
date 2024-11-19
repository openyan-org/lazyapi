package lazyapi

type DatabaseEngine string

const (
	PostgreSQL DatabaseEngine = "postgresql"
	MySQL      DatabaseEngine = "mysql"
	SQLite     DatabaseEngine = "sqlite"
	None       DatabaseEngine = "none"
)

type Language string

const (
	Python     Language = "python"
	TypeScript Language = "typescript"
	Go         Language = "go"
)

type WebFramework string

const (
	NetHTTP WebFramework = "net/http"
	Hono    WebFramework = "hono"
	Flask   WebFramework = "flask"
)

type FieldType string

const (
	Integer   FieldType = "integer"
	Float     FieldType = "float"
	Text      FieldType = "text"
	Timestamp FieldType = "timestamp"
	Boolean   FieldType = "boolean"
	UUID      FieldType = "uuid"
)

type ActionType string

const (
	InsertRecord ActionType = "insert_record"
	UpdateRecord ActionType = "update_record"
	DeleteRecord ActionType = "delete_record"
	ReadRecord   ActionType = "read_record"
	ListRecord   ActionType = "list_records"
)

type HTTPMethod string

const (
	HttpPost   HTTPMethod = "Post"
	HttpGet    HTTPMethod = "Get"
	HttpPut    HTTPMethod = "Put"
	HTTPPatch  HTTPMethod = "Patch"
	HttpDelete HTTPMethod = "Delete"
)
