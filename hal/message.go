package hal

type Message struct {
	Message string `json:"message"`
	Resource
}

func NewMessage(msg, path string, embedded ...any) Message {
	return Message{
		Message:  msg,
		Resource: NewResource(path, embedded...),
	}
}
