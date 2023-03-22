package datasource

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Get makes an HTTP Get request to the url provided
func Get(url string, resp interface{}) error {
	log.Printf("INFO: making a request to TransportAPI with URL: %s\n", url)

	// make a Get request to the url
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make http request to url: %s. Error: %v", url, err)
		return err
	}

	// read all the response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read body from response. Error: %v", err)
		return err
	}

	// close the body after reading
	err = res.Body.Close()
	if err != nil {
		log.Printf("Failed to close response body. Error: %v", err)
		return err
	}

	// unmarshal the raw byte data into the interface
	err = json.Unmarshal(data, resp)
	if err != nil {
		log.Printf("Error unmarshalling data into interface. Error: %v", err)
		return err
	}

	return nil
}
