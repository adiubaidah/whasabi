package service

import (
	"adiubaidah/adi-bot/app"
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/middleware"
	"adiubaidah/adi-bot/model"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/generative-ai-go/genai"
	"gorm.io/gorm"
)

type AiServiceImpl struct {
	Client   *genai.Client
	Validate *validator.Validate
	DB       *gorm.DB
}

func NewAiService(client *genai.Client, db *gorm.DB, validate *validator.Validate) AiService {
	return &AiServiceImpl{
		Client:   client,
		Validate: validate,
		DB:       db,
	}
}

func (service *AiServiceImpl) GetModel(ctx context.Context) *model.Ai {
	// Extract user information from context
	user := ctx.Value(middleware.UserContext).(jwt.MapClaims)
	userID := uint(user["id"].(float64)) // Convert ID to uint from float64

	// Find AI model by user ID
	aiModel := &model.Ai{}
	result := service.DB.Where("user_id = ?", userID).Take(&aiModel)

	// Check if record was not found
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// Panic on unexpected errors
		panic(result.Error)
	}

	return aiModel

}

func (service *AiServiceImpl) UpsertModel(ctx context.Context, modelAi model.CreateAIModel) *model.Ai {
	// Extract user information from context
	user := ctx.Value(middleware.UserContext).(jwt.MapClaims)
	userID := uint(user["id"].(float64)) // Convert ID to uint from float64

	// Find AI model by user ID
	aiModel := &model.Ai{}
	result := service.DB.Where("user_id = ?", userID).Take(&aiModel)

	// Check if record was not found
	var err error
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// Panic on unexpected errors
		panic(result.Error)
	}

	if result.RowsAffected == 0 {
		// If no AI model exists for this user, create a new one
		fmt.Println("Create new AI model")
		newAiModel := &model.Ai{
			UserID: userID,
			CreateAIModel: model.CreateAIModel{
				Name:        modelAi.Name,
				Phone:       modelAi.Phone,
				Instruction: modelAi.Instruction,
				Temperature: modelAi.Temperature,
				TopK:        modelAi.TopK,
				TopP:        modelAi.TopP,
			},
			IsActive:        false,
			IsAuthenticated: false,
		}

		err = service.DB.Create(&newAiModel).Error
		aiModel = newAiModel // Assign new AI model to return later
	} else {
		// Update existing AI model
		fmt.Println("Update existing AI model")
		err = service.DB.Model(&aiModel).Where("phone = ?", modelAi.Phone).Updates(&model.CreateAIModel{
			Name:        modelAi.Name,
			Phone:       modelAi.Phone,
			Instruction: modelAi.Instruction,
			Temperature: modelAi.Temperature,
			TopK:        modelAi.TopK,
			TopP:        modelAi.TopP,
		}).Error
	}

	helper.PanicIfError("Error while upsert", err)

	// Return the updated/new AI model
	return aiModel
}

func (service *AiServiceImpl) GenerateResponse(ctx context.Context, modelAi *model.Ai, histories *[]model.History, input string) string {
	var sessionHistory []*genai.Content
	for _, history := range *histories {
		sessionHistory = append(sessionHistory, &genai.Content{
			Role:  history.RoleAs,
			Parts: []genai.Part{genai.Text(history.Content)},
		})
	}

	// Generate response
	option := app.AiModelOption{
		Instruction: modelAi.Instruction,
		TopK:        modelAi.TopK,
		TopP:        modelAi.TopP,
		Temperature: modelAi.Temperature,
	}

	model := app.GetAIModel(service.Client, &option)
	session := model.StartChat()

	if (sessionHistory != nil) && (len(sessionHistory) > 0) {
		session.History = sessionHistory
	}

	resp, err := session.SendMessage(ctx, genai.Text(input))
	helper.PanicIfError("Error saat mengambil respon:", err)

	return string(resp.Candidates[0].Content.Parts[0].(genai.Text))
}
