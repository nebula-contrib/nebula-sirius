package statement

import (
	"reflect"
	"testing"
)

func TestGenerateInsertEdgeStatementWhereVidStrings(t *testing.T) {
	testCases := GetTestCasesForGenerateInsertEdgeStatementWhereVidString()
	for _, testcase := range testCases {
		actual, err := GenerateInsertEdgeStatement(testcase.Given)

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

func TestGenerateInsertEdgeStatementWhereVidInt64(t *testing.T) {
	testCases := GetTestCasesForGenerateInsertEdgeStatementWhereVidInt64()
	for _, testcase := range testCases {
		actual, err := GenerateInsertEdgeStatement(testcase.Given)

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
