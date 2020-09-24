package srpc

// Service can implement Healthy interface to provide service health check
type Healthy interface {
	Healthy() error
}
