package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

type TaskAnalyzer struct {
	claudeService *ClaudeService
	geminiService *GeminiService
	logger        domain.Logger
}

func NewTaskAnalyzer(logger domain.Logger) *TaskAnalyzer {
	return &TaskAnalyzer{
		claudeService: NewClaudeService(logger),
		geminiService: NewGeminiService(logger),
		logger:        logger,
	}
}

// AnalyzeRequirement breaks down a development requirement into tasks
func (ta *TaskAnalyzer) AnalyzeRequirement(req domain.TaskBreakdownRequest) (*domain.TaskBreakdownResponse, error) {
	ctx := context.Background()
	
	// Try AI services in order of preference, fall back to rule-based analysis
	
	// 1. Try Claude first (most accurate for code analysis)
	if ta.claudeService.IsConfigured() {
		ta.logger.Info("Using Claude AI for task analysis")
		result, err := ta.claudeService.AnalyzeRequirement(ctx, req)
		if err == nil {
			return result, nil
		}
		ta.logger.Error("Claude analysis failed, trying Gemini", "error", err)
	}
	
	// 2. Try Gemini as fallback
	if ta.geminiService.IsConfigured() {
		ta.logger.Info("Using Gemini AI for task analysis")
		result, err := ta.geminiService.AnalyzeRequirement(ctx, req)
		if err == nil {
			return result, nil
		}
		ta.logger.Error("Gemini analysis failed, using rule-based fallback", "error", err)
	}
	
	// 3. Fall back to rule-based analysis
	ta.logger.Info("Using rule-based task analysis (no AI configured)")
	return ta.ruleBasedAnalysis(req)
}

// ruleBasedAnalysis provides fallback analysis when AI services are unavailable
func (ta *TaskAnalyzer) ruleBasedAnalysis(req domain.TaskBreakdownRequest) (*domain.TaskBreakdownResponse, error) {
	tasks := ta.generateTasks(req.Requirement, req.ProjectType)
	
	// Calculate estimates based on task complexity
	totalEstimate := 0.0
	for i := range tasks {
		tasks[i].EstimateHours = ta.estimateTaskTime(tasks[i])
		totalEstimate += tasks[i].EstimateHours
	}

	// Recommend team members based on skills
	recommendedTeam := ta.recommendTeam(tasks, req.TeamSkills)

	return &domain.TaskBreakdownResponse{
		Tasks:           tasks,
		TotalEstimate:   totalEstimate,
		RecommendedTeam: recommendedTeam,
		CriticalPath:    ta.identifyCriticalPath(tasks),
		RiskFactors:     ta.identifyRiskFactors(req.Requirement),
		Confidence:      0.75, // Rule-based confidence level
	}, nil
}

func (ta *TaskAnalyzer) generateTasks(requirement, projectType string) []domain.Task {
	req := strings.ToLower(requirement)
	tasks := []domain.Task{}
	
	// Authentication system detection
	if strings.Contains(req, "auth") || strings.Contains(req, "login") || strings.Contains(req, "oauth") {
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "OAuth Provider Integration",
			Description: "Integrate OAuth providers (GitHub, Google)",
			Category:    "backend",
			Priority:    1,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "JWT Token Management",
			Description: "Implement JWT token generation and validation",
			Category:    "backend",
			Priority:    1,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Login/Signup UI",
			Description: "Create authentication user interface",
			Category:    "frontend",
			Priority:    2,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Authentication Testing",
			Description: "Test authentication flows and security",
			Category:    "qa",
			Priority:    3,
		})
	}

	// API development detection
	if strings.Contains(req, "api") || strings.Contains(req, "endpoint") {
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "API Design & Documentation",
			Description: "Design REST API endpoints and OpenAPI spec",
			Category:    "backend",
			Priority:    1,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "API Implementation",
			Description: "Implement REST API endpoints",
			Category:    "backend",
			Priority:    2,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "API Testing",
			Description: "Unit and integration tests for API",
			Category:    "qa",
			Priority:    3,
		})
	}

	// Database work detection
	if strings.Contains(req, "database") || strings.Contains(req, "data") || strings.Contains(req, "store") {
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Database Schema Design",
			Description: "Design database schema and relationships",
			Category:    "backend",
			Priority:    1,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Database Migration",
			Description: "Create database migration scripts",
			Category:    "devops",
			Priority:    2,
		})
	}

	// Frontend work detection
	if strings.Contains(req, "ui") || strings.Contains(req, "interface") || strings.Contains(req, "frontend") {
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "UI/UX Design",
			Description: "Create mockups and user interface design",
			Category:    "frontend",
			Priority:    1,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Frontend Implementation",
			Description: "Implement user interface components",
			Category:    "frontend",
			Priority:    2,
		})
	}

	// If no specific patterns matched, create generic tasks
	if len(tasks) == 0 {
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Requirement Analysis",
			Description: "Analyze and break down requirements",
			Category:    "backend",
			Priority:    1,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Implementation",
			Description: requirement,
			Category:    "backend",
			Priority:    2,
		})
		
		tasks = append(tasks, domain.Task{
			ID:          generateID(),
			Title:       "Testing & Validation",
			Description: "Test implementation and validate requirements",
			Category:    "qa",
			Priority:    3,
		})
	}

	return tasks
}

