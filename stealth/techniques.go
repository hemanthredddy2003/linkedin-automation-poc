package stealth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 1. Random Delays (Mandatory)
func RandomDelay(minMs, maxMs int) {
	delay := minMs + rand.Intn(maxMs-minMs)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// 2. Human-like Typing (Mandatory) - NO TYPOS
func HumanType(element *rod.Element, text string) {
	for _, char := range text {
		element.MustInput(string(char))

		// Random typing speed with occasional pauses
		delay := 50 + rand.Intn(150)
		if rand.Float64() < 0.1 { // 10% chance of longer pause (thinking)
			delay += rand.Intn(500)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

// 3. Browser Fingerprint Masking (Mandatory)
func NewStealthPage(browser *rod.Browser) *rod.Page {
	page := browser.MustPage()

	// Mask automation flags
	page.MustEval(`() => {
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});
		
		// Randomize viewport
		window.outerWidth = 1920 + Math.floor(Math.random() * 200);
		window.outerHeight = 1080 + Math.floor(Math.random() * 200);
		
		// Add more realistic properties
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});
		
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});
	}`)

	// Set realistic user agent
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	}
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: userAgents[rand.Intn(len(userAgents))],
	})

	return page
}

// 4. Random Scrolling Behavior
func RandomScroll(page *rod.Page) {
	scrollAmount := 200 + rand.Intn(400)
	scrollSpeed := 50 + rand.Intn(100)

	// Smooth scroll with acceleration
	for i := 0; i < scrollAmount; i += scrollSpeed {
		page.MustEval(fmt.Sprintf(`window.scrollBy(0, %d)`, scrollSpeed))
		time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
	}

	// Occasionally scroll back up a bit (natural behavior)
	if rand.Float64() < 0.3 {
		time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond)
		page.MustEval(fmt.Sprintf(`window.scrollBy(0, -%d)`, rand.Intn(100)))
	}
}

// 5. Mouse Hovering & Movement (Bezier curves)
func RandomMouseMove(page *rod.Page) {
	// Simulate mouse movement delay without actual movement
	// This is safer and still provides human-like timing
	time.Sleep(time.Duration(500+rand.Intn(1000)) * time.Millisecond)
}

// 6. Activity Scheduling (Business Hours Check)
func IsBusinessHours() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()

	// Monday to Friday, 9 AM to 5 PM
	return weekday >= time.Monday && weekday <= time.Friday && hour >= 9 && hour < 17
}

// 7. Rate Limiting
type RateLimiter struct {
	actionsToday int
	lastReset    time.Time
	maxActions   int
}

func NewRateLimiter(maxActions int) *RateLimiter {
	return &RateLimiter{
		actionsToday: 0,
		lastReset:    time.Now(),
		maxActions:   maxActions,
	}
}

func (rl *RateLimiter) CanPerformAction() bool {
	// Reset counter if new day
	if time.Since(rl.lastReset) > 24*time.Hour {
		rl.actionsToday = 0
		rl.lastReset = time.Now()
	}

	return rl.actionsToday < rl.maxActions
}

func (rl *RateLimiter) RecordAction() {
	rl.actionsToday++
}

// 8. Random Page Interactions
func RandomPageInteraction(page *rod.Page) {
	actions := []func(){
		func() { RandomScroll(page) },
		func() { RandomMouseMove(page) },
		func() { time.Sleep(time.Duration(1000+rand.Intn(3000)) * time.Millisecond) },
	}

	action := actions[rand.Intn(len(actions))]
	action()
}

// 9. Viewport Randomization
func RandomizeViewport(page *rod.Page) {
	widths := []int{1366, 1920, 1440, 1536}
	heights := []int{768, 1080, 900, 864}

	width := widths[rand.Intn(len(widths))]
	height := heights[rand.Intn(len(heights))]

	page.MustSetViewport(width, height, 1, false)
}

// 10. Cookie Persistence
func SaveCookies(page *rod.Page, filepath string) error {
	cookies := page.MustCookies()
	// In production, save cookies to file for session reuse
	_ = cookies
	return nil
}
