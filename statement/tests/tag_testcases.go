package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/tag_create"
	"github.com/nebula-contrib/nebula-sirius/statement/tag_delete"
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

type TestCaseGenerateDeleteTagStatement[TVidType string | int64] struct {
	Description   string
	Given         tag_delete.DeleteTagStatement[TVidType]
	Expected      string
	IsErrExpected bool
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

func GetTestCasesForGenerateDeleteTagStatementWhereVidString() []TestCaseGenerateDeleteTagStatement[string] {
	return []TestCaseGenerateDeleteTagStatement[string]{
		{
			Description:   "",
			Given:         tag_delete.NewDeleteTagStatement[string]([]string{"tag1", "tag2"}, []string{"vid1", "vid2"}),
			Expected:      `DELETE TAG tag1,tag2 FROM "vid1","vid2";`,
			IsErrExpected: false,
		},
		{
			Description:   "",
			Given:         tag_delete.NewDeleteTagStatement[string]([]string{"tag1"}, []string{"vid1"}),
			Expected:      `DELETE TAG tag1 FROM "vid1";`,
			IsErrExpected: false,
		},
		{
			Description:   "",
			Given:         tag_delete.NewDeleteTagStatement[string]([]string{}, []string{"vid1"}, tag_delete.WithAllTags[string]()),
			Expected:      `DELETE TAG * FROM "vid1";`,
			IsErrExpected: false,
		},
		{
			Description:   "vidList not specified",
			Given:         tag_delete.NewDeleteTagStatement[string]([]string{"tag1"}, []string{}),
			Expected:      "",
			IsErrExpected: true,
		},
		{
			Description:   "tagList and WithAllTags are mutually exclusive, cannot be used together",
			Given:         tag_delete.NewDeleteTagStatement[string]([]string{"tag1"}, []string{"vid1"}, tag_delete.WithAllTags[string]()),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateDeleteTagStatementWhereVidInt64() []TestCaseGenerateDeleteTagStatement[int64] {
	return []TestCaseGenerateDeleteTagStatement[int64]{
		{
			Description:   "",
			Given:         tag_delete.NewDeleteTagStatement[int64]([]string{"tag1", "tag2"}, []int64{100, 200}),
			Expected:      `DELETE TAG tag1,tag2 FROM 100,200;`,
			IsErrExpected: false,
		},
		{
			Description:   "",
			Given:         tag_delete.NewDeleteTagStatement[int64]([]string{"tag1"}, []int64{100}),
			Expected:      `DELETE TAG tag1 FROM 100;`,
			IsErrExpected: false,
		},
		{
			Description:   "",
			Given:         tag_delete.NewDeleteTagStatement[int64]([]string{}, []int64{100}, tag_delete.WithAllTags[int64]()),
			Expected:      `DELETE TAG * FROM 100;`,
			IsErrExpected: false,
		},
		{
			Description:   "vidList not specified",
			Given:         tag_delete.NewDeleteTagStatement[int64]([]string{"tag1"}, []int64{}),
			Expected:      "",
			IsErrExpected: true,
		},
		{
			Description:   "tagList and WithAllTags are mutually exclusive, cannot be used together",
			Given:         tag_delete.NewDeleteTagStatement[int64]([]string{"tag1"}, []int64{100}, tag_delete.WithAllTags[int64]()),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}
