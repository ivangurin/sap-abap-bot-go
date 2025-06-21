package agent

const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleUser      = "user"

	ToolChoiceAuto = "auto"

	Temperature = 0.5
)

type Answer struct {
	Answer string `json:"answer"`
}
