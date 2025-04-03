package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/vertex_delete"
	"github.com/nebula-contrib/nebula-sirius/statement/vertex_insert"
)

type PersonV2 struct {
	Vid              string  `nebula_vid:"true"`
	f_bool           bool    `nebula_field:"f_bool"`
	f_int64          int64   `nebula_field:"f_int64"`
	f_string         string  `nebula_field:"f_string"`
	f_fixed_string   string  `nebula_field:"f_fixed_string"`
	f_double         float64 `nebula_field:"f_double"`
	f_int32          int32   `nebula_field:"f_int32"`
	f_int16          int16   `nebula_field:"f_int16"`
	f_int8           int8    `nebula_field:"f_int8"`
	f_date           string  `nebula_field:"f_date" nebula_field_type:"date"`
	f_time           string  `nebula_field:"f_time" nebula_field_type:"time"`
	f_datetime       string  `nebula_field:"f_datetime" nebula_field_type:"datetime"`
	t_ts             string  `nebula_field:"f_ts" nebula_field_type:"timestamp"`
	f_geo            string  `nebula_field:"f_geo" nebula_field_type:"geography"`
	f_geo_linestring string  `nebula_field:"f_geo_linestring" nebula_field_type:"geography"`
	f_geo_polygon    string  `nebula_field:"f_geo_polygon" nebula_field_type:"geography"`
	f_duration       string  `nebula_field:"f_duration" nebula_field_type:"duration"`
}

func (p *PersonV2) GetTagName() string {
	return "PersonV2"
}

func (p *PersonV2) InsertIfNotExists() bool {
	return false
}

type PersonV1 struct {
	Vid              *string  `nebula_vid:"true"`
	f_bool           *bool    `nebula_field:"f_bool"`
	f_int64          *int64   `nebula_field:"f_int64"`
	f_string         *string  `nebula_field:"f_string"`
	f_fixed_string   *string  `nebula_field:"f_fixed_string"`
	f_double         *float64 `nebula_field:"f_double"`
	f_int32          *int32   `nebula_field:"f_int32"`
	f_int16          *int16   `nebula_field:"f_int16"`
	f_int8           *int8    `nebula_field:"f_int8"`
	f_date           *string  `nebula_field:"f_date" nebula_field_type:"date"`
	f_time           *string  `nebula_field:"f_time" nebula_field_type:"time"`
	f_datetime       *string  `nebula_field:"f_datetime" nebula_field_type:"datetime"`
	t_ts             *string  `nebula_field:"f_ts" nebula_field_type:"timestamp"`
	f_geo            *string  `nebula_field:"f_geo" nebula_field_type:"geography"`
	f_geo_linestring *string  `nebula_field:"f_geo_linestring" nebula_field_type:"geography"`
	f_geo_polygon    *string  `nebula_field:"f_geo_polygon" nebula_field_type:"geography"`
	f_duration       *string  `nebula_field:"f_duration" nebula_field_type:"duration"`
}

func (p *PersonV1) GetTagName() string {
	return "PersonV1"
}

func (p *PersonV1) InsertIfNotExists() bool {
	return false
}

type PersonV3 struct {
	Vid *string `nebula_vid:"true"`
}

func (p *PersonV3) GetTagName() string {
	return "PersonV3"
}

func (p *PersonV3) InsertIfNotExists() bool {
	return true
}

var (
	vid            = "4001"
	fBool          = true
	fInt64         = int64(1234567890)
	fString        = "your text here"
	fFixedString   = "fix"
	fDouble        = 3.1213
	fInt32         = int32(23125425)
	fInt16         = int16(23767)
	fInt8          = int8(127)
	fDate          = "2020-01-01"
	fTime          = "14:30:00"
	fDatetime      = "2017-03-04T22:30:40.003000[Asia/Shanghai]"
	tTs            = "1988-03-01T08:00:00"
	fGeo           = "POINT(1 1)"
	fGeoLineString = "LINESTRING(0 0, 1 1, 2 2)"
	fGeoPolygon    = "POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))"
	fDuration      = "{years: 12, days: 14, hours: 99, minutes: 12}"
)

type TestCaseGenerateBatchedInsertVertexStatements struct {
	Description        string
	GivenVerticesArray []vertex_insert.IInsertableVertex
	GivenBatchSize     int
	Expected           []string
	IsErrExpected      bool
}

type TestCaseGenerateInsertVertexStatement struct {
	Description        string
	GivenVerticesArray []vertex_insert.IInsertableVertex
	Expected           string
	IsErrExpected      bool
}

