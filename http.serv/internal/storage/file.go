package storage

import (
	"encoding/json"
	"log"
	"os"

	"my-go-project/http.serv/internal/model"
)

func SaveTasks(tasks []model.TaskDTO) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Println("json error:", err)
		return err
	}

	return os.WriteFile("tasks.json", data, 0644)

}

func LoadTasks() ([]model.TaskDTO, error) {
	data, err := os.ReadFile("tasks.json")

	if err != nil {
		if os.IsNotExist(err) {
			return []model.TaskDTO{}, nil
		}
		return nil, err

	}
	var tasks []model.TaskDTO

	err = json.Unmarshal(data, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
