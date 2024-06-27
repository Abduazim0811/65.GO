package method

import (
	"Task/taskpb"
	"database/sql"
	"log"
	"time"
)

func StoreNewTask(db *sql.DB, req *taskpb.TaskRequest) (*taskpb.TaskResponse, error) {
	time := time.Now().Format(time.ANSIC)

	query := "INSERT INTO tasks(task, started_time) VALUES($1,$2) RETURNING id,task"
	var (
		id int32
		task string
	)
	err := db.QueryRow(query, req.TaskDescription, time).Scan(&id, &task)
	if err != nil {
		log.Println("unable to insert task:", err)
		return nil, err
	}

	return &taskpb.TaskResponse{
		TaskId: id,
		Status: "Done",
	}, nil
}

func GetAllTask(db *sql.DB)(*taskpb.TaskList,error){
	query:= "SELECT id, task FROM tasks;"

	rows, err := db.Query(query)
	if err!=nil{
		log.Println("unable to getall task:", err)
		return nil, err
	}

	var  taskList []*taskpb.TaskResponse

	for rows.Next(){
		var (
			id int32
			task string
		)

		if err:= rows.Scan(&id, &task); err!=nil{
			log.Println(err)
			return nil, err
		}

		taskList= append(taskList, &taskpb.TaskResponse{
			TaskId: id,
			Status: task,
		})
	}

	return &taskpb.TaskList{Tasks: taskList}, nil
}

func DeleteTask(db *sql.DB, req *taskpb.CancelRequest)error{
	query := "DELETE FROM tasks WHERE id=$1"
	_, err:= db.Exec(query, req.TaskId)
	if err!=nil{
		log.Println("unable to delete tasks:", err)
		return err
	}

	return nil
}