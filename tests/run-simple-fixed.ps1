#!/usr/bin/env pwsh

Write-Host "SIMPLE FIXED BROWSER TEST" -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green
Write-Host ""

Write-Host "CRITICAL FIX APPLIED:" -ForegroundColor Yellow
Write-Host "  - Updated spoof.go to use AddScriptToEvaluateOnNewDocument" -ForegroundColor Green
Write-Host "  - Script now runs BEFORE page load (pre-load injection)" -ForegroundColor Green
Write-Host "  - Added automatic test page opening" -ForegroundColor Green
Write-Host ""

Write-Host "Building simple fixed test..." -ForegroundColor Yellow
go build -o test-simple-fixed.exe test-simple-fixed.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "Build successful" -ForegroundColor Green

Write-Host ""
Write-Host "Running FIXED browser test..." -ForegroundColor Yellow
Write-Host "This will:" -ForegroundColor Gray
Write-Host "  1. Create test profile with random fingerprint" -ForegroundColor Gray
Write-Host "  2. Launch Edge with FIXED pre-load injection" -ForegroundColor Green
Write-Host "  3. Automatically open browserleaks.com and creepjs.com" -ForegroundColor Gray
Write-Host "  4. Run for 60 seconds for testing" -ForegroundColor Gray
Write-Host ""

.\test-simple-fixed.exe

Write-Host ""
Write-Host "Cleaning up..." -ForegroundColor Yellow
Remove-Item -Path "test-simple-fixed.exe" -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "SIMPLE FIXED TEST COMPLETED" -ForegroundColor Green
Write-Host ""
Write-Host "KEY QUESTION:" -ForegroundColor Cyan
Write-Host "Are the fingerprint values now spoofed correctly on browserleaks.com?" -ForegroundColor Yellow
Write-Host ""
Write-Host "If YES: The fix works and can be integrated into main app" -ForegroundColor Green
Write-Host "If NO: Need to investigate further or use alternative approach" -ForegroundColor Red