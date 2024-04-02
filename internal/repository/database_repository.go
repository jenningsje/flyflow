package repository

import (
	"encoding/json"
	"github.com/flyflow-devs/flyflow/internal/logger"
	"github.com/flyflow-devs/flyflow/internal/models"
	"github.com/flyflow-devs/flyflow/internal/requests"
	"gorm.io/gorm"
	"strings"
	"sync"
	"github.com/pandodao/tokenizer-go"

)

type DatabaseRepository struct {
	DB          *gorm.DB
	APIKeyCache sync.Map
	repo        Repository
}

func NewDatabaseRepository(db *gorm.DB, r Repository) *DatabaseRepository {
	repo := &DatabaseRepository{
		DB:   db,
		repo: r,
	}
	return repo
}

func (dr *DatabaseRepository) SaveQueryRecord(req *requests.CompletionRequest, resp *requests.CompletionResponse) error {
	jsonData, err := json.Marshal(req.Cr)
	if err != nil {
		return err
	}

	if !req.Cr.Stream {
		var openAIResp requests.OpenAIResponse
		err = json.Unmarshal([]byte(resp.Response), &openAIResp)
		if err != nil {
			return err
		}

		var responseText string
		if len(openAIResp.Choices) > 0 {
			responseText = openAIResp.Choices[0].Message.Content
		}

		queryRecord := &models.QueryRecord{
			Request:         string(jsonData),
			Response:        responseText,
			RequestedModel:  req.Cr.Model,
			MaxTokens:       req.Cr.MaxTokens,
			InputTokens:     openAIResp.Usage.PromptTokens,
			OutputTokens:    openAIResp.Usage.CompletionTokens,
			Stream:          req.Cr.Stream,
			APIKey:          req.APIKey,
			Tags:            req.Cr.Tags,
		}

		return dr.DB.Create(queryRecord).Error
	} else {
		// Streaming version
		// Initialize variables to store the response data
		var responseBuilder strings.Builder
		var inputTokens, outputTokens int

		// Split the response data by newline
		lines := strings.Split(resp.Response, "\n")

		// Iterate over each line of the response
		for _, line := range lines {
			// Check if the line starts with "data: "
			if strings.HasPrefix(line, "data: ") {
				// Remove the "data: " prefix
				jsonData := strings.TrimPrefix(line, "data: ")

				// Check if the line indicates the end of the stream
				if jsonData == "[DONE]" {
					break
				}

				// Unmarshal the JSON data into a StreamingData struct
				var streamingData requests.StreamingData
				err := json.Unmarshal([]byte(jsonData), &streamingData)
				if err != nil {
					logger.S.Error("error unmarshalling StreamingData in SaveQueryRecord", err)
					continue
				}

				// Append the content to the responseBuilder
				if len(streamingData.Choices) > 0 {
					responseBuilder.WriteString(streamingData.Choices[0].Delta.Content)
				}

				outputTokens += 1
			}
		}

		var totalMessage string
		for _, msg := range req.Cr.Messages {
			totalMessage += msg.Content.(string)
		}

		queryRecord := &models.QueryRecord{
			Request:        string(jsonData),
			Response:       responseBuilder.String(),
			RequestedModel: req.Cr.Model,
			MaxTokens:      req.Cr.MaxTokens,
			InputTokens:    tokenizer.MustCalToken(totalMessage),
			OutputTokens:   outputTokens,
			Stream:         req.Cr.Stream,
			APIKey:         req.APIKey,
			Tags:           req.Cr.Tags,
		}

		return dr.DB.Create(queryRecord).Error
	}
}

func (dr *DatabaseRepository) ProxyRequest(r *requests.ProxyRequest) error {
	return dr.repo.ProxyRequest(r)
}

func (dr *DatabaseRepository) ChatCompletion(r *requests.CompletionRequest) (*requests.CompletionResponse, error) {
	resp, err := dr.repo.ChatCompletion(r)
	if err != nil {
		return nil, err
	}

	if resp.ShouldSave {
		err = dr.SaveQueryRecord(r, resp)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
