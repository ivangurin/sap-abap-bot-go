package github

const (
	host            = "https://models.inference.ai.azure.com/"
	chatCompletions = "chat/completions"
)

type ChatCompletionRequest struct {
	Model       string                          `json:"model"`
	Messages    []*ChatCompletionRequestMessage `json:"messages"`
	Tools       []*ChatCompletionRequestTool    `json:"tools,omitempty"`
	ToolChoice  string                          `json:"tool_choice,omitempty"`
	Temperature float64                         `json:"temperature,omitempty"`
}

type ChatCompletionRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequestTool struct {
	Type     string                         `json:"type"`
	Function *ChatCompletionRequestFunction `json:"function"`
}

type ChatCompletionRequestFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ChatCompletionResponse struct {
	Error   *ChatCompletionResponseError    `json:"error,omitempty"`
	Choices []*ChatCompletionResponseChoice `json:"choices"`
}

type ChatCompletionResponseChoice struct {
	Message *ChatCompletionResponseMessage `json:"message"`
}

type ChatCompletionResponseError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code"`
}

type ChatCompletionResponseMessage struct {
	Role      string                            `json:"role"`
	Content   string                            `json:"content"`
	ToolCalls []*ChatCompletionResponseToolCall `json:"tool_calls,omitempty"`
}

type ChatCompletionResponseToolCall struct {
	ID       string                          `json:"id"`
	Type     string                          `json:"type"`
	Function *ChatCompletionResponseFunction `json:"function"`
}

type ChatCompletionResponseFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
