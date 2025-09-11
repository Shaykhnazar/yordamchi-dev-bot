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

// ClaudeService handles integration with Claude.ai API
type ClaudeService struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
	logger     domain.Logger
}

// ClaudeRequest represents request to Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeMessage represents a message in Claude API format
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents response from Claude API
type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
	ID      string          `json:"id"`
	Model   string          `json:"model"`
	Role    string          `json:"role"`
	Usage   ClaudeUsage     `json:"usage"`
}

// ClaudeContent represents content in Claude response
type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeUsage represents token usage in Claude response
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// NewClaudeService creates a new Claude service
func NewClaudeService(logger domain.Logger) *ClaudeService {
	return &ClaudeService{
		apiKey:  os.Getenv("CLAUDE_API_KEY"),
		baseURL: "https://api.anthropic.com/v1/messages",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

// IsConfigured returns true if Claude API is properly configured
func (c *ClaudeService) IsConfigured() bool {
	return c.apiKey != ""
}

// AnalyzeRequirement sends requirement to Claude for task breakdown
func (c *ClaudeService) AnalyzeRequirement(ctx context.Context, req domain.TaskBreakdownRequest) (*domain.TaskBreakdownResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("Claude API key not configured")
	}

	prompt := c.buildAnalysisPrompt(req)
	
	response, err := c.sendRequest(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("Claude API request failed: %w", err)
	}

	result, err := c.parseTaskBreakdown(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Claude response: %w", err)
	}

	c.logger.Info("Claude analysis completed", 
		"tasks_count", len(result.Tasks),
		"confidence", result.Confidence,
		"total_estimate", result.TotalEstimate)

	return result, nil
}

// buildAnalysisPrompt creates a prompt for task analysis
func (c *ClaudeService) buildAnalysisPrompt(req domain.TaskBreakdownRequest) string {
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

// sendRequest sends request to Claude API
func (c *ClaudeService) sendRequest(ctx context.Context, prompt string) (string, error) {
	reqData := ClaudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 4000,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
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

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude")
	}

	return claudeResp.Content[0].Text, nil
}

// parseTaskBreakdown parses Claude's JSON response into task breakdown
func (c *ClaudeService) parseTaskBreakdown(response string) (*domain.TaskBreakdownResponse, error) {
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