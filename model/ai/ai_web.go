package ai

type AiRequest struct {
	Phone       string  `json:"phone" validate:"required,number"`
	Instruction string  `json:"instruction" validate:"required"`
	TopK        int32   `json:"topK" validate:"required"`
	TopP        float32 `json:"topP" validate:"required"`
	Temperature float64 `json:"temperature" validate:"required"`
}
