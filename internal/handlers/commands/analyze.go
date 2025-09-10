package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// AnalyzeCommand handles AI-powered task analysis
type AnalyzeCommand struct {
	taskAnalyzer       *services.TaskAnalyzer
	logger             domain.Logger
	fileExtractor      *services.FileExtractor
	telegramFileService *services.TelegramFileService
}

// NewAnalyzeCommand creates a new analyze command handler
func NewAnalyzeCommand(taskAnalyzer *services.TaskAnalyzer, logger domain.Logger, fileExtractor *services.FileExtractor, telegramFileService *services.TelegramFileService) *AnalyzeCommand {
	return &AnalyzeCommand{
		taskAnalyzer:       taskAnalyzer,
		logger:             logger,
		fileExtractor:      fileExtractor,
		telegramFileService: telegramFileService,
	}
}

// Handle processes the analyze command for both text and file analysis
func (c *AnalyzeCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing analyze command", "user_id", cmd.User.TelegramID)

	// Check if message contains a file attachment
	if cmd.Document != nil {
		return c.handleFileAnalysis(ctx, cmd)
	}

	// Handle text-based analysis
	return c.handleTextAnalysis(ctx, cmd)
}

// handleFileAnalysis processes uploaded files for analysis
func (c *AnalyzeCommand) handleFileAnalysis(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing file analysis", 
		"user_id", cmd.User.TelegramID, 
		"filename", cmd.Document.FileName,
		"file_size", cmd.Document.FileSize)

	// 1. Validate file
	if err := c.fileExtractor.ValidateFile(cmd.Document); err != nil {
		return &domain.Response{
			Text: fmt.Sprintf("❌ **File validation failed:** %s\n\n"+
				"**Supported formats:** %s\n"+
				"**Maximum size:** 20MB",
				err.Error(), 
				strings.Join(c.fileExtractor.GetSupportedFormats(), ", ")),
			ParseMode: "Markdown",
		}, nil
	}

	// 2. Download file temporarily
	tempFile, err := c.telegramFileService.DownloadFile(cmd.Document)
	if err != nil {
		c.logger.Error("Failed to download file", "error", err)
		return &domain.Response{
			Text: "❌ **Download failed.** Please try uploading the file again.\n\n" +
				"If the problem persists, try:\n" +
				"• Reducing file size\n" +
				"• Converting to a simpler format (TXT, MD)\n" +
				"• Checking your internet connection",
			ParseMode: "Markdown",
		}, nil
	}

	// 3. Ensure cleanup
	defer func() {
		c.telegramFileService.CleanupFile(tempFile)
	}()

	// 4. Extract content from file
	content, err := c.fileExtractor.ExtractContent(tempFile, cmd.Document.FileName)
	if err != nil {
		c.logger.Error("Failed to extract file content", "error", err, "filename", cmd.Document.FileName)
		return &domain.Response{
			Text: fmt.Sprintf("❌ **Content extraction failed:** %s\n\n"+
				"**Troubleshooting:**\n"+
				"• Ensure file is not corrupted\n"+
				"• Try saving in a different format\n"+
				"• For PDFs, ensure text is selectable (not scanned image)",
				err.Error()),
			ParseMode: "Markdown",
		}, nil
	}

	// 5. Check if content was extracted
	if strings.TrimSpace(content) == "" {
		return &domain.Response{
			Text: fmt.Sprintf("❌ **No readable content found** in `%s`\n\n"+
				"**Possible causes:**\n"+
				"• File contains only images/graphics\n"+
				"• File is corrupted or password-protected\n"+
				"• Text is embedded in images (OCR not supported yet)\n\n"+
				"**Suggestion:** Try uploading a plain text file with your requirements.",
				cmd.Document.FileName),
			ParseMode: "Markdown",
		}, nil
	}

	// 6. Analyze extracted content
	req := domain.TaskBreakdownRequest{
		Requirement: content,
		TeamSkills:  []string{"go", "react", "python", "docker", "postgresql", "javascript", "typescript", "kubernetes"},
		ProjectType: "web",
	}

	result, err := c.taskAnalyzer.AnalyzeRequirement(req)
	if err != nil {
		c.logger.Error("File content analysis failed", "error", err, "filename", cmd.Document.FileName)
		return &domain.Response{
			Text: "❌ **Analysis failed.** The file content might be too complex or unclear.\n\n" +
				"**Try:**\n" +
				"• Simplifying the requirements document\n" +
				"• Using more specific technical language\n" +
				"• Breaking down into smaller sections",
			ParseMode: "Markdown",
		}, nil
	}

	// 7. Format results with file context
	responseText := c.formatFileAnalysisResults(result, cmd.Document)

	c.logger.Info("File analysis completed",
		"user_id", cmd.User.TelegramID,
		"filename", cmd.Document.FileName,
		"content_length", len(content),
		"tasks_count", len(result.Tasks),
		"total_estimate", result.TotalEstimate,
		"confidence", result.Confidence)

	return &domain.Response{
		Text:      responseText,
		ParseMode: "Markdown",
	}, nil
}

