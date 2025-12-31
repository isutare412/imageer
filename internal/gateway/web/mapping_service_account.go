package web

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

func ServiceAccountToWeb(sa domain.ServiceAccount) ServiceAccount {
	return ServiceAccount{
		ID:          sa.ID,
		CreatedAt:   sa.CreatedAt,
		UpdatedAt:   sa.UpdatedAt,
		ExpireAt:    sa.ExpireAt,
		Name:        sa.Name,
		AccessScope: sa.AccessScope,
		Projects: lo.Map(sa.Projects, func(r domain.ProjectReference, _ int) ProjectReference {
			return ProjectReferenceToWeb(r)
		}),
	}
}

func ServiceAccountWithAPIKeyToWeb(sa domain.ServiceAccountWithAPIKey) ServiceAccountWithAPIKey {
	return ServiceAccountWithAPIKey{
		ID:          sa.ID,
		CreatedAt:   sa.CreatedAt,
		UpdatedAt:   sa.UpdatedAt,
		ExpireAt:    sa.ExpireAt,
		Name:        sa.Name,
		AccessScope: sa.AccessScope,
		Projects: lo.Map(sa.Projects, func(r domain.ProjectReference, _ int) ProjectReference {
			return ProjectReferenceToWeb(r)
		}),
		APIKey: sa.APIKey,
	}
}

func ServiceAccountsToWeb(sas domain.ServiceAccounts) ServiceAccounts {
	return ServiceAccounts{
		Items: lo.Map(sas.Items, func(sa domain.ServiceAccount, _ int) ServiceAccount {
			return ServiceAccountToWeb(sa)
		}),
		Total: sas.Total,
	}
}

func ListServiceAccountsAdminParamsToDomain(
	params ListServiceAccountsAdminParams,
) domain.ListServiceAccountsParams {
	var offset *int
	if params.Offset != nil {
		v := int(*params.Offset)
		offset = &v
	}

	var limit *int
	if params.Limit != nil {
		v := int(*params.Limit)
		limit = &v
	}

	return domain.ListServiceAccountsParams{
		Offset: offset,
		Limit:  limit,
	}
}

func CreateServiceAccountAdminRequestToDomain(
	req CreateServiceAccountAdminRequest,
) domain.CreateServiceAccountRequest {
	return domain.CreateServiceAccountRequest{
		Name:        req.Name,
		AccessScope: req.AccessScope,
		ProjectIDs:  req.ProjectIDs,
		ExpireAt:    req.ExpireAt,
	}
}

func UpdateServiceAccountAdminRequestToDomain(
	serviceAccountID string,
	req UpdateServiceAccountAdminRequest,
) domain.UpdateServiceAccountRequest {
	return domain.UpdateServiceAccountRequest{
		ID:          serviceAccountID,
		Name:        req.Name,
		AccessScope: req.AccessScope,
		ProjectIDs:  req.ProjectIDs,
		ExpireAt:    req.ExpireAt,
	}
}
