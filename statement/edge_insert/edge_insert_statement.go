package edge_insert

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"sort"
	"strings"
)

type InsertEdgeStatement[TVidType statement.VidType] struct {
	SourceVid   TVidType
	TargetVid   TVidType
	Properties  map[string]interface{}
	Rank        int
	EdgeType    string
	IfNotExists bool
}

type InsertEdgeStatementOption[TVidType statement.VidType] func(*InsertEdgeStatement[TVidType])

// NewInsertEdgeStatement creates a new InsertEdgeStatement with the given source
// vertex ID, target vertex ID, and edge type. It also allows for additional configuration
// through a variadic list of options. The function initializes the properties map
// and sets the IfNotExists flag to false by default. It applies each provided option
// to the statement before returning it.
//
// Parameters:
//   - sourceVid: The ID of the source vertex.
//   - targetVid: The ID of the target vertex.
//   - edgeType: The type of the edge.
//   - options: A variadic list of functions that can modify the InsertEdgeStatement.
//
// Returns:
//
//	An initialized InsertEdgeStatement configured with the provided parameters and options.
func NewInsertEdgeStatement[TVidType statement.VidType](sourceVid, targetVid TVidType, edgeType string, options ...InsertEdgeStatementOption[TVidType]) InsertEdgeStatement[TVidType] {
	statement := InsertEdgeStatement[TVidType]{
		SourceVid:   sourceVid,
		TargetVid:   targetVid,
		EdgeType:    edgeType,
		Properties:  make(map[string]interface{}),
		IfNotExists: false,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithRank sets the rank of the InsertEdgeStatement to the provided value.
func WithRank[TVidType statement.VidType](rank int) func(*InsertEdgeStatement[TVidType]) {
	return func(stmt *InsertEdgeStatement[TVidType]) {
		stmt.Rank = rank
	}
}

// WithProperties sets the properties map of the InsertEdgeStatement to the provided map.
func WithProperties[TVidType statement.VidType](properties map[string]interface{}) func(*InsertEdgeStatement[TVidType]) {
	return func(stmt *InsertEdgeStatement[TVidType]) {
		stmt.Properties = properties
	}
}

// WithIfNotExists sets the IfNotExists flag of the InsertEdgeStatement to the provided value.
func WithIfNotExists[TVidType statement.VidType](ifNotExists bool) func(*InsertEdgeStatement[TVidType]) {
	return func(stmt *InsertEdgeStatement[TVidType]) {
		stmt.IfNotExists = ifNotExists
	}
}

// GenerateInsertEdgeStatement generates a string representation of the InsertEdgeStatement.
// The function constructs the INSERT EDGE statement with the provided source and target vertex IDs,
// edge type, rank, properties, and IfNotExists flag. It encodes the vertex IDs and properties
// into the appropriate format and returns the resulting string.
func GenerateInsertEdgeStatement[TVidType statement.VidType](input InsertEdgeStatement[TVidType]) (string, error) {
	var sb strings.Builder

	if input.IfNotExists {
		sb.WriteString(`INSERT EDGE IF NOT EXISTS `)
	} else {
		sb.WriteString(`INSERT EDGE `)
	}

	sb.WriteString(fmt.Sprintf(`%s`, input.EdgeType))

	sortedProperties := make([]string, 0, len(input.Properties))
	for k, _ := range input.Properties {
		sortedProperties = append(sortedProperties, k)
	}
	sort.Strings(sortedProperties)

	if input.Properties == nil || len(input.Properties) == 0 {
		sb.WriteString(` () `)
	} else {
		sb.WriteString(` (`)
		for i, key := range sortedProperties {
			if i > 0 {
				sb.WriteString(`,`)
			}
			sb.WriteString(fmt.Sprintf(`%s`, key))
		}
		sb.WriteString(`) `)
	}

	sourceVidValue, err := statement.EncodeVidFieldValueAsStr(input.SourceVid)
	if err != nil {
		return "", err
	}

	targetVidValue, err := statement.EncodeVidFieldValueAsStr(input.TargetVid)
	if err != nil {
		return "", err
	}

	sb.WriteString(fmt.Sprintf(`VALUES %s->%s@%d:`, sourceVidValue, targetVidValue, input.Rank))

	if input.Properties == nil || len(input.Properties) == 0 {
		sb.WriteString(`()`)
	} else {
		sb.WriteString(`(`)
		for i, key := range sortedProperties {
			if i > 0 {
				sb.WriteString(`,`)
			}
			encodedVal, err := statement.EncodeNebulaFieldValue(input.Properties[key])
			if err != nil {
				return "", err
			}
			sb.WriteString(fmt.Sprintf(`%v`, encodedVal))
		}
		sb.WriteString(`)`)
	}

	sb.WriteString(`;`)

	return sb.String(), nil
}
