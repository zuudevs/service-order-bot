/**

 filename  : main.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Entry point for the Service Order Telegram Bot

 copyright Copyright (c) 2026

**/

package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/zuudevs/service-order-bot/internal/bot"
	"github.com/zuudevs/service-order-bot/internal/client"
)

func main() {
	// Load .env if present (ignored in production where env vars are set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("[main] No .env file found, using environment variables")
	}

	// Required env vars
	token := mustEnv("TELEGRAM_BOT_TOKEN")
	apiURL := mustEnv("SERVICE_ORDER_API_URL")
	apiToken := mustEnv("SERVICE_ORDER_API_TOKEN")

	// Optional
	debugStr := os.Getenv("BOT_DEBUG")
	debug, _ := strconv.ParseBool(debugStr)

	// Create API client
	apiClient := client.New(apiURL, apiToken)

	// Check API health on startup
	log.Printf("[main] Connecting to API at %s ...", apiURL)
	if err := apiClient.Health(); err != nil {
		log.Printf("[main] WARNING: API health check failed: %v", err)
		log.Println("[main] Continuing anyway — API may come online later")
	} else {
		log.Println("[main] API is healthy ✅")
	}

	// Create and start bot
	b, err := bot.New(token, apiClient, debug)
	if err != nil {
		log.Fatalf("[main] Failed to initialize bot: %v", err)
	}

	log.Println("[main] Bot started. Press Ctrl+C to stop.")
	b.Run()
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[main] Required env var %q is not set", key)
	}
	return v
}