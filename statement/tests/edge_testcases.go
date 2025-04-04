package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/edge_alter"
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

type TestCaseGenerateBatchEdgeStatement[TVidType string | int64] struct {
	Description     string
	GivenStatements []statement.IEdgeStatementOperation[TVidType]
	GivenBatchSize  int
	Expected        []string
	IsErrExpected   bool
}

type TestCaseGenerateEdgeAlterStatement struct {
	Description   string
	Given         edge_alter.AlterEdgeStatement
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
			Description:   "A simple insert edge statement without properties and with default settings",
			Given:         edge_delete.NewDeleteEdgeStatement("John", "Alive", "Friend"),
			Expected:      `DELETE EDGE Friend "John"->"Alive"@0;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_delete.NewDeleteEdgeStatement("John", "Alive", "Friend",
				edge_delete.WithRank[string](99),
			),
			Expected:      `DELETE EDGE Friend "John"->"Alive"@99;`,
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateDeleteEdgeStatementWhereVidInt64() []TestCaseGenerateDeleteEdgeStatement[int64] {
	return []TestCaseGenerateDeleteEdgeStatement[int64]{
		{
			Description:   "A simple insert edge statement without properties and with default settings",
			Given:         edge_delete.NewDeleteEdgeStatement[int64](100, 200, "Friend"),
			Expected:      `DELETE EDGE Friend 100->200@0;`,
			IsErrExpected: false,
		},
		{
			Description: "A simple insert edge statement without properties and with default settings",
			Given: edge_delete.NewDeleteEdgeStatement[int64](100, 200, "Friend",
				edge_delete.WithRank[int64](99),
			),
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

func GetTestCasesForGenerateBatchEdgeStatementWhereVidString() []TestCaseGenerateBatchEdgeStatement[string] {
	return []TestCaseGenerateBatchEdgeStatement[string]{
		{
			Description: "Edge batch with 4 statements and batch size of 1",
			GivenStatements: []statement.IEdgeStatementOperation[string]{
				edge_delete.NewDeleteEdgeStatement("John", "Alive", "Friend"),
				edge_delete.NewDeleteEdgeStatement("John", "Bob", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Alive", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Bob", "Friend"),
			},
			GivenBatchSize: 1,
			Expected: []string{
				`DELETE EDGE Friend "John"->"Alive"@0;`,
				`DELETE EDGE Friend "John"->"Bob"@0;`,
				`INSERT EDGE Friend () VALUES "John"->"Alive"@0:();`,
				`INSERT EDGE Friend () VALUES "John"->"Bob"@0:();`,
			},
			IsErrExpected: false,
		},
		{
			Description: "Edge batch with 4 statements and batch size of 2",
			GivenStatements: []statement.IEdgeStatementOperation[string]{
				edge_delete.NewDeleteEdgeStatement("John", "Alive", "Friend"),
				edge_delete.NewDeleteEdgeStatement("John", "Bob", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Alive", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Bob", "Friend"),
			},
			GivenBatchSize: 2,
			Expected: []string{
				`DELETE EDGE Friend "John"->"Alive"@0;` + `DELETE EDGE Friend "John"->"Bob"@0;`,
				`INSERT EDGE Friend () VALUES "John"->"Alive"@0:();` + `INSERT EDGE Friend () VALUES "John"->"Bob"@0:();`,
			},
			IsErrExpected: false,
		},
		{
			Description: "Edge batch with 4 statements and batch size of 3",
			GivenStatements: []statement.IEdgeStatementOperation[string]{
				edge_delete.NewDeleteEdgeStatement("John", "Alive", "Friend"),
				edge_delete.NewDeleteEdgeStatement("John", "Bob", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Alive", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Bob", "Friend"),
			},
			GivenBatchSize: 3,
			Expected: []string{
				`DELETE EDGE Friend "John"->"Alive"@0;` +
					`DELETE EDGE Friend "John"->"Bob"@0;` +
					`INSERT EDGE Friend () VALUES "John"->"Alive"@0:();`,
				`INSERT EDGE Friend () VALUES "John"->"Bob"@0:();`,
			},
			IsErrExpected: false,
		},
		{
			Description: "Edge batch with 4 statements and batch size of 4",
			GivenStatements: []statement.IEdgeStatementOperation[string]{
				edge_delete.NewDeleteEdgeStatement("John", "Alive", "Friend"),
				edge_delete.NewDeleteEdgeStatement("John", "Bob", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Alive", "Friend"),
				edge_insert.NewInsertEdgeStatement("John", "Bob", "Friend"),
			},
			GivenBatchSize: 4,
			Expected: []string{
				`DELETE EDGE Friend "John"->"Alive"@0;` +
					`DELETE EDGE Friend "John"->"Bob"@0;` +
					`INSERT EDGE Friend () VALUES "John"->"Alive"@0:();` +
					`INSERT EDGE Friend () VALUES "John"->"Bob"@0:();`,
			},
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateBatchEdgeStatementWhereVidInt64() []TestCaseGenerateBatchEdgeStatement[int64] {
	return []TestCaseGenerateBatchEdgeStatement[int64]{
		{
			Description: "Edge batch with 4 statements and batch size of 1",
			GivenStatements: []statement.IEdgeStatementOperation[int64]{
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(200), "Friend"),
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(300), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(200), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(300), "Friend"),
			},
			GivenBatchSize: 1,
			Expected: []string{
				`DELETE EDGE Friend 100->200@0;`,
				`DELETE EDGE Friend 100->300@0;`,
				`INSERT EDGE Friend () VALUES 100->200@0:();`,
				`INSERT EDGE Friend () VALUES 100->300@0:();`,
			},
			IsErrExpected: false,
		},
		{
			Description: "Edge batch with 4 statements and batch size of 2",
			GivenStatements: []statement.IEdgeStatementOperation[int64]{
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(200), "Friend"),
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(300), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(200), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(300), "Friend"),
			},
			GivenBatchSize: 2,
			Expected: []string{
				`DELETE EDGE Friend 100->200@0;` +
					`DELETE EDGE Friend 100->300@0;`,
				`INSERT EDGE Friend () VALUES 100->200@0:();` +
					`INSERT EDGE Friend () VALUES 100->300@0:();`,
			},
			IsErrExpected: false,
		},
		{
			Description: "Edge batch with 4 statements and batch size of 3",
			GivenStatements: []statement.IEdgeStatementOperation[int64]{
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(200), "Friend"),
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(300), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(200), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(300), "Friend"),
			},
			GivenBatchSize: 3,
			Expected: []string{
				`DELETE EDGE Friend 100->200@0;` +
					`DELETE EDGE Friend 100->300@0;` +
					`INSERT EDGE Friend () VALUES 100->200@0:();`,
				`INSERT EDGE Friend () VALUES 100->300@0:();`,
			},
			IsErrExpected: false,
		},
		{
			Description: "Edge batch with 4 statements and batch size of 4",
			GivenStatements: []statement.IEdgeStatementOperation[int64]{
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(200), "Friend"),
				edge_delete.NewDeleteEdgeStatement(int64(100), int64(300), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(200), "Friend"),
				edge_insert.NewInsertEdgeStatement(int64(100), int64(300), "Friend"),
			},
			GivenBatchSize: 4,
			Expected: []string{
				`DELETE EDGE Friend 100->200@0;` +
					`DELETE EDGE Friend 100->300@0;` +
					`INSERT EDGE Friend () VALUES 100->200@0:();` +
					`INSERT EDGE Friend () VALUES 100->300@0:();`,
			},
			IsErrExpected: false,
		},
	}
}

func GetTestCasesForGenerateEdgeAlterStatement() []TestCaseGenerateEdgeAlterStatement {
	return []TestCaseGenerateEdgeAlterStatement{
		{
			Description: "A simple tag alter statement with add definition",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				}),
			Expected: `ALTER EDGE tag1 ADD (prop1 string);`,
		},
		{
			Description: "A simple tag alter statement with add definition and tag comment",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				},
				edge_alter.WithTagComment("test comment")),
			Expected: `ALTER EDGE tag1 ADD (prop1 string) COMMENT 'test comment';`,
		},
		{
			Description: "A simple tag alter statement with add definition and ttl definition",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				},
				edge_alter.WithTtlDefinitions([]edge_alter.TTLDefinition{
					edge_alter.NewTTLDefinition(100, "created_at"),
				})),
			Expected: `ALTER EDGE tag1 ADD (prop1 string) TTL_DURATION = 100, TTL_COL = "created_at";`,
		},
		{
			Description: "A simple tag alter statement with add definition and two ttl definitions",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				},
				edge_alter.WithTtlDefinitions([]edge_alter.TTLDefinition{
					edge_alter.NewTTLDefinition(100, "created_at"),
					edge_alter.NewTTLDefinition(200, "updated_at"),
				})),
			Expected: `ALTER EDGE tag1 ADD (prop1 string) TTL_DURATION = 100, TTL_COL = "created_at", TTL_DURATION = 200, TTL_COL = "updated_at";`,
		},
		{
			Description: "A simple tag alter statement with two add definitions",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
					edge_alter.NewAlterTypeAddDefinition("prop2", statement.PropertyTypeInt),
				}),
			Expected: `ALTER EDGE tag1 ADD (prop1 string), ADD (prop2 int);`,
		},
		{
			Description: "A simple tag alter statement with change definition",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString),
				}),
			Expected: `ALTER EDGE tag1 CHANGE (prop1 string);`,
		},
		{
			Description: "A simple tag alter statement with two change definitions",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString),
					edge_alter.NewAlterTypeChangeDefinition("prop2", statement.PropertyTypeString),
				}),
			Expected: `ALTER EDGE tag1 CHANGE (prop1 string), CHANGE (prop2 string);`,
		},
		{
			Description: "A simple tag alter statement with drop definition",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeDropDefinition("prop1"),
				}),
			Expected: `ALTER EDGE tag1 DROP (prop1);`,
		},
		{
			Description: "A simple tag alter statement with two drop definitions",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeDropDefinition("prop1"),
					edge_alter.NewAlterTypeDropDefinition("prop2"),
				}),
			Expected: `ALTER EDGE tag1 DROP (prop1), DROP (prop2);`,
		},
		{
			Description: "A tag alter statement with add, change and drop definitions",
			Given: edge_alter.NewAlterEdgeStatement("tag1",
				[]edge_alter.IAlterEdgeTypeDefinition{
					edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
					edge_alter.NewAlterTypeChangeDefinition("prop2", statement.PropertyTypeString),
					edge_alter.NewAlterTypeDropDefinition("prop3"),
				}),
			Expected: `ALTER EDGE tag1 ADD (prop1 string), CHANGE (prop2 string), DROP (prop3);`,
		},
	}
}
