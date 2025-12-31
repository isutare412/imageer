package web

import "github.com/isutare412/imageer/internal/gateway/domain"

func ProjectReferenceToWeb(r domain.ProjectReference) ProjectReference {
	return ProjectReference{
		ID:   r.ID,
		Name: r.Name,
	}
}
