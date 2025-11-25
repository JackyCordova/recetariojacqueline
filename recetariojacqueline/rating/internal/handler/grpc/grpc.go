package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"recetariojacqueline.com/rating/internal/controller/rating"
	"recetariojacqueline.com/rating/pkg/model"
	"recetariojacqueline.com/src/gen"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler { return &Handler{ctrl: ctrl} }

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	avg, count, err := h.ctrl.GetAverage(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &gen.GetAggregatedRatingResponse{Avg: avg, Count: int32(count)}, nil
}

func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if err := h.ctrl.Put(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), req.UserId, req.Value); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &gen.PutRatingResponse{}, nil
}
