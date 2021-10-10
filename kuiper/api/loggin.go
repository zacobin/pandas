package api

// +build !test

import (
	"context"
	"fmt"
	"time"

	"github.com/cloustone/pandas/kuiper"
	log "github.com/cloustone/pandas/pkg/logger"
)

var _ kuiper.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    kuiper.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc kuiper.Service, logger log.Logger) kuiper.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateStreams(ctx context.Context, token string, ths ...kuiper.Stream) (saved []kuiper.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_streams for token %s and streams %s took %s to complete", token, saved, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateStreams(ctx, token, ths...)
}

func (lm *loggingMiddleware) UpdateStream(ctx context.Context, token string, stream kuiper.Stream) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_stream for token %s and stream %s took %s to complete", token, stream.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateStream(ctx, token, stream)
}

func (lm *loggingMiddleware) ViewStream(ctx context.Context, token, id string) (stream kuiper.Stream, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_stream for token %s and stream %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewStream(ctx, token, id)
}

func (lm *loggingMiddleware) ListStreams(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (_ kuiper.StreamsPage, err error) {
	defer func(begin time.Time) {
		nlog := ""
		if name != "" {
			nlog = fmt.Sprintf("with name %s ", name)
		}
		message := fmt.Sprintf("Method list_streams %sfor token %s took %s to complete", nlog, token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListStreams(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) RemoveStream(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_stream for token %s and stream %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveStream(ctx, token, id)
}

func (lm *loggingMiddleware) CreateRules(ctx context.Context, token string, rules ...kuiper.Rule) (saved []kuiper.Rule, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method create_rules for token %s and rules %s took %s to complete", token, saved, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.CreateRules(ctx, token, rules...)
}

func (lm *loggingMiddleware) UpdateRule(ctx context.Context, token string, rule kuiper.Rule) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method update_rule for token %s and rule %s took %s to complete", token, rule.ID, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.UpdateRule(ctx, token, rule)
}

func (lm *loggingMiddleware) ViewRule(ctx context.Context, token, id string) (rule kuiper.Rule, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_rule for token %s and rule %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewRule(ctx, token, id)
}

func (lm *loggingMiddleware) ListRules(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (_ kuiper.RulesPage, err error) {
	defer func(begin time.Time) {
		nlog := ""
		if name != "" {
			nlog = fmt.Sprintf("with name %s ", name)
		}
		message := fmt.Sprintf("Method list_rules %sfor token %s took %s to complete", nlog, token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListRules(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) RemoveRule(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method remove_rule for token %s and rule %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.RemoveRule(ctx, token, id)
}

func (lm *loggingMiddleware) StartRule(ctx context.Context, token string, ruleId string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method start for token %s, rules %s to complete", token, ruleId, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.StartRule(ctx, token, ruleId)
}

func (lm *loggingMiddleware) StopRule(ctx context.Context, token string, ruleId string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method stop for token %s, rules %s to complete", token, ruleId, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.StartRule(ctx, token, ruleId)
}

func (lm *loggingMiddleware) RestartRule(ctx context.Context, token string, ruleId string) (err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method restart for token %s, rules %s to complete", token, ruleId, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.StartRule(ctx, token, ruleId)
}

func (lm *loggingMiddleware) ViewPluginSource(ctx context.Context, token, id string) (source kuiper.PluginSource, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_plugin_source for token %s and plugin source %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewPluginSource(ctx, token, id)
}

func (lm *loggingMiddleware) ListPluginSources(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (_ kuiper.PluginSourcesPage, err error) {
	defer func(begin time.Time) {
		nlog := ""
		if name != "" {
			nlog = fmt.Sprintf("with name %s ", name)
		}
		message := fmt.Sprintf("Method list_plugin_sources %sfor token %s took %s to complete", nlog, token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListPluginSources(ctx, token, offset, limit, name, metadata)
}

func (lm *loggingMiddleware) ViewPluginSink(ctx context.Context, token, id string) (source kuiper.PluginSink, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method view_plugin_sink for token %s and plugin source %s took %s to complete", token, id, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ViewPluginSink(ctx, token, id)
}

func (lm *loggingMiddleware) ListPluginSinks(ctx context.Context, token string, offset, limit uint64, name string, metadata kuiper.Metadata) (_ kuiper.PluginSinksPage, err error) {
	defer func(begin time.Time) {
		nlog := ""
		if name != "" {
			nlog = fmt.Sprintf("with name %s ", name)
		}
		message := fmt.Sprintf("Method list_plugin_Sinks %sfor token %s took %s to complete", nlog, token, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))
			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.ListPluginSinks(ctx, token, offset, limit, name, metadata)
}
