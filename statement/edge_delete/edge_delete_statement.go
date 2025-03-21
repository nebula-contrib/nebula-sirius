package edge_delete

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

type DeleteEdgeStatement[TVidType statement.VidType] struct {
	SourceVid TVidType
	TargetVid TVidType
	Rank      int
	EdgeType  string
}

// GenerateDeleteEdgeStatement takes a struct DeleteEdgeStatement and generates the corresponding
// to DELETE EDGE scripts
func GenerateDeleteEdgeStatement[TVidType statement.VidType](input DeleteEdgeStatement[TVidType]) (string, error) {
	var sb strings.Builder

	sourceVidValue, err := statement.EncodeVidFieldValueAsStr(input.SourceVid)
	if err != nil {
		return "", err
	}

	targetVidValue, err := statement.EncodeVidFieldValueAsStr(input.TargetVid)
	if err != nil {
		return "", err
	}

	sb.WriteString(fmt.Sprintf(`DELETE EDGE %s %v->%v@%d;`, input.EdgeType, sourceVidValue, targetVidValue, input.Rank))

	return sb.String(), nil
}
