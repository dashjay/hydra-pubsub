package server

import "sync"

// viaHTTP, a hub manages a lot of topics by http request and response
type viaTcpServer struct {
	// port means which port this server uses
	port   int
	topics map[string]*Dispatcher
	sync.Mutex
}
