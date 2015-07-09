package mockwebserver_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/segmentio/mockwebserver"

	"github.com/bmizerany/assert"
)

func TestMockWebServer(t *testing.T) {
	s := mockwebserver.New()
	url := s.Start()
	defer s.Stop()

	assert.Equal(t, 0, len(s.Requests))
	assert.Equal(t, 0, len(s.Handlers))

	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "https://giphy.com/gifs/sloth-WVLZLE4yGCQFi", http.StatusInternalServerError)
	})
	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is the response you are looking for.")
	})

	assert.Equal(t, 0, len(s.Requests))
	assert.Equal(t, 2, len(s.Handlers))

	{
		resp, err := http.Get(url)
		assert.Equal(t, nil, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assertBodyEqual(t, resp, "https://giphy.com/gifs/sloth-WVLZLE4yGCQFi\n")

		assert.Equal(t, 1, len(s.Requests))
		assert.Equal(t, 1, len(s.Handlers))

		request := s.TakeRequest()
		assert.Equal(t, "/", request.URL.String())
	}

	assert.Equal(t, 0, len(s.Requests))
	assert.Equal(t, 1, len(s.Handlers))

	{
		resp, err := http.Get(url + "/foo")
		assert.Equal(t, nil, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assertBodyEqual(t, resp, "This is the response you are looking for.\n")

		assert.Equal(t, len(s.Requests), 1)
		assert.Equal(t, len(s.Handlers), 0)

		request := s.TakeRequest()
		assert.Equal(t, "/foo", request.URL.String())
	}

	assert.Equal(t, 0, len(s.Requests))
	assert.Equal(t, 0, len(s.Handlers))
}

func TestNoRegisteredHandlers(t *testing.T) {
	s := mockwebserver.New()
	url := s.Start()
	defer s.Stop()

	resp, err := http.Get(url)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assertBodyEqual(t, resp, "")

	request := s.TakeRequest()
	assert.Equal(t, "/", request.URL.String())
}

func TestNoRecordedResponses(t *testing.T) {
	s := mockwebserver.New()
	s.Start()
	defer s.Stop()

	request := s.TakeRequest()
	if request != nil {
		t.Errorf("request != nil")
	}
}

func ExampleMockWebServer() {
	// Start the server.
	s := mockwebserver.New()
	url := s.Start()
	defer s.Stop()

	// Enqeue a response
	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	})

	// Excercise your HTTP code.
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	// Inspect the request.
	req := s.TakeRequest()
	fmt.Println(req.Method)
	fmt.Println(req.URL)

	// Output:
	// Hello World!
	//
	// GET
	// /
}

func assertBodyEqual(t *testing.T, resp *http.Response, exp string) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	assert.Equal(t, exp, string(body))
}
