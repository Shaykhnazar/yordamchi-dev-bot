# 📅 Week 3: AI Integration & Smart Features

## 🎯 Learning Objectives

By the end of Week 3, you will master:
- OpenAI API integration and prompt engineering
- Natural language processing for code understanding
- Machine learning model integration in Go
- Advanced analytics and user behavior tracking
- Real-time recommendation systems
- AI-powered code review and suggestions
- Personalized learning path generation

## 📊 Week Overview

| Day | Focus | Key Features | AI Integration |
|-----|-------|-------------|----------------|
| 15 | OpenAI Integration | Code explanation, generation | GPT-4/3.5 API |
| 16 | AI Code Review | Automated reviews, suggestions | Custom prompts |
| 17 | Learning Tracker | Progress tracking, recommendations | ML analytics |
| 18 | Tech News AI | Curated news, summarization | NLP processing |
| 19 | Smart Recommendations | Personalized suggestions | Recommendation engine |
| 20 | Advanced AI Features | Context awareness, memory | Long-term context |
| 21 | AI Integration Testing | End-to-end AI workflows | Model validation |

---

## 📅 Day 15: OpenAI API Integration

### 🎯 Goals
- Implement OpenAI API client with streaming support
- Create code explanation and generation features
- Add intelligent error analysis and debugging suggestions
- Build prompt engineering system for optimal results

### 🔧 OpenAI Service Implementation

#### `internal/services/ai/openai.go`
```go
package ai

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

type OpenAIClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
    logger     Logger
}

type Logger interface {
    Printf(format string, args ...interface{})
    Println(args ...interface{})
}

// OpenAI API request/response types
type ChatCompletionRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    MaxTokens   int       `json:"max_tokens,omitempty"`
    Temperature float64   `json:"temperature,omitempty"`
    Stream      bool      `json:"stream,omitempty"`
    TopP        float64   `json:"top_p,omitempty"`
    Stop        []string  `json:"stop,omitempty"`
}

type Message struct {
    Role    string `json:"role"`    // "system", "user", "assistant"
    Content string `json:"content"`
}

type ChatCompletionResponse struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}

type Choice struct {
    Index        int     `json:"index"`
    Message      Message `json:"message"`
    FinishReason string  `json:"finish_reason"`
}

type Usage struct {
    PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`
}

type StreamResponse struct {
    ID      string        `json:"id"`
    Object  string        `json:"object"`
    Created int64         `json:"created"`
    Model   string        `json:"model"`
    Choices []StreamChoice `json:"choices"`
}

type StreamChoice struct {
    Index int           `json:"index"`
    Delta StreamMessage `json:"delta"`
    FinishReason *string `json:"finish_reason"`
}

type StreamMessage struct {
    Role    string `json:"role,omitempty"`
    Content string `json:"content,omitempty"`
}

// AI Service response types
type CodeExplanation struct {
    Summary     string              `json:"summary"`
    StepByStep  []ExplanationStep   `json:"step_by_step"`
    KeyConcepts []string            `json:"key_concepts"`
    Language    string              `json:"language"`
    Complexity  string              `json:"complexity"`
    Suggestions []string            `json:"suggestions"`
    Examples    []CodeExample       `json:"examples"`
}

type ExplanationStep struct {
    LineNumbers string `json:"line_numbers"`
    Code        string `json:"code"`
    Explanation string `json:"explanation"`
    Purpose     string `json:"purpose"`
}

type CodeExample struct {
    Title       string `json:"title"`
    Code        string `json:"code"`
    Description string `json:"description"`
}

type CodeGeneration struct {
    Code        string   `json:"code"`
    Language    string   `json:"language"`
    Description string   `json:"description"`
    Examples    []string `json:"examples"`
    Tests       string   `json:"tests"`
    Documentation string `json:"documentation"`
}

type ErrorAnalysis struct {
    ErrorType     string   `json:"error_type"`
    PossibleCause string   `json:"possible_cause"`
    Solutions     []string `json:"solutions"`
    Prevention    []string `json:"prevention"`
    RelatedTopics []string `json:"related_topics"`
}

func NewOpenAIClient(apiKey string, logger Logger) *OpenAIClient {
    return &OpenAIClient{
        apiKey:  apiKey,
        baseURL: "https://api.openai.com/v1",
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
        logger: logger,
    }
}

