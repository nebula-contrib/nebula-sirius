package statement

type TestCaseGenerateInsertEdgeStatement[TVidType string | int64] struct {
	Description   string
	Given         InsertEdgeStatement[TVidType]
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateDeleteEdgeStatement[TVidType string | int64] struct {
	Description   string
	Given         DeleteEdgeStatement[TVidType]
	Expected      string
	IsErrExpected bool
}

func GetTestCasesForGenerateInsertEdgeStatementWhereVidString() []TestCaseGenerateInsertEdgeStatement[string] {
	return []TestCaseGenerateInsertEdgeStatement[string]{
		{
			Description:   "A simple insert edge statement without properties and with default settings",
			Given:         NewInsertEdgeStatement[string]("John", "Alive", "Friend"),
			Expected:      `INSERT EDGE Friend () VALUES "John"->"Alive"@0:();`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: NewInsertEdgeStatement[string]("John", "Alive", "Friend",
				WithProperties[string](map[string]interface{}{"key1": "strval1", "key2": 121}),
				WithIfNotExists[string](true),
				WithRank[string](100)),
			Expected:      `INSERT EDGE IF NOT EXISTS Friend (key1,key2) VALUES "John"->"Alive"@100:("strval1",121);`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateInsertEdgeStatementWhereVidInt64() []TestCaseGenerateInsertEdgeStatement[int64] {
	return []TestCaseGenerateInsertEdgeStatement[int64]{
		{
			Description:   "A simple insert edge statement without properties and with default settings",
			Given:         NewInsertEdgeStatement[int64](100, 200, "Friend"),
			Expected:      `INSERT EDGE Friend () VALUES 100->200@0:();`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: NewInsertEdgeStatement[int64](100, 200, "Friend",
				WithProperties[int64](map[string]interface{}{"key1": "strval1", "key2": 121}),
				WithIfNotExists[int64](true),
				WithRank[int64](100)),
			Expected:      `INSERT EDGE IF NOT EXISTS Friend (key1,key2) VALUES 100->200@100:("strval1",121);`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateDeleteEdgeStatementWhereVidString() []TestCaseGenerateDeleteEdgeStatement[string] {
	return []TestCaseGenerateDeleteEdgeStatement[string]{
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: DeleteEdgeStatement[string]{
				SourceVid: "John",
				TargetVid: "Alive",
				EdgeType:  "Friend",
			},
			Expected:      `DELETE EDGE Friend "John"->"Alive"@0;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: DeleteEdgeStatement[string]{
				SourceVid: "John",
				TargetVid: "Alive",
				EdgeType:  "Friend",
				Rank:      99,
			},
			Expected:      `DELETE EDGE Friend "John"->"Alive"@99;`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateDeleteEdgeStatementWhereVidInt64() []TestCaseGenerateDeleteEdgeStatement[int64] {
	return []TestCaseGenerateDeleteEdgeStatement[int64]{
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: DeleteEdgeStatement[int64]{
				SourceVid: 100,
				TargetVid: 200,
				EdgeType:  "Friend",
			},
			Expected:      `DELETE EDGE Friend 100->200@0;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: DeleteEdgeStatement[int64]{
				SourceVid: 100,
				TargetVid: 200,
				EdgeType:  "Friend",
				Rank:      99,
			},
			Expected:      `DELETE EDGE Friend 100->200@99;`,
			IsErrExpected: false,
		},
	}
}
