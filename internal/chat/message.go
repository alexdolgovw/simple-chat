package chat

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}
