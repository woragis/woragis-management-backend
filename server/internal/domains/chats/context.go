package chats

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	experiencesdomain "woragis-management-service/internal/domains/experiences"
	projectsdomain "woragis-management-service/internal/domains/projects"
	userprofilesdomain "woragis-management-service/internal/domains/userprofiles"
	// Note: The following domains are implemented in the jobs service
	// jobapplicationsdomain, jobwebsitesdomain, resumesdomain
	// Note: The following domains are not yet implemented in this service
	// skillsdomain, casestudiesdomain, technicalwritingsdomain, postsdomain, problemsolutionsdomain
)

// ContextBuilder builds context for chat conversations.
type ContextBuilder struct {
	jobApplicationService   interface{} // Placeholder for job applications service (from jobs service)
	resumeService           interface{} // Placeholder for resumes service (from jobs service)
	userProfileService      userprofilesdomain.Service
	projectService          projectsdomain.Service
	skillService            interface{} // Placeholder for skills service
	caseStudyService        interface{} // Placeholder for case studies service
	technicalWritingService interface{} // Placeholder for technical writings service
	postService             interface{} // Placeholder for posts service
	problemSolutionService  interface{} // Placeholder for problem solutions service
	experienceService       experiencesdomain.Service
	logger                  *slog.Logger
}

// NewContextBuilder creates a new context builder.
func NewContextBuilder(
	jobApplicationService interface{}, // Placeholder for job applications service (from jobs service)
	resumeService interface{}, // Placeholder for resumes service (from jobs service)
	userProfileService userprofilesdomain.Service,
	projectService projectsdomain.Service,
	skillService interface{}, // Placeholder for skills service
	caseStudyService interface{}, // Placeholder for case studies service
	technicalWritingService interface{}, // Placeholder for technical writings service
	postService interface{}, // Placeholder for posts service
	problemSolutionService interface{}, // Placeholder for problem solutions service
	experienceService experiencesdomain.Service,
	logger *slog.Logger,
) *ContextBuilder {
	return &ContextBuilder{
		jobApplicationService:   jobApplicationService,
		resumeService:           resumeService,
		userProfileService:      userProfileService,
		projectService:          projectService,
		skillService:            skillService,
		caseStudyService:        caseStudyService,
		technicalWritingService: technicalWritingService,
		postService:             postService,
		problemSolutionService:  problemSolutionService,
		experienceService:       experienceService,
		logger:                  logger,
	}
}

// BuildContextOptions specifies what context to include.
type BuildContextOptions struct {
	IncludeJobApplication bool
	IncludeResume         bool
	IncludeUserProfile    bool
	IncludeProjects       bool
	IncludeCaseStudies    bool
	IncludeTechnicalWritings bool
	IncludePosts          bool
	IncludeProblemSolutions bool
	IncludeSkills         bool
	IncludeExperiences    bool
}

// BuildContext builds context string for a conversation.
func (cb *ContextBuilder) BuildContext(ctx context.Context, userID uuid.UUID, conv *Conversation, options BuildContextOptions) (string, error) {
	var parts []string

	// User Profile (About Me)
	if options.IncludeUserProfile {
		profile, err := cb.userProfileService.GetProfile(ctx, userID)
		if err == nil && profile != nil && profile.AboutMe != "" {
			parts = append(parts, fmt.Sprintf("## About Me\n%s\n", profile.AboutMe))
		}
	}

	// Job Application Context
	// Note: Job applications and resumes are now managed in the jobs service
	// TODO: Call jobs service API to fetch job application and resume data when needed
	if options.IncludeJobApplication && conv.JobApplicationID != nil {
		// Job application context would be retrieved from the jobs service
		_ = cb.jobApplicationService
		parts = append(parts, "## Job Application Context\n(Job application context from jobs service would be included here)\n\n")
	}

	// Resume Context
	if options.IncludeResume {
		// Resume context would be retrieved from the jobs service
		_ = cb.resumeService
		parts = append(parts, "## Resume Information\n(Resume information from jobs service would be included here)\n\n")
	}

	// Projects
	if options.IncludeProjects {
		projects, err := cb.projectService.ListProjects(ctx, userID)
		if err == nil && len(projects) > 0 {
			parts = append(parts, "## Projects\n")
			for i, p := range projects {
				if i >= 10 { // Limit to 10 most recent
					parts = append(parts, fmt.Sprintf("... and %d more projects\n", len(projects)-10))
					break
				}
				parts = append(parts, fmt.Sprintf("- **%s**: %s (Status: %s)\n", p.Name, p.Description, p.Status))
			}
			parts = append(parts, "\n")
		}
	}

	// Skills
	if options.IncludeSkills {
		// TODO: Implement skills service
		_ = cb.skillService
		_ = ctx
		// skills, err := cb.skillService.ListSkills(ctx)
		// if err == nil && len(skills) > 0 {
		// 	parts = append(parts, "## Skills\n")
		// 	skillNames := make([]string, 0, len(skills))
		// 	for _, s := range skills {
		// 		if len(skillNames) >= 20 { // Limit to 20 skills
		// 			break
		// 		}
		// 		skillNames = append(skillNames, s.Name)
		// 	}
		// 	parts = append(parts, strings.Join(skillNames, ", "))
		// 	parts = append(parts, "\n\n")
		// }
	}

	// Case Studies
	if options.IncludeCaseStudies {
		// TODO: Implement case studies service
		_ = cb.caseStudyService
		_ = ctx
		_ = userID
		// TODO: Implement case studies service
		// userIDPtr := &userID
		// caseStudies, err := cb.caseStudyService.ListCaseStudies(ctx, ...)
		// var caseStudies []interface{}
		// if len(caseStudies) > 0 {
		// 	parts = append(parts, "## Case Studies\n")
		// 	for _, cs := range caseStudies {
		// 		// Use Problem as the description since Description field doesn't exist
		// 		description := cs.Problem
		// 		if len(description) > 100 {
		// 			description = description[:100] + "..."
		// 		}
		// 		parts = append(parts, fmt.Sprintf("- **%s**: %s\n", cs.Title, description))
		// 	}
		// 	parts = append(parts, "\n")
		// }
	}

	// Technical Writings
	if options.IncludeTechnicalWritings {
		// Note: Would need to add ListTechnicalWritings method
		parts = append(parts, "## Technical Writings\n(Technical writings would be included here)\n\n")
	}

	// Posts
	if options.IncludePosts {
		// Note: Would need to add ListPosts method with user filter
		parts = append(parts, "## Posts\n(Posts would be included here)\n\n")
	}

	// Problem Solutions
	if options.IncludeProblemSolutions {
		// Note: Would need to add ListProblemSolutions method
		parts = append(parts, "## Problem Solutions\n(Problem solutions would be included here)\n\n")
	}

	// Experiences
	if options.IncludeExperiences {
		// Note: Would need to add ListExperiences method with user filter
		parts = append(parts, "## Experiences\n(Experiences would be included here)\n\n")
	}

	return strings.Join(parts, ""), nil
}

// GetDefaultContextOptions returns default context options for job application chats.
func GetDefaultContextOptions() BuildContextOptions {
	return BuildContextOptions{
		IncludeJobApplication:  true,
		IncludeResume:          true,
		IncludeUserProfile:     true,
		IncludeProjects:        true,
		IncludeCaseStudies:     true,
		IncludeTechnicalWritings: true,
		IncludePosts:           true,
		IncludeProblemSolutions: true,
		IncludeSkills:          true,
		IncludeExperiences:     true,
	}
}

