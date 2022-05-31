# batchproc
[![RELEASE](https://github.com/aryannr97/batchproc/actions/workflows/release.yml/badge.svg)](https://github.com/aryannr97/batchproc/actions/workflows/release.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aryannr97/batchproc)
[![Codecov](https://img.shields.io/codecov/c/github/aryannr97/batchproc)](https://app.codecov.io/gh/aryannr97/batchproc)
[![Go Report Card](https://goreportcard.com/badge/github.com/aryannr97/batchproc)](https://goreportcard.com/report/github.com/aryannr97/batchproc)
[![Go Report Card](https://img.shields.io/badge/Linter-golangci--lint-informational)](https://golangci-lint.run)
[![MIT license](https://img.shields.io/github/license/aryannr97/batchproc)](https://github.com/aryannr97/batchproc/blob/main/LICENSE)
## Generic batch processor

batchproc is developed to serve the purpose of enabling a uniform way to process large volumes of different types of data & run repititive computations as per requirement 

batchproc essentially provides
- Ability to process any type of **indexed**, **iterative** data
- **Customization** for batch steps
- [**Rungroups**](https://github.com/bharat-rajani/rungroup) to execute batch units concurrently

### :floppy_disk: Installation
```
go get -u github.com/aryannr97/batchproc
```

### :notebook_with_decorative_cover: Documentation
Detailed documentation can be found [here](https://github.com/aryannr97/batchproc/wiki/batchproc-wiki)

### :technologist: Usage

##### 1. **New** 
- Returns a new instance of batch processing executor.
```go
// Creation stage
executor := batchproc.New(context.Background(), "main", len(data), data, getBatchUnit)
```
- **GetBatchUnitFunc** is a predefined type `type GetBatchUnitFunc func() BatchUnit`
- User should define & pass this function while creating executor via *New* to load batches with user defined type implementing *BatchUnit* interface
```go
// getBatchUnit returns CustomBatchUnit that implements BatchUnit interface
func getBatchUnit() batchproc.BatchUnit {
	return &CustomBatchUnit{}
}
```

##### 2. **BatchUnit**
- It's an interface defining two methods Compute & GetResult
```go
type BatchUnit interface {
	Compute(int, int, interface{}) error
	GetResult() interface{}
}

```
- This provides a way for user to have customize batch steps as per use-case
```go
// CustomBatchUnit : define your own struct type
type CustomBatchUnit struct {
	result int
}

// Compute : define custom implementation for this method as per use case & store result in struct attribute
func (t *CustomBatchUnit) Compute(start, end int, collection interface{}) error {
	if data, ok := collection.([]int); ok {
		data = data[start:end]
		for _, value := range data {
			t.result += value
		}

		return nil
	}

	return fmt.Errorf("type assertion error")
}

// GetResult : return result attribute from defined struct type
func (t *CustomBatchUnit) GetResult() interface{} {
	return t.result
}
```

##### 3. **Run**
- Executes batch processing with loaded batch units
- Internally each batch will process under a separate rungroup (concurrent routine)
- Any single batch failure will interrrupt other batches & cancel overall batch processing
```go
// Computation stage
if err := executor.Run(); err != nil {
	log.Println(err)
	return
}
```

##### 4. **Aggregate**
- Aggregate performs final stage of the batch processing, where results from different batches are 
aggregated to form overall result.
- User should pass own aggregation function to process batch result as desired
```go
// Aggregation stage
result := executor.Aggregate(aggregation)
```
- **AggregationFunc** is a predefined type, to serve different case-to-case implementations `type AggregationFunc func([]interface{}) interface{}`
```go
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
```

### :rocket: Benchmark results

```
goos: linux
goarch: amd64
pkg: github.com/aryannr97/batchproc
cpu: Intel(R) Xeon(R) Platinum 8272CL CPU @ 2.60GHz
```

|                           Test case                               | b.N iterations | Avg processsing time |
| ----------------------------------------------------------------- | -------------- | -------------------- |
| Addition_of_first_10_million_integers_without_batch_processing-2  |      100       |    10665257 ns/op    |
| Addition_of_first_10_million_integers_with_batch_processing-2     |      251       |     4620470 ns/op    |


```
goos: linux
goarch: amd64
pkg: github.com/aryannr97/batchproc
cpu: Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz
```
|                           Test case                               | b.N iterations | Avg processsing time |
| ----------------------------------------------------------------- | -------------- | -------------------- |
| Addition_of_first_10_million_integers_without_batch_processing-12 |      184       |     6394545 ns/op    |
| Addition_of_first_10_million_integers_with_batch_processing-12    |      295       |     4074257 ns/op    |

As observed, batch processing can nearly run for **twice the number of iterations** and **reduce processing time by half**

### Example
To execute the demo program, run command `make demo`
