#!/usr/bin/env pwsh

Write-Host "Ghost Browser FINAL FIXED TEST" -ForegroundColor Green
Write-Host ""

Write-Host "CRITICAL FIXES APPLIED:" -ForegroundColor Yellow
Write-Host "  - Changed injection timing: POST-load to PRE-load" -ForegroundColor Green
Write-Host "  - Updated spoof.go: MustEvaluate to AddScriptToEvaluateOnNewDocument" -ForegroundColor Green
Write-Host "  - Added automatic test page opening" -ForegroundColor Green
Write-Host "  - Created fallback CDP direct approach" -ForegroundColor Green
Write-Host ""

Write-Host "Available Tests:" -ForegroundColor Cyan
Write-Host "  1. Fixed Rod-based test (may be blocked by antivirus)" -ForegroundColor White
Write-Host "  2. Fixed CDP direct test (bypasses antivirus)" -ForegroundColor Green
Write-Host ""

$choice = Read-Host "Choose test (1 for Rod, 2 for CDP, or Enter for both)"

if ($choice -eq "1" -or $choice -eq "") {
    Write-Host ""
    Write-Host "RUNNING: Fixed Rod-based Test" -ForegroundColor Cyan
    Write-Host "====================================" -ForegroundColor Cyan
    
    go build -o test-rod-fixed.exe test-rod-fixed.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Rod test built successfully" -ForegroundColor Green
        .\test-rod-fixed.exe
        Remove-Item -Path "test-rod-fixed.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "Rod test build failed" -ForegroundColor Red
    }
    
    if ($choice -eq "") {
        Write-Host ""
        Read-Host "Press Enter to continue to CDP test"
    }
}

if ($choice -eq "2" -or $choice -eq "") {
    Write-Host ""
    Write-Host "RUNNING: Fixed CDP Direct Test" -ForegroundColor Cyan
    Write-Host "==============================" -ForegroundColor Cyan
    
    go build -o test-cdp-fixed.exe test-cdp-fixed.go
    if ($LASTEXITCODE -eq 0) {
        Write-Host "CDP test built successfully" -ForegroundColor Green
        .\test-cdp-fixed.exe
        Remove-Item -Path "test-cdp-fixed.exe" -ErrorAction SilentlyContinue
    } else {
        Write-Host "CDP test build failed" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "FINAL TEST SUMMARY" -ForegroundColor Green
Write-Host "==================" -ForegroundColor Green

Write-Host ""
Write-Host "CRITICAL VERIFICATION CHECKLIST:" -ForegroundColor Cyan
Write-Host "1. Did Edge browser open successfully? ___" -ForegroundColor White
Write-Host "2. On browserleaks.com/javascript:" -ForegroundColor White
Write-Host "   - hardwareConcurrency = spoofed value (NOT real CPU)? ___" -ForegroundColor Yellow
Write-Host "   - deviceMemory = spoofed value (NOT real RAM)? ___" -ForegroundColor Yellow
Write-Host "3. Console messages in DevTools (F12):" -ForegroundColor White
Write-Host "   - [Ghost Browser] Navigator spoofed? ___" -ForegroundColor Yellow
Write-Host "   - [Ghost Browser] Fingerprint spoofing active? ___" -ForegroundColor Yellow
Write-Host "4. Any errors in terminal? ___" -ForegroundColor White
Write-Host "5. Which test worked better? (Rod/CDP) ___" -ForegroundColor White

Write-Host ""
Write-Host "NEXT STEPS:" -ForegroundColor Green
Write-Host "If fingerprint spoofing now works correctly:" -ForegroundColor White
Write-Host "  - Integration into main Ghost Browser app" -ForegroundColor Green
Write-Host "  - Replace old injection method with fixed version" -ForegroundColor Green
Write-Host "  - Choose Rod (if not blocked) or CDP direct approach" -ForegroundColor Green

Write-Host ""
Write-Host "Testing complete! Report results above." -ForegroundColor Green