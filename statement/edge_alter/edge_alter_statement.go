package edge_alter

import (
	"fmt"
	"strings"
)

// AlterEdgeStatement represents the ALTER EDGE statement in Nebula Graph.
type AlterEdgeStatement struct {
	edgeName string                     //required
	alterDef []IAlterEdgeTypeDefinition //required
	ttlDef   []TTLDefinition            // optional
	comment  string                     //optional
}

// AlterEdgeStatementOption is a functional option for configuring an AlterEdgeStatement.
type AlterEdgeStatementOption func(*AlterEdgeStatement)

// NewAlterEdgeStatement creates a new AlterEdgeStatement with the given edge name and alteration definitions.
// It also allows for additional configuration through a variadic list of options.
// The function initializes the statement with the provided edge name and alteration definitions,
// and applies each provided option to the statement before returning it.
// Parameters:
//   - edgeName: The name of the edge to be altered.
//   - alterDef: A slice of alteration definitions to be applied to the edge.
//   - options: A variadic list of functions that can modify the AlterEdgeStatement.
//
// Returns:
//   - An initialized AlterEdgeStatement configured with the provided parameters and options.
//
// Example usage:
//
//	```
//	alterEdgeStmt := NewAlterEdgeStatement("Person", []IAlterTypeDefinition{
//	    NewAlterTypeAddDefinition("age", statement.IntType),
//	    NewAlterTypeDropDefinition("address"),
//	}, WithTtlDefinitions([]TTLDefinition{
//	    NewTTLDefinition("age", 30),
//	}), WithTagComment("This is a comment"))
//	```
func NewAlterEdgeStatement(edgeName string, alterDef []IAlterEdgeTypeDefinition, options ...AlterEdgeStatementOption) AlterEdgeStatement {
	statement := AlterEdgeStatement{
		edgeName: edgeName,
		alterDef: alterDef,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithTtlDefinitions sets the TTL definitions for the AlterEdgeStatement.
// It takes a slice of TTLDefinition as its argument.
// This option is optional and can be used to specify the TTL settings for the edge.
func WithTtlDefinitions(ttlDef []TTLDefinition) func(*AlterEdgeStatement) {
	return func(stmt *AlterEdgeStatement) {
		stmt.ttlDef = ttlDef
	}
}

// WithTagComment sets the comment for the AlterEdgeStatement.
// It takes a string as its argument.
// This option is optional and can be used to add a comment to the edge.
func WithTagComment(comment string) func(*AlterEdgeStatement) {
	return func(stmt *AlterEdgeStatement) {
		stmt.comment = comment
	}
}

// GenerateAlterEdgeStatement generates the ALTER EDGE statement based on the provided AlterEdgeStatement.
// It checks for required fields and constructs the statement string.
// If any required fields are missing, it returns an error.
func GenerateAlterEdgeStatement(input AlterEdgeStatement) (string, error) {
	if input.edgeName == "" {
		return "", fmt.Errorf("edge name is required")
	}

	if len(input.alterDef) == 0 {
		return "", fmt.Errorf("at least one alter definition is required")
	}

	var sb strings.Builder
	sb.WriteString("ALTER EDGE ")
	sb.WriteString(input.edgeName)
	sb.WriteString(" ")

	for i, alterDef := range input.alterDef {
		if i > 0 {
			sb.WriteString(", ")
		}
		stmt, err := alterDef.GenerateStatement()
		if err != nil {
			return "", err
		}
		sb.WriteString(stmt)
	}

	if len(input.ttlDef) > 0 {
		sb.WriteString(" ")
		for i, ttl := range input.ttlDef {
			if i > 0 {
				sb.WriteString(", ")
			}
			stmt, err := GenerateTTlDefinitionStatement(ttl)
			if err != nil {
				return "", err
			}
			sb.WriteString(stmt)
		}
	}

	if input.comment != "" {
		sb.WriteString(fmt.Sprintf(` COMMENT '%s'`, input.comment))
	}

	// Add a semicolon at the end of the statement.
	sb.WriteString(";")

	return sb.String(), nil
}
