package Query

import (
	"baitapweek2/Db"
	"errors"
)

func InsertTask(task Db.Task) (int64, error) {
	result := Db.DB.Create(&task)
	if result.Error != nil {
		return 0, nil

	}
	return result.RowsAffected, nil

}

func GetTasks() ([]Db.Task, error) {
	var taskQ []Db.Task
	result := Db.DB.Preload("Category").Find(&taskQ)
	if result.Error != nil {
		return nil, errors.New("lấy dữ liệu tasks thất bại")
	}
	return taskQ, nil
}

func GetTasksById(id int) (*Db.Task, error) {
	var taskQ Db.Task
	result := Db.DB.Where("task_id = ?", id).First(&taskQ)
	if result.Error != nil {
		return nil, errors.New("lấy dữ liệu tasks thất bại")
	}
	return &taskQ, nil
}

func DeleteTask(Id int) (int, error) {

	result := Db.DB.Delete(&Db.Task{}, Id)

	if result.Error != nil {
		return 0, errors.New("Dtelete tasks thất bại")
	}
	return int(result.RowsAffected), nil

}

func EditTask(task Db.Task) (int, error) {
	// result := Db.DB.Model(&Db.Task{}).Where("task_id = ? ", task.TaskID).Updates(task)
	updates := map[string]interface{}{
		"title":       task.Title,
		"description": task.Description,
		"category_id": task.CategoryID,
		"due_date":    task.DueDate,
		"status":      task.Status,
	}

	// Thực hiện update
	result := Db.DB.Model(&Db.Task{}).Where("task_id = ?", task.TaskID).Updates(updates)

	if result.Error != nil {
		return 0, errors.New("Update tasks thất bại")
	}

	return int(result.RowsAffected), nil
}
