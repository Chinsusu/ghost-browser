package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("🧪 === SIMPLE ChromeDP Test ===")
	fmt.Println()

	// Find Edge
	edgePath, err := findEdgePath()
	if err != nil {
		log.Fatal("❌ Edge not found:", err)
	}
	fmt.Println("✅ Found Edge:", edgePath)

	// Create temp profile
	tempDir, err := os.MkdirTemp("", "ghost-simple-*")
	if err != nil {
		log.Fatal("❌ Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir)

	// Setup Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(edgePath),
		chromedp.UserDataDir(tempDir),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Disable logging to reduce noise
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Simple spoofing script
	spoofScript := `
(function() {
	console.log('[Ghost] Starting spoofing...');
	
	Object.defineProperty(Navigator.prototype, 'hardwareConcurrency', {
		get: () => 8,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'deviceMemory', {
		get: () => 16,
		configurable: true
	});
	
	Object.defineProperty(Navigator.prototype, 'webdriver', {
		get: () => false,
		configurable: true
	});
	
	console.log('[Ghost] ✅ Spoofing active - hardwareConcurrency:', navigator.hardwareConcurrency, 'deviceMemory:', navigator.deviceMemory);
})();
`

	fmt.Println("🚀 Launching browser...")
	fmt.Println("⚠️  Check the browser window manually!")
	fmt.Println()

	// Launch with pre-load script
	err = chromedp.Run(ctx,
		// Add script BEFORE navigation
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(spoofScript).Do(ctx)
			if err != nil {
				fmt.Printf("⚠️  Script injection failed: %v\n", err)
			} else {
				fmt.Println("✅ Script injected successfully!")
			}
			return err
		}),
		
		// Navigate to test page
		chromedp.Navigate("https://browserleaks.com/javascript"),
	)

	if err != nil {
		log.Printf("❌ Failed: %v", err)
	} else {
		fmt.Println("✅ Browser launched and navigated!")
	}

	fmt.Println()
	fmt.Println("📋 MANUAL VERIFICATION:")
	fmt.Println("1. Check the Edge browser window")
	fmt.Println("2. On browserleaks.com/javascript page:")
	fmt.Println("   - hardwareConcurrency should be 8 (NOT 16)")
	fmt.Println("   - deviceMemory should be 16 (NOT 8)")
	fmt.Println("3. Open DevTools (F12) → Console tab")
	fmt.Println("4. Look for:")
	fmt.Println("   [Ghost] Starting spoofing...")
	fmt.Println("   [Ghost] ✅ Spoofing active - hardwareConcurrency: 8 deviceMemory: 16")
	fmt.Println()
	fmt.Println("⏳ Browser will stay open for 60 seconds...")
	
	// Keep browser open for manual verification
	time.Sleep(60 * time.Second)

	fmt.Println("\n🔍 RESULTS:")
	fmt.Println("Did you see spoofed values on browserleaks.com?")
	fmt.Println("- hardwareConcurrency = 8? (Yes/No)")
	fmt.Println("- deviceMemory = 16? (Yes/No)")
	fmt.Println("- Console messages visible? (Yes/No)")
}

func findEdgePath() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("Windows only")
	}

	paths := []string{
		filepath.Join(os.Getenv("ProgramFiles(x86)"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join(os.Getenv("LocalAppData"), "Microsoft", "Edge", "Application", "msedge.exe"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	if path, err := exec.LookPath("msedge"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("not found")
}