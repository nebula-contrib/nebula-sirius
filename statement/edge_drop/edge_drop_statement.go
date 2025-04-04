package edge_drop

import (
	"fmt"
	"strings"
)

// DropEdgeStatement represents a Drop EDGE statement in a graph database.
type DropEdgeStatement struct {
	name     string // required
	ifExists bool   // optional
}

// DropEdgeStatementOption is a functional option for configuring a DropEdgeStatement.
// It takes a pointer to a DropEdgeStatement as its argument.
type DropEdgeStatementOption func(*DropEdgeStatement)

// NewDropEdgeStatement Drops a new DropEdgeStatement with the given options.
// It applies each provided option to the statement before returning it.
func NewDropEdgeStatement(name string, options ...DropEdgeStatementOption) DropEdgeStatement {
	statement := DropEdgeStatement{
		name: name,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithIfExists sets the ifNotExists flag of the DropEdgeStatement to the provided value.
func WithIfExists(ifExists bool) func(*DropEdgeStatement) {
	return func(stmt *DropEdgeStatement) {
		stmt.ifExists = ifExists
	}
}

// GenerateDropEdgeStatement generates the DROP EDGE statement based on the provided DropEdgeStatement.
func GenerateDropEdgeStatement(edge DropEdgeStatement) (string, error) {
	if edge.name == "" {
		return "", fmt.Errorf("edge name cannot be empty")
	}
	var sb strings.Builder

	sb.WriteString("DROP EDGE ")
	if edge.ifExists {
		sb.WriteString("IF EXISTS ")
	}

	sb.WriteString(edge.name)
	sb.WriteString(";")
	return sb.String(), nil
}
