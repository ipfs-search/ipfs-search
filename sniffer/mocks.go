package sniffer

import (
	t "github.com/ipfs-search/ipfs-search/types"
	"log"
	"time"
)

type mockLogger struct {
	wait time.Duration
	msgs chan map[string]interface{}
	errc chan error
}

func (m mockLogger) Next() (map[string]interface{}, error) {
	log.Printf("Call to Next()")

	time.Sleep(m.wait)

	select {
	case msg := <-m.msgs:
		return msg, nil
	case err := <-m.errc:
		return nil, err
	}
}

func (m mockLogger) Close() error {
	return nil
}

func newMockLogger() mockLogger {
	msgs := make(chan map[string]interface{}, 1)
	errc := make(chan error, 1)

	return mockLogger{
		msgs: msgs,
		errc: errc,
	}
}

type mockExtractor struct {
	provider *t.Provider
	err      error
}

func (m mockExtractor) Extract(map[string]interface{}) (*t.Provider, error) {
	return m.provider, m.err
}

type mockQueue struct {
	err        error
	pubs       chan interface{}
	priorities chan uint8
}

func (m mockQueue) Publish(pub interface{}, priority uint8) error {
	log.Printf("Mock publishing %v (priority %d)", pub, priority)

	m.pubs <- pub
	m.priorities <- priority

	return m.err
}

func mockProvider() t.Provider {
	resource := &t.Resource{
		Protocol: "ipfs",
		Id:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
	}

	return t.Provider{
		Resource: resource,
		Date:     time.Now(),
		Provider: "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
	}
}
