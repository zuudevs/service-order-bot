/**

 filename  : task_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Task & Transaction services - calls service-order-api

 copyright Copyright (c) 2026

**/

package services

import (
	"fmt"

	"github.com/zuudevs/service-order-bot/internal/client"
	"github.com/zuudevs/service-order-bot/internal/models"
)

// ========================= TaskService =========================

type TaskService struct {
	api *client.APIClient
}

func NewTaskService(api *client.APIClient) *TaskService {
	return &TaskService{api: api}
}

func (s *TaskService) List() ([]models.Task, error) {
	var tasks []models.Task
	if err := s.api.GET("/tasks", &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetByID(id uint64) (*models.Task, error) {
	var task models.Task
	if err := s.api.GET(fmt.Sprintf("/tasks/%d", id), &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) GetByStatus(status models.TaskStatus) ([]models.Task, error) {
	var tasks []models.Task
	if err := s.api.GET(fmt.Sprintf("/tasks?status=%d", status), &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) Create(req models.CreateTaskRequest) error {
	return s.api.POST("/tasks", req, nil)
}

func (s *TaskService) Patch(id uint64, req models.PatchTaskRequest) error {
	return s.api.PATCH(fmt.Sprintf("/tasks/%d", id), req, nil)
}

func (s *TaskService) Delete(id uint64) error {
	return s.api.DELETE(fmt.Sprintf("/tasks/%d", id))
}

// ========================= TransactionService =========================

type TransactionService struct {
	api *client.APIClient
}

func NewTransactionService(api *client.APIClient) *TransactionService {
	return &TransactionService{api: api}
}

func (s *TransactionService) List() ([]models.Transaction, error) {
	var txs []models.Transaction
	if err := s.api.GET("/transactions", &txs); err != nil {
		return nil, err
	}
	return txs, nil
}

func (s *TransactionService) GetByID(id uint64) (*models.Transaction, error) {
	var tx models.Transaction
	if err := s.api.GET(fmt.Sprintf("/transactions/%d", id), &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (s *TransactionService) Create(req models.CreateTransactionRequest) error {
	return s.api.POST("/transactions", req, nil)
}

func (s *TransactionService) Delete(id uint64) error {
	return s.api.DELETE(fmt.Sprintf("/transactions/%d", id))
}

// ========================= ContactService =========================

type ContactService struct {
	api *client.APIClient
}

func NewContactService(api *client.APIClient) *ContactService {
	return &ContactService{api: api}
}

func (s *ContactService) List() ([]models.Contact, error) {
	var contacts []models.Contact
	if err := s.api.GET("/contacts", &contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *ContactService) Create(req models.CreateContactRequest) error {
	return s.api.POST("/contacts", req, nil)
}

func (s *ContactService) Delete(id uint64) error {
	return s.api.DELETE(fmt.Sprintf("/contacts/%d", id))
}