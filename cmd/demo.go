package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aryannr97/batchproc"
)

func main() {
	data := getData()

	// Creation stage
	executor := batchproc.New(context.Background(), "main", len(data), data, getBatchUnit)
	fmt.Printf("Created batch executor loaded with number of batches: %v\n", executor.GetBatchCount())

	// Computation stage
	fmt.Println("Batch processor execution starting")
	if err := executor.Run(); err != nil {
		log.Println(err)
		return
	}
	// Aggregation stage
	fmt.Println("Batch processor result aggregation starting")
	result := executor.Aggregate(aggregation)

	fmt.Printf("Addition of first 2000 integers is: %v\n", result)
	fmt.Printf("Batch processor took %vns for execution\n", executor.ElapsedDuration.Nanoseconds())

}

func getData() []int {
	data := []int{}
	for i := 1; i <= 2000; i++ {
		data = append(data, i)
	}

	return data
}

// CustomBatchUnit : define your own struct type
type CustomBatchUnit struct {
	result int
}

// Compute : define custom implementation for this method as per use case & store result in struct attribute
func (t *CustomBatchUnit) Compute(start, end int, collection interface{}) error {
	if data, ok := collection.([]int); ok {
		data = data[start:end]
		count := 0
		for _, value := range data {
			count += value
		}

		t.result += count

		return nil
	}

	return fmt.Errorf("type assertion error")
}

// GetResult : return result attribute from defined struct type
func (t *CustomBatchUnit) GetResult() interface{} {
	return t.result
}

// getBatchUnit returns CustomBatchUnit that implements BatchUnit interface
func getBatchUnit() batchproc.BatchUnit {
	return &CustomBatchUnit{}
}

// aggregation : treat results gathered from all batches as required
func aggregation(results []interface{}) interface{} {
	var res int

	for _, result := range results {
		if total, ok := result.(int); ok {
			res += total
		}
	}

	return res
}
