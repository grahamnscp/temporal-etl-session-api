package main

import (
	"context"
	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	"etlfile"
	"etlfile/utils"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	clientOptions, err := utils.LoadClientOptions()
	if err != nil {
		log.Fatalln("Failed to load Temporal Cloud environment:", err)
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	fileID := uuid.New()
	workflowOptions := client.StartWorkflowOptions{
		ID:        "etl_file_" + fileID,
		TaskQueue: "ETLFile",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, etlfile.ETLFileProcessingWorkflow, fileID)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
