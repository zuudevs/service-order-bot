/**

 filename  : task_handler.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Telegram handlers for tasks and transactions

 copyright Copyright (c) 2026

**/

package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zuudevs/service-order-bot/internal/keyboards"
	"github.com/zuudevs/service-order-bot/internal/models"
	"github.com/zuudevs/service-order-bot/internal/session"
	"github.com/zuudevs/service-order-bot/internal/services"
)

// ========================= TaskHandler =========================

type TaskHandler struct {
	svc      *services.TaskService
	sessions *session.Manager
	bot      *tgbotapi.BotAPI
}

func NewTaskHandler(
	svc *services.TaskService,
	sessions *session.Manager,
	bot *tgbotapi.BotAPI,
) *TaskHandler {
	return &TaskHandler{svc: svc, sessions: sessions, bot: bot}
}

func (h *TaskHandler) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	h.bot.Send(msg)
}

func (h *TaskHandler) HandleMenu(chatID int64) {
	h.send(chatID, "✅ *Task Management*\nChoose an action:", keyboards.TaskMenu())
}

func (h *TaskHandler) HandleList(chatID int64, tasks []models.Task) {
	if tasks == nil {
		var err error
		tasks, err = h.svc.List()
		if err != nil {
			h.send(chatID, "❌ Failed to fetch tasks: "+err.Error(), keyboards.TaskMenu())
			return
		}
	}

	if len(tasks) == 0 {
		h.send(chatID, "📭 No tasks found.", keyboards.TaskMenu())
		return
	}

	var sb strings.Builder
	sb.WriteString("✅ *Tasks*\n\n")
	for _, t := range tasks {
		dueStr := t.Due.Format("02 Jan 2006")
		sb.WriteString(fmt.Sprintf("• `#%d` *%s* — %s — Due: %s — 💰 Rp%s\n",
			t.ID, t.Subject, t.Status.String(), dueStr, formatPrice(t.Price)))
	}
	sb.WriteString(fmt.Sprintf("\n_Total: %d_", len(tasks)))

	h.send(chatID, sb.String(), keyboards.TaskMenu())
}

func (h *TaskHandler) HandleCreate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateCreateTaskSubject)
	h.send(chatID, "➕ *New Task*\n\nEnter task *subject*:", keyboards.CancelOnly())
}

func (h *TaskHandler) HandleUpdate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateUpdateTaskID)
	h.send(chatID, "✏️ Enter the *Task ID* to update status:", keyboards.CancelOnly())
}

func (h *TaskHandler) HandleDelete(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateDeleteTaskID)
	h.send(chatID, "🗑️ Enter the *Task ID* to delete:", keyboards.CancelOnly())
}

func (h *TaskHandler) HandleMessage(userID, chatID int64, text string, state session.State) bool {
	text = strings.TrimSpace(text)

	switch state {

	case session.StateCreateTaskSubject:
		if text == "" {
			h.send(chatID, "❌ Subject cannot be empty.")
			return true
		}
		h.sessions.SetData(userID, "subject", text)
		h.sessions.Set(userID, session.StateCreateTaskDescription)
		h.send(chatID, "Enter *description* (or skip):", keyboards.SkipCancel())
		return true

	case session.StateCreateTaskDescription:
		h.sessions.SetData(userID, "description", text)
		h.sessions.Set(userID, session.StateCreateTaskPrice)
		h.send(chatID, "Enter *price* (e.g. 50000):", keyboards.CancelOnly())
		return true

	case session.StateCreateTaskPrice:
		price, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid price. Enter a number.")
			return true
		}
		h.sessions.SetData(userID, "price", price)
		h.sessions.Set(userID, session.StateCreateTaskDue)
		h.send(chatID, "Enter *due date* (format: DD-MM-YYYY, e.g. 30-06-2026):", keyboards.CancelOnly())
		return true

	case session.StateCreateTaskDue:
		due, err := time.Parse("02-01-2006", text)
		if err != nil {
			h.send(chatID, "❌ Invalid date. Use format DD-MM-YYYY.")
			return true
		}

		subjectRaw, _ := h.sessions.GetData(userID, "subject")
		descRaw, _ := h.sessions.GetData(userID, "description")
		priceRaw, _ := h.sessions.GetData(userID, "price")

		subject, _ := subjectRaw.(string)
		desc, _ := descRaw.(string)
		price, _ := priceRaw.(uint64)

		req := models.CreateTaskRequest{
			Subject:     subject,
			Description: desc,
			Price:       price,
			Due:         due,
		}

		if err := h.svc.Create(req); err != nil {
			h.send(chatID, "❌ Failed to create task: "+err.Error(), keyboards.TaskMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Task *%s* created!", subject), keyboards.TaskMenu())
		}
		h.sessions.Clear(userID)
		return true

	case session.StateUpdateTaskID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.SetData(userID, "task_id", id)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("Updating Task #%d\nSelect new *status*:", id),
			keyboards.TaskStatusMenu("update_task_status"))
		return true

	case session.StateDeleteTaskID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("⚠️ Delete Task #%d?", id),
			keyboards.ConfirmDelete(fmt.Sprintf("task_confirm_delete_%d", id)))
		return true
	}

	return false
}

