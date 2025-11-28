package session

type ConnectionState string

const (
	ConnectionStateNew          ConnectionState = "new"
	ConnectionStateConnecting   ConnectionState = "connecting"
	ConnectionStateConnected    ConnectionState = "connected"
	ConnectionStateDisconnected ConnectionState = "disconnected"
	ConnectionStateFailed       ConnectionState = "failed"
	ConnectionStateClosed       ConnectionState = "closed"
)
