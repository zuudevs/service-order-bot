/**

 filename  : bot.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Telegram bot dispatcher — routes updates to handlers

 copyright Copyright (c) 2026

**/

package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zuudevs/service-order-bot/internal/client"
	"github.com/zuudevs/service-order-bot/internal/handlers"
	"github.com/zuudevs/service-order-bot/internal/keyboards"
	"github.com/zuudevs/service-order-bot/internal/middlewares"
	"github.com/zuudevs/service-order-bot/internal/session"
	"github.com/zuudevs/service-order-bot/internal/services"
)

// Bot wraps the Telegram bot and all its handlers
type Bot struct {
	api            *tgbotapi.BotAPI
	auth           *middlewares.AuthMiddleware
	sessions       *session.Manager
	apiClient      *client.APIClient

	personHandler      *handlers.PersonHandler
	orderHandler       *handlers.OrderHandler
	taskHandler        *handlers.TaskHandler
	transactionHandler *handlers.TransactionHandler
	contactHandler     *handlers.ContactHandler
	statsHandler       *handlers.StatsHandler
}

// New creates and wires up all bot components
func New(token string, apiClient *client.APIClient, debug bool) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create bot api: %w", err)
	}

	api.Debug = debug
	log.Printf("[bot] Authorized as @%s", api.Self.UserName)

	sessions := session.NewManager(30 * 60 * 1e9) // 30 min TTL

	personSvc := services.NewPersonService(apiClient)
	orderSvc := services.NewOrderService(apiClient)
	taskSvc := services.NewTaskService(apiClient)
	txSvc := services.NewTransactionService(apiClient)
	contactSvc := services.NewContactService(apiClient)

	return &Bot{
		api:            api,
		auth:           middlewares.NewAuthMiddleware(),
		sessions:       sessions,
		apiClient:      apiClient,

		personHandler:      handlers.NewPersonHandler(personSvc, sessions, api),
		orderHandler:       handlers.NewOrderHandler(orderSvc, sessions, api),
		taskHandler:        handlers.NewTaskHandler(taskSvc, sessions, api),
		transactionHandler: handlers.NewTransactionHandler(txSvc, sessions, api),
		contactHandler:     handlers.NewContactHandler(contactSvc, sessions, api),
		statsHandler:       handlers.NewStatsHandler(personSvc, orderSvc, taskSvc, txSvc, api),
	}, nil
}

// Run starts the polling loop
func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	log.Println("[bot] Listening for updates...")

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

// ========================= Message Dispatcher =========================

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID
	chatID := msg.Chat.ID

	if !b.auth.IsAllowed(userID) {
		b.send(chatID, "⛔ You are not authorized to use this bot.")
		return
	}

	text := strings.TrimSpace(msg.Text)

	// Commands
	if msg.IsCommand() {
		b.handleCommand(userID, chatID, msg.Command())
		return
	}

	// Route to active session state
	sess := b.sessions.Get(userID)
	state := sess.State

	if state == session.StateIdle {
		b.sendMainMenu(chatID)
		return
	}

	// Try each handler
	handled := b.personHandler.HandleMessage(userID, chatID, text, state) ||
		b.orderHandler.HandleMessage(userID, chatID, text, state) ||
		b.taskHandler.HandleMessage(userID, chatID, text, state) ||
		b.transactionHandler.HandleMessage(userID, chatID, text, state) ||
		b.contactHandler.HandleMessage(userID, chatID, text, state)

	if !handled {
		b.send(chatID, "❓ I didn't understand that. Use the menu below.", keyboards.BackToMain())
	}
}

func (b *Bot) handleCommand(userID int64, chatID int64, cmd string) {
	switch cmd {
	case "start", "menu":
		b.sendMainMenu(chatID)
	case "cancel":
		b.sessions.Clear(userID)
		b.send(chatID, "❎ Cancelled.", keyboards.BackToMain())
	case "health":
		if err := b.apiClient.Health(); err != nil {
			b.send(chatID, "🔴 API is *DOWN*: "+err.Error())
		} else {
			b.send(chatID, "🟢 API is *UP* and healthy!")
		}
	case "stats":
		b.statsHandler.HandleStats(chatID)
	case "help":
		b.sendHelp(chatID)
	default:
		b.send(chatID, "❓ Unknown command. Use /menu to start.")
	}
}

