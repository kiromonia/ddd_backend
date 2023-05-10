package repository

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"backend/model"
)

type ITaskRepository interface {
	GetAllTasks(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	if err := tr.db.Joins("User").Where("user_id = ?", userId).Order("created_at").Find(tasks).Error; err != nil {
		return err
	}

	return nil
}

func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	if err := tr.db.Joins("User").Where("user_id = ?", userId).First(task, taskId).Error; err != nil {
		return err
	}

	return nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}

	return nil
}

func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	// Calauseを使うとTaskモデルにUpdateした内容を書き込んでくれる
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id = ? AND user_id = ?", taskId, userId).Update("title", task.Title)
	if result.Error != nil {
		return result.Error
	}

	// 更新された行数が0の場合はエラーを返す
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	result := tr.db.Where("id = ? AND user_id = ?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}

	// 削除された行数が0の場合はエラーを返す
	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
