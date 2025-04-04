package edge_alter

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

// AlterTypeAddDefinition represents an ALTER TAG statement for adding a new property.
// It contains the property name, type, and optional attributes such as
type AlterTypeAddDefinition struct {
	PropName string                 // required
	Type     statement.PropertyType // required
}

// AlterAddDefinitionOption is a functional option for configuring an AlterTypeAddDefinition.
type AlterAddDefinitionOption func(*AlterTypeAddDefinition)

// NewAlterTypeAddDefinition creates a new AlterTypeAddDefinition with the given property name and type.
// It also allows for additional configuration through a variadic list of options.
// The function initializes the statement with the provided property name and type,
// and applies each provided option to the statement before returning it.
// Parameters:
//   - propName: The name of the property to be added.
//   - propType: The type of the property to be added.
//   - options: A variadic list of functions that can modify the AlterTypeAddDefinition.
//
// Returns:
//   - An initialized AlterTypeAddDefinition configured with the provided parameters and options.
func NewAlterTypeAddDefinition(propName string, propType statement.PropertyType, options ...AlterAddDefinitionOption) AlterTypeAddDefinition {
	statement := AlterTypeAddDefinition{
		PropName: propName,
		Type:     propType,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// GetAlterType returns the type of alteration being performed.
func (def AlterTypeAddDefinition) GetAlterType() AlterType {
	return AlterTypeAdd
}

// GenerateStatement generates the NGQL statement for adding a new property to a edge.
// It constructs the statement based on the provided property name, type, and optional attributes.
// The function returns the generated statement as a string and an error if any required fields are missing.
func (def AlterTypeAddDefinition) GenerateStatement() (string, error) {

	if def.PropName == "" {
		return "", fmt.Errorf("property name is required")
	}

	var sb strings.Builder
	sb.WriteString(string(def.GetAlterType()))
	sb.WriteString(" (")
	sb.WriteString(def.PropName)
	sb.WriteString(" ")
	sb.WriteString(string(def.Type))
	sb.WriteString(")")
	return sb.String(), nil
}
