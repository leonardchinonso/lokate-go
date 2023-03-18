package interfaces

type RequestServiceInterface interface {
	Get(url string) error
}
