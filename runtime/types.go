package runtime

// ID is used to Get the Result of a Request
type ID int64

// Request contains the information related to a request
type Request []byte

type innerRequest struct {
	id       ID
	request  Request
	response Response
}

// Response contains the information related to a result
type Response interface {
	Read() ([]byte, error)
	write([]byte) error
}

type innerResponse struct {
	id     ID
	result Response
}

// RequestExecutor knows how to Put Requests and Get Responses
type RequestExecutor interface {
	Put(request Request) (Response, error)
}
