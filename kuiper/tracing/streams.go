package tracing

import (
	"context"

	"github.com/cloustone/pandas/kuiper"

	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveStreamOp               = "save_stream"
	saveStreamsOp              = "save_streams"
	updateStreamOp             = "update_stream"
	updateStreamKeyOp          = "update_stream_by_key"
	retrieveStreamByIDOp       = "retrieve_stream_by_id"
	retrieveStreamByKeyOp      = "retrieve_stream_by_key"
	retrieveAllStreamsOp       = "retrieve_all_streams"
	retrieveStreamsByChannelOp = "retrieve_streams_by_chan"
	removeStreamOp             = "remove_stream"
	retrieveStreamIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ kuiper.StreamRepository = (*streamRepositoryMiddleware)(nil)
	_ kuiper.StreamCache      = (*streamCacheMiddleware)(nil)
)

type streamRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   kuiper.StreamRepository
}

// StreamRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func StreamRepositoryMiddleware(tracer opentracing.Tracer, repo kuiper.StreamRepository) kuiper.StreamRepository {
	return streamRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm streamRepositoryMiddleware) Save(ctx context.Context, ths ...kuiper.Stream) ([]kuiper.Stream, error) {
	span := createSpan(ctx, trm.tracer, saveStreamsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm streamRepositoryMiddleware) Update(ctx context.Context, th kuiper.Stream) error {
	span := createSpan(ctx, trm.tracer, updateStreamOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm streamRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (kuiper.Stream, error) {
	span := createSpan(ctx, trm.tracer, retrieveStreamByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm streamRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.StreamsPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllStreamsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm streamRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeStreamOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type streamCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  kuiper.StreamCache
}

// StreamCacheMiddleware tracks request and their latency, and adds spans
// to context.
func StreamCacheMiddleware(tracer opentracing.Tracer, cache kuiper.StreamCache) kuiper.StreamCache {
	return streamCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm streamCacheMiddleware) Save(ctx context.Context, streamKey string, streamID string) error {
	span := createSpan(ctx, tcm.tracer, saveStreamOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, streamKey, streamID)
}

func (tcm streamCacheMiddleware) ID(ctx context.Context, streamKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveStreamIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, streamKey)
}

func (tcm streamCacheMiddleware) Remove(ctx context.Context, streamID string) error {
	span := createSpan(ctx, tcm.tracer, removeStreamOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, streamID)
}