// handleTextAnalysis handles traditional text-based analysis
func (c *AnalyzeCommand) handleTextAnalysis(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Parse command arguments (everything after /analyze)
	parts := strings.Fields(cmd.Text)
	if len(parts) < 2 {
		return &domain.Response{
			Text: "📋 **AI Requirements Analysis**\n\n" +
				"**Text Analysis:**\n" +
				"`/analyze Build user authentication with OAuth`\n\n" +
				"**File Analysis:**\n" +
				"Upload any document (PDF, DOCX, TXT, MD, XLSX) with your requirements\n\n" +
				"**Supported formats:** " + strings.Join(c.fileExtractor.GetSupportedFormats(), ", ") + "\n" +
				"**Maximum size:** 20MB\n\n" +
				"**Tips for better analysis:**\n" +
				"• Be specific about technologies (React, Go, PostgreSQL)\n" +
				"• Include project scope (backend, frontend, full-stack)\n" +
				"• Mention integrations (GitHub, Stripe, etc.)\n" +
				"• Describe user stories and acceptance criteria",
			ParseMode: "Markdown",
		}, nil
	}

	requirement := strings.Join(parts[1:], " ")

	// Create analysis request with default team skills
	req := domain.TaskBreakdownRequest{
		Requirement: requirement,
		TeamSkills:  []string{"go", "react", "python", "docker", "postgresql", "javascript", "typescript", "kubernetes"},
		ProjectType: "web", // Default to web project
	}

	// Analyze with TaskAnalyzer
	result, err := c.taskAnalyzer.AnalyzeRequirement(req)
	if err != nil {
		c.logger.Error("Task analysis failed", "error", err, "requirement", requirement)
		return &domain.Response{
			Text: "❌ **Analysis failed.** Please try again with a clearer requirement.\n\n" +
				"**Make sure to:**\n" +
				"• Use clear, specific language\n" +
				"• Include technology stack details\n" +
				"• Specify project scope and goals\n" +
				"• Provide concrete user stories",
			ParseMode: "Markdown",
		}, nil
	}

	// Format and send results
	responseText := c.formatTaskBreakdown(result)

	c.logger.Info("Text analysis completed",
		"user_id", cmd.User.TelegramID,
		"tasks_count", len(result.Tasks),
		"total_estimate", result.TotalEstimate,
		"confidence", result.Confidence)

	return &domain.Response{
		Text:      responseText,
		ParseMode: "Markdown",
	}, nil
}

