package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/edge_create"
	"reflect"
	"testing"
)

func TestGenerateCreateEdgeStatement(t *testing.T) {
	testCases := GetTestCasesForGenerateCreateEdgeStatement()
	for _, testcase := range testCases {
		actual, err := edge_create.GenerateCreateEdgeStatement(testcase.Given)

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
