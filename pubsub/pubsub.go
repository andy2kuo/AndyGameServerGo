package pubsub

import (
	"context"
	"sync"
)

type Hub struct {
	sync.Mutex
	subs map[*Subscriber]struct{}
}

func (h *Hub) Subscribe(ctx context.Context, s *Subscriber) error {
	h.Lock()
	h.subs[s] = struct{}{}
	h.Unlock()

	go func() {
		select {
		case <-s.quit:
		case <-ctx.Done():
			h.Lock()
			delete(h.subs, s)
			h.Unlock()
		}
	}()

	go s.run(ctx)

	return nil
}

func (h *Hub) Publish(ctx context.Context, data interface{}) error {
	h.Lock()

	for s := range h.subs {
		s.Publish(ctx, data)
	}
	h.Unlock()

	return nil
}

func (h *Hub) Unsubscribe(s *Subscriber) error {
	h.Lock()
	delete(h.subs, s)
	h.Unlock()
	close(s.quit)
	return nil
}

func NewHub() *Hub {
	return &Hub{
		subs: make(map[*Subscriber]struct{}),
	}
}

type message interface{}

type Subscriber struct {
	sync.Mutex

	name    string
	handler chan message
	quit    chan struct{}

	callback func(interface{})
}

func (s *Subscriber) run(ctx context.Context) {
	for {
		select {
		case msg := <-s.handler:
			s.callback(msg)
		case <-s.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *Subscriber) Publish(ctx context.Context, msg message) {
	select {
	case <-ctx.Done():
		return
	case s.handler <- msg:
	default:
	}
}

func NewSubscriber(name string, _callback func(interface{})) *Subscriber {
	return &Subscriber{
		name:     name,
		handler:  make(chan message, 100),
		quit:     make(chan struct{}),
		callback: _callback,
	}
}
