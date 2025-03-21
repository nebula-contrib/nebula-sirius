package edge_batch

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

// GenerateBatchedEdgeStatements generates a batch of edge statements.
// The function takes a slice of IEdgeStatementOperation and a batch size as input.
// It iterates over the input slice and generates a batch of statements with the provided batch size.
// The function returns a slice of strings, where each string represents a batch of statements.

// Parameters:
//   - statements: A slice of IEdgeStatementOperation.
//   - batchSize: The size of each batch.
// Returns:
//   - A slice of strings representing the batched statements.
//   - An error if there was an issue generating the statements.

func GenerateBatchedEdgeStatements[T statement.VidType](statements []statement.IEdgeStatementOperation[T], batchSize int) ([]string, error) {
	scripts := make([]string, 0)
	for i := 0; i < len(statements); i = i + batchSize {
		st := i
		end := i + batchSize
		if end > len(statements) {
			end = len(statements)
		}

		var sb strings.Builder

		for ; st < end; st++ {
			statement, err := statements[st].GenerateStatement()
			if err != nil {
				return nil, err
			}
			sb.WriteString(statement)
		}

		scripts = append(scripts, sb.String())
	}

	return scripts, nil
}
