package main

import (
	"client/taskpb"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("could not connect", err)
	}
	defer conn.Close()

	client := taskpb.NewTaskServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createRes, err := client.CreateTask(ctx, &taskpb.TaskRequest{TaskDescription: "lyuboy task"})
	if err != nil {
		log.Fatalf("Could not create task: %v", err)
	}

	fmt.Printf("created Task: %d, Status: %s\n", createRes.TaskId, createRes.Status)

	listRes, err := client.ListTasks(context.Background(), &taskpb.Empty{})
	if err != nil {
		log.Fatalf("Could not list tasks: %v", err)
	}
	for _, task := range listRes.Tasks {
		fmt.Printf("Task ID: %d, Status: %s\n", task.TaskId, task.Status)
	}

	cancelRes, err := client.CancelTask(ctx, &taskpb.CancelRequest{TaskId: createRes.TaskId})
	if err != nil {
		log.Fatalf("Could not cancel task: %v", err)
	}
	fmt.Printf("Cancelled Task Status: %s\n", cancelRes.Status)
}
