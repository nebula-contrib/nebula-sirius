package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement"
	"github.com/nebula-contrib/nebula-sirius/statement/tag_alter"
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
	Description   string
	Given         tag_drop.DropTagStatement
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateDeleteTagStatement[TVidType string | int64] struct {
	Description   string
	Given         tag_delete.DeleteTagStatement[TVidType]
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateTagTTLDefinitionStatement struct {
	Description   string
	Given         tag_alter.TTLDefinition
	Expected      string
	IsErrExpected bool
}

type TestCaseGenerateTagAlterStatement struct {
	Description   string
	Given         tag_alter.AlterTagStatement
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
					tag_create.NewTagProperty("name", statement.PropertyTypeString, false),
					tag_create.NewTagProperty("email", statement.PropertyTypeString, true),
					tag_create.NewTagProperty("phone", statement.PropertyTypeString, true),
				})),
			Expected:      `CREATE TAG IF NOT EXISTS account (name string NOT NULL, email string NULL, phone string NULL);`,
			IsErrExpected: false,
		},
		{
			Description: "A simple create tag statement with default type",
			Given: tag_create.NewCreateTagStatement("account",
				tag_create.WithIfNotExists(true),
				tag_create.WithProperties([]tag_create.TagProperty{
					tag_create.NewTagProperty("name", statement.PropertyTypeString, false),
					tag_create.NewTagProperty("email", statement.PropertyTypeString, true),
					tag_create.NewTagProperty("phone", statement.PropertyTypeString, true),
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
					tag_create.NewTagProperty("name", statement.PropertyTypeString, false),
					tag_create.NewTagProperty("email", statement.PropertyTypeString, true),
					tag_create.NewTagProperty("phone", statement.PropertyTypeString, true),
					tag_create.NewTagProperty("created_at", statement.PropertyTypeTimestamp, true),
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
					tag_create.NewTagProperty("name", statement.PropertyTypeString, false),
					tag_create.NewTagProperty("email", statement.PropertyTypeString, true),
					tag_create.NewTagProperty("phone", statement.PropertyTypeString, true),
					tag_create.NewTagProperty("created_at", statement.PropertyTypeTimestamp, true),
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
		{
			Description:   "A error case with empty tag name",
			Given:         tag_drop.NewDropTagStatement(""),
			Expected:      "",
			IsErrExpected: true,
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

func GetTestCasesForGenerateTagTTLDefinitionStatement() []TestCaseGenerateTagTTLDefinitionStatement {
	return []TestCaseGenerateTagTTLDefinitionStatement{
		{
			Description: "A simple tag ttl definition statement",
			Given:       tag_alter.NewTTLDefinition(100, "created_at"),
			Expected:    `TTL_DURATION = 100, TTL_COL = "created_at"`,
		},
		{
			Description: "A simple tag ttl definition statement",
			Given:       tag_alter.NewTTLDefinition(0, "created_at"),
			Expected:    `TTL_DURATION = 0, TTL_COL = "created_at"`,
		},
		{
			Description:   "An error case with negative ttl duration",
			Given:         tag_alter.NewTTLDefinition(-1, "created_at"),
			Expected:      "",
			IsErrExpected: true,
		},
		{
			Description:   "An error case missing ttl column",
			Given:         tag_alter.NewTTLDefinition(100, ""),
			Expected:      "",
			IsErrExpected: true,
		},
	}
}

func GetTestCasesForGenerateTagAlterStatement() []TestCaseGenerateTagAlterStatement {
	return []TestCaseGenerateTagAlterStatement{
		{
			Description: "A simple tag alter statement with add definition",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				}),
			Expected: `ALTER TAG tag1 ADD (prop1 string NULL);`,
		},
		{
			Description: "A simple tag alter statement with add definition and tag comment",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				},
				tag_alter.WithTagComment("test comment")),
			Expected: `ALTER TAG tag1 ADD (prop1 string NULL) COMMENT 'test comment';`,
		},
		{
			Description: "A simple tag alter statement with add definition and ttl definition",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				},
				tag_alter.WithTtlDefinitions([]tag_alter.TTLDefinition{
					tag_alter.NewTTLDefinition(100, "created_at"),
				})),
			Expected: `ALTER TAG tag1 ADD (prop1 string NULL) TTL_DURATION = 100, TTL_COL = "created_at";`,
		},
		{
			Description: "A simple tag alter statement with add definition and two ttl definitions",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
				},
				tag_alter.WithTtlDefinitions([]tag_alter.TTLDefinition{
					tag_alter.NewTTLDefinition(100, "created_at"),
					tag_alter.NewTTLDefinition(200, "updated_at"),
				})),
			Expected: `ALTER TAG tag1 ADD (prop1 string NULL) TTL_DURATION = 100, TTL_COL = "created_at", TTL_DURATION = 200, TTL_COL = "updated_at";`,
		},
		{
			Description: "A simple tag alter statement with two add definitions",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
					tag_alter.NewAlterTypeAddDefinition("prop2", statement.PropertyTypeInt),
				}),
			Expected: `ALTER TAG tag1 ADD (prop1 string NULL), ADD (prop2 int NULL);`,
		},
		{
			Description: "A simple tag alter statement with change definition",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString),
				}),
			Expected: `ALTER TAG tag1 CHANGE (prop1 string NULL);`,
		},
		{
			Description: "A simple tag alter statement with two change definitions",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeChangeDefinition("prop1", statement.PropertyTypeString),
					tag_alter.NewAlterTypeChangeDefinition("prop2", statement.PropertyTypeString),
				}),
			Expected: `ALTER TAG tag1 CHANGE (prop1 string NULL), CHANGE (prop2 string NULL);`,
		},
		{
			Description: "A simple tag alter statement with drop definition",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeDropDefinition("prop1"),
				}),
			Expected: `ALTER TAG tag1 DROP (prop1);`,
		},
		{
			Description: "A simple tag alter statement with two drop definitions",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeDropDefinition("prop1"),
					tag_alter.NewAlterTypeDropDefinition("prop2"),
				}),
			Expected: `ALTER TAG tag1 DROP (prop1), DROP (prop2);`,
		},
		{
			Description: "A tag alter statement with add, change and drop definitions",
			Given: tag_alter.NewAlterTagStatement("tag1",
				[]tag_alter.IAlterTagTypeDefinition{
					tag_alter.NewAlterTypeAddDefinition("prop1", statement.PropertyTypeString),
					tag_alter.NewAlterTypeChangeDefinition("prop2", statement.PropertyTypeString),
					tag_alter.NewAlterTypeDropDefinition("prop3"),
				}),
			Expected: `ALTER TAG tag1 ADD (prop1 string NULL), CHANGE (prop2 string NULL), DROP (prop3);`,
		},
	}
}
