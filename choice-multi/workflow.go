package choice_multi

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	OrderChoiceApple  = "apple"
	OrderChoiceBanana = "banana"
	OrderChoiceCherry = "cherry"
	OrderChoiceOrange = "orange"
)

func MultiChoiceWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var orderActivities *OrderActivities 

	var choices []string
	err := workflow.ExecuteActivity(ctx, orderActivities.GetBasketOrder).Get(ctx, &choices)
	if err != nil {
		return err
	}
	logger := workflow.GetLogger(ctx)

	var futures []workflow.Future
	for _, item := range choices {
		var f workflow.Future
		switch item {
		case OrderChoiceApple:
			f = workflow.ExecuteActivity(ctx, orderActivities.OrderApple, item)
		case OrderChoiceBanana:
			f = workflow.ExecuteActivity(ctx, orderActivities.OrderBanana, item)
		case OrderChoiceCherry:
			f = workflow.ExecuteActivity(ctx, orderActivities.OrderCherry, item)
		case OrderChoiceOrange:
			f = workflow.ExecuteActivity(ctx, orderActivities.OrderOrange, item)
		default:
			logger.Error("Unexpected order.", "Order", item)
			return errors.New("invalid choice-multi")
		}
		futures = append(futures, f)
	}

	for _, future := range futures {
		_ = future.Get(ctx, nil)
	}

	logger.Info("Workflow completed.")
	return nil
}
