/**

 filename  : contact_handler.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Telegram handlers for contacts and stats

 copyright Copyright (c) 2026

**/

package handlers

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zuudevs/service-order-bot/internal/keyboards"
	"github.com/zuudevs/service-order-bot/internal/models"
	"github.com/zuudevs/service-order-bot/internal/session"
	"github.com/zuudevs/service-order-bot/internal/services"
)

// ========================= ContactHandler =========================

type ContactHandler struct {
	svc      *services.ContactService
	sessions *session.Manager
	bot      *tgbotapi.BotAPI
}

func NewContactHandler(
	svc *services.ContactService,
	sessions *session.Manager,
	bot *tgbotapi.BotAPI,
) *ContactHandler {
	return &ContactHandler{svc: svc, sessions: sessions, bot: bot}
}

func (h *ContactHandler) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	h.bot.Send(msg)
}

func (h *ContactHandler) HandleMenu(chatID int64) {
	h.send(chatID, "📞 *Contact Management*\nChoose an action:", keyboards.ContactMenu())
}

func (h *ContactHandler) HandleList(chatID int64) {
	contacts, err := h.svc.List()
	if err != nil {
		h.send(chatID, "❌ Failed to fetch contacts: "+err.Error(), keyboards.ContactMenu())
		return
	}

	if len(contacts) == 0 {
		h.send(chatID, "📭 No contacts found.", keyboards.ContactMenu())
		return
	}

	var sb strings.Builder
	sb.WriteString("📞 *Contacts*\n\n")
	for _, c := range contacts {
		mainTag := ""
		if c.IsMain {
			mainTag = " ⭐"
		}
		sb.WriteString(fmt.Sprintf("• `#%d` [Person #%d] %s: `%s`%s\n",
			c.ID, c.PersonID, c.ContactType.String(), c.Value, mainTag))
	}
	sb.WriteString(fmt.Sprintf("\n_Total: %d_", len(contacts)))

	h.send(chatID, sb.String(), keyboards.ContactMenu())
}

func (h *ContactHandler) HandleCreate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateCreateContactPersonID)
	h.send(chatID, "➕ *Add Contact*\n\nEnter the *Person ID*:", keyboards.CancelOnly())
}

func (h *ContactHandler) HandleDelete(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateDeleteContactID)
	h.send(chatID, "🗑️ Enter the *Contact ID* to delete:", keyboards.CancelOnly())
}

func (h *ContactHandler) HandleMessage(userID, chatID int64, text string, state session.State) bool {
	text = strings.TrimSpace(text)

	switch state {

	case session.StateCreateContactPersonID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid person ID.")
			return true
		}
		h.sessions.SetData(userID, "person_id", id)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, "Select *contact type*:", keyboards.ContactTypeMenu())
		return true

	case session.StateCreateContactValue:
		h.sessions.SetData(userID, "value", text)
		h.finishCreate(userID, chatID)
		return true

	case session.StateDeleteContactID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("⚠️ Delete Contact #%d?", id),
			keyboards.ConfirmDelete(fmt.Sprintf("contact_confirm_delete_%d", id)))
		return true
	}

	return false
}