func (ta *TaskAnalyzer) estimateTaskTime(task domain.Task) float64 {
	baseHours := map[string]float64{
		"backend":  4.0,
		"frontend": 3.0,
		"qa":       2.0,
		"devops":   3.0,
	}

	base := baseHours[task.Category]
	if base == 0 {
		base = 3.0 // default
	}

	// Adjust based on priority and complexity keywords
	multiplier := 1.0
	desc := strings.ToLower(task.Description)
	
	if strings.Contains(desc, "oauth") || strings.Contains(desc, "security") {
		multiplier = 1.5
	}
	if strings.Contains(desc, "integration") || strings.Contains(desc, "api") {
		multiplier = 1.3
	}
	if strings.Contains(desc, "testing") || strings.Contains(desc, "validation") {
		multiplier = 0.8
	}
	
	return base * multiplier
}

func (ta *TaskAnalyzer) recommendTeam(tasks []domain.Task, teamSkills []string) []string {
	skillMap := make(map[string]bool)
	for _, skill := range teamSkills {
		skillMap[strings.ToLower(skill)] = true
	}

	recommendations := []string{}
	
	for _, task := range tasks {
		switch task.Category {
		case "backend":
			if skillMap["go"] || skillMap["backend"] || skillMap["api"] {
				recommendations = append(recommendations, "@backend-dev")
			}
		case "frontend":
			if skillMap["react"] || skillMap["frontend"] || skillMap["ui"] {
				recommendations = append(recommendations, "@frontend-dev")
			}
		case "qa":
			if skillMap["testing"] || skillMap["qa"] {
				recommendations = append(recommendations, "@qa-engineer")
			}
		case "devops":
			if skillMap["devops"] || skillMap["docker"] || skillMap["kubernetes"] {
				recommendations = append(recommendations, "@devops-engineer")
			}
		}
	}

	return removeDuplicates(recommendations)
}

func (ta *TaskAnalyzer) identifyCriticalPath(tasks []domain.Task) []string {
	// Simple critical path: tasks with priority 1 and high estimates
	critical := []string{}
	
	for _, task := range tasks {
		if task.Priority == 1 {
			critical = append(critical, task.ID)
		}
	}
	
	return critical
}

func (ta *TaskAnalyzer) identifyRiskFactors(requirement string) []string {
	risks := []string{}
	req := strings.ToLower(requirement)
	
	if strings.Contains(req, "oauth") || strings.Contains(req, "auth") {
		risks = append(risks, "Authentication security complexity")
	}
	if strings.Contains(req, "integration") {
		risks = append(risks, "Third-party API dependency")
	}
	if strings.Contains(req, "database") && strings.Contains(req, "migration") {
		risks = append(risks, "Data migration complexity")
	}
	if len(strings.Fields(requirement)) > 20 {
		risks = append(risks, "Large scope - consider breaking down further")
	}
	
	return risks
}

func generateID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}