// Package mockwebserver contains a scriptable web server for testing HTTP
// clients.
package mockwebserver

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// A scriptable server. It wraps an httptest.Server and lets you lets you
// specify which responses to return and then verify that requests were made as
// expected.
type Server struct {
	testServer *httptest.Server
	handlers   []http.HandlerFunc
	requests   chan *http.Request
	sync.Mutex
}

// New returns a new mock web server.
func New() *Server {
	return &Server{
		requests: make(chan *http.Request),
	}
}

// Start starts the server. The caller should call Stop when finished, to shut
// it down.
func (s *Server) Start() string {
	s.testServer = httptest.NewServer(s)
	return s.testServer.URL
}

// Stop shuts down the server and blocks until all outstanding requests on this
// server have completed.
func (s *Server) Stop() {
	s.testServer.Close()
}

// ServeHTTP satisifies the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	go func() {
		s.requests <- r
	}()

	if len(s.handlers) == 0 {
		w.WriteHeader(200)
		return
	}

	h := s.handlers[0]
	s.handlers = s.handlers[1:]
	h.ServeHTTP(w, r)
}

// Enqueue adds an `http.HandlerFunc` that will be executed to a request made in
// sequence. The first request is served by the first enqueued handler; the
// second request by the second enqueued handler; and so on.
func (s *Server) Enqueue(h http.HandlerFunc) {
	s.Lock()
	defer s.Unlock()

	s.handlers = append(s.handlers, h)
}

// TakeRequest gets the first HTTP request, removes it, and returns it. Callers
// should use this to verify the request was sent as intended. This method will
// block until the request is available, possibly forever.
func (s *Server) TakeRequest() *http.Request {
	return <-s.requests
}

// TakeRequest gets the first HTTP request (waiting up to the specified wait
// time if necessary), removes it, and returns it. Callers should use this to
// verify the request was sent as intended.
func (s *Server) TakeRequestWithTimeout(duration time.Duration) *http.Request {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(duration)
		timeout <- true
	}()
	select {
	case r := <-s.requests:
		return r
	case <-timeout:
		return nil
	}
}
