package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// OpenAIService handles integration with OpenAI ChatGPT API
type OpenAIService struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
	logger     domain.Logger
}

// OpenAIRequest represents request to OpenAI API
type OpenAIRequest struct {
	Model       string           `json:"model"`
	Messages    []OpenAIMessage  `json:"messages"`
	MaxTokens   int              `json:"max_tokens"`
	Temperature float64          `json:"temperature"`
}

// OpenAIMessage represents a message in OpenAI API format
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents response from OpenAI API
type OpenAIResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []OpenAIChoice   `json:"choices"`
	Usage   OpenAIUsage      `json:"usage"`
}

// OpenAIChoice represents a choice in OpenAI response
type OpenAIChoice struct {
	Index   int           `json:"index"`
	Message OpenAIMessage `json:"message"`
	Finish  string        `json:"finish_reason"`
}

// OpenAIUsage represents token usage in OpenAI response
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(logger domain.Logger) *OpenAIService {
	// Default to GPT-3.5-turbo if no model specified
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	
	return &OpenAIService{
		apiKey:  os.Getenv("OPENAI_API_KEY"),
		model:   model,
		baseURL: "https://api.openai.com/v1/chat/completions",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

// IsConfigured returns true if OpenAI API is properly configured
func (o *OpenAIService) IsConfigured() bool {
	return o.apiKey != ""
}

// AnalyzeRequirement sends requirement to OpenAI for task breakdown
func (o *OpenAIService) AnalyzeRequirement(ctx context.Context, req domain.TaskBreakdownRequest) (*domain.TaskBreakdownResponse, error) {
	if !o.IsConfigured() {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	prompt := o.buildAnalysisPrompt(req)
	
	response, err := o.sendRequest(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}

	result, err := o.parseTaskBreakdown(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	o.logger.Info("OpenAI analysis completed", 
		"tasks_count", len(result.Tasks),
		"confidence", result.Confidence,
		"total_estimate", result.TotalEstimate)

	return result, nil
}

// buildAnalysisPrompt creates a prompt for task analysis
func (o *OpenAIService) buildAnalysisPrompt(req domain.TaskBreakdownRequest) string {
	skillsStr := strings.Join(req.TeamSkills, ", ")
	
	return fmt.Sprintf(`You are an expert software project manager and technical architect with deep experience in project planning and estimation.

Break down this development requirement into actionable tasks:

**Requirement:** %s
**Project Type:** %s
**Team Skills:** %s

Please provide a detailed task breakdown in the following JSON format:

{
  "tasks": [
    {
      "id": "task_1",
      "title": "Task title",
      "description": "Detailed description of what needs to be done",
      "category": "backend|frontend|qa|devops",
      "estimate_hours": 4.5,
      "priority": 1,
      "dependencies": []
    }
  ],
  "total_estimate": 40.5,
  "recommended_team": ["Backend Developer", "Frontend Developer", "DevOps Engineer"],
  "critical_path": ["task_1", "task_2"],
  "risk_factors": ["Potential complexity in authentication", "Third-party API dependencies"],
  "confidence": 0.85
}

Guidelines:
- Break down into 3-15 specific, actionable tasks
- Estimate hours realistically considering complexity and potential blockers
- Use priority: 1 (high/critical), 2 (medium), 3 (low)
- Categories: backend, frontend, qa, devops
- Include task dependencies where one task blocks another
- Confidence: 0.6-1.0 based on requirement clarity and your certainty
- Consider the team's available skills when making recommendations
- Think about integration points, testing requirements, and deployment considerations

Respond ONLY with valid JSON, no additional text or formatting.`, req.Requirement, req.ProjectType, skillsStr)
}

// sendRequest sends request to OpenAI API
func (o *OpenAIService) sendRequest(ctx context.Context, prompt string) (string, error) {
	reqData := OpenAIRequest{
		Model: o.model,
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: "You are an expert software project manager and technical architect. Always respond with valid JSON only.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   4000,
		Temperature: 0.3, // Lower temperature for more consistent, focused responses
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error %d: %s", resp.StatusCode, string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

// parseTaskBreakdown parses OpenAI's JSON response into task breakdown
func (o *OpenAIService) parseTaskBreakdown(response string) (*domain.TaskBreakdownResponse, error) {
	// Clean the response to extract JSON
	response = strings.TrimSpace(response)
	
	// Remove code block markers if present
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
	}
	
	response = strings.TrimSpace(response)

	var result domain.TaskBreakdownResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response from OpenAI: %w", err)
	}

	// Validate the response
	if len(result.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks found in OpenAI response")
	}

	// Generate task IDs if missing and validate task data
	for i := range result.Tasks {
		if result.Tasks[i].ID == "" {
			result.Tasks[i].ID = fmt.Sprintf("openai_task_%d_%d", time.Now().UnixNano(), i)
		}
		
		// Set defaults for missing fields
		if result.Tasks[i].Priority == 0 {
			result.Tasks[i].Priority = 2 // Default to medium priority
		}
		if result.Tasks[i].EstimateHours == 0 {
			result.Tasks[i].EstimateHours = 4.0 // Default estimate
		}
		if result.Tasks[i].Category == "" {
			result.Tasks[i].Category = "backend" // Default category
		}
	}

	// Validate confidence score
	if result.Confidence == 0 {
		result.Confidence = 0.8 // Default confidence for OpenAI
	} else if result.Confidence > 1.0 {
		result.Confidence = 1.0
	} else if result.Confidence < 0.1 {
		result.Confidence = 0.1
	}

	return &result, nil
}