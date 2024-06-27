package taskservice

import (
	"Task/internal/method"
	"Task/taskpb"
	"context"
	"database/sql"
	"log"
	"time"
)


type TaskService struct{
	taskpb.UnimplementedTaskServiceServer
	db *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService{
	return &TaskService{db:db}
}

func (t *TaskService)CreateTask(ctx context.Context, req *taskpb.TaskRequest)(*taskpb.TaskResponse, error){
	task, err:=method.StoreNewTask(t.db, req)  
	if err!=nil{
	  log.Println("failed to store task:", err)
	  return nil, err
	}
	for i:=0; i<2; i++{
		select{
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1* time.Second):
			log.Println("Working")
		}
	}
	return task, nil
  }

func (t *TaskService)ListTasks(ctx context.Context, req *taskpb.Empty)(*taskpb.TaskList, error){
	taskList, err:=method.GetAllTask(t.db)
	if err != nil {
		log.Println("failed to get all tasks:", err)
		return nil, err
	}

	return taskList, nil
}

func (t *TaskService)CancelTask(ctx context.Context, req *taskpb.CancelRequest)(*taskpb.CancelResponse, error){
	err:=method.DeleteTask(t.db, req)
	if err != nil {
		return nil, err
	}

	return &taskpb.CancelResponse{Status: "deleted"}, nil
}



