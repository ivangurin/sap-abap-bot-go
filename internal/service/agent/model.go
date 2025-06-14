package agent

type SendAnswer struct {
	CorrectQuestion bool   `json:"correct_question"`
	Answer          string `json:"answer"`
}
