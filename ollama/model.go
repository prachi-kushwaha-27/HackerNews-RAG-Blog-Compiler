package ollama

type InstructRequest struct {
	Content     string   `json:"content"`
	ImagesB64   []string `json:"imagesB64"`
	Instruction string   `json:"instruction"`
}

type InstructResponse struct {
	Result string `json:"result"`
}
