package request

// Fetcher defines an object that can make fetch requests
type Fetcher interface {
	Fetch(options Request) Response
}

// Request contains the information necessary to make a request to a graphql endpoint
type Request struct {
	payload []byte
}

// Response contains the response
type Response struct {
	Payload []byte
}

// NewRequest creates an request struct
func NewRequest(payload []byte) Request {
	return Request{payload: payload}
}
