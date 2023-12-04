package pubsub

import "sync"

type PubSub struct {
	mu sync.Mutex
}
