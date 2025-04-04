package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/edge_create"
	"github.com/nebula-contrib/nebula-sirius/statement/edge_drop"
)

type TestCaseGenerateCreateEdgeStatement struct {
	Description   string
	Given         edge_create.CreateEdgeStatement
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateDropEdgeStatement struct {
	Description   string
	Given         edge_drop.DropEdgeStatement
	Expected      string
	IsErrExpected bool
}

func GetTestCasesForGenerateCreateEdgeStatement() []TestCaseGenerateCreateEdgeStatement {
	return []TestCaseGenerateCreateEdgeStatement{
		{
			Description:   "A simple create tag statement without IfNotExists",
			Given:         edge_create.NewCreateEdgeStatement("account"),
			Expected:      `CREATE EDGE account ();`,
			IsErrExpected: false,
		},
		{
			Description:   "A simple create tag statement without properties",
			Given:         edge_create.NewCreateEdgeStatement("account", edge_create.WithIfNotExists(true)),
			Expected:      `CREATE EDGE IF NOT EXISTS account ();`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement properties",
			Given: edge_create.NewCreateEdgeStatement("account",
				edge_create.WithIfNotExists(true),
				edge_create.WithProperties([]edge_create.EdgeProperty{
					edge_create.NewEdgeProperty("name", statement.PropertyTypeString, false),
					edge_create.NewEdgeProperty("email", statement.PropertyTypeString, true),
					edge_create.NewEdgeProperty("phone", statement.PropertyTypeString, true),
				})),
			Expected:      `CREATE EDGE IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL);`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with default type",
			Given: edge_create.NewCreateEdgeStatement("account",
				edge_create.WithIfNotExists(true),
				edge_create.WithProperties([]edge_create.EdgeProperty{
					edge_create.NewEdgeProperty("name", statement.PropertyTypeString, false),
					edge_create.NewEdgeProperty("email", statement.PropertyTypeString, true),
					edge_create.NewEdgeProperty("phone", statement.PropertyTypeString, true),
				})),
			Expected:      `CREATE EDGE IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL);`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with default type and TTL",
			Given: edge_create.NewCreateEdgeStatement("account",
				edge_create.WithIfNotExists(true),
				edge_create.WithTtlCol("created_at"),
				edge_create.WithTtlDuration(100),
				edge_create.WithProperties([]edge_create.EdgeProperty{
					edge_create.NewEdgeProperty("name", statement.PropertyTypeString, false),
					edge_create.NewEdgeProperty("email", statement.PropertyTypeString, true),
					edge_create.NewEdgeProperty("phone", statement.PropertyTypeString, true),
					edge_create.NewEdgeProperty("created_at", statement.PropertyTypeTimestamp, true),
				})),
			Expected:      `CREATE EDGE IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL, created_at timestamp NULL) TTL_DURATION = 100, TTL_COL = "created_at";`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with invalid ttl column",
			Given: edge_create.NewCreateEdgeStatement("account",
				edge_create.WithIfNotExists(true),
				edge_create.WithTtlCol("not_existed_field"),
				edge_create.WithTtlDuration(100),
				edge_create.WithProperties([]edge_create.EdgeProperty{
					edge_create.NewEdgeProperty("name", statement.PropertyTypeString, false),
					edge_create.NewEdgeProperty("email", statement.PropertyTypeString, true),
					edge_create.NewEdgeProperty("phone", statement.PropertyTypeString, true),
					edge_create.NewEdgeProperty("created_at", statement.PropertyTypeTimestamp, true),
				})),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateDropEdgeStatement() []TestCaseGenerateDropEdgeStatement {
	return []TestCaseGenerateDropEdgeStatement{
		{
			Description: "A simple drop tag statement without IfExists",
			Given:       edge_drop.NewDropEdgeStatement("edge1"),
			Expected:    `DROP EDGE edge1;`,
		},
		{
			Description: "A simple drop tag statement with IfExists",
			Given:       edge_drop.NewDropEdgeStatement("edge1", edge_drop.WithIfExists(true)),
			Expected:    `DROP EDGE IF EXISTS edge1;`,
		},
		{
			Description:   "A error case with empty tag name",
			Given:         edge_drop.NewDropEdgeStatement(""),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}
