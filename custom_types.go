package batchproc

// BatchUnit interface defines functionalities of a single batch
type BatchUnit interface {
	Compute(int, int, interface{}) error
	GetResult() interface{}
}

// GetBatchUnitFunc defines func signature for returning user defined type which implements BatchUnit interface
type GetBatchUnitFunc func() BatchUnit

// AggrerationFunc defines func signature for result aggregation
type AggregationFunc func([]interface{}) interface{}
