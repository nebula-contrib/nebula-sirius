package tag_alter

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

// AlterTypeChangeDefinition represents an ALTER TAG statement for changing a property type.
type AlterTypeChangeDefinition struct {
	propName    string                 // required
	propType    statement.PropertyType // required
	notNullable bool                   // optional
	defaultVal  string                 // optional
	comment     string                 // optional
}

// AlterChangeDefinitionOption is a functional option for configuring an AlterTypeChangeDefinition.
type AlterChangeDefinitionOption func(*AlterTypeChangeDefinition)

// NewAlterTypeChangeDefinition creates a new AlterTypeChangeDefinition with the given property name and type.
// It also allows for additional configuration through a variadic list of options.
// The function initializes the statement with the provided property name and type,
// and applies each provided option to the statement before returning it.
// Parameters:
//   - propName: The name of the property to be changed.
//   - propType: The new type of the property.
//   - options: A variadic list of functions that can modify the AlterTypeChangeDefinition.
//
// Returns:
//   - An initialized AlterTypeChangeDefinition configured with the provided parameters and options.
func NewAlterTypeChangeDefinition(propName string, propType statement.PropertyType, options ...AlterChangeDefinitionOption) AlterTypeChangeDefinition {
	statement := AlterTypeChangeDefinition{
		propName: propName,
		propType: propType,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithAlterTypeChangeNotNullable sets the NotNullable attribute of the AlterTypeChangeDefinition to the provided value.
func WithAlterTypeChangeNotNullable(notNullable bool) func(*AlterTypeChangeDefinition) {
	return func(stmt *AlterTypeChangeDefinition) {
		stmt.notNullable = notNullable
	}
}

// WithAlterTypeChangeDefault sets the Default attribute of the AlterTypeChangeDefinition to the provided value.
func WithAlterTypeChangeDefault(defaultValue string) func(*AlterTypeChangeDefinition) {
	return func(stmt *AlterTypeChangeDefinition) {
		stmt.defaultVal = defaultValue
	}
}

// WithAlterTypeChangeComment sets the Comment attribute of the AlterTypeChangeDefinition to the provided value.
func WithAlterTypeChangeComment(comment string) func(*AlterTypeChangeDefinition) {
	return func(stmt *AlterTypeChangeDefinition) {
		stmt.comment = comment
	}
}

// GetAlterType returns the type of alteration being performed.
func (def AlterTypeChangeDefinition) GetAlterType() AlterType {
	return AlterTypeChange
}

// GenerateStatement generates the NGQL statement for altering the tag type.
func (def AlterTypeChangeDefinition) GenerateStatement() (string, error) {
	if def.propName == "" {
		return "", fmt.Errorf("property name is required")
	}

	var sb strings.Builder
	sb.WriteString(string(def.GetAlterType()))
	sb.WriteString(" (")
	sb.WriteString(def.propName)
	sb.WriteString(" ")
	sb.WriteString(string(def.propType))
	if def.notNullable {
		sb.WriteString(" NOT NULL")

	} else {
		sb.WriteString(" NULL")
	}
	if def.defaultVal != "" {
		sb.WriteString(" DEFAULT ")
		sb.WriteString(fmt.Sprintf("'%s'", def.defaultVal))
	}
	if def.comment != "" {
		sb.WriteString(" COMMENT '")
		sb.WriteString(def.comment)
		sb.WriteString("'")
	}
	sb.WriteString(")")

	return sb.String(), nil
}