// formatTaskBreakdown formats the analysis results for display
func (c *AnalyzeCommand) formatTaskBreakdown(result *domain.TaskBreakdownResponse) string {
	var response strings.Builder
	
	response.WriteString("📋 **Task Breakdown Analysis**\n\n")
	
	// Group tasks by category
	categories := make(map[string][]domain.Task)
	for _, task := range result.Tasks {
		categories[task.Category] = append(categories[task.Category], task)
	}
	
	// Format each category with appropriate icons
	categoryIcons := map[string]string{
		"backend":  "🔐",
		"frontend": "🎨", 
		"qa":       "🧪",
		"devops":   "⚙️",
	}
	
	categoryNames := map[string]string{
		"backend":  "Backend Development",
		"frontend": "Frontend Development", 
		"qa":       "Quality Assurance",
		"devops":   "DevOps & Infrastructure",
	}
	
	for category, tasks := range categories {
		icon := categoryIcons[category]
		if icon == "" {
			icon = "📝"
		}
		
		categoryName := categoryNames[category]
		if categoryName == "" {
			categoryName = strings.Title(category) + " Tasks"
		}
		
		categoryTotal := 0.0
		response.WriteString(fmt.Sprintf("%s **%s** (Est: %.1fh)\n", icon, categoryName, getCategoryTotal(tasks)))
		
		for _, task := range tasks {
			priorityIcon := getPriorityIcon(task.Priority)
			response.WriteString(fmt.Sprintf("├── %s %s - %.1fh\n", priorityIcon, task.Title, task.EstimateHours))
			categoryTotal += task.EstimateHours
		}
		
		response.WriteString(fmt.Sprintf("└── **Subtotal: %.1f hours**\n\n", categoryTotal))
	}
	
	// Total estimate with developer days calculation
	devDays := result.TotalEstimate / 8
	response.WriteString(fmt.Sprintf("⏱️ **Total Estimate: %.1f hours (%.1f developer days)**\n\n", 
		result.TotalEstimate, devDays))
	
	// Recommended team
	if len(result.RecommendedTeam) > 0 {
		response.WriteString("👥 **Recommended Team:**\n")
		for _, member := range result.RecommendedTeam {
			response.WriteString(fmt.Sprintf("• %s\n", member))
		}
		response.WriteString("\n")
	}
	
	// Critical path tasks
	if len(result.CriticalPath) > 0 {
		response.WriteString("🎯 **Critical Path Tasks:** ")
		response.WriteString(fmt.Sprintf("%d high-priority items\n\n", len(result.CriticalPath)))
	}
	
	// Risk factors
	if len(result.RiskFactors) > 0 {
		response.WriteString("⚠️ **Risk Factors & Considerations:**\n")
		for _, risk := range result.RiskFactors {
			response.WriteString(fmt.Sprintf("• %s\n", risk))
		}
		response.WriteString("\n")
	}
	
	// Analysis confidence and next steps
	confidenceEmoji := getConfidenceEmoji(result.Confidence)
	response.WriteString(fmt.Sprintf("📊 **Analysis Confidence:** %s %.0f%%\n\n", confidenceEmoji, result.Confidence*100))
	
	response.WriteString("**Next Steps:**\n")
	response.WriteString("• Use `/create_project <name>` to create a project\n")
	response.WriteString("• Use `/add_member @user skills` to build your team\n")
	response.WriteString("• Use `/workload` to check team capacity")
	
	return response.String()
}

