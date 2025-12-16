package search

import (
	"fmt"
	"linkedin-automation/stealth"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func SearchProfiles(page *rod.Page, jobTitle string, maxProfiles int) []string {
	var profiles []string
	visited := make(map[string]bool)

	// Navigate to search page
	searchURL := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s", strings.ReplaceAll(jobTitle, " ", "+"))
	log.Printf("Searching: %s", searchURL)

	err := page.Navigate(searchURL)
	if err != nil {
		log.Printf("Navigation error: %v", err)
		return profiles
	}

	page.MustWaitLoad()
	stealth.RandomDelay(5000, 7000)

	// Scroll using WaitIdle instead of eval
	page.Mouse.MustScroll(0, 500)
	time.Sleep(2 * time.Second)

	// Find all profile links
	elements, err := page.Elements("a[href*='/in/']")
	if err != nil {
		log.Printf("Error finding profile links: %v", err)
		return profiles
	}

	log.Printf("Found %d potential profile links", len(elements))

	for _, el := range elements {
		if len(profiles) >= maxProfiles {
			break
		}

		href, err := el.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		profileURL := *href

		// Clean URL and check if valid profile
		if !strings.Contains(profileURL, "/in/") || strings.Contains(profileURL, "/company/") {
			continue
		}

		// Make sure it's a full URL
		if !strings.HasPrefix(profileURL, "http") {
			profileURL = "https://www.linkedin.com" + profileURL
		}

		// Remove query parameters
		if idx := strings.Index(profileURL, "?"); idx != -1 {
			profileURL = profileURL[:idx]
		}

		// Deduplicate
		if visited[profileURL] {
			continue
		}

		visited[profileURL] = true
		profiles = append(profiles, profileURL)

		log.Printf("Found profile %d: %s", len(profiles), profileURL)
	}

	log.Printf("Collected %d profiles total", len(profiles))
	return profiles
}
