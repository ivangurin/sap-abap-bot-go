package agent

const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleUser      = "user"

	ToolChoiceAuto = "auto"

	Temperature = 1
)

type Answer struct {
	Answer string `json:"answer"`
}
