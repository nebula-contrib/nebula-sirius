package tests

import (
	"github.com/nebula-contrib/nebula-sirius/statement/edge_batch"
	"reflect"
	"testing"
)

func TestGenerateBatchEdgeStatementWhereVidStrings(t *testing.T) {
	testCases := GetTestCasesForGenerateBatchEdgeStatementWhereVidString()
	for _, testcase := range testCases {
		actual, err := edge_batch.GenerateBatchedEdgeStatements(testcase.GivenStatements, testcase.GivenBatchSize)

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
					testcase.GivenStatements, testcase.GivenBatchSize,
					testcase.Expected, len(testcase.Expected),
					actual, len(actual))
			}
		}
	}
}

func TestGenerateBatchEdgeStatementWhereVidInt64(t *testing.T) {
	testCases := GetTestCasesForGenerateBatchEdgeStatementWhereVidInt64()
	for _, testcase := range testCases {
		actual, err := edge_batch.GenerateBatchedEdgeStatements(testcase.GivenStatements, testcase.GivenBatchSize)

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
					testcase.GivenStatements, testcase.GivenBatchSize,
					testcase.Expected, len(testcase.Expected),
					actual, len(actual))
			}
		}
	}
}
