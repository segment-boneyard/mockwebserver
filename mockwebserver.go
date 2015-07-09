// Package mockwebserver contains a scriptable web server for testing HTTP
// clients.
package mockwebserver

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

// A scriptable server. It wraps an httptest.Server and lets you lets you
// specify which responses to return and then verify that requests were made as
// expected.
type Server struct {
	TestServer *httptest.Server
	Handlers   []http.HandlerFunc
	Requests   []*http.Request
	sync.Mutex
}

// New returns a new mock web server.
func New() *Server {
	return &Server{}
}

// Start starts the server. The caller should call Stop when finished, to shut
// it down.
func (s *Server) Start() string {
	s.TestServer = httptest.NewServer(s)
	return s.TestServer.URL
}

// Stop shuts down the server and blocks until all outstanding requests on this
// server have completed.
func (s *Server) Stop() {
	s.TestServer.Close()
}

// ServeHTTP satisifies the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	s.Requests = append(s.Requests, r)

	if len(s.Handlers) == 0 {
		w.WriteHeader(200)
		return
	}

	h := s.Handlers[0]
	s.Handlers = s.Handlers[1:]
	h.ServeHTTP(w, r)
}

// Enqueue adds an `http.HandlerFunc` that will be executed to a request made in
// sequence. The first request is served by the first enqueued handler; the
// second request by the second enqueued handler; and so on.
func (s *Server) Enqueue(h http.HandlerFunc) {
	s.Lock()
	defer s.Unlock()

	s.Handlers = append(s.Handlers, h)
}

// TakeRequest gets the first HTTP request, removes it, and returns it.
// Callers should use this to verify the request was sent as intended.
func (s *Server) TakeRequest() *http.Request {
	s.Lock()
	defer s.Unlock()

	if len(s.Requests) == 0 {
		return nil
	}

	r := s.Requests[0]
	s.Requests = s.Requests[1:]
	return r
}
