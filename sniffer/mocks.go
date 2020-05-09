package sniffer

import (
	t "github.com/ipfs-search/ipfs-search/types"
	"time"
)

type mockLogger struct {
	wait   time.Duration
	msg    map[string]interface{}
	err    error
	closed bool
}

func (m mockLogger) Next() (map[string]interface{}, error) {
	time.Sleep(m.wait)

	return m.msg, m.err
}

func (m mockLogger) Close() error {
	m.closed = true
	return nil
}

type mockExtractor struct {
	provider *t.Provider
	err      error
}

func (m mockExtractor) Extract(map[string]interface{}) (*t.Provider, error) {
	return m.provider, m.err
}

type mockQueue struct {
	msg       interface{}
	priority  uint8
	err       error
	publishes int
}

func (m mockQueue) Publish(msg interface{}, priority uint8) error {
	m.msg = msg
	m.priority = priority
	m.publishes++
	return m.err
}
