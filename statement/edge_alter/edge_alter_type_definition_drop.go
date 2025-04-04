package edge_alter

import (
	"fmt"
	"strings"
)

// AlterTypeDropDefinition represents an ALTER TAG statement for dropping a property.
type AlterTypeDropDefinition struct {
	PropName string // required
}

// NewAlterTypeDropDefinition creates a new AlterTypeDropDefinition with the given property name.
func NewAlterTypeDropDefinition(propName string) AlterTypeDropDefinition {
	return AlterTypeDropDefinition{
		PropName: propName,
	}
}

// GetAlterType returns the type of alteration for this definition.
func (def AlterTypeDropDefinition) GetAlterType() AlterType {
	return AlterTypeDrop
}

// GenerateStatement generates the ALTER TAG statement for dropping a property.
func (def AlterTypeDropDefinition) GenerateStatement() (string, error) {
	if def.PropName == "" {
		return "", fmt.Errorf("property name is required")
	}

	var sb strings.Builder
	sb.WriteString(string(def.GetAlterType()))
	sb.WriteString(" (")
	sb.WriteString(def.PropName)
	sb.WriteString(")")
	return sb.String(), nil
}
