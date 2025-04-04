package edge_create

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

// CreateEdgeStatement represents a CREATE EDGE statement in a graph database.
// It contains the name of the edge, properties, TTL duration, TTL column, and ifNotExists flag.
type CreateEdgeStatement struct {
	name        string         // required
	properties  []EdgeProperty // optional
	ttlDuration uint           // optional
	ttlCol      string         // optional
	ifNotExists bool           // optional
}

// EdgeProperty represents a property in a edge.
// It contains the name, type, and nullable flag.
type EdgeProperty struct {
	field    string
	ttype    statement.PropertyType
	nullable bool
}

func NewEdgeProperty(field string, propertyType statement.PropertyType, nullable bool) EdgeProperty {
	return EdgeProperty{
		field:    field,
		ttype:    propertyType,
		nullable: nullable,
	}
}

// CreateEdgeStatementOption is a functional option for configuring a CreateEdgeStatement.
// It takes a pointer to a CreateEdgeStatement as its argument.
type CreateEdgeStatementOption func(*CreateEdgeStatement)

// NewCreateEdgeStatement creates a new CreateEdgeStatement with the given options.
// It applies each provided option to the statement before returning it.
func NewCreateEdgeStatement(name string, options ...CreateEdgeStatementOption) CreateEdgeStatement {
	statement := CreateEdgeStatement{
		name: name,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithIfNotExists sets the ifNotExists flag of the CreateEdgeStatement to the provided value.
func WithIfNotExists(ifNotExists bool) func(*CreateEdgeStatement) {
	return func(stmt *CreateEdgeStatement) {
		stmt.ifNotExists = ifNotExists
	}
}

// WithTtlCol sets the ttlCol of the CreateEdgeStatement to the provided value.
func WithTtlCol(ttlCol string) func(*CreateEdgeStatement) {
	return func(stmt *CreateEdgeStatement) {
		stmt.ttlCol = ttlCol
	}
}

// WithTtlDuration sets the ttlDuration of the CreateEdgeStatement to the provided value.
func WithTtlDuration(ttlDuration uint) func(*CreateEdgeStatement) {
	return func(stmt *CreateEdgeStatement) {
		stmt.ttlDuration = ttlDuration
	}
}

// WithProperties sets the properties of the CreateEdgeStatement to the provided value.
func WithProperties(properties []EdgeProperty) func(*CreateEdgeStatement) {
	return func(stmt *CreateEdgeStatement) {
		stmt.properties = properties
	}
}

// GenerateCreateEdgeStatement generates a string representation of the CreateEdgeStatement.
// The function checks if the TTL column exists in the properties and returns an error if it doesn't.
// Otherwise, it returns a string representation of the CreateEdgeStatement.
func GenerateCreateEdgeStatement(edge CreateEdgeStatement) (string, error) {
	err := isTTLColValid(edge.ttlCol, edge.properties)
	if err != nil {
		return "", fmt.Errorf("TTL column %s does not exist in the properties", edge.ttlCol)
	}

	var sb strings.Builder
	sb.Grow(len(edge.name) + 10 + len(edge.properties)*(len(edge.name)+10))

	sb.WriteString("CREATE EDGE ")
	if edge.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(edge.name)
	sb.WriteString(" (")

	for i, field := range edge.properties {
		if i > 0 {
			sb.WriteString(", ")
		}

		t := string(field.ttype)
		if t == "" {
			t = "string"
		}
		n := "NULL"
		if !field.nullable {
			n = "NOT NULL"
		}
		sb.WriteString(field.field)
		sb.WriteString(" ")
		sb.WriteString(t)
		sb.WriteString(" ")
		sb.WriteString(n)
	}

	if edge.ttlCol != "" {
		sb.WriteString(") ")
		sb.WriteString(fmt.Sprintf(`TTL_DURATION = %d, TTL_COL = "%s"`, edge.ttlDuration, edge.ttlCol))
		sb.WriteString(";")
	} else {
		sb.WriteString(");")
	}

	return sb.String(), nil
}

func isTTLColValid(ttlCol string, properties []EdgeProperty) error {
	if ttlCol == "" {
		// no ttl column is valid
		return nil
	}

	for _, field := range properties {
		if field.field == ttlCol {
			return nil
		}
	}

	return fmt.Errorf("TTL column %s does not exist in the fields", ttlCol)
}
