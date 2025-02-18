package statement

import (
	"reflect"
	"testing"
)

func TestGenerateDeleteEdgeStatementWhereVidStrings(t *testing.T) {
	testCases := GetTestCasesForGenerateDeleteEdgeStatementWhereVidString()
	for _, testcase := range testCases {
		actual, err := GenerateDeleteEdgeStatement(testcase.Given)

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

func TestGenerateDeleteEdgeStatementWhereVidInt64(t *testing.T) {
	testCases := GetTestCasesForGenerateDeleteEdgeStatementWhereVidInt64()
	for _, testcase := range testCases {
		actual, err := GenerateDeleteEdgeStatement(testcase.Given)

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
