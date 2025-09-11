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

// GeminiService handles integration with Google Gemini API
type GeminiService struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
	logger     domain.Logger
}

// GeminiRequest represents request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content in Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents response from Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate response
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// NewGeminiService creates a new Gemini service
func NewGeminiService(logger domain.Logger) *GeminiService {
	// Default to Gemini Pro if no model specified
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-pro"
	}
	
	// Build the full URL with model
	baseURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)
	
	return &GeminiService{
		apiKey:  os.Getenv("GEMINI_API_KEY"),
		model:   model,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

// IsConfigured returns true if Gemini API is properly configured
func (g *GeminiService) IsConfigured() bool {
	return g.apiKey != ""
}

// AnalyzeRequirement sends requirement to Gemini for task breakdown
func (g *GeminiService) AnalyzeRequirement(ctx context.Context, req domain.TaskBreakdownRequest) (*domain.TaskBreakdownResponse, error) {
	if !g.IsConfigured() {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	prompt := g.buildAnalysisPrompt(req)
	
	response, err := g.sendRequest(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("Gemini API request failed: %w", err)
	}

	result, err := g.parseTaskBreakdown(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	g.logger.Info("Gemini analysis completed", 
		"tasks_count", len(result.Tasks),
		"confidence", result.Confidence,
		"total_estimate", result.TotalEstimate)

	return result, nil
}

// buildAnalysisPrompt creates a prompt for task analysis
func (g *GeminiService) buildAnalysisPrompt(req domain.TaskBreakdownRequest) string {
	skillsStr := strings.Join(req.TeamSkills, ", ")
	
	return fmt.Sprintf(`You are an expert software project manager and technical architect.

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
      "description": "Detailed description",
      "category": "backend|frontend|qa|devops",
      "estimate_hours": 4.5,
      "priority": 1,
      "dependencies": []
    }
  ],
  "total_estimate": 40.5,
  "recommended_team": ["Backend Developer", "Frontend Developer"],
  "critical_path": ["task_1", "task_2"],
  "risk_factors": ["Potential complexity in authentication"],
  "confidence": 0.85
}

Guidelines:
- Break down into 3-15 specific, actionable tasks
- Estimate hours realistically (consider complexity)
- Use priority 1 (high), 2 (medium), 3 (low)
- Categories: backend, frontend, qa, devops
- Include dependencies between tasks
- Confidence: 0.6-1.0 based on requirement clarity
- Consider the team's available skills

Respond only with valid JSON.`, req.Requirement, req.ProjectType, skillsStr)
}

// sendRequest sends request to Gemini API
func (g *GeminiService) sendRequest(ctx context.Context, prompt string) (string, error) {
	reqData := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", g.baseURL, g.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// parseTaskBreakdown parses Gemini's JSON response into task breakdown
func (g *GeminiService) parseTaskBreakdown(response string) (*domain.TaskBreakdownResponse, error) {
	// Clean the response to extract JSON
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	}
	response = strings.TrimSpace(response)

	var result domain.TaskBreakdownResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}

	// Validate the response
	if len(result.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks found in response")
	}

	// Generate task IDs if missing
	for i := range result.Tasks {
		if result.Tasks[i].ID == "" {
			result.Tasks[i].ID = fmt.Sprintf("task_%d_%d", time.Now().UnixNano(), i)
		}
	}

	return &result, nil
}