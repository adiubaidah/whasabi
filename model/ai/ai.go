package ai

type Ai struct {
	Phone       string  `json:"id"`
	Instruction string  `json:"instruction" validate:"required"`
	TopK        int32   `json:"topK" validate:"required"`
	TopP        float32 `json:"topP" validate:"required"`
	Temperature float64 `json:"temperature" validate:"required"`
}
