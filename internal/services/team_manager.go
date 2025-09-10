package services

import (
	"sort"
	"strings"

	"yordamchi-dev-bot/internal/domain"
)

type TeamManager struct {
	// Database repository will be injected
}

func NewTeamManager() *TeamManager {
	return &TeamManager{}
}

// AnalyzeWorkload calculates current team workload and availability
func (tm *TeamManager) AnalyzeWorkload(teamID string, members []domain.TeamMember, tasks []domain.Task) *domain.TeamWorkload {
	memberWorkloads := make([]domain.MemberWorkload, 0, len(members))
	totalAvailable := 0.0
	totalAllocated := 0.0

	for _, member := range members {
		workload := tm.calculateMemberWorkload(member, tasks)
		memberWorkloads = append(memberWorkloads, workload)
		
		totalAvailable += member.Capacity
		totalAllocated += workload.Current
	}

	utilization := 0.0
	if totalAvailable > 0 {
		utilization = totalAllocated / totalAvailable
	}

	return &domain.TeamWorkload{
		TeamID:      teamID,
		Members:     memberWorkloads,
		Available:   totalAvailable,
		Allocated:   totalAllocated,
		Utilization: utilization,
	}
}

// RecommendAssignment suggests optimal task assignments based on skills and workload
func (tm *TeamManager) RecommendAssignment(task domain.Task, members []domain.TeamMember, currentTasks []domain.Task) *domain.TeamMember {
	// Filter members by skill match
	candidates := tm.filterBySkills(task, members)
	if len(candidates) == 0 {
		candidates = members // fallback to all members
	}

	// Score each candidate
	type candidate struct {
		member *domain.TeamMember
		score  float64
	}

	scored := make([]candidate, 0, len(candidates))
	
	for i := range candidates {
		member := &candidates[i]
		score := tm.calculateAssignmentScore(task, member, currentTasks)
		scored = append(scored, candidate{member: member, score: score})
	}

	// Sort by score (higher is better)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if len(scored) > 0 {
		return scored[0].member
	}

	return nil
}

// OptimizeWorkload redistributes tasks to balance team workload
func (tm *TeamManager) OptimizeWorkload(teamID string, members []domain.TeamMember, tasks []domain.Task) []domain.Task {
	optimized := make([]domain.Task, len(tasks))
	copy(optimized, tasks)

	// Calculate current workloads
	memberWorkloads := make(map[string]float64)
	for _, member := range members {
		memberWorkloads[member.ID] = tm.calculateCurrentWorkload(member.ID, tasks)
	}

	// Find overloaded and underloaded members
	overloaded := []string{}
	underloaded := []string{}
	
	for _, member := range members {
		utilization := memberWorkloads[member.ID] / member.Capacity
		
		if utilization > 0.9 {
			overloaded = append(overloaded, member.ID)
		} else if utilization < 0.6 {
			underloaded = append(underloaded, member.ID)
		}
	}

	// Redistribute tasks from overloaded to underloaded members
	for _, overloadedID := range overloaded {
		for i, task := range optimized {
			if task.AssignedTo == overloadedID && task.Status == "todo" {
				// Find best alternative assignment
				for _, underloadedID := range underloaded {
					if tm.canAssignTask(task, underloadedID, members) {
						optimized[i].AssignedTo = underloadedID
						
						// Update workload tracking
						memberWorkloads[overloadedID] -= task.EstimateHours
						memberWorkloads[underloadedID] += task.EstimateHours
						
						break
					}
				}
			}
		}
	}

	return optimized
}

func (tm *TeamManager) calculateMemberWorkload(member domain.TeamMember, tasks []domain.Task) domain.MemberWorkload {
	current := tm.calculateCurrentWorkload(member.ID, tasks)
	utilization := 0.0
	
	if member.Capacity > 0 {
		utilization = current / member.Capacity
	}

	status := "available"
	if utilization > 0.9 {
		status = "overloaded"
	} else if utilization > 0.75 {
		status = "busy"
	}

	return domain.MemberWorkload{
		MemberID:    member.ID,
		Username:    member.Username,
		Capacity:    member.Capacity,
		Current:     current,
		Utilization: utilization,
		Status:      status,
	}
}

