/**

 filename  : auth.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Authorization middleware - whitelist Telegram user IDs

 copyright Copyright (c) 2026

**/

package middlewares

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// AuthMiddleware checks if a Telegram user is allowed to use the bot
type AuthMiddleware struct {
	allowedIDs map[int64]bool
	allowAll   bool
}

// NewAuthMiddleware creates an AuthMiddleware from ALLOWED_TELEGRAM_USER_IDS env var.
// If the env var is empty or "*", all users are allowed.
func NewAuthMiddleware() *AuthMiddleware {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_TELEGRAM_USER_IDS"))

	if raw == "" || raw == "*" {
		log.Println("[auth] WARNING: no user whitelist set — all users can access this bot")
		return &AuthMiddleware{allowAll: true}
	}

	allowed := make(map[int64]bool)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			log.Printf("[auth] invalid user ID in whitelist: %q", part)
			continue
		}
		allowed[id] = true
	}

	log.Printf("[auth] whitelist loaded: %d user(s)", len(allowed))
	return &AuthMiddleware{allowedIDs: allowed}
}

// IsAllowed returns true if the given Telegram user ID is permitted
func (a *AuthMiddleware) IsAllowed(userID int64) bool {
	if a.allowAll {
		return true
	}
	return a.allowedIDs[userID]
}