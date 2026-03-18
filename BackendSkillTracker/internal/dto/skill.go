package dto

import "time"

type SkillRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type SkillResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RecommendedEmployeeResponse struct {
	ID            int             `json:"id"`
	Username      string          `json:"username"`
	Name          string          `json:"name"`
	Skills        []SkillResponse `json:"skills"`
	MatchScore    int             `json:"match_score"`
	MatchedSkills []string        `json:"matched_skills"`
	MissingSkills []string        `json:"missing_skills"`
}
