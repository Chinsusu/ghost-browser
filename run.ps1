#!/usr/bin/env pwsh

Write-Host "=== Ghost Browser Launcher ===" -ForegroundColor Green
Write-Host ""

# Check available executables in priority order
if (Test-Path "ghost-browser-api-v3.exe") {
    Write-Host "🚀 Starting Ghost Browser API v3 (Latest)..." -ForegroundColor Yellow
    Write-Host "✅ Features: Full API, Profile Management, ChromeDP Integration" -ForegroundColor Green
    Write-Host "🌐 Access: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "Press Ctrl+C to stop the application" -ForegroundColor Gray
    Write-Host ""
    
    .\ghost-browser-api-v3.exe
} elseif (Test-Path "ghost-browser-chromedp-demo.exe") {
    Write-Host "🧪 Starting Ghost Browser ChromeDP Demo..." -ForegroundColor Yellow
    Write-Host "✅ Features: Fingerprint Spoofing Demo, browserleaks.com test" -ForegroundColor Green
    Write-Host "Press Ctrl+C to stop the application" -ForegroundColor Gray
    Write-Host ""
    
    .\ghost-browser-chromedp-demo.exe
} elseif (Test-Path "ghost-browser-backend.exe") {
    Write-Host "🔧 Starting Ghost Browser Backend..." -ForegroundColor Yellow
    Write-Host "Note: This is backend-only mode. For full UI, install Node.js and run 'wails dev'" -ForegroundColor Cyan
    Write-Host "Press Ctrl+C to stop the application" -ForegroundColor Gray
    Write-Host ""
    
    .\ghost-browser-backend.exe
} elseif (Test-Path "build/bin/Ghost Browser.exe") {
    Write-Host "🖥️ Starting Ghost Browser (Desktop)..." -ForegroundColor Yellow
    & "build/bin/Ghost Browser.exe"
} else {
    Write-Host "❌ Ghost Browser executable not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Available options:" -ForegroundColor Yellow
    Write-Host "  1. Run Ghost Browser v4 (Advanced):" -ForegroundColor White
    Write-Host "     cd ghost-browser-v4; .\run.ps1" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  2. Build from source:" -ForegroundColor White
    Write-Host "     .\build_and_test.ps1" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  3. Run specific version:" -ForegroundColor White
    Write-Host "     .\ghost-browser-api-v3.exe" -ForegroundColor Cyan
    Write-Host "     .\ghost-browser-chromedp-demo.exe" -ForegroundColor Cyan
    exit 1
}