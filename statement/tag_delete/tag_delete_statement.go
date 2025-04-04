package tag_delete

import (
	"fmt"
	"github.com/nebula-contrib/nebula-sirius/statement"
	"strings"
)

type DeleteTagStatement[TVidType statement.VidType] struct {
	nameList []string   // required
	vidList  []TVidType // required
	allTags  bool       // optional
}

// DeleteTagStatementOption is a function that configures an DeleteEdgeStatement.
type DeleteTagStatementOption[TVidType statement.VidType] func(*DeleteTagStatement[TVidType])

func NewDeleteTagStatement[TVidType statement.VidType](nameList []string, vidList []TVidType, options ...DeleteTagStatementOption[TVidType]) DeleteTagStatement[TVidType] {
	statement := DeleteTagStatement[TVidType]{
		nameList: nameList,
		vidList:  vidList,
	}

	// Apply all the functional options to configure the statement.
	for _, opt := range options {
		opt(&statement)
	}

	return statement
}

// WithAllTags sets the rank of the DeleteTagStatement to the provided value.
func WithAllTags[TVidType statement.VidType]() func(*DeleteTagStatement[TVidType]) {
	return func(stmt *DeleteTagStatement[TVidType]) {
		stmt.allTags = true
	}
}

// DELETE TAG <tag_name_list> FROM <VID_list>;
func GenerateDeleteTagStatement[TVidType statement.VidType](input DeleteTagStatement[TVidType]) (string, error) {
	if len(input.nameList) == 0 && !input.allTags {
		return "", fmt.Errorf("tag name list is required or allTags is set")
	}

	if len(input.nameList) > 0 && input.allTags {
		return "", fmt.Errorf("cannot specify tag name list when allTags is set")
	}

	if len(input.vidList) == 0 {
		return "", fmt.Errorf("vid list is required")
	}

	var sb strings.Builder
	sb.WriteString("DELETE TAG ")

	if input.allTags {
		sb.WriteString("*")
	} else {
		for i, tagName := range input.nameList {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(tagName)
		}
	}

	sb.WriteString(" FROM ")

	for i, vid := range input.vidList {
		if i > 0 {
			sb.WriteString(",")
		}
		vidStr, err := statement.EncodeVidFieldValueAsStr(vid)
		if err != nil {
			return "", err
		}
		sb.WriteString(vidStr)
	}

	sb.WriteString(";")

	query := sb.String()
	return query, nil
}
