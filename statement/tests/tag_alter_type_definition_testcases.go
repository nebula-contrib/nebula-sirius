package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/tag_alter"
)

type TestCaseTagAlterTypeDefinitionAddStatement struct {
	Description   string
	Given         tag_alter.AlterTypeAddDefinition
	Expected      string
	IsErrExpected bool
}

type TestCaseTagAlterTypeDefinitionChangeStatement struct {
	Description   string
	Given         tag_alter.AlterTypeChangeDefinition
	Expected      string
	IsErrExpected bool
}

type TestCaseTagAlterTypeDefinitionDropStatement struct {
	Description   string
	Given         tag_alter.AlterTypeDropDefinition
	Expected      string
	IsErrExpected bool
}

func GetTestCasesForGenerateTagAlterTypeDefinitionAddStatement() []TestCaseTagAlterTypeDefinitionAddStatement {
	return []TestCaseTagAlterTypeDefinitionAddStatement{
		{
			Description: "A simple add tag statement",
			Given:       tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
			Expected:    `ADD (prop1 string NULL)`,
		},
		{
			Description: "A simple add tag statement with not null",
			Given:       tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString, tag_alter.WithAlterTypeAddNotNullable(true)),
			Expected:    `ADD (prop1 string NOT NULL)`,
		},
		{
			Description: "A simple add tag statement with comment, not null and default value",
			Given: tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString,
				tag_alter.WithAlterTypeAddNotNullable(true),
				tag_alter.WithAlterTypeAddDefault("default_value"),
				tag_alter.WithAlterTypeAddComment("my_comment")),
			Expected: `ADD (prop1 string NOT NULL DEFAULT 'default_value' COMMENT 'my_comment')`,
		},
		{
			Description:   "A error case for empty property name",
			Given:         tag_alter.NewAlterTypeAddDefinition("", statement.PropertyTypeString),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateTagAlterTypeDefinitionChangeStatement() []TestCaseTagAlterTypeDefinitionChangeStatement {
	return []TestCaseTagAlterTypeDefinitionChangeStatement{
		{
			Description: "A simple change tag statement",
			Given:       tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString),
			Expected:    `CHANGE (prop1 string NULL)`,
		},
		{
			Description: "A simple change tag statement with not null",
			Given: tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString,
				tag_alter.WithAlterTypeChangeNotNullable(true)),
			Expected: `CHANGE (prop1 string NOT NULL)`,
		},
		{
			Description: "A simple change tag statement with comment",
			Given: tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString,
				tag_alter.WithAlterTypeChangeComment("my_comment")),
			Expected: `CHANGE (prop1 string NULL COMMENT 'my_comment')`,
		},
		{
			Description: "A simple change tag statement default value",
			Given: tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString,
				tag_alter.WithAlterTypeChangeDefault("default_value_set")),
			Expected: `CHANGE (prop1 string NULL DEFAULT 'default_value_set')`,
		},
		{
			Description: "A simple change tag statement with not null, comment and default value",
			Given: tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString,
				tag_alter.WithAlterTypeChangeNotNullable(true),
				tag_alter.WithAlterTypeChangeComment("my_comment"),
				tag_alter.WithAlterTypeChangeDefault("default_value_set")),
			Expected: `CHANGE (prop1 string NOT NULL DEFAULT 'default_value_set' COMMENT 'my_comment')`,
		},
		{
			Description:   "A error case for empty property name",
			Given:         tag_alter.NewAlterTypeChangeDefinition("", statement.PropertyTypeString),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateTagAlterTypeDefinitionDropStatement() []TestCaseTagAlterTypeDefinitionDropStatement {
	return []TestCaseTagAlterTypeDefinitionDropStatement{
		{
			Description: "A simple change tag statement",
			Given:       tag_alter.NewAlterTypeDropDefinition("prop1"),
			Expected:    `DROP (prop1)`,
		},
		{
			Description:   "A error case for empty property name",
			Given:         tag_alter.NewAlterTypeDropDefinition(""),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}
