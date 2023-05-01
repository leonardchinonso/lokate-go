package interfaces

// RequestServiceInterface represents an interface for the requests
type RequestServiceInterface interface {
	Get(url string) error
}