func (c *OpenAIClient) ExplainCode(ctx context.Context, code, language string) (*CodeExplanation, error) {
    prompt := c.buildCodeExplanationPrompt(code, language)
    
    messages := []Message{
        {
            Role:    "system",
            Content: "You are an expert programming instructor who explains code clearly and comprehensively. Always provide structured, educational explanations.",
        },
        {
            Role:    "user",
            Content: prompt,
        },
    }

    response, err := c.chatCompletion(ctx, messages, "gpt-4", 2000, 0.3)
    if err != nil {
        return nil, fmt.Errorf("OpenAI API error: %w", err)
    }

    explanation, err := c.parseCodeExplanation(response.Choices[0].Message.Content, language)
    if err != nil {
        return nil, fmt.Errorf("parse explanation: %w", err)
    }

    c.logger.Printf("🤖 Code explained: %s (%d tokens)", language, response.Usage.TotalTokens)
    return explanation, nil
}

func (c *OpenAIClient) GenerateCode(ctx context.Context, description, language string) (*CodeGeneration, error) {
    prompt := c.buildCodeGenerationPrompt(description, language)
    
    messages := []Message{
        {
            Role:    "system",
            Content: "You are an expert programmer. Generate clean, well-documented, production-ready code with proper error handling and best practices.",
        },
        {
            Role:    "user",
            Content: prompt,
        },
    }

    response, err := c.chatCompletion(ctx, messages, "gpt-4", 2500, 0.4)
    if err != nil {
        return nil, fmt.Errorf("OpenAI API error: %w", err)
    }

    generation, err := c.parseCodeGeneration(response.Choices[0].Message.Content, language)
    if err != nil {
        return nil, fmt.Errorf("parse generation: %w", err)
    }

    c.logger.Printf("🤖 Code generated: %s (%d tokens)", language, response.Usage.TotalTokens)
    return generation, nil
}

func (c *OpenAIClient) AnalyzeError(ctx context.Context, errorMsg, code, language string) (*ErrorAnalysis, error) {
    prompt := c.buildErrorAnalysisPrompt(errorMsg, code, language)
    
    messages := []Message{
        {
            Role:    "system",
            Content: "You are a debugging expert. Analyze errors thoroughly and provide actionable solutions with clear explanations.",
        },
        {
            Role:    "user",
            Content: prompt,
        },
    }

    response, err := c.chatCompletion(ctx, messages, "gpt-3.5-turbo", 1500, 0.2)
    if err != nil {
        return nil, fmt.Errorf("OpenAI API error: %w", err)
    }

    analysis, err := c.parseErrorAnalysis(response.Choices[0].Message.Content)
    if err != nil {
        return nil, fmt.Errorf("parse analysis: %w", err)
    }

    c.logger.Printf("🤖 Error analyzed: %s", errorMsg[:min(50, len(errorMsg))])
    return analysis, nil
}

func (c *OpenAIClient) StreamCompletion(ctx context.Context, messages []Message, callback func(string)) error {
    req := ChatCompletionRequest{
        Model:       "gpt-3.5-turbo",
        Messages:    messages,
        MaxTokens:   1500,
        Temperature: 0.7,
        Stream:      true,
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return fmt.Errorf("marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", strings.NewReader(string(jsonData)))
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }

    c.setHeaders(httpReq)

    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return fmt.Errorf("execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("OpenAI API error: %d", resp.StatusCode)
    }

    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "data: ") {
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                break
            }

            var streamResp StreamResponse
            if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
                continue
            }

            if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
                callback(streamResp.Choices[0].Delta.Content)
            }
        }
    }

    return scanner.Err()
}

func (c *OpenAIClient) chatCompletion(ctx context.Context, messages []Message, model string, maxTokens int, temperature float64) (*ChatCompletionResponse, error) {
    req := ChatCompletionRequest{
        Model:       model,
        Messages:    messages,
        MaxTokens:   maxTokens,
        Temperature: temperature,
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    c.setHeaders(httpReq)

    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(body))
    }

    var response ChatCompletionResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }

    return &response, nil
}

func (c *OpenAIClient) setHeaders(req *http.Request) {
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("User-Agent", "DevMate-Bot/1.0")
}

