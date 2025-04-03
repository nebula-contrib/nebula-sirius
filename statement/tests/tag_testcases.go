package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/tag_create"
	"github.com/nebula-contrib/nebula-sirius/statement/tag_drop"
)

type TestCaseGenerateCreateTagStatement struct {
	Description   string
	Given         tag_create.CreateTagStatement
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateDropTagStatement struct {
	Description string
	Given       tag_drop.DropTagStatement
	Expected    string
}

func GetTestCasesForGenerateCreateTagStatement() []TestCaseGenerateCreateTagStatement {
	return []TestCaseGenerateCreateTagStatement{
		{
			Description:   "A simple create tag statement without IfNotExists",
			Given:         tag_create.NewCreateTagStatement("account"),
			Expected:      `CREATE TAG account ();`,
			IsErrExpected: false,
		},
		{
			Description:   "A simple create tag statement without properties",
			Given:         tag_create.NewCreateTagStatement("account", tag_create.WithIfNotExists(true)),
			Expected:      `CREATE TAG IF NOT EXISTS account ();`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement properties",
			Given: tag_create.NewCreateTagStatement("account",
				tag_create.WithIfNotExists(true),
				tag_create.WithProperties([]tag_create.TagProperty{
					{Field: "name", Type: statement.PropertyTypeString, Nullable: false},
					{Field: "email", Type: statement.PropertyTypeString, Nullable: true},
					{Field: "phone", Type: statement.PropertyTypeString, Nullable: true},
				})),
			Expected:      `CREATE TAG IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL);`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with default type",
			Given: tag_create.NewCreateTagStatement("account",
				tag_create.WithIfNotExists(true),
				tag_create.WithProperties([]tag_create.TagProperty{
					{Field: "name", Nullable: false},
					{Field: "email", Nullable: true},
					{Field: "phone", Nullable: true},
				})),
			Expected:      `CREATE TAG IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL);`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with default type and TTL",
			Given: tag_create.NewCreateTagStatement("account",
				tag_create.WithIfNotExists(true),
				tag_create.WithTtlCol("created_at"),
				tag_create.WithTtlDuration(100),
				tag_create.WithProperties([]tag_create.TagProperty{
					{Field: "name", Nullable: false},
					{Field: "email", Nullable: true},
					{Field: "phone", Nullable: true},
					{Field: "created_at", Type: "timestamp", Nullable: true},
				})),
			Expected:      `CREATE TAG IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL, created_at timestamp NULL) TTL_DURATION = 100, TTL_COL = "created_at";`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with invalid ttl column",
			Given: tag_create.NewCreateTagStatement("account",
				tag_create.WithIfNotExists(true),
				tag_create.WithTtlCol("not_existed_field"),
				tag_create.WithTtlDuration(100),
				tag_create.WithProperties([]tag_create.TagProperty{
					{Field: "name", Nullable: false},
					{Field: "email", Nullable: true},
					{Field: "phone", Nullable: true},
					{Field: "created_at", Type: "timestamp", Nullable: true},
				})),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateDropTagStatement() []TestCaseGenerateDropTagStatement {
	return []TestCaseGenerateDropTagStatement{
		{
			Description: "A simple drop tag statement without IfExists",
			Given:       tag_drop.NewDropTagStatement("account"),
			Expected:    `DROP TAG account;`,
		},
		{
			Description: "A simple drop tag statement with IfExists",
			Given:       tag_drop.NewDropTagStatement("account", tag_drop.WithIfExists(true)),
			Expected:    `DROP TAG IF EXISTS account;`,
		},
	}
}
