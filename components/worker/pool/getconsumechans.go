package pool

import (
	"context"

	samqp "github.com/rabbitmq/amqp091-go"
)

type consumeFunc func(context.Context) (<-chan samqp.Delivery, error)

func (p *Pool) getConsumeChans(ctx context.Context) (*consumeChans, error) {
	queues, err := p.getQueues(ctx)
	if err != nil {
		return nil, err
	}

	// Note: Manually adjust indexCount whenever the amount of indexes change
	const consumeChanCnt = 3
	var (
		consumeFuncs = [consumeChanCnt]consumeFunc{queues.Files.Consume, queues.Directories.Consume, queues.Hashes.Consume}
		chans        [consumeChanCnt]<-chan samqp.Delivery
	)

	for i, f := range consumeFuncs {
		chans[i], err = f(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Note: Manually adjust order here!
	return &consumeChans{
		Files:       chans[0],
		Directories: chans[1],
		Hashes:      chans[2],
	}, nil
}