// ========================= Callback Dispatcher =========================

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	userID := cb.From.ID
	chatID := cb.Message.Chat.ID
	data := cb.Data

	// Acknowledge the callback
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))

	if !b.auth.IsAllowed(userID) {
		return
	}

	// Global callbacks
	switch data {
	case "menu_main":
		b.sessions.Clear(userID)
		b.sendMainMenu(chatID)
		return

	case "cancel":
		b.sessions.Clear(userID)
		b.send(chatID, "❎ Cancelled.", keyboards.MainMenu())
		return

	case "skip":
		sess := b.sessions.Get(userID)
		b.handleSkip(userID, chatID, sess.State)
		return

	case "menu_health":
		if err := b.apiClient.Health(); err != nil {
			b.send(chatID, "🔴 API is *DOWN*: "+err.Error(), keyboards.BackToMain())
		} else {
			b.send(chatID, "🟢 API is *UP* and healthy!", keyboards.BackToMain())
		}
		return

	case "menu_stats":
		b.statsHandler.HandleStats(chatID)
		return
	}

	// Route to domain handlers
	_ = b.personHandler.HandleCallback(userID, chatID, data) ||
		b.orderHandler.HandleCallback(userID, chatID, data) ||
		b.taskHandler.HandleCallback(userID, chatID, data) ||
		b.transactionHandler.HandleCallback(userID, chatID, data) ||
		b.contactHandler.HandleCallback(userID, chatID, data)
}

// handleSkip handles the "skip" button — advances the state machine
func (b *Bot) handleSkip(userID, chatID int64, state session.State) {
	switch state {
	case session.StateCreatePersonMiddleName:
		b.sessions.SetData(userID, "middlename", "")
		b.sessions.Set(userID, session.StateCreatePersonLastName)
		b.send(chatID, "Enter *last name* (or skip):", keyboards.SkipCancel())

	case session.StateCreatePersonLastName:
		b.sessions.SetData(userID, "lastname", "")
		// Trigger finishCreate by sending empty — person handler doesn't export finishCreate,
		// so we simulate the message flow
		b.personHandler.HandleMessage(userID, chatID, "", state)

	case session.StateCreateOrderPersonID:
		b.sessions.SetData(userID, "person_id", (*uint64)(nil))
		b.sessions.Set(userID, session.StateIdle)
		b.send(chatID, "Select order *status*:", keyboards.OrderStatusMenu("new_order_status"))

	case session.StateCreateTaskDescription:
		b.taskHandler.HandleMessage(userID, chatID, "", state)

	case session.StateCreateTxOrderID:
		b.sessions.SetData(userID, "order_id", (*uint64)(nil))
		b.sessions.Set(userID, session.StateCreateTxAmount)
		b.send(chatID, "Enter *amount* (e.g. 150000):", keyboards.CancelOnly())

	case session.StateCreateTxEvidence:
		b.transactionHandler.HandleMessage(userID, chatID, "", state)

	default:
		b.sessions.Clear(userID)
		b.sendMainMenu(chatID)
	}
}

// ========================= Helpers =========================

func (b *Bot) send(chatID int64, text string, markup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("[bot] send error: %v", err)
	}
}

func (b *Bot) sendMainMenu(chatID int64) {
	b.send(chatID, "🏠 *Service Order Bot*\nWelcome! Choose a module:", keyboards.MainMenu())
}

func (b *Bot) sendHelp(chatID int64) {
	help := `*Service Order Bot — Help*

*Commands:*
/start or /menu — Show main menu
/stats — Dashboard statistics
/health — Check API status
/cancel — Cancel current action
/help — Show this message

*Features:*
👥 *Persons* — Manage customers/service providers
📋 *Orders* — Create and track service orders
✅ *Tasks* — Manage work tasks
💰 *Transactions* — Record payments
📞 *Contacts* — Phone numbers and emails

_Use the inline buttons to navigate._`

	b.send(chatID, help, keyboards.BackToMain())
}