// Prompt engineering functions
func (c *OpenAIClient) buildCodeExplanationPrompt(code, language string) string {
    return fmt.Sprintf(`Please explain this %s code in detail:

```%s
%s
```

Provide a comprehensive explanation including:
1. A brief summary of what the code does
2. Step-by-step breakdown of each significant part
3. Key programming concepts demonstrated
4. Complexity level (beginner/intermediate/advanced)
5. Suggestions for improvement if any
6. Related examples or variations

Format your response as structured text that's easy to read and understand.`, language, language, code)
}

func (c *OpenAIClient) buildCodeGenerationPrompt(description, language string) string {
    return fmt.Sprintf(`Generate %s code for the following requirements:

Requirements: %s

Please provide:
1. Clean, well-commented code following best practices
2. Proper error handling where applicable
3. Brief explanation of the approach
4. Usage examples
5. Basic unit tests if relevant
6. Documentation comments

Make the code production-ready and follow %s conventions.`, language, description, language)
}

func (c *OpenAIClient) buildErrorAnalysisPrompt(errorMsg, code, language string) string {
    return fmt.Sprintf(`Analyze this %s error:

Error: %s

Code context:
```%s
%s
```

Please provide:
1. What type of error this is
2. Most likely cause of the error
3. Step-by-step solutions to fix it
4. How to prevent similar errors in the future
5. Related concepts to learn

Be specific and actionable in your recommendations.`, language, errorMsg, language, code)
}

// Parsing functions for structured responses
func (c *OpenAIClient) parseCodeExplanation(content, language string) (*CodeExplanation, error) {
    // Simplified parsing - in production, use more sophisticated NLP
    explanation := &CodeExplanation{
        Language:    language,
        Summary:     c.extractSection(content, "summary", "brief summary"),
        Complexity:  c.extractComplexity(content),
        KeyConcepts: c.extractKeyConcepts(content),
        Suggestions: c.extractSuggestions(content),
        StepByStep:  c.extractSteps(content),
    }

    return explanation, nil
}

func (c *OpenAIClient) parseCodeGeneration(content, language string) (*CodeGeneration, error) {
    generation := &CodeGeneration{
        Language:      language,
        Code:          c.extractCodeBlock(content),
        Description:   c.extractSection(content, "explanation", "approach"),
        Examples:      c.extractExamples(content),
        Tests:         c.extractTests(content),
        Documentation: c.extractDocumentation(content),
    }

    return generation, nil
}

func (c *OpenAIClient) parseErrorAnalysis(content string) (*ErrorAnalysis, error) {
    analysis := &ErrorAnalysis{
        ErrorType:     c.extractSection(content, "type of error", "error type"),
        PossibleCause: c.extractSection(content, "cause", "likely cause"),
        Solutions:     c.extractSolutions(content),
        Prevention:    c.extractPrevention(content),
        RelatedTopics: c.extractRelatedTopics(content),
    }

    return analysis, nil
}

// Helper parsing functions
func (c *OpenAIClient) extractSection(content, keyword1, keyword2 string) string {
    // Simplified text extraction - in production use better NLP
    lines := strings.Split(content, "\n")
    for i, line := range lines {
        lowerLine := strings.ToLower(line)
        if strings.Contains(lowerLine, keyword1) || strings.Contains(lowerLine, keyword2) {
            if i+1 < len(lines) {
                return strings.TrimSpace(lines[i+1])
            }
        }
    }
    return ""
}

func (c *OpenAIClient) extractComplexity(content string) string {
    content = strings.ToLower(content)
    if strings.Contains(content, "beginner") {
        return "beginner"
    } else if strings.Contains(content, "advanced") {
        return "advanced"
    } else if strings.Contains(content, "intermediate") {
        return "intermediate"
    }
    return "intermediate"
}

func (c *OpenAIClient) extractKeyConcepts(content string) []string {
    // Extract concepts from structured response
    var concepts []string
    lines := strings.Split(content, "\n")
    
    inConceptsSection := false
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.Contains(strings.ToLower(line), "concept") {
            inConceptsSection = true
            continue
        }
        
        if inConceptsSection && strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
            concept := strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*")
            concepts = append(concepts, strings.TrimSpace(concept))
        } else if inConceptsSection && line == "" {
            break
        }
    }
    
    return concepts
}

