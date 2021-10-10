package api

// +build !test

import (
	"context"
	"time"

	"github.com/cloustone/pandas/kuiper"
	"github.com/go-kit/kit/metrics"
)

var _ kuiper.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     kuiper.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc kuiper.Service, counter metrics.Counter, latency metrics.Histogram) kuiper.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateStreams(ctx context.Context, token string, ths ...kuiper.Stream) (saved []kuiper.Stream, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_kuiper").Add(1)
		ms.latency.With("method", "create_kuiper").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateStreams(ctx, token, ths...)
}

func (ms *metricsMiddleware) UpdateStream(ctx context.Context, token string, stream kuiper.Stream) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_stream").Add(1)
		ms.latency.With("method", "update_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateStream(ctx, token, stream)
}

func (ms *metricsMiddleware) ViewStream(ctx context.Context, token, id string) (kuiper.Stream, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_stream").Add(1)
		ms.latency.With("method", "view_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewStream(ctx, token, id)
}

func (ms *metricsMiddleware) ListStreams(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.StreamsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_kuiper").Add(1)
		ms.latency.With("method", "list_kuiper").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListStreams(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) RemoveStream(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_stream").Add(1)
		ms.latency.With("method", "remove_stream").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveStream(ctx, token, id)
}

func (ms *metricsMiddleware) CreateRules(ctx context.Context, token string, rules ...kuiper.Rule) (saved []kuiper.Rule, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_rules").Add(1)
		ms.latency.With("method", "create_rules").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateRules(ctx, token, rules...)
}

func (ms *metricsMiddleware) UpdateRule(ctx context.Context, token string, rule kuiper.Rule) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_rule").Add(1)
		ms.latency.With("method", "update_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateRule(ctx, token, rule)
}

func (ms *metricsMiddleware) ViewRule(ctx context.Context, token, id string) (kuiper.Rule, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_rule").Add(1)
		ms.latency.With("method", "view_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewRule(ctx, token, id)
}

func (ms *metricsMiddleware) StartRule(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "start_rule").Add(1)
		ms.latency.With("method", "start_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.StartRule(ctx, token, id)
}

func (ms *metricsMiddleware) StopRule(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "stop_rule").Add(1)
		ms.latency.With("method", "stop_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.StopRule(ctx, token, id)
}

func (ms *metricsMiddleware) RestartRule(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "restart_rule").Add(1)
		ms.latency.With("method", "restart_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RestartRule(ctx, token, id)
}

func (ms *metricsMiddleware) ListRules(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.RulesPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_rules").Add(1)
		ms.latency.With("method", "list_rules").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListRules(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) RemoveRule(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_rule").Add(1)
		ms.latency.With("method", "remove_rule").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveRule(ctx, token, id)
}

func (ms *metricsMiddleware) ViewPluginSource(ctx context.Context, token, id string) (kuiper.PluginSource, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_plugin_source").Add(1)
		ms.latency.With("method", "view_plugin_source").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewPluginSource(ctx, token, id)
}

func (ms *metricsMiddleware) ListPluginSources(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.PluginSourcesPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_plugin_sources").Add(1)
		ms.latency.With("method", "list_plugin_sources").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListPluginSources(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) ViewPluginSink(ctx context.Context, token, id string) (kuiper.PluginSink, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_plugin_sink").Add(1)
		ms.latency.With("method", "view_plugin_sink").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewPluginSink(ctx, token, id)
}

func (ms *metricsMiddleware) ListPluginSinks(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.PluginSinksPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_plugin_sources").Add(1)
		ms.latency.With("method", "list_plugin_sources").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListPluginSinks(ctx, token, offset, limit, name, metadata)
}
