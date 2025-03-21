package edge_upsert

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"sort"
	"strings"
)

// UpsertEdgeStatement represents a UPSERT EDGE statement in Nebula Graph.
type UpsertEdgeStatement[TVidType statement.VidType] struct {
	edgeType   string                 // required
	srcVid     TVidType               // required
	dstVid     TVidType               // required
	rank       int                    // optional
	updateProp map[string]interface{} //required
	condition  string                 //optional
	yield      string                 //optional
}

// UpsertEdgeStatementOption is a function that configures an UpsertEdgeStatement.
type UpsertEdgeStatementOption[TVidType statement.VidType] func(*UpsertEdgeStatement[TVidType])

func (s UpsertEdgeStatement[TVidType]) GetSrcVid() TVidType {
	return s.srcVid
}

func (s UpsertEdgeStatement[TVidType]) GetDstVid() TVidType {
	return s.dstVid
}

func (s UpsertEdgeStatement[TVidType]) GetOperationType() statement.OperationTypeStatement {
	return statement.UpsertStatement
}

func (s UpsertEdgeStatement[TVidType]) GenerateStatement() (string, error) {
	return GenerateUpsertEdgeStatement(s)
}

// NewUpsertEdgeStatement creates a new UpsertEdgeStatement with the given source
// vertex ID, target vertex ID, edge type, and properties. It also allows for additional configuration
// through a variadic list of options. The function initializes the properties map
// and sets the rank to 0 by default. It applies each provided option
// to the statement before returning it.
//
// Parameters:
//   - edgeType: The type of the edge.
//   - srcVid: The ID of the source vertex.
//   - dstVid: The ID of the target vertex.
//   - updateProp: The properties to update.
//   - options: A variadic list of functions that can modify the UpsertEdgeStatement.
func NewUpsertEdgeStatement[TVidType statement.VidType](edgeType string, srcVid TVidType, dstVid TVidType, updateProp map[string]interface{}, options ...UpsertEdgeStatementOption[TVidType]) UpsertEdgeStatement[TVidType] {
	if updateProp == nil {
		updateProp = make(map[string]interface{})
	}

	statement := UpsertEdgeStatement[TVidType]{
		edgeType:   edgeType,
		srcVid:     srcVid,
		dstVid:     dstVid,
		updateProp: updateProp,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithWhen sets the rank of the UpsertEdgeStatement to the provided value.
func WithWhen[TVidType statement.VidType](condition string) func(*UpsertEdgeStatement[TVidType]) {
	return func(stmt *UpsertEdgeStatement[TVidType]) {
		stmt.condition = condition
	}
}

// WithYield sets the rank of the UpsertEdgeStatement to the provided value.
func WithYield[TVidType statement.VidType](yield string) func(*UpsertEdgeStatement[TVidType]) {
	return func(stmt *UpsertEdgeStatement[TVidType]) {
		stmt.yield = yield
	}
}

// WithRank sets the rank of the UpsertEdgeStatement to the provided value.
func WithRank[TVidType statement.VidType](rank int) func(*UpsertEdgeStatement[TVidType]) {
	return func(stmt *UpsertEdgeStatement[TVidType]) {
		stmt.rank = rank
	}
}

// GenerateUpsertEdgeStatement generates a string representation of the UpsertEdgeStatement.
// The function constructs the UPSERT EDGE statement with the provided source and target vertex IDs,
// edge type, rank, properties, and condition and yield flag. It encodes the vertex IDs, properties,
// condition and yield into the appropriate format and returns the resulting string.
func GenerateUpsertEdgeStatement[TVidType statement.VidType](input UpsertEdgeStatement[TVidType]) (string, error) {

	if len(input.updateProp) == 0 {
		return "", fmt.Errorf("update properties are required")
	}

	var sb strings.Builder

	sb.WriteString(`UPSERT EDGE ON `)
	sb.WriteString(input.edgeType)
	sb.WriteString(` `)
	srcVidValue, err := statement.EncodeVidFieldValueAsStr(input.srcVid)
	if err != nil {
		return "", err
	}
	dstVidValue, err := statement.EncodeVidFieldValueAsStr(input.dstVid)
	if err != nil {
		return "", err
	}
	sb.WriteString(fmt.Sprintf(`%s->%s@%d`, srcVidValue, dstVidValue, input.rank))
	sb.WriteString(` SET `)

	sortedPropertKeys := make([]string, 0, len(input.updateProp))
	for k := range input.updateProp {
		sortedPropertKeys = append(sortedPropertKeys, k)
	}
	sort.Strings(sortedPropertKeys)

	firstProp := true
	for _, k := range sortedPropertKeys {
		if !firstProp {
			sb.WriteString(`, `)
		}
		sb.WriteString(k)
		sb.WriteString(`=`)

		v := input.updateProp[k]
		val, err := statement.EncodeNebulaFieldValue(v)
		if err != nil {
			return "", err
		}
		sb.WriteString(val)
		firstProp = false
	}

	if input.condition != "" {
		sb.WriteString(` WHEN `)
		sb.WriteString(input.condition)
	}

	if input.yield != "" {
		sb.WriteString(` YIELD `)
		sb.WriteString(input.yield)
	}

	sb.WriteString(`;`)
	return sb.String(), nil
}