// formatFileAnalysisResults formats analysis results with file context
func (c *AnalyzeCommand) formatFileAnalysisResults(result *domain.TaskBreakdownResponse, document *domain.TelegramDocument) string {
	var response strings.Builder
	
	// File header with metadata
	response.WriteString("📄 **File Analysis Complete**\n\n")
	response.WriteString(fmt.Sprintf("**File:** `%s`\n", document.FileName))
	response.WriteString(fmt.Sprintf("**Size:** %s\n", c.telegramFileService.GetFileSize(document.FileSize)))
	response.WriteString(fmt.Sprintf("**Type:** %s\n\n", document.MimeType))
	
	// Analysis summary
	response.WriteString("🤖 **AI Analysis Summary:**\n")
	response.WriteString(fmt.Sprintf("├── **Tasks Generated:** %d\n", len(result.Tasks)))
	response.WriteString(fmt.Sprintf("├── **Total Estimate:** %.1f hours (%.1f days)\n", result.TotalEstimate, result.TotalEstimate/8))
	confidence := getConfidenceEmoji(result.Confidence)
	response.WriteString(fmt.Sprintf("└── **Confidence:** %s %.0f%%\n\n", confidence, result.Confidence*100))

	// Task breakdown by category
	response.WriteString("📋 **Task Breakdown:**\n\n")
	
	// Group tasks by category
	categories := make(map[string][]domain.Task)
	for _, task := range result.Tasks {
		categories[task.Category] = append(categories[task.Category], task)
	}
	
	categoryIcons := map[string]string{
		"backend":  "🔐",
		"frontend": "🎨", 
		"qa":       "🧪",
		"devops":   "⚙️",
	}
	
	for category, tasks := range categories {
		icon := categoryIcons[category]
		if icon == "" {
			icon = "📝"
		}
		
		categoryName := strings.Title(category)
		categoryTotal := getCategoryTotal(tasks)
		response.WriteString(fmt.Sprintf("%s **%s** (%.1fh)\n", icon, categoryName, categoryTotal))
		
		// Show up to 3 tasks per category to keep response manageable
		maxTasks := 3
		for i, task := range tasks {
			if i >= maxTasks {
				response.WriteString(fmt.Sprintf("├── ... and %d more %s tasks\n", len(tasks)-maxTasks, category))
				break
			}
			
			priority := getPriorityIcon(task.Priority)
			response.WriteString(fmt.Sprintf("├── %s %s (%.1fh)\n", priority, task.Title, task.EstimateHours))
		}
		response.WriteString("\n")
	}
	
	// Project insights
	if len(result.RecommendedTeam) > 0 {
		response.WriteString("👥 **Recommended Team Skills:**\n")
		for _, skill := range result.RecommendedTeam[:min(len(result.RecommendedTeam), 5)] {
			response.WriteString(fmt.Sprintf("• %s\n", skill))
		}
		response.WriteString("\n")
	}
	
	// Risk factors (if any)
	if len(result.RiskFactors) > 0 {
		response.WriteString("⚠️ **Key Risks:**\n")
		for _, risk := range result.RiskFactors[:min(len(result.RiskFactors), 3)] {
			response.WriteString(fmt.Sprintf("• %s\n", risk))
		}
		response.WriteString("\n")
	}
	
	// Next steps
	response.WriteString("🚀 **Next Steps:**\n")
	response.WriteString("• Use `/create_project <name>` to create a project\n")
	response.WriteString("• Use `/add_member @user skills` to build your team\n")
	response.WriteString("• Use `/workload` to analyze team capacity\n")
	response.WriteString("• Use `/list_projects` to track progress")
	
	return response.String()
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper functions for formatting
func getCategoryTotal(tasks []domain.Task) float64 {
	total := 0.0
	for _, task := range tasks {
		total += task.EstimateHours
	}
	return total
}

func getPriorityIcon(priority int) string {
	switch priority {
	case 1:
		return "🔴" // High priority
	case 2:
		return "🟡" // Medium priority
	case 3:
		return "🟢" // Low priority
	default:
		return "⚪" // Unknown priority
	}
}

func getConfidenceEmoji(confidence float64) string {
	if confidence >= 0.9 {
		return "🎯" // Very high confidence
	} else if confidence >= 0.75 {
		return "✅" // High confidence
	} else if confidence >= 0.6 {
		return "⚠️"  // Medium confidence
	} else {
		return "❓" // Low confidence
	}
}

// CanHandle checks if this handler can process the command
func (c *AnalyzeCommand) CanHandle(command string) bool {
	return command == "/analyze"
}

// Description returns the command description
func (c *AnalyzeCommand) Description() string {
	return "🔍 Break down development requirements into actionable tasks with AI analysis"
}

// Usage returns the command usage instructions
func (c *AnalyzeCommand) Usage() string {
	return "/analyze <requirement> - Analyze development requirements and break them down into tasks"
}