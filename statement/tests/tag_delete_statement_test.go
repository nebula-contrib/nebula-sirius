package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/tag_delete"
	"reflect"
	"testing"
)

func TestGenerateDeleteTagStatementWhereVidStrings(t *testing.T) {
	testCases := GetTestCasesForGenerateDeleteTagStatementWhereVidString()
	for _, testcase := range testCases {
		actual, err := tag_delete.GenerateDeleteTagStatement(testcase.Given)

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

func TestGenerateDeleteTagStatementWhereVidInt64(t *testing.T) {
	testCases := GetTestCasesForGenerateDeleteTagStatementWhereVidInt64()
	for _, testcase := range testCases {
		actual, err := tag_delete.GenerateDeleteTagStatement(testcase.Given)

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
