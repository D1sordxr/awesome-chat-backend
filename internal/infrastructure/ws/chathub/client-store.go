package chathub

import "sync"

type ClientStore interface {
	Add(client *Client)
	Remove(client *Client)
	GetClients(chatID string) (map[string]*Client, bool)
}

type InMemoryClientStoreImpl struct {
	chatClients sync.Map // chatID -> *chatEntry
}

func NewInMemoryClientStoreImpl() *InMemoryClientStoreImpl {
	return &InMemoryClientStoreImpl{
		chatClients: sync.Map{},
	}
}

type chatEntry struct {
	mu      sync.Mutex
	clients map[string]*Client
}

func (i *InMemoryClientStoreImpl) Add(client *Client) {
	for _, chat := range client.chats {
		entry, _ := i.chatClients.LoadOrStore(chat, &chatEntry{clients: make(map[string]*Client)})
		e := entry.(*chatEntry)
		e.mu.Lock()
		e.clients[client.id] = client
		e.mu.Unlock()
	}
}

func (i *InMemoryClientStoreImpl) Remove(client *Client) {
	for _, chat := range client.chats {
		entry, loaded := i.chatClients.Load(chat)
		if !loaded {
			continue
		}
		e := entry.(*chatEntry)
		e.mu.Lock()
		delete(e.clients, client.id)
		if len(e.clients) == 0 {
			i.chatClients.Delete(chat)
		}
		e.mu.Unlock()
	}
}

func (i *InMemoryClientStoreImpl) GetClients(chatID string) (map[string]*Client, bool) {
	entry, loaded := i.chatClients.Load(chatID)
	if !loaded {
		return nil, false
	}
	e := entry.(*chatEntry)
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.clients, true
}
