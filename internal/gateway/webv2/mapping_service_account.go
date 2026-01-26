package webv2

import (
	"github.com/samber/lo"

	"github.com/isutare412/imageer/internal/gateway/domain"
	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
)

func ServiceAccountToWeb(sa domain.ServiceAccount) gen.ServiceAccount {
	return gen.ServiceAccount{
		ID:          sa.ID,
		CreatedAt:   sa.CreatedAt,
		UpdatedAt:   sa.UpdatedAt,
		ExpireAt:    sa.ExpireAt,
		Name:        sa.Name,
		AccessScope: sa.AccessScope,
		Projects: lo.Map(sa.Projects, func(r domain.ProjectReference, _ int) gen.ProjectReference {
			return ProjectReferenceToWeb(r)
		}),
	}
}

func ServiceAccountWithAPIKeyToWeb(sa domain.ServiceAccountWithAPIKey,
) gen.ServiceAccountWithAPIKey {
	return gen.ServiceAccountWithAPIKey{
		ID:          sa.ID,
		CreatedAt:   sa.CreatedAt,
		UpdatedAt:   sa.UpdatedAt,
		ExpireAt:    sa.ExpireAt,
		Name:        sa.Name,
		AccessScope: sa.AccessScope,
		Projects: lo.Map(sa.Projects, func(r domain.ProjectReference, _ int) gen.ProjectReference {
			return ProjectReferenceToWeb(r)
		}),
		APIKey: sa.APIKey,
	}
}

func ServiceAccountsToWeb(sas domain.ServiceAccounts) gen.ServiceAccounts {
	return gen.ServiceAccounts{
		Items: lo.Map(sas.Items, func(sa domain.ServiceAccount, _ int) gen.ServiceAccount {
			return ServiceAccountToWeb(sa)
		}),
		Total: sas.Total,
	}
}

func ListServiceAccountsAdminParamsToDomain(
	params gen.ListServiceAccountsAdminParams,
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
	req gen.CreateServiceAccountAdminRequest,
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
	req gen.UpdateServiceAccountAdminRequest,
) domain.UpdateServiceAccountRequest {
	return domain.UpdateServiceAccountRequest{
		ID:          serviceAccountID,
		Name:        req.Name,
		AccessScope: req.AccessScope,
		ProjectIDs:  req.ProjectIDs,
		ExpireAt:    req.ExpireAt,
	}
}
