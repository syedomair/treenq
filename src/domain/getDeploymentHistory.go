package domain

import (
	"context"

	"github.com/treenq/treenq/pkg/vel"
)

type GetDeploymentHistoryRequest struct {
	AppID string
}

type GetDeploymentHistoryResponse struct {
	History []AppDefinition
}

func (h *Handler) GetDeploymentHistory(ctx context.Context, req GetDeploymentHistoryRequest) (GetDeploymentHistoryResponse, *vel.Error) {
	history, err := h.db.GetDeploymentHistory(ctx, req.AppID)
	if err != nil {
		return GetDeploymentHistoryResponse{}, &vel.Error{
			Code:    "UNKNOWN",
			Message: err.Error(),
		}
	}

	return GetDeploymentHistoryResponse{
		History: history,
	}, nil
}
