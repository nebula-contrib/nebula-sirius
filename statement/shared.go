package statement

import "fmt"

type PropertyType string

const (
	PropertyTypeString       PropertyType = "string"
	PropertyTypeFixedStringF PropertyType = "fixed_string(%d)"

	PropertyTypeInt   PropertyType = "int"
	PropertyTypeInt64 PropertyType = "int64"
	PropertyTypeInt32 PropertyType = "int32"
	PropertyTypeInt16 PropertyType = "int16"
	PropertyTypeInt8  PropertyType = "int8"

	PropertyTypeFloat  PropertyType = "float"
	PropertyTypeDouble PropertyType = "double"

	PropertyTypeBoolean PropertyType = "bool"

	PropertyTypeDate      PropertyType = "date"
	PropertyTypeTime      PropertyType = "time"
	PropertyTypeDateTime  PropertyType = "datetime"
	PropertyTypeTimestamp PropertyType = "timestamp"
	PropertyTypeDuration  PropertyType = "duration"

	PropertyTypeGeography PropertyType = "geography"
)

type IEdgeStatementOperation[T VidType] interface {
	GetSrcVid() T
	GetDstVid() T
	GetOperationType() OperationTypeStatement
	GenerateStatement() (string, error)
}

type OperationTypeStatement int

const (
	InsertStatement OperationTypeStatement = iota
	UpdateStatement
	DeleteStatement
	UpsertStatement
)

// VidType is a type constraint for vertex ID types.
type VidType interface {
	string | int64
}

// Go tags representing Nebula related customizations
const (
	VID_GO_TAG               = "nebula_vid"
	NEBULA_FIELD_GO_TAG      = "nebula_field"
	NEBULA_FIELD_TYPE_GO_TAG = "nebula_field_type"
)

// EncodeVidFieldValueAsStr encodes a vertex ID field value into a string representation.
// It handles both string and int64 types.
// The function returns the encoded string and an error if the type is unsupported.
// Supported types:
// - string
// - int64
// The function formats the string with double quotes for string types and without quotes for int64 types.
func EncodeVidFieldValueAsStr(vidValue any) (string, error) {
	switch v := vidValue.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v), nil
	case int64:
		return fmt.Sprintf(`%v`, v), nil
	default:
		return "", fmt.Errorf("unsupported vid type: %T", v)
	}
}

// EncodeNebulaFieldValue encodes a Nebula field value into a string representation.
// It handles various types such as string, int, float, and bool.
// The function returns the encoded string and an error if the type is unsupported.
// Supported types:
// - string
// - int (and its variants)
// - float (and its variants)
// - bool
func EncodeNebulaFieldValue(fieldValue any) (string, error) {
	switch v := fieldValue.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf(`%v`, v), nil
	case float32, float64:
		return fmt.Sprintf(`%v`, v), nil
	case bool:
		return fmt.Sprintf(`%v`, v), nil
	default:
		return "", fmt.Errorf("unsupported vid type: %T", v)
	}
}
