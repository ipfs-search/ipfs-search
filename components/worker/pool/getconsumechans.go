package pool

import (
	"context"

	samqp "github.com/rabbitmq/amqp091-go"
)

type consumeFunc func(context.Context) (<-chan samqp.Delivery, error)

const consumerCount = 3

func (p *Pool) getConsumeChans(ctx context.Context) (*consumeChans, error) {
	queues, err := p.getQueues(ctx)
	if err != nil {
		return nil, err
	}

	var consumeFuncs = [consumerCount]consumeFunc{queues.Files.Consume, queues.Directories.Consume, queues.Hashes.Consume}
	var chans [consumerCount]<-chan samqp.Delivery

	for i, f := range consumeFuncs {
		chans[i], err = f(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &consumeChans{
		Files:       chans[0],
		Directories: chans[1],
		Hashes:      chans[2],
	}, nil
}
