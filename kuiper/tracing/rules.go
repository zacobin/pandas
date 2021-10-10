package tracing

import (
	"context"

	"github.com/cloustone/pandas/kuiper"

	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveRuleOp               = "save_rule"
	saveRulesOp              = "save_rules"
	updateRuleOp             = "update_rule"
	updateRuleKeyOp          = "update_rule_by_key"
	retrieveRuleByIDOp       = "retrieve_rule_by_id"
	retrieveRuleByKeyOp      = "retrieve_rule_by_key"
	retrieveAllRulesOp       = "retrieve_all_rules"
	retrieveRulesByChannelOp = "retrieve_rules_by_chan"
	removeRuleOp             = "remove_rule"
	retrieveRuleIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ kuiper.RuleRepository = (*ruleRepositoryMiddleware)(nil)
	_ kuiper.RuleCache      = (*ruleCacheMiddleware)(nil)
)

type ruleRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   kuiper.RuleRepository
}

// RuleRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func RuleRepositoryMiddleware(tracer opentracing.Tracer, repo kuiper.RuleRepository) kuiper.RuleRepository {
	return ruleRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm ruleRepositoryMiddleware) Save(ctx context.Context, ths ...kuiper.Rule) ([]kuiper.Rule, error) {
	span := createSpan(ctx, trm.tracer, saveRulesOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm ruleRepositoryMiddleware) Update(ctx context.Context, th kuiper.Rule) error {
	span := createSpan(ctx, trm.tracer, updateRuleOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm ruleRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (kuiper.Rule, error) {
	span := createSpan(ctx, trm.tracer, retrieveRuleByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm ruleRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.RulesPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllRulesOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm ruleRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeRuleOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type ruleCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  kuiper.RuleCache
}

// RuleCacheMiddleware tracks request and their latency, and adds spans
// to context.
func RuleCacheMiddleware(tracer opentracing.Tracer, cache kuiper.RuleCache) kuiper.RuleCache {
	return ruleCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm ruleCacheMiddleware) Save(ctx context.Context, ruleKey string, ruleID string) error {
	span := createSpan(ctx, tcm.tracer, saveRuleOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, ruleKey, ruleID)
}

func (tcm ruleCacheMiddleware) ID(ctx context.Context, ruleKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveRuleIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, ruleKey)
}

func (tcm ruleCacheMiddleware) Remove(ctx context.Context, ruleID string) error {
	span := createSpan(ctx, tcm.tracer, removeRuleOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, ruleID)
}

func createSpan(ctx context.Context, tracer opentracing.Tracer, opName string) opentracing.Span {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return tracer.StartSpan(
			opName,
			opentracing.ChildOf(parentSpan.Context()),
		)
	}
	return tracer.StartSpan(opName)
}
