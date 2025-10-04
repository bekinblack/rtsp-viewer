package model

type Status int

const (
	Disconnected Status = iota
	Connected
)

type Connector interface {
	Connect(url string) error
	Disconnect()
}

type Viewer interface {
	View()
}

type Validator interface {
	Validate() error
}

type Logger interface {
	Log(msgs ...any)
}
