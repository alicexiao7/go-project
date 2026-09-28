package service

import (
	"errors"

	"my-go-project/http.serv/internal/model"

)


type TaskStorage interface {
    LoadTasks() ([]model.TaskDTO, error)
    SaveTasks(tasks []model.TaskDTO) error
}


func CreateTask(title string, storage TaskStorage) (model.TaskDTO, error) {

	if title ==""{
		return model.TaskDTO{},errors.New("title empty")
	}

	tasks, err := storage.LoadTasks()
	if err != nil{
		return model.TaskDTO{}, err
	}


	nextID:=1

	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}

	task := model.TaskDTO{
	ID:          nextID,
	Title:       title,
	IsCompleted: false,
	}

	tasks = append(tasks, task)
	

	err = storage.SaveTasks(tasks)
	if err != nil {
		return model.TaskDTO{}, err
	}




	return task, nil

}