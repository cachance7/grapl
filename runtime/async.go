package runtime

import (
	"github.com/cachance7/grapl/request"
	log "github.com/sirupsen/logrus"
)

// AsyncResponse knows how to retrieve the result of a call
type AsyncResponse struct {
	ch chan ([]byte)
}

func newAsyncResponse() AsyncResponse {
	return AsyncResponse{ch: make(chan []byte, 1)}
}

func (res AsyncResponse) Read() ([]byte, error) {
	return <-res.ch, nil
}

func (res AsyncResponse) write(data []byte) error {
	res.ch <- data
	return nil
}

// AsyncRequestExecutor will make requests and store results
type AsyncRequestExecutor struct {
	next    ID
	fetcher request.Fetcher
	in      chan innerRequest
}

// NewAsyncRequestExecutor creates and returns a new AsyncRequestExecutor
func NewAsyncRequestExecutor(fetcher request.Fetcher) AsyncRequestExecutor {
	return AsyncRequestExecutor{
		next:    0,
		fetcher: fetcher,
		in:      make(chan innerRequest, 100),
	}
}

func (executor AsyncRequestExecutor) loop() {
	log.Println("starting loop")
	for {
		msg, ok := <-executor.in
		if !ok {
			log.Println("could not read from in queue; terminating")
			break
		}
		res := executor.fetcher.Fetch(request.NewRequest(msg.request))
		msg.response.write(res.Payload)
	}
}

// Start begins executor processing
func (executor AsyncRequestExecutor) Start() {
	go executor.loop()
}

// Stop terminates executor processing
func (executor AsyncRequestExecutor) Stop() {
	// Do something here. Maybe enqueue terminate message
}

func (executor *AsyncRequestExecutor) nextID() ID {
	executor.next++
	nextID := executor.next
	return nextID
}

// Put enqueues a request with this executor and returns an ID. This function will not block.
func (executor AsyncRequestExecutor) Put(req Request) (Response, error) {
	newID := executor.nextID()
	res := newAsyncResponse()
	executor.in <- innerRequest{id: newID, request: req, response: res}
	return res, nil
}
