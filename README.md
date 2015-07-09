# MockWebServer

`MockWebServer` is a scriptable web server for testing HTTP clients.

This library makes it easy to test that your app Does The Right Thing when it makes HTTP and HTTPS calls. It lets you specify which responses to return and then verify that requests were made as expected.

Because it exercises your full HTTP stack, you can be confident that you're testing everything. You can even copy & paste HTTP responses from your real web server to create representative test cases. Or test that your code survives in awkward-to-reproduce situations like 500 errors or slow-loading responses.

Inspired by [`MockWebServer` in OkHttp](https://github.com/square/okhttp/tree/master/mockwebserver).

## Example

Use mockwebserver the same way that you use mocking frameworks:

1. Script the mocks.
2. Run application code.
3. Verify that the expected requests were made.


## Example

```go
// Create a MockWebServer. These are lean enough that you can create a new
// instance for every unit test.
s := mockwebserver.New()

// Start the server and get it's URL. You'll need this to make HTTP requests.
url := s.Start()
// Stop the server after the test finishes.
defer s.Stop()

// Enqueue some responses.
s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "https://giphy.com/gifs/sloth-WVLZLE4yGCQFi")
})

// Exercise your application code, which should make those HTTP requests.
// Responses are returned in the same order that they are enqueued.
resp, err := http.Get(url + "/foo")

// Assert your application code.
assert.Equal(t, err, nil)
assert.Equal(t, resp.StatusCode, 200)
assertBodyEqual(t, resp, "https://giphy.com/gifs/sloth-WVLZLE4yGCQFi\n")

// Optional: confirm that your app made the HTTP requests you were expecting.
request := s.TakeRequest()
assert.Equal(t, request.URL.String(), "/foo")
```

## Usage

#### type Server

```go
type Server struct {
    TestServer *httptest.Server
    Handlers   []http.HandlerFunc
    Requests   []*http.Request
    sync.Mutex
}
```

A scriptable web server.

#### func  New

```go
func New() *Server
```
New returns a new mock web server.

#### func (*Server) Start

```go
func (c *Server) Start() string
```
Start the mock web server and return the URL it is running on.

#### func (*Server) Stop

```go
func (c *Server) Stop()
```
Stop the server.

#### func (*Server) Enqueue

```go
func (s *Server) Enqueue(h http.HandlerFunc)
```
Enqueue a `http.HandlerFunc` that will be executed to a request made in sequence. The first request is served by the first enqueued handler; the second request by the second enqueued handler; and so on.

#### func (*Server) TakeRequest

```go
func (s *Server) TakeRequest() *http.Request
```
Takes the first HTTP request, removes it, and returns it. Callers should use this to verify the request was sent as intended.

# License

 MIT