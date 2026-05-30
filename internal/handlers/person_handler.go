/**

 filename  : person_handler.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Telegram handler for person-related flows

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

type PersonHandler struct {
	svc      *services.PersonService
	sessions *session.Manager
	bot      *tgbotapi.BotAPI
}

func NewPersonHandler(
	svc *services.PersonService,
	sessions *session.Manager,
	bot *tgbotapi.BotAPI,
) *PersonHandler {
	return &PersonHandler{svc: svc, sessions: sessions, bot: bot}
}

// send is a helper to send a message
func (h *PersonHandler) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	h.bot.Send(msg)
}

// HandleMenu sends the person sub-menu
func (h *PersonHandler) HandleMenu(chatID int64) {
	h.send(chatID, "👥 *Person Management*\nChoose an action:", keyboards.PersonMenu())
}

// HandleList lists all persons
func (h *PersonHandler) HandleList(chatID int64) {
	persons, err := h.svc.List()
	if err != nil {
		h.send(chatID, "❌ Failed to fetch persons: "+err.Error(), keyboards.BackToMain())
		return
	}

	if len(persons) == 0 {
		h.send(chatID, "📭 No persons found.", keyboards.PersonMenu())
		return
	}

	var sb strings.Builder
	sb.WriteString("👥 *Persons List*\n\n")
	for _, p := range persons {
		sb.WriteString(fmt.Sprintf("• `#%d` — *%s*\n", p.ID, p.FullName()))
	}
	sb.WriteString(fmt.Sprintf("\n_Total: %d_", len(persons)))

	h.send(chatID, sb.String(), keyboards.PersonMenu())
}

// HandleView shows details for a single person
func (h *PersonHandler) HandleView(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateViewPersonID)
	h.send(chatID, "🔍 Enter the *Person ID* to view:", keyboards.CancelOnly())
}

// HandleViewByID fetches and displays a person
func (h *PersonHandler) HandleViewByID(chatID int64, idStr string) {
	id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil {
		h.send(chatID, "❌ Invalid ID. Please enter a number.", keyboards.CancelOnly())
		return
	}

	person, err := h.svc.GetByID(id)
	if err != nil {
		h.send(chatID, "❌ Person not found: "+err.Error(), keyboards.PersonMenu())
		return
	}

	contacts, _ := h.svc.GetContacts(id)
	orders, _ := h.svc.GetOrders(id)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👤 *Person #%d*\n\n", person.ID))
	sb.WriteString(fmt.Sprintf("First Name: `%s`\n", person.FirstName))
	if person.MiddleName != nil {
		sb.WriteString(fmt.Sprintf("Middle Name: `%s`\n", *person.MiddleName))
	}
	if person.LastName != nil {
		sb.WriteString(fmt.Sprintf("Last Name: `%s`\n", *person.LastName))
	}
	sb.WriteString(fmt.Sprintf("Created: `%s`\n", person.CreatedAt.Format("02 Jan 2006")))

	if len(contacts) > 0 {
		sb.WriteString("\n📞 *Contacts:*\n")
		for _, c := range contacts {
			mainTag := ""
			if c.IsMain {
				mainTag = " ⭐"
			}
			sb.WriteString(fmt.Sprintf("  • %s: `%s`%s\n", c.ContactType.String(), c.Value, mainTag))
		}
	}

	if len(orders) > 0 {
		sb.WriteString(fmt.Sprintf("\n📋 *Orders: %d total*\n", len(orders)))
	}

	h.send(chatID, sb.String(), keyboards.PersonMenu())
}

// HandleCreate starts the create person flow
func (h *PersonHandler) HandleCreate(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateCreatePersonFirstName)
	h.send(chatID, "➕ *Create New Person*\n\nEnter *first name*:", keyboards.CancelOnly())
}

// HandleEdit starts the edit person flow
func (h *PersonHandler) HandleEdit(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateEditPersonID)
	h.send(chatID, "✏️ Enter the *Person ID* to edit:", keyboards.CancelOnly())
}

// HandleDelete starts the delete person flow
func (h *PersonHandler) HandleDelete(userID, chatID int64) {
	h.sessions.Clear(userID)
	h.sessions.Set(userID, session.StateDeletePersonID)
	h.send(chatID, "🗑️ Enter the *Person ID* to delete:", keyboards.CancelOnly())
}

// HandleMessage routes text input based on the current session state
func (h *PersonHandler) HandleMessage(userID, chatID int64, text string, state session.State) bool {
	switch state {

	case session.StateViewPersonID:
		h.HandleViewByID(chatID, text)
		h.sessions.Clear(userID)
		return true

	case session.StateCreatePersonFirstName:
		text = strings.TrimSpace(text)
		if text == "" {
			h.send(chatID, "❌ First name cannot be empty.")
			return true
		}
		h.sessions.SetData(userID, "firstname", text)
		h.sessions.Set(userID, session.StateCreatePersonMiddleName)
		h.send(chatID, "Enter *middle name* (or skip):", keyboards.SkipCancel())
		return true

	case session.StateCreatePersonMiddleName:
		text = strings.TrimSpace(text)
		h.sessions.SetData(userID, "middlename", text)
		h.sessions.Set(userID, session.StateCreatePersonLastName)
		h.send(chatID, "Enter *last name* (or skip):", keyboards.SkipCancel())
		return true

	case session.StateCreatePersonLastName:
		text = strings.TrimSpace(text)
		h.sessions.SetData(userID, "lastname", text)
		h.finishCreate(userID, chatID)
		return true

	case session.StateDeletePersonID:
		id, err := strconv.ParseUint(strings.TrimSpace(text), 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.SetData(userID, "delete_id", id)
		h.sessions.Set(userID, session.StateIdle)
		h.send(chatID, fmt.Sprintf("⚠️ Delete Person #%d? This cannot be undone.", id),
			keyboards.ConfirmDelete(fmt.Sprintf("person_confirm_delete_%d", id)))
		return true

	case session.StateEditPersonID:
		id, err := strconv.ParseUint(strings.TrimSpace(text), 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.")
			return true
		}
		h.sessions.SetData(userID, "edit_id", id)
		h.sessions.Set(userID, session.StateEditPersonField)
		h.send(chatID,
			fmt.Sprintf("Editing Person #%d\nWhich field to update?", id),
			tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("First Name", "edit_person_field_firstname"),
					tgbotapi.NewInlineKeyboardButtonData("Middle Name", "edit_person_field_middlename"),
					tgbotapi.NewInlineKeyboardButtonData("Last Name", "edit_person_field_lastname"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
				),
			),
		)
		return true

	case session.StateEditPersonValue:
		field, _ := h.sessions.GetData(userID, "edit_field")
		idRaw, _ := h.sessions.GetData(userID, "edit_id")
		id, _ := idRaw.(uint64)

		req := models.PatchPersonRequest{}
		value := strings.TrimSpace(text)

		switch field {
		case "firstname":
			req.FirstName = &value
		case "middlename":
			req.MiddleName = &value
		case "lastname":
			req.LastName = &value
		}

		if err := h.svc.Patch(id, req); err != nil {
			h.send(chatID, "❌ Update failed: "+err.Error(), keyboards.PersonMenu())
		} else {
			h.send(chatID, fmt.Sprintf("✅ Person #%d updated successfully!", id), keyboards.PersonMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	return false
}

// HandleCallback processes callback queries related to persons
func (h *PersonHandler) HandleCallback(userID, chatID int64, data string) bool {
	switch data {
	case "menu_persons":
		h.HandleMenu(chatID)
		return true
	case "person_list":
		h.HandleList(chatID)
		return true
	case "person_create":
		h.HandleCreate(userID, chatID)
		return true
	case "person_edit":
		h.HandleEdit(userID, chatID)
		return true
	case "person_delete":
		h.HandleDelete(userID, chatID)
		return true
	case "person_view":
		h.HandleView(userID, chatID)
		return true
	}

	// Edit field selection
	if strings.HasPrefix(data, "edit_person_field_") {
		field := strings.TrimPrefix(data, "edit_person_field_")
		h.sessions.SetData(userID, "edit_field", field)
		h.sessions.Set(userID, session.StateEditPersonValue)
		h.send(chatID, fmt.Sprintf("Enter new value for *%s*:", field), keyboards.CancelOnly())
		return true
	}

	// Confirm delete
	if strings.HasPrefix(data, "person_confirm_delete_") {
		idStr := strings.TrimPrefix(data, "person_confirm_delete_")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			h.send(chatID, "❌ Invalid ID.", keyboards.PersonMenu())
			return true
		}
		if err := h.svc.Delete(id); err != nil {
			h.send(chatID, "❌ Delete failed: "+err.Error(), keyboards.PersonMenu())
		} else {
			h.send(chatID, fmt.Sprintf("🗑️ Person #%d deleted.", id), keyboards.PersonMenu())
		}
		h.sessions.Clear(userID)
		return true
	}

	return false
}

// finishCreate completes the create person flow
func (h *PersonHandler) finishCreate(userID, chatID int64) {
	firstnameRaw, _ := h.sessions.GetData(userID, "firstname")
	middlenameRaw, _ := h.sessions.GetData(userID, "middlename")
	lastnameRaw, _ := h.sessions.GetData(userID, "lastname")

	firstname, _ := firstnameRaw.(string)

	req := models.CreatePersonRequest{FirstName: firstname}

	if mn, ok := middlenameRaw.(string); ok && mn != "" {
		req.MiddleName = &mn
	}
	if ln, ok := lastnameRaw.(string); ok && ln != "" {
		req.LastName = &ln
	}

	if err := h.svc.Create(req); err != nil {
		h.send(chatID, "❌ Failed to create person: "+err.Error(), keyboards.PersonMenu())
	} else {
		h.send(chatID, fmt.Sprintf("✅ Person *%s* created successfully!", firstname), keyboards.PersonMenu())
	}

	h.sessions.Clear(userID)
}