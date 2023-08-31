package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"etlfile"
	"etlfile/utils"
)

func main() {

	// The client and worker are heavyweight objects that should be created once per process.
	clientOptions, err := utils.LoadClientOptions()
	if err != nil {
		log.Fatalln("Failed to load Temporal Cloud environment:", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	log.Println("Starting worker with EnableSessionWorker option..")

	workerOptions := worker.Options{
		EnableSessionWorker: true, // Important for a worker to participate in the session
	}
	w := worker.New(c, "ETLFile", workerOptions)

	log.Println("Registering for workflow and activites..")

	w.RegisterWorkflow(etlfile.ETLFileProcessingWorkflow)
	w.RegisterActivity(&etlfile.Activities{BlobStore: &etlfile.BlobStore{}})

	log.Println("Running worker..")

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
