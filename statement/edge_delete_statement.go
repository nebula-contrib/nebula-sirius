package statement

import (
	"fmt"
	"strings"
)

type DeleteEdgeStatement[TVidType VidType] struct {
	SourceVid TVidType
	TargetVid TVidType
	Rank      int
	EdgeType  string
}

// GenerateDeleteEdgeStatement takes a struct DeleteEdgeStatement and generates the corresponding
// to DELETE EDGE scripts
func GenerateDeleteEdgeStatement[TVidType VidType](statement DeleteEdgeStatement[TVidType]) (string, error) {
	var sb strings.Builder

	sourceVidValue, err := encodeVidFieldValueAsStr(statement.SourceVid)
	if err != nil {
		return "", err
	}

	targetVidValue, err := encodeVidFieldValueAsStr(statement.TargetVid)
	if err != nil {
		return "", err
	}

	sb.WriteString(fmt.Sprintf(`DELETE EDGE %s %v->%v@%d;`, statement.EdgeType, sourceVidValue, targetVidValue, statement.Rank))

	return sb.String(), nil
}