type TestCaseGenerateDeleteVertexStatement[Tvid string | int64] struct {
	Description   string
	Given         vertex_delete.DeleteVertexStatement[Tvid]
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateBatchedDeleteVertexStatement[Tvid string | int64] struct {
	Description    string
	Given          vertex_delete.DeleteVertexStatement[Tvid]
	GivenBatchSize int
	Expected       []string
	IsErrExpected  bool
}

func GetTestCasesForBatchedDeleteVertexStatementWhereVidString() []TestCaseGenerateBatchedDeleteVertexStatement[string] {
	return []TestCaseGenerateBatchedDeleteVertexStatement[string]{
		{
			Description: "Given Delete Statement with batch, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[string]{
				VertexIds: []string{"1001", "1002", "1003"},
				WithEdge:  false,
			},
			GivenBatchSize: 1,
			Expected:       []string{`DELETE VERTEX "1001";`, `DELETE VERTEX "1002";`, `DELETE VERTEX "1003";`},
			IsErrExpected:  false,
		},
		{
			Description: "Given Delete Statement with batch, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[string]{
				VertexIds: []string{"1001", "1002", "1003"},
				WithEdge:  true,
			},
			GivenBatchSize: 2,
			Expected:       []string{`DELETE VERTEX "1001", "1002" WITH EDGE;`, `DELETE VERTEX "1003" WITH EDGE;`},
			IsErrExpected:  false,
		},
		{
			Description: "Given Delete Statement with batch, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[string]{
				VertexIds: []string{"1001", "1002", "1003"},
				WithEdge:  true,
			},
			GivenBatchSize: 3,
			Expected:       []string{`DELETE VERTEX "1001", "1002", "1003" WITH EDGE;`},
			IsErrExpected:  false,
		},
	}
}

func GetTestCasesForBatchedDeleteVertexStatementWhereVidInt64() []TestCaseGenerateBatchedDeleteVertexStatement[int64] {
	return []TestCaseGenerateBatchedDeleteVertexStatement[int64]{
		{
			Description: "Given Delete Statement with batch, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[int64]{
				VertexIds: []int64{1001, 1002, 1003},
				WithEdge:  false,
			},
			GivenBatchSize: 1,
			Expected:       []string{`DELETE VERTEX 1001;`, `DELETE VERTEX 1002;`, `DELETE VERTEX 1003;`},
			IsErrExpected:  false,
		},
		{
			Description: "Given Delete Statement with batch, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[int64]{
				VertexIds: []int64{1001, 1002, 1003},
				WithEdge:  true,
			},
			GivenBatchSize: 2,
			Expected:       []string{`DELETE VERTEX 1001, 1002 WITH EDGE;`, `DELETE VERTEX 1003 WITH EDGE;`},
			IsErrExpected:  false,
		},
		{
			Description: "Given Delete Statement with batch, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[int64]{
				VertexIds: []int64{1001, 1002, 1003},
				WithEdge:  true,
			},
			GivenBatchSize: 3,
			Expected:       []string{`DELETE VERTEX 1001, 1002, 1003 WITH EDGE;`},
			IsErrExpected:  false,
		},
	}
}

