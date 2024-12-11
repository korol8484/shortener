package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/korol8484/shortener/internal/app/grpc/service"
	"github.com/korol8484/shortener/internal/app/grpc/service/contract"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/util"
)

// Handler Grpc service
type Handler struct {
	service.InternalServer
	usecase *usecase.Usecase
}

// NewHandler factory
func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

// Ping - Health check
func (h *Handler) Ping(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {
	if h.usecase.Ping() {
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(codes.Unavailable, "service unavailable")
}

// HandleShort - Shorten url
func (h *Handler) HandleShort(ctx context.Context, in *contract.RequestShort) (*contract.ResponseShort, error) {
	userID, ok := util.ReadUserIDFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	if in.GetData() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "data can't by empty")
	}

	u, err := h.usecase.CreateURL(ctx, in.GetData(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	return &contract.ResponseShort{Data: h.usecase.FormatAlias(u)}, nil
}

// HandleGet Load url by alias
func (h *Handler) HandleGet(ctx context.Context, in *contract.RequestFindByAlias) (*contract.ResponseShort, error) {
	if in.GetAlias() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "data can't by empty")
	}

	u, err := h.usecase.LoadByAlias(ctx, in.GetAlias())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	return &contract.ResponseShort{Data: h.usecase.FormatAlias(u)}, nil
}

// HandleBatch Handler for a collection of shortened links
func (h *Handler) HandleBatch(ctx context.Context, in *contract.RequestBatch) (*contract.ResponseBatch, error) {
	userID, ok := util.ReadUserIDFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	if len(in.GetBatch()) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "bad request")
	}

	resp := &contract.ResponseBatch{Batch: make([]*contract.Batch, 0, len(in.GetBatch()))}
	for _, b := range in.GetBatch() {
		u, err := h.usecase.CreateURL(ctx, b.GetOriginalUrl(), userID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "bad request")
		}

		resp.Batch = append(resp.Batch, &contract.Batch{
			CorrelationId: b.GetCorrelationId(),
			OriginalUrl:   h.usecase.FormatAlias(u),
		})
	}

	return resp, nil
}

// UserURL Handler for list user shortened links
func (h *Handler) UserURL(ctx context.Context, in *empty.Empty) (*contract.ResponseUserUrl, error) {
	userID, ok := util.ReadUserIDFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	list, err := h.usecase.LoadAllUserURL(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	resp := &contract.ResponseUserUrl{Url: make([]*contract.Url, 0, len(list))}
	for _, u := range list {
		resp.Url = append(resp.Url, &contract.Url{
			ShortUrl: u.URL,
			Alias:    h.usecase.FormatAlias(u),
		})
	}

	return resp, nil
}

// BatchDelete handler for a collection of delete user shorten URLs
func (h *Handler) BatchDelete(ctx context.Context, in *contract.RequestBatchDelete) (*empty.Empty, error) {
	userID, ok := util.ReadUserIDFromCtx(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "bad request")
	}

	h.usecase.AddToDelete(in.GetAlias(), userID)

	return &empty.Empty{}, nil
}

// Stats - return service statistic
func (h *Handler) Stats(ctx context.Context, in *empty.Empty) (*contract.ResponseStats, error) {
	stats, err := h.usecase.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error load stats: %s", err)
	}

	return &contract.ResponseStats{
		Urls:  stats.Urls,
		Users: stats.Users,
	}, nil
}
