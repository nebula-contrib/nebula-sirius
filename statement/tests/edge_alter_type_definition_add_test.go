package tests

import (
	"reflect"
	"testing"
)

func TestGenerateEdgeAlterTypeDefinitionAddStatement(t *testing.T) {
	testCases := GetTestCasesForGenerateEdgeAlterTypeDefinitionAddStatement()
	for _, testcase := range testCases {
		actual, err := testcase.Given.GenerateStatement()

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
