package edge_insert

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"sort"
	"strings"
)

type InsertEdgeStatement[TVidType statement.VidType] struct {
	srcVid      TVidType               // required
	dstVid      TVidType               // required
	properties  map[string]interface{} // optional
	rank        int                    // optional
	edgeType    string                 // required
	ifNotExists bool                   // optional
}

type InsertEdgeStatementOption[TVidType statement.VidType] func(*InsertEdgeStatement[TVidType])

func (s InsertEdgeStatement[TVidType]) GetSrcVid() TVidType {
	return s.srcVid
}

func (s InsertEdgeStatement[TVidType]) GetDstVid() TVidType {
	return s.dstVid
}

func (s InsertEdgeStatement[TVidType]) GetOperationType() statement.OperationTypeStatement {
	return statement.StatementInsert
}

func (s InsertEdgeStatement[TVidType]) GenerateStatement() (string, error) {
	return GenerateInsertEdgeStatement(s)
}

// NewInsertEdgeStatement creates a new InsertEdgeStatement with the given source
// vertex ID, target vertex ID, and edge type. It also allows for additional configuration
// through a variadic list of options. The function initializes the properties map
// and sets the ifNotExists flag to false by default. It applies each provided option
// to the statement before returning it.
//
// Parameters:
//   - srcVid: The ID of the source vertex.
//   - dstVid: The ID of the target vertex.
//   - edgeType: The type of the edge.
//   - options: A variadic list of functions that can modify the InsertEdgeStatement.
//
// Returns:
//
//	An initialized InsertEdgeStatement configured with the provided parameters and options.
func NewInsertEdgeStatement[TVidType statement.VidType](srcVid, dstVid TVidType, edgeType string, options ...InsertEdgeStatementOption[TVidType]) InsertEdgeStatement[TVidType] {
	statement := InsertEdgeStatement[TVidType]{
		srcVid:   srcVid,
		dstVid:   dstVid,
		edgeType: edgeType,
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
		stmt.rank = rank
	}
}

// WithProperties sets the properties map of the InsertEdgeStatement to the provided map.
func WithProperties[TVidType statement.VidType](properties map[string]interface{}) func(*InsertEdgeStatement[TVidType]) {
	return func(stmt *InsertEdgeStatement[TVidType]) {
		stmt.properties = properties
	}
}

// WithIfNotExists sets the ifNotExists flag of the InsertEdgeStatement to the provided value.
func WithIfNotExists[TVidType statement.VidType](ifNotExists bool) func(*InsertEdgeStatement[TVidType]) {
	return func(stmt *InsertEdgeStatement[TVidType]) {
		stmt.ifNotExists = ifNotExists
	}
}

// GenerateInsertEdgeStatement generates a string representation of the InsertEdgeStatement.
// The function constructs the INSERT EDGE statement with the provided source and target vertex IDs,
// edge type, rank, properties, and ifNotExists flag. It encodes the vertex IDs and properties
// into the appropriate format and returns the resulting string.
func GenerateInsertEdgeStatement[TVidType statement.VidType](input InsertEdgeStatement[TVidType]) (string, error) {
	var sb strings.Builder

	if input.ifNotExists {
		sb.WriteString(`INSERT EDGE IF NOT EXISTS `)
	} else {
		sb.WriteString(`INSERT EDGE `)
	}

	sb.WriteString(input.edgeType)

	sortedProperties := make([]string, 0, len(input.properties))
	for k := range input.properties {
		sortedProperties = append(sortedProperties, k)
	}
	sort.Strings(sortedProperties)

	if len(input.properties) == 0 {
		sb.WriteString(` () `)
	} else {
		sb.WriteString(` (`)
		for i, key := range sortedProperties {
			if i > 0 {
				sb.WriteString(`,`)
			}
			sb.WriteString(key)
		}
		sb.WriteString(`) `)
	}

	sourceVidValue, err := statement.EncodeVidFieldValueAsStr(input.srcVid)
	if err != nil {
		return "", err
	}

	targetVidValue, err := statement.EncodeVidFieldValueAsStr(input.dstVid)
	if err != nil {
		return "", err
	}

	sb.WriteString(fmt.Sprintf(`VALUES %s->%s@%d:`, sourceVidValue, targetVidValue, input.rank))

	if len(input.properties) == 0 {
		sb.WriteString(`()`)
	} else {
		sb.WriteString(`(`)
		for i, key := range sortedProperties {
			if i > 0 {
				sb.WriteString(`,`)
			}
			encodedVal, err := statement.EncodeNebulaFieldValue(input.properties[key])
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
