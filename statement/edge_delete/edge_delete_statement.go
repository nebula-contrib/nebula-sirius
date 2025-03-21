package edge_delete

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

type DeleteEdgeStatement[TVidType statement.VidType] struct {
	srcVid   TVidType // required
	dstVid   TVidType // required
	rank     int      // optional
	edgeType string   // required
}

func (s DeleteEdgeStatement[TVidType]) GetSrcVid() TVidType {
	return s.srcVid
}

func (s DeleteEdgeStatement[TVidType]) GetDstVid() TVidType {
	return s.dstVid
}

func (s DeleteEdgeStatement[TVidType]) GetOperationType() statement.OperationTypeStatement {
	return statement.DeleteStatement
}

func (s DeleteEdgeStatement[TVidType]) GenerateStatement() (string, error) {
	return GenerateDeleteEdgeStatement(s)
}

// DeleteEdgeStatementOption is a function that configures an DeleteEdgeStatement.
type DeleteEdgeStatementOption[TVidType statement.VidType] func(*DeleteEdgeStatement[TVidType])

// NewDeleteEdgeStatement creates a new DeleteEdgeStatement with the given source
// vertex ID, target vertex ID, edge type, and properties. It also allows for additional configuration
// through a variadic list of options. The function initializes the properties map
// and sets the rank to 0 by default. It applies each provided option
// to the statement before returning it.
//
// Parameters:
//   - srcVid: The ID of the source vertex.
//   - dstVid: The ID of the target vertex.
//   - edgeType: The type of the edge.
//   - options: A variadic list of functions that can modify the DeleteEdgeStatement.
func NewDeleteEdgeStatement[TVidType statement.VidType](srcVid TVidType, dstVid TVidType, edgeType string, options ...DeleteEdgeStatementOption[TVidType]) DeleteEdgeStatement[TVidType] {

	statement := DeleteEdgeStatement[TVidType]{
		edgeType: edgeType,
		srcVid:   srcVid,
		dstVid:   dstVid,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithRank sets the rank of the DeleteEdgeStatement to the provided value.
func WithRank[TVidType statement.VidType](rank int) func(*DeleteEdgeStatement[TVidType]) {
	return func(stmt *DeleteEdgeStatement[TVidType]) {
		stmt.rank = rank
	}
}

// GenerateDeleteEdgeStatement generates a statement for deleting an edge.
// It takes in a DeleteEdgeStatement and returns a string representation of the statement.
func GenerateDeleteEdgeStatement[TVidType statement.VidType](input DeleteEdgeStatement[TVidType]) (string, error) {
	var sb strings.Builder

	sourceVidValue, err := statement.EncodeVidFieldValueAsStr(input.srcVid)
	if err != nil {
		return "", err
	}

	targetVidValue, err := statement.EncodeVidFieldValueAsStr(input.dstVid)
	if err != nil {
		return "", err
	}

	sb.WriteString(fmt.Sprintf(`DELETE EDGE %s %v->%v@%d;`, input.edgeType, sourceVidValue, targetVidValue, input.rank))

	return sb.String(), nil
}
