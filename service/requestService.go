package service

import (
	"fmt"
	"github.com/leonardchinonso/lokate-go/models/interfaces"
	"io"
	"log"
	"net/http"
)

type req struct{}

// NewRequestService returns an interface for all external request service methods
func NewRequestService() interfaces.RequestServiceInterface {
	return &req{}
}

func (r *req) Get(url string) error {
	// make the http request to the url
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error sending get request to url: %v. Error: %v", url, err)
		return err
	}

	// read the body of the response
	data, _ := io.ReadAll(res.Body)

	// close response body to avoid leak
	res.Body.Close()

	fmt.Printf("%s\n", data)

	return nil
}
