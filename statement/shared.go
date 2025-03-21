package statement

import "fmt"

type VidType interface {
	string | int64
}

// Go tags representing Nebula related customizations
const (
	VID_GO_TAG               = "nebula_vid"
	NEBULA_FIELD_GO_TAG      = "nebula_field"
	NEBULA_FIELD_TYPE_GO_TAG = "nebula_field_type"
)

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
