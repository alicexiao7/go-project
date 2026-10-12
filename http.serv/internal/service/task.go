package service

import (
	"errors"

	"my-go-project/http.serv/internal/model"

)


type TaskStorage interface {
    LoadTasks() ([]model.TaskDTO, error)
    SaveTasks(tasks []model.TaskDTO) error
}

type TaskService struct{
	storage TaskStorage
}

func NewTaskService(storage TaskStorage)*TaskService{
	return &TaskService{
		storage:storage,
	}

}



func (s*TaskService)CreateTask(title string) (model.TaskDTO, error) {

	if title ==""{
		return model.TaskDTO{},errors.New("title empty")
	}

	tasks, err := s.storage.LoadTasks()
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
	

	err = s.storage.SaveTasks(tasks)
	if err != nil {
		return model.TaskDTO{}, err
	}




	return task, nil

}

func (s*TaskService) GetTasks()([]model.TaskDTO, error){
	tasks, err := s.storage.LoadTasks()
	if err != nil{
		return nil, err
	}

	return tasks, nil

}	

func (s*TaskService)GetTaskByID(id int)(model.TaskDTO, error){

	tasks, err := s.storage.LoadTasks()
	if err != nil{
		return model.TaskDTO{}, err
	}
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return model.TaskDTO{},errors.New("task not found")

}

func (s*TaskService)UpdateTask(id int, updatedTask model.TaskDTO)(model.TaskDTO, error){

	tasks, err := s.storage.LoadTasks()
	if err != nil{
		return model.TaskDTO{}, err
	}
	
	
	if updatedTask.Title == "" {

		return  model.TaskDTO{},errors.New("task title is empty")
	}

	for i, task := range tasks {
		if task.ID == id {
			updatedTask.ID = id
			tasks[i] = updatedTask

			err = s.storage.SaveTasks(tasks)
			
		
			if err != nil {
				return  model.TaskDTO{},err
			}
			return updatedTask, nil
		}

		
		
	}
	return model.TaskDTO{}, errors.New("task id not found")

}

func (s*TaskService)DeleteTask(id int) error {
    tasks, err := s.storage.LoadTasks()
    if err != nil {
        return err
    }

    for i, task := range tasks {
        if task.ID == id {
            tasks = append(tasks[:i], tasks[i+1:]...)

            err = s.storage.SaveTasks(tasks)
            if err != nil {
                return err
            }

            return nil
        }
    }

    return errors.New("task id not found")
}