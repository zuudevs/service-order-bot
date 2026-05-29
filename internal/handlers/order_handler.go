/**

 filename  : order_handler.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Telegram handler for order-related flows

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

type OrderHandler struct {
	svc      *services.OrderService
	sessions *session.Manager
	bot      *tgbotapi.BotAPI
}

func NewOrderHandler(
	svc *services.OrderService,
	sessions *session.Manager,
	bot *tgbotapi.BotAPI,
) *OrderHandler {
	return &OrderHandler{svc: svc, sessions: sessions, bot: bot}
}

func (h *OrderHandler) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	h.bot.Send(msg)
}

func (h *OrderHandler) HandleMenu(chatID int64) {
	h.send(chatID, "📋 *Order Management*\nChoose an action:", keyboards.OrderMenu())
}

func (h *OrderHandler) HandleList(chatID int64, orders []models.Order) {
	if orders == nil {
		var err error
		orders, err = h.svc.List()
		if err != nil {
			h.send(chatID, "❌ Failed to fetch orders: "+err.Error(), keyboards.OrderMenu())
			return
		}
	}

	if len(orders) == 0 {
		h.send(chatID, "📭 No orders found.", keyboards.OrderMenu())
		return
	}

	var sb strings.Builder
	sb.WriteString("📋 *Orders*\n\n")
	for _, o := range orders {
		personStr := "—"
		if o.PersonID != nil {
			personStr = fmt.Sprintf("#%d", *o.PersonID)
		}
		sb.WriteString(fmt.Sprintf(
			"• `#%d` %s — Person: %s — 💰 Rp%s\n",
			o.ID,
			o.Status.Emoji(),
			personStr,
			formatPrice(o.TotalPrice),
		))
	}
	sb.WriteString(fmt.Sprintf("\n_Total: %d_", len(orders)))

	h.send(chatID, sb.String(), keyboards.OrderMenu())
}

func (h *OrderHandler) HandleCreate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateCreateOrderPersonID)
	h.send(chatID,
		"➕ *New Order*\n\nEnter the *Person ID* for this order (or skip to create without person):",
		keyboards.SkipCancel(),
	)
}

func (h *OrderHandler) HandleUpdate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateUpdateOrderID)
	h.send(chatID, "✏️ Enter the *Order ID* to update:", keyboards.CancelOnly())
}

func (h *OrderHandler) HandleDelete(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateDeleteOrderID)
	h.send(chatID, "🗑️ Enter the *Order ID* to delete:", keyboards.CancelOnly())
}

// HandleMessage routes text input based on the current session state
func (h *OrderHandler) HandleMessage(userID, chatID int64, text string, state session.State) bool {
	text = strings.TrimSpace(text)

	switch state {

	case session.StateCreateOrderPersonID:
		var personID *uint64
		id, err := strconv.ParseUint(text, 10, 64)
		if err == nil {
			personID = &id
		}
		h.sessions.SetData(userID, "person_id", personID)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, "Select order *status*:", keyboards.OrderStatusMenu("new_order_status"))
		return true

	case session.StateUpdateOrderID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.SetData(userID, "order_id", id)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("Updating Order #%d\nSelect new *status*:", id),
			keyboards.OrderStatusMenu("update_order_status"))
		return true

	case session.StateUpdateOrderPrice:
		price, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid price. Enter a number.")
			return true
		}
		idRaw, _ := h.sessions.GetData(userID, "order_id")
		id, _ := idRaw.(uint64)

		if err := h.svc.Patch(id, models.PatchOrderRequest{TotalPrice: &price}); err != nil {
			h.send(chatID, "❌ Update failed: "+err.Error(), keyboards.OrderMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Order #%d total price updated to Rp%s!", id, formatPrice(price)), keyboards.OrderMenu())
		}
		h.sessions.Clear(userID)
		return true

	case session.StateDeleteOrderID:
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.SetData(userID, "delete_id", id)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("⚠️ Delete Order #%d?", id),
			keyboards.ConfirmDelete(fmt.Sprintf("order_confirm_delete_%d", id)))
		return true
	}

	return false
}

// HandleCallback processes callback queries related to orders
func (h *OrderHandler) HandleCallback(userID, chatID int64, data string) bool {
	switch data {
	case "menu_orders":
		h.HandleMenu(chatID)
		return true
	case "order_list":
		h.HandleList(chatID, nil)
		return true
	case "order_create":
		h.HandleCreate(userID, chatID)
		return true
	case "order_update":
		h.HandleUpdate(userID, chatID)
		return true
	case "order_delete":
		h.HandleDelete(userID, chatID)
		return true
	}

	// Filter by status
	if strings.HasPrefix(data, "order_list_") {
		statusStr := strings.TrimPrefix(data, "order_list_")
		statusInt, err := strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			return false
		}
		orders, err := h.svc.GetByStatus(models.OrderStatus(statusInt))
		if err != nil {
			h.send(chatID, "❌ "+err.Error(), keyboards.OrderMenu())
			return true
		}
		h.HandleList(chatID, orders)
		return true
	}

	// New order status selection
	if strings.HasPrefix(data, "new_order_status_") {
		statusStr := strings.TrimPrefix(data, "new_order_status_")
		statusInt, _ := strconv.ParseUint(statusStr, 10, 8)
		status := models.OrderStatus(statusInt)

		personIDRaw, _ := h.sessions.GetData(userID, "person_id")
		personID, _ := personIDRaw.(*uint64)

		req := models.CreateOrderRequest{
			Status:   status,
			PersonID: personID,
		}

		if err := h.svc.Create(req); err != nil {
			h.send(chatID, "❌ Failed to create order: "+err.Error(), keyboards.OrderMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Order created with status *%s*!", status.String()), keyboards.OrderMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	// Update order status
	if strings.HasPrefix(data, "update_order_status_") {
		statusStr := strings.TrimPrefix(data, "update_order_status_")
		statusInt, _ := strconv.ParseUint(statusStr, 10, 8)
		status := models.OrderStatus(statusInt)

		idRaw, _ := h.sessions.GetData(userID, "order_id")
		id, _ := idRaw.(uint64)

		if err := h.svc.Patch(id, models.PatchOrderRequest{Status: &status}); err != nil {
			h.send(chatID, "❌ Update failed: "+err.Error(), keyboards.OrderMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Order #%d status updated to *%s*!", id, status.String()), keyboards.OrderMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	// Confirm delete
	if strings.HasPrefix(data, "order_confirm_delete_") {
		idStr := strings.TrimPrefix(data, "order_confirm_delete_")
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if err := h.svc.Delete(id); err != nil {
			h.send(chatID, "❌ Delete failed: "+err.Error(), keyboards.OrderMenu())
		} else {
			h.send(chatID, fmt.Sprintf("🗑️ Order #%d deleted.", id), keyboards.OrderMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	return false
}

// formatPrice formats a price with thousand separators
func formatPrice(price uint64) string {
	s := strconv.FormatUint(price, 10)
	n := len(s)
	var result strings.Builder
	for i, ch := range s {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(ch)
	}
	return result.String()
}