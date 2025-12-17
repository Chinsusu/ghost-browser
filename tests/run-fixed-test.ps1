#!/usr/bin/env pwsh

Write-Host "🧪 === FIXED CDP Browser Test (Pre-load Injection) ===" -ForegroundColor Green
Write-Host ""

Write-Host "🔧 CRITICAL FIX APPLIED:" -ForegroundColor Yellow
Write-Host "  • Changed from Runtime.evaluate (post-load)" -ForegroundColor Red
Write-Host "  • To Page.addScriptToEvaluateOnNewDocument (pre-load)" -ForegroundColor Green
Write-Host "  • Script now runs BEFORE page reads navigator values" -ForegroundColor Green
Write-Host ""

Write-Host "🔧 Building FIXED test executable..." -ForegroundColor Yellow
go build -o test-cdp-fixed.exe test-cdp-fixed.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ FIXED test executable built successfully" -ForegroundColor Green

Write-Host "`n🚀 Running FIXED CDP test..." -ForegroundColor Yellow
Write-Host "This will:" -ForegroundColor Gray
Write-Host "  1. Create a test profile with random fingerprint" -ForegroundColor Gray
Write-Host "  2. Launch Edge with CDP enabled" -ForegroundColor Gray
Write-Host "  3. Add spoofing script to NEW DOCUMENT evaluation" -ForegroundColor Green
Write-Host "  4. Navigate to test pages (script runs BEFORE page load)" -ForegroundColor Green
Write-Host "  5. Verify fingerprint values are spoofed" -ForegroundColor Gray
Write-Host ""

.\test-cdp-fixed.exe

Write-Host "`n🧹 Cleaning up..." -ForegroundColor Yellow
Remove-Item -Path "test-cdp-fixed.exe" -ErrorAction SilentlyContinue

Write-Host "FIXED CDP test completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 CRITICAL CHECK - Please verify:" -ForegroundColor Cyan
Write-Host "1. Did Edge browser open successfully?" -ForegroundColor White
Write-Host "2. Was CDP connection established?" -ForegroundColor White
Write-Host "3. On browserleaks.com/javascript:" -ForegroundColor White
Write-Host "   - hardwareConcurrency = spoofed value (NOT your real CPU cores)?" -ForegroundColor Yellow
Write-Host "   - deviceMemory = spoofed value (NOT your real RAM)?" -ForegroundColor Yellow
Write-Host "4. Did you see console messages in DevTools?" -ForegroundColor White
Write-Host "5. Any errors in the terminal?" -ForegroundColor White
Write-Host ""
Write-Host "💡 If fingerprint values are now spoofed correctly," -ForegroundColor Green
Write-Host "   this fix can be integrated into the main app!" -ForegroundColor Green