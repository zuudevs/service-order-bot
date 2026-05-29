/**

 filename  : models.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Data models mirroring the service-order-api

 copyright Copyright (c) 2026

**/

package models

import "time"

// ========================= Person =========================

type Person struct {
	ID         uint64    `json:"id"`
	FirstName  string    `json:"firstname"`
	MiddleName *string   `json:"middlename"`
	LastName   *string   `json:"lastname"`
	CreatedAt  time.Time `json:"created_at"`
}

func (p Person) FullName() string {
	name := p.FirstName
	if p.MiddleName != nil && *p.MiddleName != "" {
		name += " " + *p.MiddleName
	}
	if p.LastName != nil && *p.LastName != "" {
		name += " " + *p.LastName
	}
	return name
}

type CreatePersonRequest struct {
	FirstName  string  `json:"firstname"`
	MiddleName *string `json:"middlename,omitempty"`
	LastName   *string `json:"lastname,omitempty"`
}

type PatchPersonRequest struct {
	FirstName  *string `json:"firstname,omitempty"`
	MiddleName *string `json:"middlename,omitempty"`
	LastName   *string `json:"lastname,omitempty"`
}

// ========================= Contact =========================

type ContactType uint8

const (
	ContactTypeEmail ContactType = iota
	ContactTypePhone
)

func (c ContactType) String() string {
	switch c {
	case ContactTypeEmail:
		return "Email"
	case ContactTypePhone:
		return "Phone"
	default:
		return "Unknown"
	}
}

type Contact struct {
	ID          uint64      `json:"id"`
	Value       string      `json:"value"`
	ContactType ContactType `json:"contact_type"`
	IsMain      bool        `json:"is_main"`
	PersonID    uint64      `json:"person_id"`
	CreatedAt   time.Time   `json:"created_at"`
}

type CreateContactRequest struct {
	Value       string      `json:"value"`
	ContactType ContactType `json:"contact_type"`
	IsMain      bool        `json:"is_main"`
	PersonID    uint64      `json:"person_id"`
}

// ========================= Order =========================

type OrderStatus uint8

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusAccepted  OrderStatus = 1
	OrderStatusRejected  OrderStatus = 2
	OrderStatusRevised   OrderStatus = 3
	OrderStatusCompleted OrderStatus = 4
)

func (o OrderStatus) String() string {
	switch o {
	case OrderStatusPending:
		return "⏳ Pending"
	case OrderStatusAccepted:
		return "✅ Accepted"
	case OrderStatusRejected:
		return "❌ Rejected"
	case OrderStatusRevised:
		return "🔄 Revised"
	case OrderStatusCompleted:
		return "🏁 Completed"
	default:
		return "Unknown"
	}
}

func (o OrderStatus) Emoji() string {
	switch o {
	case OrderStatusPending:
		return "⏳"
	case OrderStatusAccepted:
		return "✅"
	case OrderStatusRejected:
		return "❌"
	case OrderStatusRevised:
		return "🔄"
	case OrderStatusCompleted:
		return "🏁"
	default:
		return "❓"
	}
}

type Order struct {
	ID           uint64      `json:"id"`
	Status       OrderStatus `json:"status"`
	OrderDate    time.Time   `json:"order_date"`
	LastModified time.Time   `json:"last_modified"`
	TotalPrice   uint64      `json:"total_price"`
	PersonID     *uint64     `json:"person_id"`
}

type CreateOrderRequest struct {
	Status   OrderStatus `json:"status"`
	PersonID *uint64     `json:"person_id,omitempty"`
}

type PatchOrderRequest struct {
	Status     *OrderStatus `json:"status,omitempty"`
	TotalPrice *uint64      `json:"total_price,omitempty"`
}

// ========================= Task =========================

type TaskStatus uint8

const (
	TaskStatusPending    TaskStatus = 0
	TaskStatusOnProgress TaskStatus = 1
	TaskStatusComplete   TaskStatus = 2
)

func (t TaskStatus) String() string {
	switch t {
	case TaskStatusPending:
		return "⏳ Pending"
	case TaskStatusOnProgress:
		return "🔧 On Progress"
	case TaskStatusComplete:
		return "✅ Complete"
	default:
		return "Unknown"
	}
}

type Task struct {
	ID          uint64     `json:"id"`
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Price       uint64     `json:"price"`
	Due         time.Time  `json:"due"`
}

type CreateTaskRequest struct {
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
	Price       uint64     `json:"price"`
	Due         time.Time  `json:"due"`
}

type PatchTaskRequest struct {
	Subject     *string     `json:"subject,omitempty"`
	Description *string     `json:"description,omitempty"`
	Price       *uint64     `json:"price,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
}

// ========================= Transaction =========================

type TransactionStatus uint8
type TransactionMethod uint8

const (
	TransactionStatusPending   TransactionStatus = 0
	TransactionStatusCanceled  TransactionStatus = 1
	TransactionStatusSuccess   TransactionStatus = 2
)

const (
	TransactionMethodCash    TransactionMethod = 0
	TransactionMethodEWallet TransactionMethod = 1
	TransactionMethodBank    TransactionMethod = 2
)

func (t TransactionStatus) String() string {
	switch t {
	case TransactionStatusPending:
		return "⏳ Pending"
	case TransactionStatusCanceled:
		return "❌ Canceled"
	case TransactionStatusSuccess:
		return "✅ Success"
	default:
		return "Unknown"
	}
}

func (t TransactionMethod) String() string {
	switch t {
	case TransactionMethodCash:
		return "💵 Cash"
	case TransactionMethodEWallet:
		return "📱 E-Wallet"
	case TransactionMethodBank:
		return "🏦 Bank Transfer"
	default:
		return "Unknown"
	}
}

type Transaction struct {
	ID           uint64            `json:"id"`
	Timestamp    time.Time         `json:"timestamp"`
	Status       TransactionStatus `json:"status"`
	Method       TransactionMethod `json:"method"`
	Amount       uint64            `json:"amount"`
	EvidencePath string            `json:"evidance_path"`
	OrderID      *uint64           `json:"order_id"`
}

type CreateTransactionRequest struct {
	Status       TransactionStatus `json:"status"`
	Method       TransactionMethod `json:"method"`
	Amount       uint64            `json:"amount"`
	EvidencePath string            `json:"evidence_path"`
	OrderID      *uint64           `json:"order_id,omitempty"`
}

// ========================= DetailTask =========================

type DetailTask struct {
	ID     uint64 `json:"id"`
	TaskID uint64 `json:"task_id"`
}

type CreateDetailTaskRequest struct {
	TaskID uint64 `json:"task_id"`
}

// ========================= API Response =========================

type APISuccess struct {
	Success bool `json:"success"`
}

type APIError struct {
	Error string `json:"error"`
}