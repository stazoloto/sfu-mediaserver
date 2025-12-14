package ws

type Controller interface {
	Handle(data []byte) error
}