func (tm *TeamManager) calculateCurrentWorkload(memberID string, tasks []domain.Task) float64 {
	workload := 0.0
	
	for _, task := range tasks {
		if task.AssignedTo == memberID && (task.Status == "todo" || task.Status == "in_progress") {
			workload += task.EstimateHours
		}
	}
	
	return workload
}

func (tm *TeamManager) filterBySkills(task domain.Task, members []domain.TeamMember) []domain.TeamMember {
	filtered := []domain.TeamMember{}
	requiredSkills := tm.extractRequiredSkills(task)
	
	for _, member := range members {
		if tm.hasMatchingSkills(member.Skills, requiredSkills) {
			filtered = append(filtered, member)
		}
	}
	
	return filtered
}

func (tm *TeamManager) extractRequiredSkills(task domain.Task) []string {
	skills := []string{}
	desc := strings.ToLower(task.Description + " " + task.Title)
	
	// Map task categories to skills
	categorySkills := map[string][]string{
		"backend":  {"go", "backend", "api", "database"},
		"frontend": {"react", "frontend", "ui", "javascript"},
		"qa":       {"testing", "qa", "automation"},
		"devops":   {"devops", "docker", "kubernetes", "ci/cd"},
	}
	
	if taskSkills, exists := categorySkills[task.Category]; exists {
		skills = append(skills, taskSkills...)
	}
	
	// Extract specific technology mentions
	technologies := []string{"go", "react", "python", "javascript", "docker", "kubernetes", "postgres", "mongodb"}
	for _, tech := range technologies {
		if strings.Contains(desc, tech) {
			skills = append(skills, tech)
		}
	}
	
	return skills
}

func (tm *TeamManager) hasMatchingSkills(memberSkills, requiredSkills []string) bool {
	memberSkillMap := make(map[string]bool)
	for _, skill := range memberSkills {
		memberSkillMap[strings.ToLower(skill)] = true
	}
	
	for _, required := range requiredSkills {
		if memberSkillMap[strings.ToLower(required)] {
			return true
		}
	}
	
	return false
}

func (tm *TeamManager) calculateAssignmentScore(task domain.Task, member *domain.TeamMember, currentTasks []domain.Task) float64 {
	score := 0.0
	
	// Skill match bonus (0-3 points)
	requiredSkills := tm.extractRequiredSkills(task)
	skillMatches := 0
	for _, required := range requiredSkills {
		for _, memberSkill := range member.Skills {
			if strings.EqualFold(required, memberSkill) {
				skillMatches++
				break
			}
		}
	}
	score += float64(skillMatches)
	
	// Workload penalty (subtract utilization percentage)
	currentWorkload := tm.calculateCurrentWorkload(member.ID, currentTasks)
	utilization := currentWorkload / member.Capacity
	score -= utilization * 2 // penalty for high utilization
	
	// Role match bonus
	if tm.roleMatchesTask(member.Role, task.Category) {
		score += 1.0
	}
	
	return score
}

func (tm *TeamManager) roleMatchesTask(role, category string) bool {
	roleMatches := map[string][]string{
		"lead":   {"backend", "frontend", "qa", "devops"},
		"senior": {"backend", "frontend", "devops"},
		"mid":    {"backend", "frontend", "qa"},
		"junior": {"qa", "frontend"},
	}
	
	if categories, exists := roleMatches[strings.ToLower(role)]; exists {
		for _, cat := range categories {
			if cat == category {
				return true
			}
		}
	}
	
	return false
}

func (tm *TeamManager) canAssignTask(task domain.Task, memberID string, members []domain.TeamMember) bool {
	for _, member := range members {
		if member.ID == memberID {
			requiredSkills := tm.extractRequiredSkills(task)
			return tm.hasMatchingSkills(member.Skills, requiredSkills)
		}
	}
	return false
}