package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"recetariojacqueline.com/metadata/internal/controller/metadata"
	model "recetariojacqueline.com/metadata/pkg"
	"recetariojacqueline.com/src/gen"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler { return &Handler{ctrl: ctrl} }

func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	m, err := h.ctrl.Get(ctx, req.RecipeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}

func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	if err := h.ctrl.Put(ctx, model.MetadataFromProto(req.Metadata)); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &gen.PutMetadataResponse{}, nil
}
