package statement

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
)

// IInsertableVertex is an interface that must be implemented by all struct intended to store vertex information
// that are used to generate INSERT VERTEX scripts
type IInsertableVertex interface {
	GetTagName() string
	InsertIfNotExists() bool
}

// nebulaInfoPerStruct is a struct that stores the Nebula fields and the corresponding struct fields
type nebulaInfoPerStruct struct {
	NebulaFieldAndStructFieldMap map[string]reflect.StructField
	NebulaFields                 []string
	VidStructField               reflect.StructField
}

var cachedNebulaInfoPerStruct sync.Map

// GenerateInsertVertexStatement takes a slice of struct vertices and generates the corresponding
// INSERT VERTEX scripts
//
// The struct must have a field tagged with the "nebula_vid" tag, which is used as the
// vertex ID.
//
// The struct can also have fields tagged with the "nebula_field" and "nebula_field_type" tags.
// The "nebula_field" tag is used to specify the Nebula field name, and the "nebula_field_type"
// tag is used to specify the Nebula field type obviously preventing inferring golang type.
//
// Currently supported Nebula field types are:
// - date
// - time
// - datetime
// - timestamp
// - geography
// - duration
// For other nebula field types, it will infer the golang type and convert it to the appropriate nebula type.
// The function returns a string containing the INSERT VERTEX scripts separated by semicolons.
// If an error occurs, the function returns an empty string and the error.
func GenerateInsertVertexStatement(vertices []IInsertableVertex) (string, error) {
	if len(vertices) == 0 {
		return "", fmt.Errorf("no vertices provided")
	}

	var sb strings.Builder

	for i, vertex := range vertices {
		// raise error if vertex is nil
		if vertex == nil {
			return "", fmt.Errorf("vertex is nil")
		}

		v := reflect.ValueOf(vertex)

		if v.Elem().Kind() == reflect.Struct {
			v = v.Elem()
		} else {
			return "", fmt.Errorf(fmt.Sprintf("vertex is not a struct: %v", v.Kind()))
		}

		// Get the type of the first vertex to extract the struct name
		vertexType := v.Type()

		// Get the type of the first vertex to extract the struct name
		vertexTypeName := vertex.GetTagName()
		if vertex.InsertIfNotExists() {
			sb.WriteString(fmt.Sprintf("INSERT VERTEX IF NOT EXISTS %s ", vertexTypeName))
		} else {
			sb.WriteString(fmt.Sprintf("INSERT VERTEX %s ", vertexTypeName))
		}

		nebulaInfoPerStructt, err := readThroughCache(vertexTypeName, vertexType)
		if err != nil {
			return "", err
		}

		nebulaFields := nebulaInfoPerStructt.NebulaFields
		nebulaVidStructField := nebulaInfoPerStructt.VidStructField
		nebulaFieldAndStructFieldMap := nebulaInfoPerStructt.NebulaFieldAndStructFieldMap

		availableNebulaFields := make([]string, 0)
		for _, nebulaField := range nebulaFields {
			structField := nebulaFieldAndStructFieldMap[nebulaField]
			structFieldVal := reflect.Indirect(v).FieldByName(structField.Name)

			if structFieldVal.Kind() == reflect.Pointer {
				structFieldVal = structFieldVal.Elem()
			}

			if structFieldVal.IsValid() {
				availableNebulaFields = append(availableNebulaFields, nebulaField)
			}
		}

		sb.WriteString("(" + strings.Join(availableNebulaFields, ", ") + ") VALUES ")

		vidField := reflect.Indirect(v).FieldByName(nebulaVidStructField.Name)

		vidFieldValue := ""
		if vidField.Kind() == reflect.Pointer {
			vidField = vidField.Elem()
		}
		switch vidField.Kind() {
		case reflect.String:
			vidFieldValue = fmt.Sprintf(`"%v"`, vidField)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vidFieldValue = fmt.Sprintf("%v", vidField)
		case reflect.Pointer:
			vidFieldValue = fmt.Sprintf("%v", vidField.Elem())
		default:
			return "", fmt.Errorf(fmt.Sprintf("`%s` tagged field of struct is either nil or not supported", VID_GO_TAG))
		}

		var values []string
		for _, nebulaField := range availableNebulaFields {
			structField, found := nebulaFieldAndStructFieldMap[nebulaField]
			if !found {
				return "", fmt.Errorf("field not found")
			}

			structFieldVal := reflect.Indirect(v).FieldByName(structField.Name)

			if structFieldVal.Kind() == reflect.Pointer {
				structFieldVal = structFieldVal.Elem()
			}

			switch structFieldVal.Kind() {
			case reflect.String:
				{
					switch structField.Tag.Get(NEBULA_FIELD_TYPE_GO_TAG) {
					case "date":
						//		 date("2025-02-15"),
						values = append(values, fmt.Sprintf("date(\"%v\")", structFieldVal))
					case "time":
						//		  time("14:30:00"),
						values = append(values, fmt.Sprintf("time(\"%v\")", structFieldVal))
					case "datetime":
						//        datetime("2017-03-04T22:30:40.003000[Asia/Shanghai]"),
						values = append(values, fmt.Sprintf("datetime(\"%v\")", structFieldVal))
					case "timestamp":
						//        timestamp("1988-03-01T08:00:00"),
						values = append(values, fmt.Sprintf("timestamp(\"%v\")", structFieldVal))
					case "geography":
						//        ST_GeogFromText("POINT(1 1)"),
						values = append(values, fmt.Sprintf("ST_GeogFromText(\"%v\")", structFieldVal))
					case "duration":
						//        duration({years: 12, days: 14, hours: 99, minutes: 12})
						values = append(values, fmt.Sprintf("duration(%v)", structFieldVal))
					default:
						values = append(values, fmt.Sprintf("\"%v\"", structFieldVal))
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				values = append(values, fmt.Sprintf("%v", structFieldVal))
			case reflect.Float32, reflect.Float64:
				values = append(values, fmt.Sprintf("%v", structFieldVal))
			case reflect.Bool:
				values = append(values, fmt.Sprintf("%v", structFieldVal))
			default:
				return "", fmt.Errorf(fmt.Sprintf("field type not supported: %v", structFieldVal.Kind()))
			}

		}

		sb.WriteString(fmt.Sprintf("%v:(%s)", vidFieldValue, strings.Join(values, ", ")))

		if i < len(vertices)-1 {
			sb.WriteString("; ")
		}

	}

	sb.WriteString(";")

	return sb.String(), nil
}

// GenerateBatchedInsertVertexStatements takes a slice of struct vertices and generates the corresponding
// INSERT VERTEX scripts separated by semicolons. The function takes an additional parameter batchSize
// which specifies the number of vertices to process in each batch.
func GenerateBatchedInsertVertexStatements(vertices []IInsertableVertex, batchSize int) ([]string, error) {
	scripts := make([]string, 0)
	for i := 0; i < len(vertices); i = i + batchSize {
		st := i
		end := i + batchSize
		if end > len(vertices) {
			end = len(vertices)
		}

		script, err := GenerateInsertVertexStatement(vertices[st:end])
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, script)
	}

	return scripts, nil
}

func readThroughCache(structName string, vertexType reflect.Type) (nebulaInfoPerStruct, error) {
	result, ok := cachedNebulaInfoPerStruct.Load(structName)

	if !ok {
		var nebulaVidStructField reflect.StructField
		nebulaFieldAndStructFieldMap := make(map[string]reflect.StructField)

		// Collect struct fields corresponding nebula fields into map
		var nebulaFields []string

		fieldCount := vertexType.NumField()
		for i := 0; i < fieldCount; i++ {
			structField := vertexType.Field(i)
			nebulaField := structField.Tag.Get(NEBULA_FIELD_GO_TAG)
			if nebulaField != "" {
				nebulaFieldAndStructFieldMap[nebulaField] = structField
			}

			// Is it VID field?
			if structField.Tag.Get(VID_GO_TAG) != "" {
				nebulaVidStructField = structField
			}
		}
		nebulaFields = slices.Collect(maps.Keys(nebulaFieldAndStructFieldMap))
		sort.Strings(nebulaFields)

		cachedNebulaInfoPerStruct.Store(structName, nebulaInfoPerStruct{
			NebulaFieldAndStructFieldMap: nebulaFieldAndStructFieldMap,
			NebulaFields:                 nebulaFields,
			VidStructField:               nebulaVidStructField,
		})
	}

	// try to read second time
	result, ok = cachedNebulaInfoPerStruct.Load(structName)
	if !ok {
		return nebulaInfoPerStruct{}, fmt.Errorf(fmt.Sprintf("struct not found: %v in the cache", structName))
	}

	return result.(nebulaInfoPerStruct), nil
}