func (c *OpenAIClient) extractCodeBlock(content string) string {
    // Extract code from markdown code blocks
    lines := strings.Split(content, "\n")
    var codeLines []string
    inCodeBlock := false
    
    for _, line := range lines {
        if strings.HasPrefix(line, "```") {
            if inCodeBlock {
                break
            }
            inCodeBlock = true
            continue
        }
        
        if inCodeBlock {
            codeLines = append(codeLines, line)
        }
    }
    
    return strings.Join(codeLines, "\n")
}

func (c *OpenAIClient) extractSteps(content string) []ExplanationStep {
    // Extract step-by-step breakdown
    var steps []ExplanationStep
    lines := strings.Split(content, "\n")
    
    for i, line := range lines {
        if strings.Contains(strings.ToLower(line), "step") && strings.Contains(line, ":") {
            step := ExplanationStep{
                LineNumbers: fmt.Sprintf("%d", i+1),
                Explanation: strings.TrimSpace(strings.Split(line, ":")[1]),
            }
            steps = append(steps, step)
        }
    }
    
    return steps
}

func (c *OpenAIClient) extractSuggestions(content string) []string {
    return c.extractListItems(content, "suggestion", "improvement")
}

func (c *OpenAIClient) extractSolutions(content string) []string {
    return c.extractListItems(content, "solution", "fix")
}

func (c *OpenAIClient) extractPrevention(content string) []string {
    return c.extractListItems(content, "prevent", "avoid")
}

func (c *OpenAIClient) extractRelatedTopics(content string) []string {
    return c.extractListItems(content, "related", "topic")
}

func (c *OpenAIClient) extractExamples(content string) []string {
    return c.extractListItems(content, "example", "usage")
}

func (c *OpenAIClient) extractTests(content string) string {
    return c.extractSection(content, "test", "testing")
}

func (c *OpenAIClient) extractDocumentation(content string) string {
    return c.extractSection(content, "documentation", "comment")
}

func (c *OpenAIClient) extractListItems(content, keyword1, keyword2 string) []string {
    var items []string
    lines := strings.Split(content, "\n")
    
    inSection := false
    for _, line := range lines {
        lowerLine := strings.ToLower(line)
        if strings.Contains(lowerLine, keyword1) || strings.Contains(lowerLine, keyword2) {
            inSection = true
            continue
        }
        
        if inSection {
            line = strings.TrimSpace(line)
            if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•") {
                item := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "-"), "*"), "•")
                items = append(items, strings.TrimSpace(item))
            } else if line == "" && len(items) > 0 {
                break
            }
        }
    }
    
    return items
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

#### `internal/handlers/commands/ai.go` - AI Command Handler
```go
package commands

import (
    "context"
    "fmt"
    "regexp"
    "strings"
    
    "yordamchi-dev-bot/internal/domain"
    "yordamchi-dev-bot/internal/services/ai"
)

type AIHandler struct {
    aiService ai.Service
    logger    Logger
}

func NewAIHandler(aiService ai.Service, logger Logger) *AIHandler {
    return &AIHandler{
        aiService: aiService,
        logger:    logger,
    }
}

func (h *AIHandler) CanHandle(command string) bool {
    cmd := strings.ToLower(strings.TrimSpace(command))
    return strings.HasPrefix(cmd, "/explain") || 
           strings.HasPrefix(cmd, "/generate") ||
           strings.HasPrefix(cmd, "/debug") ||
           strings.HasPrefix(cmd, "/ai") ||
           strings.HasPrefix(cmd, "/ask")
}

func (h *AIHandler) Description() string {
    return "AI-powered code explanation, generation, and debugging assistance"
}

func (h *AIHandler) Usage() string {
    return `/explain <code> - Explain code functionality
