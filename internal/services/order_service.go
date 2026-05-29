/**

 filename  : order_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Order service - calls service-order-api

 copyright Copyright (c) 2026

**/

package services

import (
	"fmt"

	"github.com/zuudevs/service-order-bot/internal/client"
	"github.com/zuudevs/service-order-bot/internal/models"
)

type OrderService struct {
	api *client.APIClient
}

func NewOrderService(api *client.APIClient) *OrderService {
	return &OrderService{api: api}
}

func (s *OrderService) List() ([]models.Order, error) {
	var orders []models.Order
	if err := s.api.GET("/orders", &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) GetByID(id uint64) (*models.Order, error) {
	var order models.Order
	if err := s.api.GET(fmt.Sprintf("/orders/%d", id), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetByStatus(status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	if err := s.api.GET(fmt.Sprintf("/orders?status=%d", status), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) Create(req models.CreateOrderRequest) error {
	return s.api.POST("/orders", req, nil)
}

func (s *OrderService) Patch(id uint64, req models.PatchOrderRequest) error {
	return s.api.PATCH(fmt.Sprintf("/orders/%d", id), req, nil)
}

func (s *OrderService) Delete(id uint64) error {
	return s.api.DELETE(fmt.Sprintf("/orders/%d", id))
}

func (s *OrderService) GetTransactions(orderID uint64) ([]models.Transaction, error) {
	var txs []models.Transaction
	if err := s.api.GET(fmt.Sprintf("/transactions?order_id=%d", orderID), &txs); err != nil {
		return nil, err
	}
	return txs, nil
}