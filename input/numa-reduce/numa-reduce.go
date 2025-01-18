package numareduce

import (
	"context"
	"fmt"

	"github.com/numaproj/numaflow-go/pkg/reducer"
	"github.com/redpanda-data/benthos/v4/public/service"
)

var numareduceConfigSpec = service.NewConfigSpec().
	Summary("a Numaflow reduce input").Field(service.NewIntField("maxBatchSize").Default(10))

var myInput *numaReduceInput

type numaReduceInput struct {
	maxBatchSize    int
	batchChannelIn  chan []byte
	batchChannelOut chan []byte
}

func reduce(ctx context.Context, keys []string, inputCh <-chan reducer.Datum, md reducer.Metadata) reducer.Messages {
	// sum up values for the same keys
	intervalWindow := md.IntervalWindow()
	_ = intervalWindow
	var resultKeys = keys

	myInput.batchChannelIn = make(chan []byte, myInput.maxBatchSize)
	myInput.batchChannelOut = make(chan []byte, myInput.maxBatchSize)

	for d := range inputCh {
		myInput.batchChannelIn <- d.Value()
	}
	close(myInput.batchChannelIn)

	messages := reducer.MessagesBuilder()
	for d := range myInput.batchChannelOut {
		messages = messages.Append(reducer.NewMessage(d).WithKeys(resultKeys))
	}

	return messages
}

func (g *numaReduceInput) Connect(ctx context.Context) error {
	numaServer := reducer.NewServer(reducer.SimpleCreatorWithReduceFn(reduce))
	err := numaServer.Start(context.Background())

	return err
}

func (n *numaReduceInput) Close(ctx context.Context) error {
	return nil
}

func (n *numaReduceInput) ReadBatch(ctx context.Context) (service.MessageBatch, service.AckFunc, error) {

	batch := make(service.MessageBatch, 0)
	syncStore := make([]*service.SyncResponseStore, 0)

	for m := range n.batchChannelIn {
		msg, store := service.NewMessage(m).WithSyncResponseStore()
		syncStore = append(syncStore, store)
		batch = append(batch, msg)
	}

	return batch, func(ctx context.Context, err error) error {
		for _, store := range syncStore {
			msgs := store.Read()

			for _, msg := range msgs {
				for _, m := range msg {
					data, err := m.AsBytes()

					if err != nil {
						n.batchChannelOut <- data
					}

				}
			}
		}
		close(n.batchChannelOut)
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

	myInput = &numaReduceInput{
		maxBatchSize,
		make(chan []byte),
		make(chan []byte),
	}

	return service.InputBatchedWithMaxInFlight(maxBatchSize, myInput), nil
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