func GetTestCasesForDeleteVertexStatementWhereVidString() []TestCaseGenerateDeleteVertexStatement[string] {
	return []TestCaseGenerateDeleteVertexStatement[string]{
		{
			Description: "Given Delete Statement, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[string]{
				VertexIds: []string{"1001"},
				WithEdge:  false,
			},
			Expected:      `DELETE VERTEX "1001";`,
			IsErrExpected: false,
		},
		{
			Description: "Given Delete Statement, prepare delete scripts where vid are string and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[string]{
				VertexIds: []string{"1001", "1002"},
				WithEdge:  false,
			},
			Expected:      `DELETE VERTEX "1001", "1002";`,
			IsErrExpected: false,
		},
		{
			Description: "Given Delete Statement, prepare delete scripts where vid are string AND WithEdge true",
			Given: vertex_delete.DeleteVertexStatement[string]{
				VertexIds: []string{"1012", "1022"},
				WithEdge:  true,
			},
			Expected:      `DELETE VERTEX "1012", "1022" WITH EDGE;`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForDeleteVertexStatementWhereVidInt64() []TestCaseGenerateDeleteVertexStatement[int64] {
	return []TestCaseGenerateDeleteVertexStatement[int64]{
		{
			Description: "Given Delete Statement, prepare delete scripts where vid are int64 and WithEdge false",
			Given: vertex_delete.DeleteVertexStatement[int64]{
				VertexIds: []int64{1001, 1002},
				WithEdge:  false,
			},
			Expected:      `DELETE VERTEX 1001, 1002;`,
			IsErrExpected: false,
		},
		{
			Description: "Given Delete Statement, prepare delete scripts where vid are int64 and WithEdge true",
			Given: vertex_delete.DeleteVertexStatement[int64]{
				VertexIds: []int64{1001, 1002},
				WithEdge:  true,
			},
			Expected:      `DELETE VERTEX 1001, 1002 WITH EDGE;`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateInsertVertexStatement() []TestCaseGenerateInsertVertexStatement {
	return []TestCaseGenerateInsertVertexStatement{
		{
			Description: "Given Struct with all reference fields, expect insert script containing all fields",
			GivenVerticesArray: []vertex_insert.IInsertableVertex{
				&PersonV1{
					Vid:              &vid,
					f_bool:           &fBool,
					f_int64:          &fInt64,
					f_string:         &fString,
					f_fixed_string:   &fFixedString,
					f_double:         &fDouble,
					f_int32:          &fInt32,
					f_int16:          &fInt16,
					f_int8:           &fInt8,
					f_date:           &fDate,
					f_time:           &fTime,
					f_datetime:       &fDatetime,
					t_ts:             &tTs,
					f_geo:            &fGeo,
					f_geo_linestring: &fGeoLineString,
					f_geo_polygon:    &fGeoPolygon,
					f_duration:       &fDuration,
				},
			},
			Expected:      `INSERT VERTEX PersonV1 (f_bool, f_date, f_datetime, f_double, f_duration, f_fixed_string, f_geo, f_geo_linestring, f_geo_polygon, f_int16, f_int32, f_int64, f_int8, f_string, f_time, f_ts) VALUES "4001":(true, date("2020-01-01"), datetime("2017-03-04T22:30:40.003000[Asia/Shanghai]"), 3.1213, duration({years: 12, days: 14, hours: 99, minutes: 12}), "fix", ST_GeogFromText("POINT(1 1)"), ST_GeogFromText("LINESTRING(0 0, 1 1, 2 2)"), ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))"), 23767, 23125425, 1234567890, 127, "your text here", time("14:30:00"), timestamp("1988-03-01T08:00:00"));`,
			IsErrExpected: false,
		},
		{
			Description: "Given Struct with all scalar fields, expect insert script containing all fields",
			GivenVerticesArray: []vertex_insert.IInsertableVertex{
				&PersonV2{
					Vid:              "4001",
					f_bool:           true,
					f_int64:          1234567890,
					f_string:         "your text here",
					f_fixed_string:   "fix",
					f_double:         3.1213,
					f_int32:          23125425,
					f_int16:          23767,
					f_int8:           127,
					f_date:           "2020-01-01",
					f_time:           "14:30:00",
					f_datetime:       "2017-03-04T22:30:40.003000[Asia/Shanghai]",
					t_ts:             "1988-03-01T08:00:00",
					f_geo:            "POINT(1 1)",
					f_geo_linestring: "LINESTRING(0 0, 1 1, 2 2)",
					f_geo_polygon:    "POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))",
					f_duration:       "{years: 12, days: 14, hours: 99, minutes: 12}",
				},
			},
			Expected:      `INSERT VERTEX PersonV2 (f_bool, f_date, f_datetime, f_double, f_duration, f_fixed_string, f_geo, f_geo_linestring, f_geo_polygon, f_int16, f_int32, f_int64, f_int8, f_string, f_time, f_ts) VALUES "4001":(true, date("2020-01-01"), datetime("2017-03-04T22:30:40.003000[Asia/Shanghai]"), 3.1213, duration({years: 12, days: 14, hours: 99, minutes: 12}), "fix", ST_GeogFromText("POINT(1 1)"), ST_GeogFromText("LINESTRING(0 0, 1 1, 2 2)"), ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))"), 23767, 23125425, 1234567890, 127, "your text here", time("14:30:00"), timestamp("1988-03-01T08:00:00"));`,
			IsErrExpected: false,
		},
		{
			Description: "Given Struct with partially reference fields, expect insert script containing non-nil fields",
			GivenVerticesArray: []vertex_insert.IInsertableVertex{
				&PersonV1{
					Vid:           &vid,
					f_bool:        &fBool,
					f_geo:         &fGeo,
					f_geo_polygon: &fGeoPolygon,
					f_duration:    &fDuration,
				},
				&PersonV1{
					Vid:        &vid,
					f_duration: &fDuration,
				},
			},
			Expected:      `INSERT VERTEX PersonV1 (f_bool, f_duration, f_geo, f_geo_polygon) VALUES "4001":(true, duration({years: 12, days: 14, hours: 99, minutes: 12}), ST_GeogFromText("POINT(1 1)"), ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))")); INSERT VERTEX PersonV1 (f_duration) VALUES "4001":(duration({years: 12, days: 14, hours: 99, minutes: 12}));`,
			IsErrExpected: false,
		},
		{
			Description: "Given Struct with partially reference fields, expect insert script containing non-nil fields",
			GivenVerticesArray: []vertex_insert.IInsertableVertex{
				&PersonV1{
					Vid:           &vid,
					f_bool:        &fBool,
					f_geo:         &fGeo,
					f_geo_polygon: &fGeoPolygon,
					f_duration:    &fDuration,
				},
				&PersonV1{
					Vid:        &vid,
					f_duration: &fDuration,
				},
			},
			Expected:      `INSERT VERTEX PersonV1 (f_bool, f_duration, f_geo, f_geo_polygon) VALUES "4001":(true, duration({years: 12, days: 14, hours: 99, minutes: 12}), ST_GeogFromText("POINT(1 1)"), ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))")); INSERT VERTEX PersonV1 (f_duration) VALUES "4001":(duration({years: 12, days: 14, hours: 99, minutes: 12}));`,
			IsErrExpected: false,
		},
		{
			Description: "Given Struct with no fields except vid field, expect insert script with no fields and insert if exists option",
			GivenVerticesArray: []vertex_insert.IInsertableVertex{
				&PersonV3{
					Vid: &vid,
				},
			},
			Expected:      `INSERT VERTEX IF NOT EXISTS PersonV3 () VALUES "4001":();`,
			IsErrExpected: false,
		},
		{
			Description: "Given Struct with no vid field, return error",
			GivenVerticesArray: []vertex_insert.IInsertableVertex{
				&PersonV3{},
			},
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateBatchedInsertVertexStatements() []TestCaseGenerateBatchedInsertVertexStatements {
	vertices := []vertex_insert.IInsertableVertex{
		&PersonV1{
			Vid:    &vid,
			f_bool: &fBool,
		},
		&PersonV1{
			Vid:        &vid,
			f_duration: &fDuration,
		},
		&PersonV1{
			Vid:   &vid,
			f_geo: &fGeoPolygon,
		},
		&PersonV1{
			Vid:    &vid,
			f_time: &fTime,
		},
	}

	return []TestCaseGenerateBatchedInsertVertexStatements{
		{
			Description:        "Given batch size is 3 and 4 vertices, expect 2 insert scripts",
			GivenVerticesArray: vertices,
			GivenBatchSize:     3,
			Expected: []string{
				`INSERT VERTEX PersonV1 (f_bool) VALUES "4001":(true); INSERT VERTEX PersonV1 (f_duration) VALUES "4001":(duration({years: 12, days: 14, hours: 99, minutes: 12})); INSERT VERTEX PersonV1 (f_geo) VALUES "4001":(ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))"));`,
				`INSERT VERTEX PersonV1 (f_time) VALUES "4001":(time("14:30:00"));`,
			},
			IsErrExpected: false,
		},
		{
			Description:        "Given batch size is 4 and 4 vertices, expect 1 insert scripts",
			GivenVerticesArray: vertices,
			GivenBatchSize:     4,
			Expected: []string{
				`INSERT VERTEX PersonV1 (f_bool) VALUES "4001":(true); INSERT VERTEX PersonV1 (f_duration) VALUES "4001":(duration({years: 12, days: 14, hours: 99, minutes: 12})); INSERT VERTEX PersonV1 (f_geo) VALUES "4001":(ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))")); INSERT VERTEX PersonV1 (f_time) VALUES "4001":(time("14:30:00"));`,
			},
			IsErrExpected: false,
		},
		{
			Description:        "Given batch size is 1 and 4 vertices, expect 4 insert scripts",
			GivenVerticesArray: vertices,
			GivenBatchSize:     1,
			Expected: []string{
				`INSERT VERTEX PersonV1 (f_bool) VALUES "4001":(true);`,
				`INSERT VERTEX PersonV1 (f_duration) VALUES "4001":(duration({years: 12, days: 14, hours: 99, minutes: 12}));`,
				`INSERT VERTEX PersonV1 (f_geo) VALUES "4001":(ST_GeogFromText("POLYGON((0 0, 0 1, 1 1, 1 0, 0 0))"));`,
				`INSERT VERTEX PersonV1 (f_time) VALUES "4001":(time("14:30:00"));`,
			},
			IsErrExpected: false,
		},
	}
}