func (h *TaskHandler) HandleCallback(userID, chatID int64, data string) bool {
	switch data {
	case "menu_tasks":
		h.HandleMenu(chatID)
		return true
	case "task_list":
		h.HandleList(chatID, nil)
		return true
	case "task_create":
		h.HandleCreate(userID, chatID)
		return true
	case "task_update":
		h.HandleUpdate(userID, chatID)
		return true
	case "task_delete":
		h.HandleDelete(userID, chatID)
		return true
	}

	if strings.HasPrefix(data, "task_list_") {
		statusStr := strings.TrimPrefix(data, "task_list_")
		statusInt, _ := strconv.ParseUint(statusStr, 10, 8)
		tasks, err := h.svc.GetByStatus(models.TaskStatus(statusInt))
		if err != nil {
			h.send(chatID, "❌ "+err.Error(), keyboards.TaskMenu())
			return true
		}
		h.HandleList(chatID, tasks)
		return true
	}

	if strings.HasPrefix(data, "update_task_status_") {
		statusStr := strings.TrimPrefix(data, "update_task_status_")
		statusInt, _ := strconv.ParseUint(statusStr, 10, 8)
		status := models.TaskStatus(statusInt)

		idRaw, _ := h.sessions.GetData(userID, "task_id")
		id, _ := idRaw.(uint64)

		if err := h.svc.Patch(id, models.PatchTaskRequest{Status: &status}); err != nil {
			h.send(chatID, "❌ Update failed: "+err.Error(), keyboards.TaskMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Task #%d status updated to *%s*!", id, status.String()), keyboards.TaskMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	if strings.HasPrefix(data, "task_confirm_delete_") {
		idStr := strings.TrimPrefix(data, "task_confirm_delete_")
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if err := h.svc.Delete(id); err != nil {
			h.send(chatID, "❌ Delete failed: "+err.Error(), keyboards.TaskMenu())
		} else {
			h.send(chatID, fmt.Sprintf("🗑️ Task #%d deleted.", id), keyboards.TaskMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	return false
}

// ========================= TransactionHandler =========================

type TransactionHandler struct {
	svc      *services.TransactionService
	sessions *session.Manager
	bot      *tgbotapi.BotAPI
}

func NewTransactionHandler(
	svc *services.TransactionService,
	sessions *session.Manager,
	bot *tgbotapi.BotAPI,
) *TransactionHandler {
	return &TransactionHandler{svc: svc, sessions: sessions, bot: bot}
}

func (h *TransactionHandler) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	h.bot.Send(msg)
}

func (h *TransactionHandler) HandleMenu(chatID int64) {
	h.send(chatID, "💰 *Transaction Management*\nChoose an action:", keyboards.TransactionMenu())
}

func (h *TransactionHandler) HandleList(chatID int64) {
	txs, err := h.svc.List()
	if err != nil {
		h.send(chatID, "❌ Failed to fetch transactions: "+err.Error(), keyboards.TransactionMenu())
		return
	}

	if len(txs) == 0 {
		h.send(chatID, "📭 No transactions found.", keyboards.TransactionMenu())
		return
	}

	var sb strings.Builder
	sb.WriteString("💰 *Transactions*\n\n")
	for _, tx := range txs {
		orderStr := "—"
		if tx.OrderID != nil {
			orderStr = fmt.Sprintf("#%d", *tx.OrderID)
		}
		sb.WriteString(fmt.Sprintf("• `#%d` %s — %s — Order: %s — 💵 Rp%s\n",
			tx.ID, tx.Status.String(), tx.Method.String(), orderStr, formatPrice(tx.Amount)))
	}
	sb.WriteString(fmt.Sprintf("\n_Total: %d_", len(txs)))

	h.send(chatID, sb.String(), keyboards.TransactionMenu())
}

func (h *TransactionHandler) HandleCreate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateCreateTxOrderID)
	h.send(chatID,
		"➕ *New Transaction*\n\nEnter *Order ID* (or skip):",
		keyboards.SkipCancel(),
	)
}

func (h *TransactionHandler) HandleDelete(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateDeleteTxID)
	h.send(chatID, "🗑️ Enter the *Transaction ID* to delete:", keyboards.CancelOnly())
}

func (h *TransactionHandler) HandleMessage(userID, chatID int64, text string, state session.State) bool {
	text = strings.TrimSpace(text)

	switch state {

	case session.StateCreateTxOrderID:
		var orderID *uint64
		id, err := strconv.ParseUint(text, 10, 64)
		if err == nil {
			orderID = &id
		}
		h.sessions.SetData(userID, "order_id", orderID)
		h.sessions.Set(userID, session.StateCreateTxAmount)
		h.send(chatID, "Enter *amount* (e.g. 150000):", keyboards.CancelOnly())
		return true

	case session.StateCreateTxAmount:
		amount, err := strconv.ParseUint(text, 10, 64)
		if err != nil || amount == 0 {
			h.send(chatID, "❌ Invalid amount. Enter a number greater than 0.")
			return true
		}
		h.sessions.SetData(userID, "amount", amount)
		h.sessions.Set(userID, session.StateCreateTxEvidence)
		h.send(chatID, "Enter *evidence path / note* (or skip):", keyboards.SkipCancel())
		return true

	case session.StateCreateTxEvidence:
		h.sessions.SetData(userID, "evidence", text)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, "Select *payment method*:", keyboards.TransactionMethodMenu())
		return true

	case session.StateDeleteTxID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("⚠️ Delete Transaction #%d?", id),
			keyboards.ConfirmDelete(fmt.Sprintf("tx_confirm_delete_%d", id)))
		return true
	}

	return false
}

