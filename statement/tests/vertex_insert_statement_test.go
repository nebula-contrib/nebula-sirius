package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/vertex_insert"
	"reflect"
	"testing"
)

func TestGenerateInsertVertexStatement(t *testing.T) {
	testCases := GetTestCasesForGenerateInsertVertexStatement()
	for _, testcase := range testCases {
		actual, err := vertex_insert.GenerateInsertVertexStatement(testcase.GivenVerticesArray)

		if err != nil {
			if !testcase.IsErrExpected {
				t.Errorf("For %s, expected no error, got %v", testcase.Description, err)
			}
			continue
		}

		if !reflect.DeepEqual(actual, testcase.Expected) {
			t.Errorf("For Case: %s "+
				"\n Given: %+q, arr len: %d "+
				"\n Expected: %s, len: %d"+
				"\n Got: %s, len: %d",
				testcase.Description,
				testcase.GivenVerticesArray, len(testcase.GivenVerticesArray),
				testcase.Expected, len(testcase.Expected),
				actual, len(actual))
		}
	}
}

func TestGenerateBatchedInsertVertexStatements(t *testing.T) {
	testCases := GetTestCasesForGenerateBatchedInsertVertexStatements()
	for _, testcase := range testCases {
		actual, err := vertex_insert.GenerateBatchedInsertVertexStatements(testcase.GivenVerticesArray, testcase.GivenBatchSize)

		if err != nil {
			if !testcase.IsErrExpected {
				t.Errorf("For %s, expected no error, got %v", testcase.Description, err)
			}
			continue
		}

		if len(actual) != len(testcase.Expected) {
			t.Errorf("For %s, expected %d insert scripts, got %d", testcase.Description, len(testcase.Expected), len(actual))
		}

		for i := 0; i < len(actual); i++ {
			if !reflect.DeepEqual(actual[i], testcase.Expected[i]) {
				t.Errorf("For Case: %s "+
					"\n Given: %v, batchSize: %d "+
					"\n Expected: %s, len: %d"+
					"\n Got: %s, len: %d",
					testcase.Description,
					testcase.GivenVerticesArray, testcase.GivenBatchSize,
					testcase.Expected, len(testcase.Expected),
					actual, len(actual))
			}
		}
	}
}
