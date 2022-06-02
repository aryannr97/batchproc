package batchproc

import (
	"fmt"
)

const (
	tenMillion = 10000000
)

var (
	testCollection = []int{}
)

// TestBatchUnit is custom type implementing BatchUnit interface
type TestBatchUnit struct {
	result int
}

// Compute performs addition portion of given slice of numbers
func (t *TestBatchUnit) Compute(start, end int, collection interface{}) error {
	if data, ok := collection.([]int); ok {
		data = data[start:end]
		count := 0
		for _, value := range data {
			count += value
		}
		t.result += count
	}

	return nil
}

// GetResult returns result stored during Compute operation
func (t *TestBatchUnit) GetResult() interface{} {
	return t.result
}

// // TestBatchUnit is custom type implementing BatchUnit interface
type TestNegativeBatchUnit struct{}

// Compute intentionally returns error for unit testing
func (t *TestNegativeBatchUnit) Compute(start, end int, data interface{}) error {
	return fmt.Errorf("test batch compute error")
}

// GetResult returns result stored during Compute operation
func (t *TestNegativeBatchUnit) GetResult() interface{} {
	return nil
}

// getValidBatchUnit returns instance of TestBatchUnit
func getValidBatchUnit() BatchUnit {
	return &TestBatchUnit{}
}

// getErrorBatchUnit returns instance of TestNegativeBatchUnit
func getErrorBatchUnit() BatchUnit {
	return &TestNegativeBatchUnit{}

}

// prepareTestCollection prepares slice of given N number of integers
func prepareTestCollection(N int) []int {
	testCollection = []int{}
	for i := 1; i <= N; i++ {
		testCollection = append(testCollection, i)
	}

	return testCollection
}

// aggregation is user defined custom AggregationFunc, which adds integer result of all batches
func aggregation(results []interface{}) interface{} {
	var res int

	for _, result := range results {
		if total, ok := result.(int); ok {
			res += total
		}
	}

	return res
}
