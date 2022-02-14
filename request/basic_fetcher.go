package request

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

// DefaultFetcher is a basic implementation fo the Fetcher interface
type DefaultFetcher struct {
	url string
}

// NewDefaultFetcher constructs a DefaultFetcher
func NewDefaultFetcher(url string) DefaultFetcher {
	return DefaultFetcher{url: url}
}

// Fetch uses a basic fetch strategy to query the server and return a response
func (fetcher DefaultFetcher) Fetch(request Request) Response {
	//Encode the data
	responseBody := bytes.NewBuffer(request.payload)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(fetcher.url, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return Response{body}
}
