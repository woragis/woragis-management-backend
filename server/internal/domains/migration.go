package management

import (
	"gorm.io/gorm"

	"woragis-management-service/internal/domains/projects"
	"woragis-management-service/internal/domains/projects/projectcasestudies"
	"woragis-management-service/internal/domains/ideas"
	"woragis-management-service/internal/domains/chats"
	"woragis-management-service/internal/domains/clients"
	"woragis-management-service/internal/domains/finances"
	"woragis-management-service/internal/domains/experiences"
	"woragis-management-service/internal/domains/userpreferences"
	"woragis-management-service/internal/domains/userprofiles"
	"woragis-management-service/internal/domains/apikeys"
	"woragis-management-service/internal/domains/languages"
	"woragis-management-service/internal/domains/scheduler"
	"woragis-management-service/internal/domains/testimonials"
)

// MigrateManagementTables runs database migrations for all management-related domains
func MigrateManagementTables(db *gorm.DB) error {
	// Enable UUID extension if not already enabled
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Enable gen_random_uuid function if not already available
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
		return err
	}

	// Migrate projects tables
	if err := db.AutoMigrate(
		&projects.Project{},
		&projects.Milestone{},
		&projects.KanbanColumn{},
		&projects.KanbanCard{},
		&projects.ProjectDependency{},
		&projects.ProjectDocumentation{},
		&projects.DocumentationSection{},
		&projects.ProjectTechnology{},
		&projects.ProjectFileStructure{},
		&projects.ProjectArchitectureDiagram{},
		&projectcasestudies.ProjectCaseStudy{},
	); err != nil {
		return err
	}

	// Migrate ideas tables
	if err := db.AutoMigrate(
		&ideas.Idea{},
		&ideas.IdeaLink{},
		&ideas.IdeaNode{},
		&ideas.IdeaNodeConnection{},
		&ideas.Document{},
	); err != nil {
		return err
	}

	// Migrate chats tables
	if err := db.AutoMigrate(
		&chats.Conversation{},
		&chats.Message{},
	); err != nil {
		return err
	}

	// Migrate clients tables
	if err := db.AutoMigrate(
		&clients.Client{},
	); err != nil {
		return err
	}

	// Migrate finances tables
	if err := db.AutoMigrate(
		&finances.Transaction{},
		&finances.RecurringTemplate{},
	); err != nil {
		return err
	}

	// Migrate experiences tables
	if err := db.AutoMigrate(
		&experiences.Experience{},
	); err != nil {
		return err
	}

	// Migrate user preferences tables
	if err := db.AutoMigrate(
		&userpreferences.UserPreferences{},
	); err != nil {
		return err
	}

	// Migrate user profiles tables
	if err := db.AutoMigrate(
		&userprofiles.UserProfile{},
	); err != nil {
		return err
	}

	// Migrate API keys tables
	if err := db.AutoMigrate(
		&apikeys.APIKey{},
	); err != nil {
		return err
	}

	// Migrate languages tables
	if err := db.AutoMigrate(
		&languages.StudySession{},
		&languages.VocabularyEntry{},
	); err != nil {
		return err
	}

	// Migrate scheduler tables
	if err := db.AutoMigrate(
		&scheduler.Schedule{},
		&scheduler.ExecutionRun{},
	); err != nil {
		return err
	}

	// Migrate testimonials tables
	if err := db.AutoMigrate(
		&testimonials.Testimonial{},
	); err != nil {
		return err
	}

	return nil
}
