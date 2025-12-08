package rp

import (
	"context"
	"sync"
)

type MemoryProducer struct {
	data []Msg

	mutex sync.Mutex

	Pinged  bool
	PingErr error
}

func (p *MemoryProducer) Ping(ctx context.Context) error {
	p.Pinged = true
	return p.PingErr
}

func (p *MemoryProducer) Publish(ctx context.Context, msg Msg) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.data = append(p.data, msg)

	return nil
}

func (p *MemoryProducer) PublishAsync(ctx context.Context, msg Msg, fn func(Msg, error)) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.data = append(p.data, msg)

	return nil
}

func (p *MemoryProducer) Flush(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.data = []Msg{}

	return nil
}

func (p *MemoryProducer) Close() {
}

func (p *MemoryProducer) Data() []Msg {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.data
}

func NewMemoryProducer() *MemoryProducer {
	return &MemoryProducer{}
}
