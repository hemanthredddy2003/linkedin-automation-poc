package main

import (
	"log"
	"time"

	"linkedin-automation/auth"
	"linkedin-automation/config"
	"linkedin-automation/messaging"
	"linkedin-automation/search"
	"linkedin-automation/stealth"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.LoadConfig()
	log.Printf("Starting LinkedIn Automation - Daily Limit: %d connections", cfg.DailyConnectionLimit)

	// Setup browser with stealth mode
	l := launcher.New().
		Headless(false).
		Devtools(false)

	browser := rod.New().
		ControlURL(l.MustLaunch()).
		MustConnect()
	defer browser.MustClose()

	// Create page with stealth settings
	page := stealth.NewStealthPage(browser)

	// Navigate to LinkedIn
	log.Println("Navigating to LinkedIn...")
	page.MustNavigate("https://www.linkedin.com/login")
	stealth.RandomDelay(2000, 3000)

	// Authenticate
	log.Println("Attempting login...")
	if err := auth.Login(page, cfg.LinkedInEmail, cfg.LinkedInPassword); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Println("✓ Login successful")

	// Wait for homepage to load
	stealth.RandomDelay(3000, 5000)

	// Search for profiles
	log.Println("Searching for profiles...")
	jobTitle := "Software Engineer"
	profiles := search.SearchProfiles(page, jobTitle, 2) // Only get 2 profiles
	log.Printf("Found %d profiles", len(profiles))

	// Send connection requests (limit to 2)
	sentCount := 0
	maxRequests := 2 // Hard limit to 2 requests

	for i, profileURL := range profiles {
		if sentCount >= maxRequests {
			log.Println("Request limit reached (2)")
			break
		}

		log.Printf("[%d/%d] Processing: %s", i+1, len(profiles), profileURL)

		// Random delay between profiles
		stealth.RandomDelay(3000, 5000)

		// Send connection request
		if err := messaging.SendConnectionRequest(page, profileURL, cfg.ConnectionNote); err != nil {
			log.Printf("Failed to send request: %v", err)
			continue
		}

		sentCount++
		log.Printf("✓ Connection request sent (%d/%d)", sentCount, maxRequests)

		// Delay after sending
		stealth.RandomDelay(5000, 8000)
	}

	log.Printf("✅ Automation complete. Sent %d connection requests.", sentCount)
	log.Println("Keeping browser open for 10 seconds...")
	time.Sleep(10 * time.Second)
}
