package tag_create

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

// CreateTagStatement represents a CREATE TAG statement in a graph database.
// It contains the name of the tag, properties, TTL duration, TTL column, and ifNotExists flag.
type CreateTagStatement struct {
	name        string        // required
	properties  []TagProperty // optional
	ttlDuration uint          // optional
	ttlCol      string        // optional
	ifNotExists bool          // optional
}

// TagProperty represents a property in a tag.
// It contains the name, type, and nullable flag.
type TagProperty struct {
	Field    string
	Type     statement.PropertyType
	Nullable bool
}

// CreateTagStatementOption is a functional option for configuring a CreateTagStatement.
// It takes a pointer to a CreateTagStatement as its argument.
type CreateTagStatementOption func(*CreateTagStatement)

// NewCreateTagStatement creates a new CreateTagStatement with the given options.
// It applies each provided option to the statement before returning it.
func NewCreateTagStatement(name string, options ...CreateTagStatementOption) CreateTagStatement {
	statement := CreateTagStatement{
		name: name,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithIfNotExists sets the ifNotExists flag of the CreateTagStatement to the provided value.
func WithIfNotExists(ifNotExists bool) func(*CreateTagStatement) {
	return func(stmt *CreateTagStatement) {
		stmt.ifNotExists = ifNotExists
	}
}

// WithTtlCol sets the ttlCol of the CreateTagStatement to the provided value.
func WithTtlCol(ttlCol string) func(*CreateTagStatement) {
	return func(stmt *CreateTagStatement) {
		stmt.ttlCol = ttlCol
	}
}

// WithTtlDuration sets the ttlDuration of the CreateTagStatement to the provided value.
func WithTtlDuration(ttlDuration uint) func(*CreateTagStatement) {
	return func(stmt *CreateTagStatement) {
		stmt.ttlDuration = ttlDuration
	}
}

// WithProperties sets the properties of the CreateTagStatement to the provided value.
func WithProperties(properties []TagProperty) func(*CreateTagStatement) {
	return func(stmt *CreateTagStatement) {
		stmt.properties = properties
	}
}

// GenerateCreateTagStatement generates a string representation of the CreateTagStatement.
// The function checks if the TTL column exists in the properties and returns an error if it doesn't.
// Otherwise, it returns a string representation of the CreateTagStatement.
func GenerateCreateTagStatement(tag CreateTagStatement) (string, error) {
	err := isTTLColValid(tag.ttlCol, tag.properties)
	if err != nil {
		return "", fmt.Errorf("TTL column %s does not exist in the properties", tag.ttlCol)
	}

	var sb strings.Builder
	sb.Grow(len(tag.name) + 10 + len(tag.properties)*(len(tag.name)+10))

	sb.WriteString("CREATE TAG ")
	if tag.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(tag.name)
	sb.WriteString(" (")

	for i, field := range tag.properties {
		if i > 0 {
			sb.WriteString(", ")
		}

		t := string(field.Type)
		if t == "" {
			t = "string"
		}
		n := "NULL"
		if !field.Nullable {
			n = "NOT NULL"
		}
		sb.WriteString(field.Field)
		sb.WriteString(" ")
		sb.WriteString(t)
		sb.WriteString(" ")
		sb.WriteString(n)
	}

	if tag.ttlCol != "" {
		sb.WriteString(") ")
		sb.WriteString(fmt.Sprintf(`TTL_DURATION = %d, TTL_COL = "%s"`, tag.ttlDuration, tag.ttlCol))
		sb.WriteString(";")
	} else {
		sb.WriteString(");")
	}

	return sb.String(), nil
}

func isTTLColValid(ttlCol string, properties []TagProperty) error {
	if ttlCol == "" {
		// no ttl column is valid
		return nil
	}

	for _, field := range properties {
		if field.Field == ttlCol {
			return nil
		}
	}

	return fmt.Errorf("TTL column %s does not exist in the fields", ttlCol)
}