/generate <description> - Generate code from description
/debug <error> <code> - Debug and fix errors
/ai <question> - Ask AI about programming`
}

func (h *AIHandler) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
    parts := strings.Fields(cmd.Text)
    if len(parts) < 2 {
        return &domain.Response{
            Text: h.getUsageMessage(),
            ParseMode: "HTML",
        }, nil
    }

    command := strings.ToLower(parts[0])
    
    switch command {
    case "/explain":
        return h.handleExplain(ctx, cmd.Text)
    case "/generate":
        return h.handleGenerate(ctx, cmd.Text)
    case "/debug":
        return h.handleDebug(ctx, cmd.Text)
    case "/ai", "/ask":
        return h.handleAIQuestion(ctx, cmd.Text)
    default:
        return h.handleAIQuestion(ctx, cmd.Text)
    }
}

func (h *AIHandler) handleExplain(ctx context.Context, message string) (*domain.Response, error) {
    code, language := h.extractCodeFromMessage(message)
    if code == "" {
        return &domain.Response{
            Text: "❌ Please provide code to explain.\n\nExample:\n<code>/explain\n```go\nfunc main() {\n  fmt.Println(\"Hello\")\n}\n```</code>",
            ParseMode: "HTML",
        }, nil
    }

    explanation, err := h.aiService.ExplainCode(ctx, code, language)
    if err != nil {
        h.logger.Printf("❌ AI explanation error: %v", err)
        return &domain.Response{
            Text: "❌ Sorry, I couldn't explain the code right now. Please try again later.",
            ParseMode: "HTML",
        }, nil
    }

    return &domain.Response{
        Text:      h.formatExplanation(explanation),
        ParseMode: "HTML",
    }, nil
}

func (h *AIHandler) handleGenerate(ctx context.Context, message string) (*domain.Response, error) {
    // Extract description from message
    parts := strings.SplitN(message, " ", 2)
    if len(parts) < 2 {
        return &domain.Response{
            Text: "❌ Please provide a description of what code to generate.\n\nExample:\n<code>/generate Create a HTTP server in Go that responds with JSON</code>",
            ParseMode: "HTML",
        }, nil
    }

    description := parts[1]
    language := h.detectLanguageFromDescription(description)

    generation, err := h.aiService.GenerateCode(ctx, description, language)
    if err != nil {
        h.logger.Printf("❌ AI generation error: %v", err)
        return &domain.Response{
            Text: "❌ Sorry, I couldn't generate the code right now. Please try again later.",
            ParseMode: "HTML",
        }, nil
    }

    return &domain.Response{
        Text:      h.formatGeneration(generation),
        ParseMode: "HTML",
    }, nil
}

func (h *AIHandler) handleDebug(ctx context.Context, message string) (*domain.Response, error) {
    // Extract error message and code
    parts := strings.SplitN(message, " ", 2)
    if len(parts) < 2 {
        return &domain.Response{
            Text: "❌ Please provide error message and code.\n\nExample:\n<code>/debug syntax error\n```go\nfunc main() {\n  fmt.Println(\"Hello\"\n}\n```</code>",
            ParseMode: "HTML",
        }, nil
    }

    content := parts[1]
    errorMsg, code, language := h.extractErrorAndCode(content)
    
    if errorMsg == "" {
        return &domain.Response{
            Text: "❌ Please provide an error message to debug.",
            ParseMode: "HTML",
        }, nil
    }

    analysis, err := h.aiService.AnalyzeError(ctx, errorMsg, code, language)
    if err != nil {
        h.logger.Printf("❌ AI debug error: %v", err)
        return &domain.Response{
            Text: "❌ Sorry, I couldn't analyze the error right now. Please try again later.",
            ParseMode: "HTML",
        }, nil
    }

    return &domain.Response{
        Text:      h.formatErrorAnalysis(analysis, errorMsg),
        ParseMode: "HTML",
    }, nil
}

func (h *AIHandler) handleAIQuestion(ctx context.Context, message string) (*domain.Response, error) {
    parts := strings.SplitN(message, " ", 2)
    if len(parts) < 2 {
        return &domain.Response{
            Text: "❌ Please ask a programming question.\n\nExample:\n<code>/ai How do I handle errors in Go?</code>",
            ParseMode: "HTML",
        }, nil
    }

    question := parts[1]
    
    // Use a simple AI response for general questions
    response, err := h.aiService.AnswerQuestion(ctx, question)
    if err != nil {
        h.logger.Printf("❌ AI question error: %v", err)
        return &domain.Response{
            Text: "❌ Sorry, I couldn't answer your question right now. Please try again later.",
            ParseMode: "HTML",
        }, nil
    }

    return &domain.Response{
        Text:      h.formatAIResponse(response, question),
        ParseMode: "HTML",
    }, nil
}

