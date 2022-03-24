package main

import (
	"context"
	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	choice_multi "github.com/temporalio/samples-go/choice-multi"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "multi_choice_" + uuid.New(),
		TaskQueue: "choice-multi",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, choice_multi.MultiChoiceWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

}
