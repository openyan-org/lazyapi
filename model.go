package lazyapi

type Field struct {
	Name        string
	Type        FieldType
	Constraints FieldConstraints
}

type FieldConstraints struct {
	Unique   bool
	Required bool
}

type Relationship struct {
	RelatedModel string
	RelationType string // e.g., "one-to-one", "one-to-many"
}

type Model struct {
	Name          string
	Fields        []Field
	Relationships []Relationship
}

func NewModel(name string, fields []Field, relationships []Relationship) *Model {
	return &Model{
		Name:          name,
		Fields:        fields,
		Relationships: relationships,
	}
}

func (m *Model) AddField(name string, fieldType FieldType, constraints FieldConstraints) {
	m.Fields = append(m.Fields, Field{
		Name:        name,
		Type:        fieldType,
		Constraints: constraints,
	})
}

func (m *Model) AddRelationship(relatedModel, relationType string) {
	m.Relationships = append(m.Relationships, Relationship{
		RelatedModel: relatedModel,
		RelationType: relationType,
	})
}
