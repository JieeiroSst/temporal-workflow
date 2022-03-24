package ctxpropagation

import (
	"context"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type (
	contextKey struct{}
	propagator struct{}
	Values struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

var PropagateKey = contextKey{}

const propagationKey = "custom-header"

func NewContextPropagator() workflow.ContextPropagator {
	return &propagator{}
}

func (s *propagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	value := ctx.Value(PropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

func (s *propagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	value := ctx.Value(PropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

func (s *propagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if value, ok := reader.Get(propagationKey); ok {
		var values Values
		if err := converter.GetDefaultDataConverter().FromPayload(value, &values); err != nil {
			return ctx, nil
		}
		ctx = context.WithValue(ctx, PropagateKey, values)
	}

	return ctx, nil
}

func (s *propagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if value, ok := reader.Get(propagationKey); ok {
		var values Values
		if err := converter.GetDefaultDataConverter().FromPayload(value, &values); err != nil {
			return ctx, nil
		}
		ctx = workflow.WithValue(ctx, PropagateKey, values)
	}

	return ctx, nil
}
