package model

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
)

// ChatMessage - провайдеро-независимое сообщение диалога, которое агент
// передает выбранному ии провайдеру.
type ChatMessage struct {
	Role ChatRole
	Text string
}
