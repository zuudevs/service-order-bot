/**

 filename  : session.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : In-memory session state for multi-step bot conversations

 copyright Copyright (c) 2026

**/

package session

import (
	"sync"
	"time"
)

// State represents the current conversation state of a user
type State string

const (
	StateIdle State = ""

	// Person states
	StateCreatePersonFirstName  State = "create_person_firstname"
	StateCreatePersonMiddleName State = "create_person_middlename"
	StateCreatePersonLastName   State = "create_person_lastname"
	StateDeletePersonID         State = "delete_person_id"
	StateEditPersonID           State = "edit_person_id"
	StateEditPersonField        State = "edit_person_field"
	StateEditPersonValue        State = "edit_person_value"

	// Order states
	StateCreateOrderPersonID State = "create_order_person_id"
	StateCreateOrderStatus   State = "create_order_status"
	StateUpdateOrderID       State = "update_order_id"
	StateUpdateOrderStatus   State = "update_order_status"
	StateUpdateOrderPrice    State = "update_order_price"
	StateDeleteOrderID       State = "delete_order_id"

	// Task states
	StateCreateTaskSubject      State = "create_task_subject"
	StateCreateTaskDescription  State = "create_task_description"
	StateCreateTaskPrice        State = "create_task_price"
	StateCreateTaskDue          State = "create_task_due"
	StateUpdateTaskID           State = "update_task_id"
	StateUpdateTaskStatus       State = "update_task_status"
	StateDeleteTaskID           State = "delete_task_id"

	// Transaction states
	StateCreateTxOrderID   State = "create_tx_order_id"
	StateCreateTxAmount    State = "create_tx_amount"
	StateCreateTxMethod    State = "create_tx_method"
	StateCreateTxEvidence  State = "create_tx_evidence"
	StateDeleteTxID        State = "delete_tx_id"

	// Contact states
	StateCreateContactPersonID State = "create_contact_person_id"
	StateCreateContactType     State = "create_contact_type"
	StateCreateContactValue    State = "create_contact_value"
	StateDeleteContactID       State = "delete_contact_id"
)

// Session holds the state and temporary data for a user conversation
type Session struct {
	State     State
	Data      map[string]any
	UpdatedAt time.Time
}

// Manager manages user sessions in memory
type Manager struct {
	mu       sync.RWMutex
	sessions map[int64]*Session
	ttl      time.Duration
}

// NewManager creates a new session manager
func NewManager(ttl time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[int64]*Session),
		ttl:      ttl,
	}

	// Start background cleanup goroutine
	go m.cleanup()

	return m
}

// Get retrieves or creates a session for the given user ID
func (m *Manager) Get(userID int64) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[userID]
	if !ok {
		s = &Session{
			State:     StateIdle,
			Data:      make(map[string]any),
			UpdatedAt: time.Now(),
		}
		m.sessions[userID] = s
	}

	return s
}

// Set updates the state for a user
func (m *Manager) Set(userID int64, state State) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[userID]
	if !ok {
		s = &Session{
			Data: make(map[string]any),
		}
		m.sessions[userID] = s
	}

	s.State = state
	s.UpdatedAt = time.Now()
}

// SetData stores a value in session data
func (m *Manager) SetData(userID int64, key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[userID]
	if !ok {
		s = &Session{
			Data: make(map[string]any),
		}
		m.sessions[userID] = s
	}

	s.Data[key] = value
	s.UpdatedAt = time.Now()
}

// GetData retrieves a value from session data
func (m *Manager) GetData(userID int64, key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[userID]
	if !ok {
		return nil, false
	}

	v, found := s.Data[key]
	return v, found
}

// Clear resets the session to idle state and clears all data
func (m *Manager) Clear(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[userID] = &Session{
		State:     StateIdle,
		Data:      make(map[string]any),
		UpdatedAt: time.Now(),
	}
}

// cleanup removes expired sessions periodically
func (m *Manager) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		for id, s := range m.sessions {
			if time.Since(s.UpdatedAt) > m.ttl {
				delete(m.sessions, id)
			}
		}
		m.mu.Unlock()
	}
}