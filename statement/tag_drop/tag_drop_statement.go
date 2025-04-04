package tag_drop

import (
	"fmt"
	"strings"
)

// DropTagStatement represents a Drop TAG statement in a graph database.
type DropTagStatement struct {
	name     string // required
	ifExists bool   // optional
}

// DropTagStatementOption is a functional option for configuring a DropTagStatement.
// It takes a pointer to a DropTagStatement as its argument.
type DropTagStatementOption func(*DropTagStatement)

// NewDropTagStatement Drops a new DropTagStatement with the given options.
// It applies each provided option to the statement before returning it.
func NewDropTagStatement(name string, options ...DropTagStatementOption) DropTagStatement {
	statement := DropTagStatement{
		name: name,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithIfExists sets the ifNotExists flag of the DropTagStatement to the provided value.
func WithIfExists(ifExists bool) func(*DropTagStatement) {
	return func(stmt *DropTagStatement) {
		stmt.ifExists = ifExists
	}
}

// GenerateDropTagStatement generates the DROP TAG statement based on the provided DropTagStatement.
func GenerateDropTagStatement(tag DropTagStatement) (string, error) {
	if tag.name == "" {
		return "", fmt.Errorf("tag name cannot be empty")
	}
	var sb strings.Builder

	sb.WriteString("DROP TAG ")
	if tag.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(tag.name)
	sb.WriteString(";")
	return sb.String(), nil
}
