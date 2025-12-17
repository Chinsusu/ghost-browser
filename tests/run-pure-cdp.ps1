#!/usr/bin/env pwsh

Write-Host "PURE CDP BROWSER TEST (No Rod Library)" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green
Write-Host ""

Write-Host "ANTIVIRUS BYPASS SOLUTION:" -ForegroundColor Yellow
Write-Host "  - Completely removes Rod library dependency" -ForegroundColor Green
Write-Host "  - Uses Edge --user-script parameter for injection" -ForegroundColor Green
Write-Host "  - Script runs at startup (before any page loads)" -ForegroundColor Green
Write-Host "  - No external executables = No antivirus blocking" -ForegroundColor Green
Write-Host ""

Write-Host "Building pure CDP test..." -ForegroundColor Yellow
go build -o test-pure-cdp.exe test-pure-cdp.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "Build successful" -ForegroundColor Green

Write-Host ""
Write-Host "Running PURE CDP test..." -ForegroundColor Yellow
Write-Host "This will:" -ForegroundColor Gray
Write-Host "  1. Create test profile with random fingerprint" -ForegroundColor Gray
Write-Host "  2. Create JavaScript spoofing script" -ForegroundColor Gray
Write-Host "  3. Launch Edge with --user-script parameter" -ForegroundColor Green
Write-Host "  4. Open browserleaks.com and creepjs.com" -ForegroundColor Gray
Write-Host "  5. Verify fingerprint spoofing works" -ForegroundColor Gray
Write-Host ""

.\test-pure-cdp.exe

Write-Host ""
Write-Host "Cleaning up..." -ForegroundColor Yellow
Remove-Item -Path "test-pure-cdp.exe" -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "PURE CDP TEST COMPLETED" -ForegroundColor Green
Write-Host ""
Write-Host "CRITICAL VERIFICATION:" -ForegroundColor Cyan
Write-Host "Check browserleaks.com/javascript for spoofed values:" -ForegroundColor Yellow
Write-Host "  - hardwareConcurrency should NOT be your real CPU cores" -ForegroundColor White
Write-Host "  - deviceMemory should NOT be your real RAM" -ForegroundColor White
Write-Host "  - Console should show [Ghost Browser] messages" -ForegroundColor White
Write-Host ""
Write-Host "If spoofing works: This approach can replace Rod library!" -ForegroundColor Green
Write-Host "If spoofing fails: Need to investigate Edge --user-script support" -ForegroundColor Red