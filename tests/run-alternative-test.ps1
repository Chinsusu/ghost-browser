#!/usr/bin/env pwsh

Write-Host "🧪 === Alternative Browser Test (No Rod Library) ===" -ForegroundColor Green
Write-Host ""

Write-Host "This test bypasses the Rod library antivirus issue by using:" -ForegroundColor Yellow
Write-Host "  • Direct Edge launch with CDP enabled" -ForegroundColor Gray
Write-Host "  • JavaScript injection via startup scripts" -ForegroundColor Gray
Write-Host "  • Manual CDP communication (if available)" -ForegroundColor Gray
Write-Host ""

Write-Host "🔧 Building alternative test executable..." -ForegroundColor Yellow
go build -o test-alternative-browser.exe test-alternative-browser.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Test executable built successfully" -ForegroundColor Green

Write-Host "`n🚀 Running alternative browser test..." -ForegroundColor Yellow
Write-Host "This will:" -ForegroundColor Gray
Write-Host "  1. Create a test profile with random fingerprint" -ForegroundColor Gray
Write-Host "  2. Launch Edge with CDP and injection scripts" -ForegroundColor Gray
Write-Host "  3. Attempt fingerprint spoofing without Rod" -ForegroundColor Gray
Write-Host "  4. Open test pages for verification" -ForegroundColor Gray
Write-Host ""

.\test-alternative-browser.exe

Write-Host "`n🧹 Cleaning up..." -ForegroundColor Yellow
Remove-Item -Path "test-alternative-browser.exe" -ErrorAction SilentlyContinue

Write-Host "Alternative test completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Please report the results:" -ForegroundColor Cyan
Write-Host "1. Did Edge browser open successfully?" -ForegroundColor White
Write-Host "2. Were fingerprint values spoofed correctly?" -ForegroundColor White
Write-Host "3. Did you see console messages in DevTools?" -ForegroundColor White
Write-Host "4. Any errors in the terminal?" -ForegroundColor White
Write-Host "5. Did the CDP injection work?" -ForegroundColor White