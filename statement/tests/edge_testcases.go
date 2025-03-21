package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/edge_delete"
	"github.com/nebula-contrib/nebula-sirius/statement/edge_insert"
	"github.com/nebula-contrib/nebula-sirius/statement/edge_upsert"
)

type TestCaseGenerateInsertEdgeStatement[TVidType string | int64] struct {
	Description   string
	Given         edge_insert.InsertEdgeStatement[TVidType]
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateDeleteEdgeStatement[TVidType string | int64] struct {
	Description   string
	Given         edge_delete.DeleteEdgeStatement[TVidType]
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateUpsertEdgeStatement[TVidType string | int64] struct {
	Description   string
	Given         edge_upsert.UpsertEdgeStatement[TVidType]
	Expected      string
	IsErrExpected bool
}

func GetTestCasesForGenerateInsertEdgeStatementWhereVidString() []TestCaseGenerateInsertEdgeStatement[string] {
	return []TestCaseGenerateInsertEdgeStatement[string]{
		{
			Description:   "A simple insert edge statement without properties and with default settings",
			Given:         edge_insert.NewInsertEdgeStatement[string]("John", "Alive", "Friend"),
			Expected:      `INSERT EDGE Friend () VALUES "John"->"Alive"@0:();`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_insert.NewInsertEdgeStatement[string]("John", "Alive", "Friend",
				edge_insert.WithProperties[string](map[string]interface{}{"key1": "strval1", "key2": 121}),
				edge_insert.WithIfNotExists[string](true),
				edge_insert.WithRank[string](100)),
			Expected:      `INSERT EDGE IF NOT EXISTS Friend (key1,key2) VALUES "John"->"Alive"@100:("strval1",121);`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateInsertEdgeStatementWhereVidInt64() []TestCaseGenerateInsertEdgeStatement[int64] {
	return []TestCaseGenerateInsertEdgeStatement[int64]{
		{
			Description:   "A simple insert edge statement without properties and with default settings",
			Given:         edge_insert.NewInsertEdgeStatement[int64](100, 200, "Friend"),
			Expected:      `INSERT EDGE Friend () VALUES 100->200@0:();`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_insert.NewInsertEdgeStatement[int64](100, 200, "Friend",
				edge_insert.WithProperties[int64](map[string]interface{}{"key1": "strval1", "key2": 121}),
				edge_insert.WithIfNotExists[int64](true),
				edge_insert.WithRank[int64](100)),
			Expected:      `INSERT EDGE IF NOT EXISTS Friend (key1,key2) VALUES 100->200@100:("strval1",121);`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateDeleteEdgeStatementWhereVidString() []TestCaseGenerateDeleteEdgeStatement[string] {
	return []TestCaseGenerateDeleteEdgeStatement[string]{
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_delete.DeleteEdgeStatement[string]{
				SourceVid: "John",
				TargetVid: "Alive",
				EdgeType:  "Friend",
			},
			Expected:      `DELETE EDGE Friend "John"->"Alive"@0;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_delete.DeleteEdgeStatement[string]{
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
			Given: edge_delete.DeleteEdgeStatement[int64]{
				SourceVid: 100,
				TargetVid: 200,
				EdgeType:  "Friend",
			},
			Expected:      `DELETE EDGE Friend 100->200@0;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_delete.DeleteEdgeStatement[int64]{
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

func GetTestCasesForGenerateUpsertEdgeStatementWhereVidString() []TestCaseGenerateUpsertEdgeStatement[string] {
	return []TestCaseGenerateUpsertEdgeStatement[string]{
		{
			Description: "A simple upsert edge statement with default settings",
			Given: edge_upsert.NewUpsertEdgeStatement[string]("Friend", "John", "Alive",
				map[string]interface{}{
					"key1": "strval1",
					"key2": 121,
				}),
			Expected:      `UPSERT EDGE ON Friend "John"->"Alive"@0 SET key1="strval1", key2=121;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple upsert edge statement configuration options",
			Given: edge_upsert.NewUpsertEdgeStatement[string]("Friend", "John", "Alive",
				map[string]interface{}{
					"key1": "strval1",
					"key2": 121,
				},
				edge_upsert.WithWhen[string]("key1 == 'strval1'"),
				edge_upsert.WithRank[string](100),
				edge_upsert.WithYield[string]("key1,key2")),
			Expected:      `UPSERT EDGE ON Friend "John"->"Alive"@100 SET key1="strval1", key2=121 WHEN key1 == 'strval1' YIELD key1,key2;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple upsert edge statement with nil update properties",
			Given: edge_upsert.NewUpsertEdgeStatement[string]("Friend", "John", "Alive",
				map[string]interface{}{}),
			Expected:      "",
			IsErrExpected: true,
		},
		{
			Description: "A simple upsert edge statement without update properties",
			Given: edge_upsert.NewUpsertEdgeStatement[string]("Friend", "John", "Alive",
				nil),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateUpsertEdgeStatementWhereVidInt64() []TestCaseGenerateUpsertEdgeStatement[int64] {
	return []TestCaseGenerateUpsertEdgeStatement[int64]{
		{
			Description: "A simple upsert edge statement with default settings",
			Given: edge_upsert.NewUpsertEdgeStatement[int64]("Friend", 100, 200,
				map[string]interface{}{
					"key1": "strval1",
					"key2": 121,
				}),
			Expected:      `UPSERT EDGE ON Friend 100->200@0 SET key1="strval1", key2=121;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple upsert edge statement configuration options",
			Given: edge_upsert.NewUpsertEdgeStatement[int64]("Friend", 100, 200,
				map[string]interface{}{
					"key1": "strval1",
					"key2": 121,
				},
				edge_upsert.WithWhen[int64]("key1 == 'strval1'"),
				edge_upsert.WithRank[int64](100),
				edge_upsert.WithYield[int64]("key1,key2")),
			Expected:      `UPSERT EDGE ON Friend 100->200@100 SET key1="strval1", key2=121 WHEN key1 == 'strval1' YIELD key1,key2;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple upsert edge statement with nil update properties",
			Given: edge_upsert.NewUpsertEdgeStatement[int64]("Friend", 100, 200,
				map[string]interface{}{}),
			Expected:      "",
			IsErrExpected: true,
		},

		{
			Description: "A simple upsert edge statement without update properties",
			Given: edge_upsert.NewUpsertEdgeStatement[int64]("Friend", 100, 200,
				nil),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}
