package branch

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SampleBranchWorkflow(ctx workflow.Context, totalBranches int) (result []string, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("SampleBranchWorkflow begin")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var futures []workflow.Future
	for i := 1; i <= totalBranches; i++ {
		activityInput := fmt.Sprintf("branch %d of %d.", i, totalBranches)
		future := workflow.ExecuteActivity(ctx, SampleActivity, activityInput)
		futures = append(futures, future)
	}
	logger.Info("Activities started")

	for _, future := range futures {
		var singleResult string
		err = future.Get(ctx, &singleResult)
		logger.Info("Activity returned with result", "resutl", singleResult)
		if err != nil {
			return
		}
		result = append(result, singleResult)
	}

	logger.Info("SampleBranchWorkflow end")
	return
}
