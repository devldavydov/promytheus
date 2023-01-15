package metrics

type Publisher interface {
	Publish(metrics Metrics) error
}

type HttpPublisher struct {
}

func NewHttpPublisher() *HttpPublisher {
	return &HttpPublisher{}
}

func (httpPublisher *HttpPublisher) Publish(Metrics Metrics) error {
	return nil
}
