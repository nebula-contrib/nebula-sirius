package vertex_delete

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"reflect"
	"strings"
)

// DeleteVertexStatement is a struct that stores the vertex IDs and whether to delete the edges associated to vertices
type DeleteVertexStatement[TVidType statement.VidType] struct {
	VertexIds []TVidType
	WithEdge  bool
}

// GenerateDeleteVertexStatement takes a slice of struct vertices and generates the corresponding
// to DELETE VERTEX scripts
func GenerateDeleteVertexStatement[TVidType statement.VidType](input DeleteVertexStatement[TVidType]) (string, error) {
	if len(input.VertexIds) == 0 {
		return "", fmt.Errorf("empty VertexIds provided ")
	}

	var sb strings.Builder

	sb.WriteString("DELETE VERTEX ")

	// DELETE VERTEX <vid> [, <vid> ...] [WITH EDGE];
	var vidsJoined string
	for i, item := range input.VertexIds {
		if i > 0 {
			vidsJoined += ", "
		}

		switch reflect.TypeOf(item).Kind() {
		case reflect.String:
			vidsJoined += fmt.Sprintf(`"%v"`, item)
		case reflect.Int64:
			vidsJoined += fmt.Sprintf(`%v`, item)
		default:
			return "", fmt.Errorf("unsupported type")
		}
	}

	sb.WriteString(vidsJoined)

	if input.WithEdge {
		sb.WriteString(" WITH EDGE;")
	} else {
		sb.WriteString(";")
	}

	return sb.String(), nil
}

// GenerateBatchedDeleteVertexStatements takes a slice of struct vertices and generates the corresponding
// to DELETE VERTEX scripts separated by semicolons. The function takes an additional parameter batchSize
// which specifies the number of vertices to process in each batch.
func GenerateBatchedDeleteVertexStatements[TVidType statement.VidType](statement DeleteVertexStatement[TVidType], batchSize int) ([]string, error) {
	scripts := make([]string, 0)
	for i := 0; i < len(statement.VertexIds); i = i + batchSize {
		st := i
		end := i + batchSize
		if end > len(statement.VertexIds) {
			end = len(statement.VertexIds)
		}

		newStatement := DeleteVertexStatement[TVidType]{
			VertexIds: statement.VertexIds[st:end],
			WithEdge:  statement.WithEdge,
		}

		script, err := GenerateDeleteVertexStatement(newStatement)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}

	return scripts, nil
}
