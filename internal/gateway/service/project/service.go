package project

import (
	"context"
	"fmt"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/port"
	"github.com/isutare412/imageer/pkg/validation"
)

type Service struct {
	projectRepo port.ProjectRepository
}

func NewService(projectRepo port.ProjectRepository) *Service {
	return &Service{
		projectRepo: projectRepo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("finding project: %w", err)
	}
	return project, nil
}

func (s *Service) List(
	ctx context.Context, params domain.ListProjectsParams,
) (domain.Projects, error) {
	projects, err := s.projectRepo.List(ctx, params)
	if err != nil {
		return domain.Projects{}, fmt.Errorf("listing projects: %w", err)
	}
	return projects, nil
}

func (s *Service) Create(ctx context.Context, req domain.CreateProjectRequest) (domain.Project, error) {
	if err := validation.Validate(req); err != nil {
		return domain.Project{}, fmt.Errorf("validating request: %w", err)
	}

	project := req.ToProject()
	project, err := s.projectRepo.Create(ctx, project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("creating project: %w", err)
	}

	return project, nil
}

func (s *Service) Update(
	ctx context.Context, req domain.UpdateProjectRequest,
) (domain.Project, error) {
	if err := validation.Validate(req); err != nil {
		return domain.Project{}, fmt.Errorf("validating request: %w", err)
	}

	project, err := s.projectRepo.Update(ctx, req)
	if err != nil {
		return domain.Project{}, fmt.Errorf("updating project: %w", err)
	}

	return project, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting project: %w", err)
	}
	return nil
}