func (h *AIHandler) extractCodeFromMessage(message string) (string, string) {
    // Try to extract code from markdown code blocks
    codeBlockRegex := regexp.MustCompile("```(\\w+)?\\n([\\s\\S]*?)```")
    matches := codeBlockRegex.FindStringSubmatch(message)
    
    if len(matches) >= 3 {
        language := matches[1]
        if language == "" {
            language = "text"
        }
        code := strings.TrimSpace(matches[2])
        return code, language
    }

    // Try inline code
    inlineCodeRegex := regexp.MustCompile("`([^`]+)`")
    matches = inlineCodeRegex.FindStringSubmatch(message)
    if len(matches) >= 2 {
        return strings.TrimSpace(matches[1]), "text"
    }

    // Extract from command arguments (fallback)
    parts := strings.SplitN(message, " ", 2)
    if len(parts) >= 2 {
        return parts[1], "text"
    }

    return "", ""
}

func (h *AIHandler) extractErrorAndCode(content string) (string, string, string) {
    // Try to find error message and code separately
    lines := strings.Split(content, "\n")
    var errorMsg string
    var codeLines []string
    var language string
    
    inCodeBlock := false
    for _, line := range lines {
        if strings.HasPrefix(line, "```") {
            if inCodeBlock {
                break
            }
            inCodeBlock = true
            // Extract language from code fence
            lang := strings.TrimPrefix(line, "```")
            if lang != "" {
                language = lang
            }
            continue
        }
        
        if inCodeBlock {
            codeLines = append(codeLines, line)
        } else if errorMsg == "" && strings.TrimSpace(line) != "" {
            errorMsg = strings.TrimSpace(line)
        }
    }
    
    code := strings.Join(codeLines, "\n")
    if language == "" {
        language = "text"
    }
    
    return errorMsg, code, language
}

func (h *AIHandler) detectLanguageFromDescription(description string) string {
    desc := strings.ToLower(description)
    
    if strings.Contains(desc, "go") || strings.Contains(desc, "golang") {
        return "go"
    } else if strings.Contains(desc, "javascript") || strings.Contains(desc, "js") || strings.Contains(desc, "node") {
        return "javascript"
    } else if strings.Contains(desc, "python") || strings.Contains(desc, "py") {
        return "python"
    } else if strings.Contains(desc, "java") && !strings.Contains(desc, "javascript") {
        return "java"
    } else if strings.Contains(desc, "c#") || strings.Contains(desc, "csharp") {
        return "csharp"
    }
    
    return "go" // Default to Go for this bot
}

func (h *AIHandler) formatExplanation(explanation *ai.CodeExplanation) string {
    var text strings.Builder
    
    text.WriteString("🤖 <b>AI Code Explanation</b>\n\n")
    
    // Summary
    if explanation.Summary != "" {
        text.WriteString(fmt.Sprintf("📋 <b>Summary:</b>\n%s\n\n", explanation.Summary))
    }
    
    // Language and complexity
    text.WriteString(fmt.Sprintf("💻 <b>Language:</b> %s\n", explanation.Language))
    text.WriteString(fmt.Sprintf("📊 <b>Complexity:</b> %s\n\n", strings.Title(explanation.Complexity)))
    
    // Key concepts
    if len(explanation.KeyConcepts) > 0 {
        text.WriteString("🎯 <b>Key Concepts:</b>\n")
        for _, concept := range explanation.KeyConcepts {
            text.WriteString(fmt.Sprintf("• %s\n", concept))
        }
        text.WriteString("\n")
    }
    
    // Step by step breakdown
    if len(explanation.StepByStep) > 0 {
        text.WriteString("📝 <b>Step-by-Step Breakdown:</b>\n")
        for i, step := range explanation.StepByStep {
            text.WriteString(fmt.Sprintf("<b>%d.</b> %s\n", i+1, step.Explanation))
            if step.Purpose != "" {
                text.WriteString(fmt.Sprintf("   <i>Purpose: %s</i>\n", step.Purpose))
            }
        }
        text.WriteString("\n")
    }
    
    // Suggestions
    if len(explanation.Suggestions) > 0 {
        text.WriteString("💡 <b>Suggestions for Improvement:</b>\n")
        for _, suggestion := range explanation.Suggestions {
            text.WriteString(fmt.Sprintf("• %s\n", suggestion))
        }
        text.WriteString("\n")
    }
    
    text.WriteString("🚀 <i>Need more help? Ask follow-up questions with /ai</i>")
    
    return text.String()
}

