package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/edge_alter"
)

type TestCaseEdgeAlterTypeDefinitionAddStatement struct {
	Description   string
	Given         edge_alter.AlterTypeAddDefinition
	Expected      string
	IsErrExpected bool
}

type TestCaseEdgeAlterTypeDefinitionChangeStatement struct {
	Description   string
	Given         edge_alter.AlterTypeChangeDefinition
	Expected      string
	IsErrExpected bool
}

type TestCaseEdgeAlterTypeDefinitionDropStatement struct {
	Description   string
	Given         edge_alter.AlterTypeDropDefinition
	Expected      string
	IsErrExpected bool
}

func GetTestCasesForGenerateEdgeAlterTypeDefinitionAddStatement() []TestCaseEdgeAlterTypeDefinitionAddStatement {
	return []TestCaseEdgeAlterTypeDefinitionAddStatement{
		{
			Description: "A simple add tag statement",
			Given:       edge_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
			Expected:    `ADD (prop1 string)`,
		},
		{
			Description:   "A error case for empty property name",
			Given:         edge_alter.NewAlterTypeAddDefinition("", statement.PropertyTypeString),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateEdgeAlterTypeDefinitionChangeStatement() []TestCaseEdgeAlterTypeDefinitionChangeStatement {
	return []TestCaseEdgeAlterTypeDefinitionChangeStatement{
		{
			Description: "A simple change tag statement",
			Given:       edge_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString),
			Expected:    `CHANGE (prop1 string)`,
		},
		{
			Description:   "A error case for empty property name",
			Given:         edge_alter.NewAlterTypeChangeDefinition("", statement.PropertyTypeString),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateEdgeAlterTypeDefinitionDropStatement() []TestCaseEdgeAlterTypeDefinitionDropStatement {
	return []TestCaseEdgeAlterTypeDefinitionDropStatement{
		{
			Description: "A simple change tag statement",
			Given:       edge_alter.NewAlterTypeDropDefinition("prop1"),
			Expected:    `DROP (prop1)`,
		},
		{
			Description:   "A error case for empty property name",
			Given:         edge_alter.NewAlterTypeDropDefinition(""),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}