func (h *TransactionHandler) HandleCallback(userID, chatID int64, data string) bool {
	switch data {
	case "menu_transactions":
		h.HandleMenu(chatID)
		return true
	case "tx_list":
		h.HandleList(chatID)
		return true
	case "tx_create":
		h.HandleCreate(userID, chatID)
		return true
	case "tx_delete":
		h.HandleDelete(userID, chatID)
		return true
	}

	if strings.HasPrefix(data, "tx_method_") {
		methodStr := strings.TrimPrefix(data, "tx_method_")
		methodInt, _ := strconv.ParseUint(methodStr, 10, 8)
		method := models.TransactionMethod(methodInt)

		orderIDRaw, _ := h.sessions.GetData(userID, "order_id")
		orderID, _ := orderIDRaw.(*uint64)
		amountRaw, _ := h.sessions.GetData(userID, "amount")
		amount, _ := amountRaw.(uint64)
		evidenceRaw, _ := h.sessions.GetData(userID, "evidence")
		evidence, _ := evidenceRaw.(string)

		req := models.CreateTransactionRequest{
			Status:       models.TransactionStatusPending,
			Method:       method,
			Amount:       amount,
			EvidencePath: evidence,
			OrderID:      orderID,
		}

		if err := h.svc.Create(req); err != nil {
			h.send(chatID, "❌ Failed to create transaction: "+err.Error(), keyboards.TransactionMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Transaction created! Method: *%s*, Amount: Rp%s",
				method.String(), formatPrice(amount)), keyboards.TransactionMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	if strings.HasPrefix(data, "tx_confirm_delete_") {
		idStr := strings.TrimPrefix(data, "tx_confirm_delete_")
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if err := h.svc.Delete(id); err != nil {
			h.send(chatID, "❌ Delete failed: "+err.Error(), keyboards.TransactionMenu())
		} else {
			h.send(chatID, fmt.Sprintf("🗑️ Transaction #%d deleted.", id), keyboards.TransactionMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	return false
}