func (h *AIHandler) formatGeneration(generation *ai.CodeGeneration) string {
    var text strings.Builder
    
    text.WriteString("🤖 <b>AI Generated Code</b>\n\n")
    
    if generation.Description != "" {
        text.WriteString(fmt.Sprintf("📋 <b>Description:</b>\n%s\n\n", generation.Description))
    }
    
    text.WriteString(fmt.Sprintf("💻 <b>Language:</b> %s\n\n", generation.Language))
    
    if generation.Code != "" {
        text.WriteString("<b>Generated Code:</b>\n")
        text.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>\n", generation.Code))
    }
    
    if len(generation.Examples) > 0 {
        text.WriteString("\n📖 <b>Usage Examples:</b>\n")
        for i, example := range generation.Examples {
            text.WriteString(fmt.Sprintf("%d. %s\n", i+1, example))
        }
    }
    
    if generation.Tests != "" {
        text.WriteString(fmt.Sprintf("\n🧪 <b>Tests:</b>\n<pre><code>%s</code></pre>", generation.Tests))
    }
    
    text.WriteString("\n\n💡 <i>Use /explain to understand the generated code better</i>")
    
    return text.String()
}

func (h *AIHandler) formatErrorAnalysis(analysis *ai.ErrorAnalysis, originalError string) string {
    var text strings.Builder
    
    text.WriteString("🤖 <b>AI Error Analysis</b>\n\n")
    
    text.WriteString(fmt.Sprintf("❌ <b>Error:</b> <code>%s</code>\n\n", originalError))
    
    if analysis.ErrorType != "" {
        text.WriteString(fmt.Sprintf("🏷 <b>Error Type:</b> %s\n\n", analysis.ErrorType))
    }
    
    if analysis.PossibleCause != "" {
        text.WriteString(fmt.Sprintf("🔍 <b>Possible Cause:</b>\n%s\n\n", analysis.PossibleCause))
    }
    
    if len(analysis.Solutions) > 0 {
        text.WriteString("🔧 <b>Solutions:</b>\n")
        for i, solution := range analysis.Solutions {
            text.WriteString(fmt.Sprintf("%d. %s\n", i+1, solution))
        }
        text.WriteString("\n")
    }
    
    if len(analysis.Prevention) > 0 {
        text.WriteString("🛡 <b>Prevention Tips:</b>\n")
        for _, tip := range analysis.Prevention {
            text.WriteString(fmt.Sprintf("• %s\n", tip))
        }
        text.WriteString("\n")
    }
    
    if len(analysis.RelatedTopics) > 0 {
        text.WriteString("📚 <b>Related Topics to Learn:</b>\n")
        for _, topic := range analysis.RelatedTopics {
            text.WriteString(fmt.Sprintf("• %s\n", topic))
        }
    }
    
    text.WriteString("\n💡 <i>Need more specific help? Use /generate to create working code</i>")
    
    return text.String()
}

func (h *AIHandler) formatAIResponse(response, question string) string {
    var text strings.Builder
    
    text.WriteString("🤖 <b>AI Assistant</b>\n\n")
    text.WriteString(fmt.Sprintf("❓ <b>Question:</b> %s\n\n", question))
    text.WriteString(fmt.Sprintf("💬 <b>Answer:</b>\n%s\n\n", response))
    text.WriteString("🚀 <i>Ask follow-up questions or use /explain, /generate, /debug for specific help</i>")
    
    return text.String()
}

func (h *AIHandler) getUsageMessage() string {
    return `🤖 <b>AI Assistant Commands</b>

<b>Code Explanation:</b>
<code>/explain</code> - Explain code functionality
<pre>```go
func main() {
    fmt.Println("Hello")
}
```</pre>

<b>Code Generation:</b>
<code>/generate Create a HTTP server in Go</code>

<b>Error Debugging:</b>
<code>/debug syntax error</code>
<pre>```go
func main() {
    fmt.Println("Hello"
}
```</pre>

<b>General Questions:</b>
<code>/ai How do I handle errors in Go?</code>
<code>/ask What are goroutines?</code>

<b>Tips:</b>
• Use code blocks for better AI understanding
• Be specific in your questions
• Follow up with related questions

🚀 <i>Powered by GPT-4 for accurate programming assistance</i>`
}
```

This covers Day 15 of Week 3 with comprehensive OpenAI integration. Would you like me to continue with the remaining days (16-21) to complete Week 3, and then proceed to Week 4?

