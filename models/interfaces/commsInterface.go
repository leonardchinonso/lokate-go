package interfaces

// CommsServiceInterface composes the contactUs interface and the AboutServiceInterface
type CommsServiceInterface interface {
	ContactUsServiceInterface
	AboutServiceInterface
}
