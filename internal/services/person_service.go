/**

 filename  : person_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Person service - calls service-order-api

 copyright Copyright (c) 2026

**/

package services

import (
	"fmt"

	"github.com/zuudevs/service-order-bot/internal/client"
	"github.com/zuudevs/service-order-bot/internal/models"
)

type PersonService struct {
	api *client.APIClient
}

func NewPersonService(api *client.APIClient) *PersonService {
	return &PersonService{api: api}
}

func (s *PersonService) List() ([]models.Person, error) {
	var persons []models.Person
	if err := s.api.GET("/persons", &persons); err != nil {
		return nil, err
	}
	return persons, nil
}

func (s *PersonService) GetByID(id uint64) (*models.Person, error) {
	var person models.Person
	if err := s.api.GET(fmt.Sprintf("/persons/%d", id), &person); err != nil {
		return nil, err
	}
	return &person, nil
}

func (s *PersonService) Create(req models.CreatePersonRequest) error {
	return s.api.POST("/persons", req, nil)
}

func (s *PersonService) Patch(id uint64, req models.PatchPersonRequest) error {
	return s.api.PATCH(fmt.Sprintf("/persons/%d", id), req, nil)
}

func (s *PersonService) Delete(id uint64) error {
	return s.api.DELETE(fmt.Sprintf("/persons/%d", id))
}

func (s *PersonService) GetContacts(personID uint64) ([]models.Contact, error) {
	var contacts []models.Contact
	if err := s.api.GET(fmt.Sprintf("/contacts?person_id=%d", personID), &contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *PersonService) GetOrders(personID uint64) ([]models.Order, error) {
	var orders []models.Order
	if err := s.api.GET(fmt.Sprintf("/orders?person_id=%d", personID), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}