func (h *ContactHandler) HandleCallback(userID, chatID int64, data string) bool {
	switch data {
	case "menu_contacts":
		h.HandleMenu(chatID)
		return true
	case "contact_list":
		h.HandleList(chatID)
		return true
	case "contact_create":
		h.HandleCreate(userID, chatID)
		return true
	case "contact_delete":
		h.HandleDelete(userID, chatID)
		return true
	}

	if strings.HasPrefix(data, "contact_type_") {
		typeStr := strings.TrimPrefix(data, "contact_type_")
		typeInt, _ := strconv.ParseUint(typeStr, 10, 8)
		h.sessions.SetData(userID, "contact_type", models.ContactType(typeInt))
		h.sessions.Set(userID, session.StateCreateContactValue)

		typeLabel := "email"
		if typeInt == 1 {
			typeLabel = "phone number"
		}
		h.send(chatID, fmt.Sprintf("Enter the *%s*:", typeLabel), keyboards.CancelOnly())
		return true
	}

	if strings.HasPrefix(data, "contact_confirm_delete_") {
		idStr := strings.TrimPrefix(data, "contact_confirm_delete_")
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if err := h.svc.Delete(id); err != nil {
			h.send(chatID, "❌ Delete failed: "+err.Error(), keyboards.ContactMenu())
		} else {
			h.send(chatID, fmt.Sprintf("🗑️ Contact #%d deleted.", id), keyboards.ContactMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	return false
}

func (h *ContactHandler) finishCreate(userID, chatID int64) {
	personIDRaw, _ := h.sessions.GetData(userID, "person_id")
	contactTypeRaw, _ := h.sessions.GetData(userID, "contact_type")
	valueRaw, _ := h.sessions.GetData(userID, "value")

	personID, _ := personIDRaw.(uint64)
	contactType, _ := contactTypeRaw.(models.ContactType)
	value, _ := valueRaw.(string)

	req := models.CreateContactRequest{
		Value:       value,
		ContactType: contactType,
		IsMain:      false,
		PersonID:    personID,
	}

	if err := h.svc.Create(req); err != nil {
		h.send(chatID, "❌ Failed to add contact: "+err.Error(), keyboards.ContactMenu())
	} else {
		h.send(chatID, fmt.Sprintf("✅ Contact `%s` added to Person #%d!", value, personID), keyboards.ContactMenu())
	}
	h.sessions.Clear(userID)
}

// ========================= StatsHandler =========================

type StatsHandler struct {
	personSvc      *services.PersonService
	orderSvc       *services.OrderService
	taskSvc        *services.TaskService
	transactionSvc *services.TransactionService
	bot            *tgbotapi.BotAPI
}

func NewStatsHandler(
	personSvc *services.PersonService,
	orderSvc *services.OrderService,
	taskSvc *services.TaskService,
	transactionSvc *services.TransactionService,
	bot *tgbotapi.BotAPI,
) *StatsHandler {
	return &StatsHandler{
		personSvc:      personSvc,
		orderSvc:       orderSvc,
		taskSvc:        taskSvc,
		transactionSvc: transactionSvc,
		bot:            bot,
	}
}

func (h *StatsHandler) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	h.bot.Send(msg)
}

func (h *StatsHandler) HandleStats(chatID int64) {
	persons, _ := h.personSvc.List()
	orders, _ := h.orderSvc.List()
	tasks, _ := h.taskSvc.List()
	txs, _ := h.transactionSvc.List()

	// Order status breakdown
	orderCounts := map[models.OrderStatus]int{}
	var totalRevenue uint64
	for _, o := range orders {
		orderCounts[o.Status]++
		totalRevenue += o.TotalPrice
	}

	// Task status breakdown
	taskCounts := map[models.TaskStatus]int{}
	for _, t := range tasks {
		taskCounts[t.Status]++
	}

	var sb strings.Builder
	sb.WriteString("📊 *Dashboard Stats*\n\n")

	sb.WriteString(fmt.Sprintf("👥 *Persons:* %d\n", len(persons)))
	sb.WriteString(fmt.Sprintf("📋 *Orders:* %d total\n", len(orders)))
	sb.WriteString(fmt.Sprintf("   • ⏳ Pending: %d\n", orderCounts[models.OrderStatusPending]))
	sb.WriteString(fmt.Sprintf("   • ✅ Accepted: %d\n", orderCounts[models.OrderStatusAccepted]))
	sb.WriteString(fmt.Sprintf("   • 🔄 Revised: %d\n", orderCounts[models.OrderStatusRevised]))
	sb.WriteString(fmt.Sprintf("   • 🏁 Completed: %d\n", orderCounts[models.OrderStatusCompleted]))
	sb.WriteString(fmt.Sprintf("   • ❌ Rejected: %d\n", orderCounts[models.OrderStatusRejected]))
	sb.WriteString(fmt.Sprintf("💰 *Total Revenue:* Rp%s\n\n", formatPrice(totalRevenue)))

	sb.WriteString(fmt.Sprintf("✅ *Tasks:* %d total\n", len(tasks)))
	sb.WriteString(fmt.Sprintf("   • ⏳ Pending: %d\n", taskCounts[models.TaskStatusPending]))
	sb.WriteString(fmt.Sprintf("   • 🔧 On Progress: %d\n", taskCounts[models.TaskStatusOnProgress]))
	sb.WriteString(fmt.Sprintf("   • ✅ Complete: %d\n", taskCounts[models.TaskStatusComplete]))

	sb.WriteString(fmt.Sprintf("\n💰 *Transactions:* %d total\n", len(txs)))

	h.send(chatID, sb.String(), keyboards.BackToMain())
}