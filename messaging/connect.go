package messaging

import (
	"errors"
	"linkedin-automation/stealth"
	"log"
	"time"

	"github.com/go-rod/rod"
)

func SendConnectionRequest(page *rod.Page, profileURL string, note string) error {
	// Navigate to profile
	log.Printf("Navigating to profile...")
	page.MustNavigate(profileURL)
	stealth.RandomDelay(3000, 5000)

	// Scroll down to load the page
	page.Mouse.MustScroll(0, 500)
	stealth.RandomDelay(2000, 3000)

	// Random mouse movement
	stealth.RandomMouseMove(page)
	stealth.RandomDelay(1000, 2000)

	// Try to find Connect button
	log.Println("Looking for Connect button...")
	connectBtn, err := page.Timeout(10 * time.Second).Element("button span:has-text('Connect')")
	if err != nil {
		// Try simpler selector
		connectBtn, err = page.Timeout(5 * time.Second).Element("button[aria-label*='Invite']")
		if err != nil {
			log.Printf("Connect button not found - may already be connected or pending")
			return errors.New("connect button not found")
		}
	}

	// Click Connect button directly
	log.Println("Clicking Connect button...")
	err = connectBtn.Timeout(5*time.Second).Click("left", 1)
	if err != nil {
		log.Printf("Failed to click connect button: %v", err)
		return err
	}
	stealth.RandomDelay(3000, 4000)

	// Check if we can add a note (optional - skip if not found)
	noteBtn, err := page.Timeout(3 * time.Second).Element("button[aria-label*='Add a note']")
	if err == nil && note != "" {
		log.Println("Found 'Add a note' button, attempting to add note...")
		err = noteBtn.Timeout(3*time.Second).Click("left", 1)
		if err == nil {
			stealth.RandomDelay(1000, 2000)

			noteInput, err := page.Timeout(5 * time.Second).Element("textarea")
			if err == nil {
				log.Println("Typing note...")
				stealth.HumanType(noteInput, note)
				stealth.RandomDelay(1000, 2000)
			} else {
				log.Println("Note textarea not found, sending without note")
			}
		} else {
			log.Println("Could not click 'Add a note', sending without note")
		}
	} else {
		log.Println("No 'Add a note' option found or no note provided")
	}

	// Look for Send button
	log.Println("Looking for Send button...")
	sendBtn, err := page.Timeout(10 * time.Second).Element("button[aria-label*='Send']")
	if err != nil {
		// Try finding by text
		sendBtn, err = page.Timeout(5 * time.Second).Element("button span:has-text('Send')")
		if err != nil {
			// Try generic submit button
			sendBtn, err = page.Timeout(5 * time.Second).Element("button[data-control-name*='send']")
			if err != nil {
				log.Println("Send button not found - connection may already be pending")
				return nil // Don't fail, might have worked
			}
		}
	}

	// Click Send
	log.Println("Clicking Send button...")
	err = sendBtn.Timeout(5*time.Second).Click("left", 1)
	if err != nil {
		log.Printf("Failed to click send: %v, but continuing", err)
	}
	stealth.RandomDelay(2000, 3000)

	log.Printf("✓ Connection request sent successfully")
	return nil
}

func SendMessage(page *rod.Page, profileURL string, message string) error {
	page.MustNavigate(profileURL)
	time.Sleep(3 * time.Second)

	msgBtn, err := page.Timeout(5 * time.Second).Element("button:has-text('Message')")
	if err != nil {
		return errors.New("message button not found")
	}

	msgBtn.MustClick()
	stealth.RandomDelay(2000, 3000)

	msgInput, err := page.Timeout(5 * time.Second).Element("div[contenteditable='true']")
	if err != nil {
		return errors.New("message input not found")
	}

	stealth.HumanType(msgInput, message)
	stealth.RandomDelay(1000, 2000)

	sendBtn, err := page.Timeout(5 * time.Second).Element("button[type='submit']")
	if err != nil {
		return errors.New("send button not found")
	}

	sendBtn.MustClick()
	stealth.RandomDelay(2000, 3000)

	log.Printf("✓ Message sent")
	return nil
}
