# MockWebServer [![GoDoc](https://godoc.org/github.com/segmentio/mockwebserver?status.svg)](https://godoc.org/github.com/segmentio/mockwebserver)

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


# License

 MIT