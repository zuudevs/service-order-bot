/**

 filename  : keyboards.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Telegram inline keyboard builders

 copyright Copyright (c) 2026

**/

package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// MainMenu returns the main menu keyboard
func MainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Persons", "menu_persons"),
			tgbotapi.NewInlineKeyboardButtonData("📋 Orders", "menu_orders"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Tasks", "menu_tasks"),
			tgbotapi.NewInlineKeyboardButtonData("💰 Transactions", "menu_transactions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Contacts", "menu_contacts"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Stats", "menu_stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏥 API Health", "menu_health"),
		),
	)
}

// PersonMenu returns the persons sub-menu
func PersonMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 List Persons", "person_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Add Person", "person_create"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Edit Person", "person_edit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 View Person", "person_view"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete Person", "person_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Main Menu", "menu_main"),
		),
	)
}

// OrderMenu returns the orders sub-menu
func OrderMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 All Orders", "order_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Pending", "order_list_0"),
			tgbotapi.NewInlineKeyboardButtonData("✅ Accepted", "order_list_1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏁 Completed", "order_list_4"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Rejected", "order_list_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ New Order", "order_create"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Update Status", "order_update"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete Order", "order_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Main Menu", "menu_main"),
		),
	)
}

// TaskMenu returns the tasks sub-menu
func TaskMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 All Tasks", "task_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Pending", "task_list_0"),
			tgbotapi.NewInlineKeyboardButtonData("🔧 On Progress", "task_list_1"),
			tgbotapi.NewInlineKeyboardButtonData("✅ Complete", "task_list_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ New Task", "task_create"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Update Status", "task_update"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete Task", "task_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Main Menu", "menu_main"),
		),
	)
}

// TransactionMenu returns the transactions sub-menu
func TransactionMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 All Transactions", "tx_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ New Transaction", "tx_create"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete", "tx_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Main Menu", "menu_main"),
		),
	)
}

// ContactMenu returns the contacts sub-menu
func ContactMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 All Contacts", "contact_list"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Add Contact", "contact_create"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Delete", "contact_delete"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Main Menu", "menu_main"),
		),
	)
}

// OrderStatusMenu returns inline buttons to pick an order status
func OrderStatusMenu(callbackPrefix string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Pending", callbackPrefix+"_0"),
			tgbotapi.NewInlineKeyboardButtonData("✅ Accepted", callbackPrefix+"_1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Rejected", callbackPrefix+"_2"),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Revised", callbackPrefix+"_3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏁 Completed", callbackPrefix+"_4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// TaskStatusMenu returns inline buttons to pick a task status
func TaskStatusMenu(callbackPrefix string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏳ Pending", callbackPrefix+"_0"),
			tgbotapi.NewInlineKeyboardButtonData("🔧 On Progress", callbackPrefix+"_1"),
			tgbotapi.NewInlineKeyboardButtonData("✅ Complete", callbackPrefix+"_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// ContactTypeMenu returns inline buttons to pick a contact type
func ContactTypeMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📧 Email", "contact_type_0"),
			tgbotapi.NewInlineKeyboardButtonData("📞 Phone", "contact_type_1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// TransactionMethodMenu returns inline buttons to pick a payment method
func TransactionMethodMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💵 Cash", "tx_method_0"),
			tgbotapi.NewInlineKeyboardButtonData("📱 E-Wallet", "tx_method_1"),
			tgbotapi.NewInlineKeyboardButtonData("🏦 Bank", "tx_method_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// SkipCancel returns a keyboard with Skip and Cancel buttons
func SkipCancel() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏩ Skip", "skip"),
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// CancelOnly returns a keyboard with just a Cancel button
func CancelOnly() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// ConfirmDelete returns a confirm/cancel keyboard for destructive actions
func ConfirmDelete(confirmCallback string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚠️ Yes, Delete", confirmCallback),
			tgbotapi.NewInlineKeyboardButtonData("❎ Cancel", "cancel"),
		),
	)
}

// BackToMain returns a single back-to-main-menu button
func BackToMain() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Main Menu", "menu_main"),
		),
	)
}