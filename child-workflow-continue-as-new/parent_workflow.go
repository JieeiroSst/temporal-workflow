package child_workflow_continue_as_new

import (
	"fmt"

	"go.temporal.io/sdk/workflow"
)

func SampleParentWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	execution := workflow.GetInfo(ctx).WorkflowExecution
	childID := fmt.Sprintf("child_workflow:%v", execution.RunID)
	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: childID,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	var result string
	err := workflow.ExecuteChildWorkflow(ctx, SampleChildWorkflow, 0, 5).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return err
	}

	logger.Info("Parent execution completed.", "Result", result)
	return nil
}

