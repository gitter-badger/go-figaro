package figp2p

import (
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
)

const (
	nodeIDTopic        = "nodeID"
	newConnectionTopic = "newConnection"
)

// Handler receives a Message
type Handler func(*Node, *floodsub.Message)

// FanoutMux handles multiplexing handlers and subscriptions for a Node
type FanoutMux struct {
	rwMutex     sync.RWMutex
	handlersMap map[string][]Handler
}

// NewFanoutMux returns a new FanoutMux
func NewFanoutMux() *FanoutMux {
	return &FanoutMux{
		rwMutex:     sync.RWMutex{},
		handlersMap: make(map[string][]Handler),
	}
}

// Delete removes a topic
func (hm *FanoutMux) Delete(topicName string) {
	hm.rwMutex.Lock()
	delete(hm.handlersMap, topicName)
	hm.rwMutex.Unlock()
}

// Get returns the registred handlers for a topic
func (hm *FanoutMux) Get(topicName string) []Handler {
	hm.rwMutex.RLock()
	handlers := hm.handlersMap[topicName]
	hm.rwMutex.RUnlock()
	return handlers
}

// Handle registers a new handler for a topic
func (hm *FanoutMux) Handle(topicName string, handler Handler) {
	hm.rwMutex.Lock()
	hm.handlersMap[topicName] = append(hm.handlersMap[topicName], handler)
	hm.rwMutex.Unlock()
}

// HandleDirectMessage attaches a Handler for and subcribes to a Topic
func (hm *FanoutMux) HandleDirectMessage(handler Handler) {
	hm.rwMutex.Lock()
	hm.handlersMap[nodeIDTopic] = append(hm.handlersMap[nodeIDTopic], handler)
	hm.rwMutex.Unlock()
}

// HandleNewConnection attaches a Handler for and subcribes to a Topic
func (hm *FanoutMux) HandleNewConnection(handler Handler) {
	hm.rwMutex.Lock()
	hm.handlersMap[newConnectionTopic] = append(hm.handlersMap[newConnectionTopic], handler)
	hm.rwMutex.Unlock()
}

// Keys returns the topics that are currectly registered
func (hm *FanoutMux) Keys() []string {
	var keys []string
	hm.rwMutex.RLock()
	for key := range hm.handlersMap {
		keys = append(keys, key)
	}
	hm.rwMutex.RUnlock()
	return keys
}

// Set replaces the handlers for a topic
func (hm *FanoutMux) Set(topicName string, handlers []Handler) {
	hm.rwMutex.Lock()
	hm.handlersMap[topicName] = handlers
	hm.rwMutex.Unlock()
}
