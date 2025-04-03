package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/edge_upsert"
	"reflect"
	"testing"
)

func TestGenerateUpsertEdgeStatementWhereVidStrings(t *testing.T) {
	testCases := GetTestCasesForGenerateUpsertEdgeStatementWhereVidString()
	for _, testcase := range testCases {
		actual, err := edge_upsert.GenerateUpsertEdgeStatement(testcase.Given)

		if err != nil {
			if !testcase.IsErrExpected {
				t.Errorf("For %s, expected no error, got %v", testcase.Description, err)
			}
			continue
		}

		if !reflect.DeepEqual(actual, testcase.Expected) {
			t.Errorf("For Case: %s "+
				"\n Given: %+v "+
				"\n Expected: %s, len: %d"+
				"\n Got: %s, len: %d",
				testcase.Description,
				testcase.Given,
				testcase.Expected, len(testcase.Expected),
				actual, len(actual))
		}
	}
}

func TestGenerateUpsertEdgeStatementWhereVidInt64(t *testing.T) {
	testCases := GetTestCasesForGenerateUpsertEdgeStatementWhereVidInt64()
	for _, testcase := range testCases {
		actual, err := edge_upsert.GenerateUpsertEdgeStatement(testcase.Given)

		if err != nil {
			if !testcase.IsErrExpected {
				t.Errorf("For %s, expected no error, got %v", testcase.Description, err)
			}
			continue
		}

		if !reflect.DeepEqual(actual, testcase.Expected) {
			t.Errorf("For Case: %s "+
				"\n Given: %+v "+
				"\n Expected: %s, len: %d"+
				"\n Got: %s, len: %d",
				testcase.Description,
				testcase.Given,
				testcase.Expected, len(testcase.Expected),
				actual, len(actual))
		}
	}
}
