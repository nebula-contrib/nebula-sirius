package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/tag_alter"
	"reflect"
	"testing"
)

func TestGenerateTagAlterStatement(t *testing.T) {
	testCases := GetTestCasesForGenerateTagAlterStatement()
	for _, testcase := range testCases {
		actual, err := tag_alter.GenerateAlterTagStatement(testcase.Given)

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
				"\n Got:      %s, len: %d",
				testcase.Description,
				testcase.Given,
				testcase.Expected, len(testcase.Expected),
				actual, len(actual))
		}
	}
}
