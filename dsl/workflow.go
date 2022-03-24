package dsl

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	Workflow struct {
		Variables map[string]string
		Root      Statement
	}

	Statement struct {
		Activity *ActivityInvocation
		Sequence *Sequence
		Parallel *Parallel
	}

	Sequence struct {
		Elements []*Statement
	}

	Parallel struct {
		Branches []*Statement
	}

	ActivityInvocation struct {
		Name      string
		Arguments []string
		Result    string
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]string) error
	}
)

func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	bindings := make(map[string]string)
	for k, v := range dslWorkflow.Variables {
		bindings[k] = v
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	err := dslWorkflow.Root.execute(ctx, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
}

func (b *Statement) execute(ctx workflow.Context, bindings map[string]string) error {
	if b.Parallel != nil {
		err := b.Parallel.execute(ctx, bindings)
		if err != nil {
			return err
		}
	}
	if b.Sequence != nil {
		err := b.Sequence.execute(ctx, bindings)
		if err != nil {
			return err
		}
	}
	if b.Activity != nil {
		err := b.Activity.execute(ctx, bindings)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a ActivityInvocation) execute(ctx workflow.Context, bindings map[string]string) error {
	inputParam := makeInput(a.Arguments, bindings)
	var result string
	err := workflow.ExecuteActivity(ctx, a.Name, inputParam).Get(ctx, &result)
	if err != nil {
		return err
	}
	if a.Result != "" {
		bindings[a.Result] = result
	}
	return nil
}

func (s Sequence) execute(ctx workflow.Context, bindings map[string]string) error {
	for _, a := range s.Elements {
		err := a.execute(ctx, bindings)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Parallel) execute(ctx workflow.Context, bindings map[string]string) error {
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	var activityErr error
	for _, s := range p.Branches {
		f := executeAsync(s, childCtx, bindings)
		selector.AddFuture(f, func(f workflow.Future) {
			err := f.Get(ctx, nil)
			if err != nil {
				cancelHandler()
				activityErr = err
			}
		})
	}

	for i := 0; i < len(p.Branches); i++ {
		selector.Select(ctx)
		if activityErr != nil {
			return activityErr
		}
	}

	return nil
}

func executeAsync(exe executable, ctx workflow.Context, bindings map[string]string) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, bindings)
		settable.Set(nil, err)
	})
	return future
}

func makeInput(argNames []string, argsMap map[string]string) []string {
	var args []string
	for _, arg := range argNames {
		args = append(args, argsMap[arg])
	}
	return args
}
