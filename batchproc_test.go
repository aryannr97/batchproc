package batchproc

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	newCtx := context.WithValue(ctx, TestRun, true)
	type args struct {
		ctx           context.Context
		id            string
		totalCount    int
		data          interface{}
		loadBatches   GetBatchUnitFunc
		batchSize     []int
		testDataCount int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Get new instance of batch process executor",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				2000,
			},
			16,
		},
		{
			"Get new instance of batch process executor with fixed batchSize",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{100},
				2000,
			},
			20,
		},
		{
			"Get new instance of batch process executor with min number of batches",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				100,
			},
			1,
		},
		{
			"Get new instance of batch process executor with min number of batches",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				100,
			},
			1,
		},
		{
			"Get new instance of batch process executor with 4 number of batches",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				500,
			},
			4,
		},
		{
			"Get new instance of batch process executor with 8 number of batches",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				1000,
			},
			8,
		},
		{
			"Get new instance of batch process executor with 16 number of batches",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				2000,
			},
			16,
		},
		{
			"Get new instance of batch process executor with max number of batches",
			args{
				ctx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				2500,
			},
			20,
		},
		{
			"Get new instance of batch process executor with modified test context",
			args{
				newCtx,
				"unit-test",
				0,
				nil,
				getValidBatchUnit,
				[]int{},
				5,
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCollection := prepareTestCollection(tt.args.testDataCount)
			tt.args.totalCount = len(testCollection)
			tt.args.data = testCollection
			if got := New(tt.args.ctx, tt.args.id, tt.args.totalCount, tt.args.data, tt.args.loadBatches, tt.args.batchSize...); got.GetBatchCount() != tt.want {
				t.Errorf("batchCount = %v want = %v", got.GetBatchCount(), tt.want)
			}
		})
	}
}

func TestExecutor_Run(t *testing.T) {
	ctx := context.Background()
	testCollection := prepareTestCollection(2000)

	type args struct {
		ctx         context.Context
		id          string
		totalCount  int
		data        interface{}
		loadBatches GetBatchUnitFunc
		batchSize   []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Successful batch processing",
			args{
				ctx,
				"unit-test",
				len(testCollection),
				testCollection,
				getValidBatchUnit,
				[]int{},
			},
			false,
		},
		{
			"Failed batch processing",
			args{
				ctx,
				"unit-test",
				len(testCollection),
				testCollection,
				getErrorBatchUnit,
				[]int{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.args.ctx, tt.args.id, tt.args.totalCount, tt.args.data, tt.args.loadBatches, tt.args.batchSize...)
			if err := e.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Executor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutor_Aggregate(t *testing.T) {
	ctx := context.Background()
	testCollection := prepareTestCollection(2000)

	type fields struct {
		ctx         context.Context
		id          string
		totalCount  int
		data        interface{}
		loadBatches GetBatchUnitFunc
		batchSize   []int
	}
	type args struct {
		aggregation AggregationFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			"Aggregate batch results to get final sum of first 2000 +ve integers",
			fields{
				ctx,
				"unit-test",
				len(testCollection),
				testCollection,
				getValidBatchUnit,
				[]int{},
			},
			args{
				aggregation,
			},
			2001000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.fields.ctx, tt.fields.id, tt.fields.totalCount, tt.fields.data, tt.fields.loadBatches, tt.fields.batchSize...)
			if err := e.Run(); err != nil {
				t.Errorf("Executor.Run() = %v", err)
			}
			if got := e.Aggregate(tt.args.aggregation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Executor.Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkBatchprocessor(b *testing.B) {
	ctx := context.Background()
	testCollection := prepareTestCollection(tenMillion)

	type args struct {
		ctx         context.Context
		id          string
		totalCount  int
		data        interface{}
		loadBatches GetBatchUnitFunc
		batchSize   []int
		batchFlag   bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Addition of first 10 million integers without batch processing",
			args{
				ctx,
				"unit-test",
				len(testCollection),
				testCollection,
				getValidBatchUnit,
				[]int{},
				false,
			},
		},
		{
			"Addition of first 10 million integers with batch processing",
			args{
				ctx,
				"unit-test",
				len(testCollection),
				testCollection,
				getValidBatchUnit,
				[]int{},
				true,
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			if tt.args.batchFlag {
				e := New(tt.args.ctx, tt.args.id, tt.args.totalCount, tt.args.data, tt.args.loadBatches, tt.args.batchSize...)
				for n := 0; n < b.N; n++ {
					_ = e.Run()
					e.Aggregate(aggregation)
				}
			} else {
				e := &TestBatchUnit{}
				for n := 0; n < b.N; n++ {
					_ = e.Compute(0, tt.args.totalCount, tt.args.data)
				}
			}
		})
	}
}
