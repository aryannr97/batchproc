package batchproc

import (
	"context"
	"fmt"
	"time"

	"github.com/bharat-rajani/rungroup"
)

// New creates, initializes & returns a new batch Executor,
//
// id is used to mark each batch with unique identification e.g {id}-batch-0,
//
// totalCount is the total size of the collection,
//
// data is actual collection need to be processed in batched fashion,
//
// getBatchUnit is used to fetch user defined type implementing BatchUnit interface and load batches,
//
// batchSize is optional & can be passed to override dynamic batch sizing.
func New(ctx context.Context, id string, totalCount int, data interface{}, getBatchUnit GetBatchUnitFunc, batchSize ...int) *Executor {
	e := &Executor{
		ctx:        ctx,
		id:         id,
		totalCount: totalCount,
		collection: data,
	}

	if len(batchSize) == 0 {
		e.prepareDynamicLoading()
	} else {
		e.prepareFixedLoading(batchSize[0])
	}

	// load batch units with user defined func
	for i := 0; i < e.batchCount; i++ {
		e.batches = append(e.batches, getBatchUnit())
	}

	return e
}

// Executor represents concurrent batch processing entity
type Executor struct {
	ctx             context.Context
	id              string
	batchSize       int
	totalCount      int
	batchCount      int
	batches         []BatchUnit
	collection      interface{}
	startTime       time.Time
	endTime         time.Time
	ElapsedDuration time.Duration
}

// prepareDynamicLoading loads number of batches, batchSize depending on totalCount,
// number of batches formed will be min:1 max:20.
func (e *Executor) prepareDynamicLoading() {
	/**
	Scale down the sample space by a factor of 100 if function is triggered through unit tests.
	Testing different number of batches can be accomplished even with smaller size of mockData.
	*/
	scaleDownFactor := 1
	if test, _ := e.ctx.Value(TestRun).(bool); test {
		scaleDownFactor = 100
	}

	if e.totalCount <= (100 / scaleDownFactor) {
		e.batchCount = 1
	} else if e.totalCount <= (500 / scaleDownFactor) {
		e.batchCount = 4
	} else if e.totalCount <= (1000 / scaleDownFactor) {
		e.batchCount = 8
	} else if e.totalCount <= (2000 / scaleDownFactor) {
		e.batchCount = 16
	} else {
		e.batchCount = 20
	}

	e.batchSize = e.totalCount / e.batchCount
}

// prepareFixedLoading loads number of batches based on fixed batchSize
func (e *Executor) prepareFixedLoading(batchSize int) {
	e.batchSize = batchSize
	e.batchCount = e.totalCount / e.batchSize
}

// Run triggers concurrent routines for computation through batch units
func (e *Executor) Run() error {
	e.startTime = time.Now()
	countRg, ctxRg := rungroup.WithContext(e.ctx)

	for i, item := range e.batches {
		index := i
		batch := item
		id := fmt.Sprintf("%v-batch-%v", e.id, index)

		countRg.GoWithFunc(func(context.Context) error {
			start := index * e.batchSize
			end := (index * e.batchSize) + e.batchSize
			// Consider all the remaining items for last batch
			if index == e.batchCount-1 {
				end = e.totalCount
			}
			// Perform computation per batch
			err := batch.Compute(start, end, e.collection)
			return err
		}, ctxRg, true, id)
	}

	// Wait (blocking) till all the batches are done processing
	if countRgErr := countRg.Wait(); countRgErr != nil {
		return countRgErr
	}

	return nil
}

// Aggregate is used for result aggregation and mark closure of batch processing,
// user should define own aggregation func to process results from different batches.
func (e *Executor) Aggregate(aggregation AggregationFunc) interface{} {
	results := []interface{}{}
	for _, batch := range e.batches {
		results = append(results, batch.GetResult())
	}

	response := aggregation(results)

	// Mark process completion post result aggregation
	e.onComplete()

	return response
}

// onComplete performs logging,reporting activities on batch processing completion
func (e *Executor) onComplete() {
	e.endTime = time.Now()
	e.ElapsedDuration = time.Duration(e.endTime.Sub(e.startTime))
}

// GetBatchCount returns number of batches for given executor
func (e *Executor) GetBatchCount() int {
	return e.batchCount
}
