package numareduce

import (
	"context"
	"fmt"

	"github.com/redpanda-data/benthos/v4/public/service"
)

var numareduceConfigSpec = service.NewConfigSpec().
	Summary("a Numaflow reduce input").Field(service.NewIntField("maxBatchSize").Default(10))

type numaReduceInput struct {
	maxBatchSize int
}

func (g *numaReduceInput) Connect(ctx context.Context) error {
	return nil
}

func (n *numaReduceInput) Close(ctx context.Context) error {
	return nil
}

func (n *numaReduceInput) ReadBatch(ctx context.Context) (service.MessageBatch, service.AckFunc, error) {

	batch := make(service.MessageBatch, 0)

	return batch, func(ctx context.Context, err error) error {
		return nil
	}, nil
}

func newNumaReduceInput(conf *service.ParsedConfig) (service.BatchInput, error) {
	maxBatchSize, err := conf.FieldInt("maxBatchSize")
	if err != nil {
		return nil, err
	}
	if maxBatchSize <= 0 {
		return nil, fmt.Errorf("maxBatchSize must be greater than 0, got: %v", maxBatchSize)
	}

	return service.InputBatchedWithMaxInFlight(maxBatchSize, &numaReduceInput{maxBatchSize}), nil
}

func init() {
	err := service.RegisterBatchInput("numareduce", numareduceConfigSpec,
		func(conf *service.ParsedConfig, mgr *service.Resources) (service.BatchInput, error) {
			return newNumaReduceInput(conf)
		})
	if err != nil {
		panic(err)
	}
}
