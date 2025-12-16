package auth

import (
	"errors"
	"linkedin-automation/stealth"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func Login(page *rod.Page, email, password string) error {
	if email == "" || password == "" {
		return errors.New("email and password are required")
	}

	log.Println("Waiting for login form...")
	time.Sleep(2 * time.Second)

	// Find and fill email field
	emailInput, err := page.Timeout(10 * time.Second).Element("#username")
	if err != nil {
		return errors.New("could not find email input field")
	}
	log.Println("Typing email...")
	stealth.HumanType(emailInput, email)
	stealth.RandomDelay(1000, 2000)

	// Find and fill password field
	passwordInput, err := page.Timeout(10 * time.Second).Element("#password")
	if err != nil {
		return errors.New("could not find password input field")
	}
	log.Println("Typing password...")
	stealth.HumanType(passwordInput, password)
	stealth.RandomDelay(1000, 2000)

	// Click login
	stealth.RandomMouseMove(page)
	stealth.RandomDelay(500, 1000)

	loginBtn, err := page.Timeout(10 * time.Second).Element("button[type='submit']")
	if err != nil {
		return errors.New("could not find login button")
	}
	log.Println("Clicking login...")
	loginBtn.MustClick()

	time.Sleep(5 * time.Second)
	currentURL := page.MustInfo().URL
	log.Printf("Current URL: %s", currentURL)

	// Handle checkpoint/verification
	if strings.Contains(currentURL, "checkpoint") || strings.Contains(currentURL, "challenge") {
		log.Println("")
		log.Println("═══════════════════════════════════════════")
		log.Println("⚠️  MANUAL VERIFICATION REQUIRED")
		log.Println("═══════════════════════════════════════════")
		log.Println("LinkedIn requires verification.")
		log.Println("Please complete it in the browser window.")
		log.Println("You have 2 MINUTES...")
		log.Println("")

		// Wait 2 minutes for verification
		for i := 0; i < 120; i++ {
			time.Sleep(1 * time.Second)
			currentURL = page.MustInfo().URL

			if !strings.Contains(currentURL, "checkpoint") && !strings.Contains(currentURL, "challenge") {
				log.Println("✓ Verification completed!")
				break
			}

			if i%15 == 0 && i > 0 {
				log.Printf("Still waiting... (%d seconds)", i)
			}
		}
	}

	time.Sleep(3 * time.Second)

	// Check success
	currentURL = page.MustInfo().URL
	if strings.Contains(currentURL, "feed") || strings.Contains(currentURL, "mynetwork") {
		log.Println("✓ Login successful!")
		return nil
	}

	// Try to find any success indicator
	_, err = page.Timeout(5 * time.Second).Element("nav")
	if err == nil {
		log.Println("✓ Login successful!")
		return nil
	}

	log.Println("Please check browser and complete any verification")
	time.Sleep(30 * time.Second)
	return nil
}
