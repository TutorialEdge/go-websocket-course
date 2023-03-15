package stream

// A pool represents all websocket clients on our service
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Channels   map[string][]*Client
	Clients    []*Client
}
