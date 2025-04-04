package tag_alter

import (
	"fmt"
	"strings"
)

// AlterTagStatement represents the ALTER TAG statement in Nebula Graph.
type AlterTagStatement struct {
	tagName  string                    //required
	alterDef []IAlterTagTypeDefinition //required
	ttlDef   []TTLDefinition           // optional
	comment  string                    //optional
}

// AlterTagStatementOption is a functional option for configuring an AlterTagStatement.
type AlterTagStatementOption func(*AlterTagStatement)

// NewAlterTagStatement creates a new AlterTagStatement with the given tag name and alteration definitions.
// It also allows for additional configuration through a variadic list of options.
// The function initializes the statement with the provided tag name and alteration definitions,
// and applies each provided option to the statement before returning it.
// Parameters:
//   - tagName: The name of the tag to be altered.
//   - alterDef: A slice of alteration definitions to be applied to the tag.
//   - options: A variadic list of functions that can modify the AlterTagStatement.
//
// Returns:
//   - An initialized AlterTagStatement configured with the provided parameters and options.
//
// Example usage:
//
//	```
//	alterTagStmt := NewAlterTagStatement("Person", []IAlterTagTypeDefinition{
//	    NewAlterTypeAddDefinition("age", statement.IntType),
//	    NewAlterTypeDropDefinition("address"),
//	}, WithTtlDefinitions([]TTLDefinition{
//	    NewTTLDefinition("age", 30),
//	}), WithTagComment("This is a comment"))
//	```
func NewAlterTagStatement(tagName string, alterDef []IAlterTagTypeDefinition, options ...AlterTagStatementOption) AlterTagStatement {
	statement := AlterTagStatement{
		tagName:  tagName,
		alterDef: alterDef,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithTtlDefinitions sets the TTL definitions for the AlterTagStatement.
// It takes a slice of TTLDefinition as its argument.
// This option is optional and can be used to specify the TTL settings for the tag.
func WithTtlDefinitions(ttlDef []TTLDefinition) func(*AlterTagStatement) {
	return func(stmt *AlterTagStatement) {
		stmt.ttlDef = ttlDef
	}
}

// WithTagComment sets the comment for the AlterTagStatement.
// It takes a string as its argument.
// This option is optional and can be used to add a comment to the tag.
func WithTagComment(comment string) func(*AlterTagStatement) {
	return func(stmt *AlterTagStatement) {
		stmt.comment = comment
	}
}

// GenerateAlterTagStatement generates the ALTER TAG statement based on the provided AlterTagStatement.
// It checks for required fields and constructs the statement string.
// If any required fields are missing, it returns an error.
func GenerateAlterTagStatement(input AlterTagStatement) (string, error) {
	if input.tagName == "" {
		return "", fmt.Errorf("tag name is required")
	}

	if len(input.alterDef) == 0 {
		return "", fmt.Errorf("at least one alter definition is required")
	}

	var sb strings.Builder
	sb.WriteString("ALTER TAG ")
	sb.WriteString(input.tagName)
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
